// Copyright (c) 2023, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package armaze represents an N-armed maze ("bandit")
// with each Arm having a distinctive CS stimulus at the start
// (could be one of multiple possibilities) and (some probability of)
// a US outcome at the end of the maze (could be either positive
// or negative, with (variable) magnitude and probability.
//
// The maze can have open or closed arms -- open arms allow
// switching to a neighboring arm anytime, while closed arms
// only allow switching at the start.
package armaze

import (
	"github.com/emer/emergent/econfig"
	"github.com/emer/emergent/env"
	"github.com/emer/emergent/erand"
	"github.com/emer/etable/etensor"
	"github.com/emer/etable/minmax"
	"github.com/goki/ki/kit"
)

// Actions is a list of mutually exclusive states
// for tracing the behavior and internal state of Emery
type Actions int

const (
	Forward Actions = iota
	Left
	Right
	Consume
	None
	ActionsN
)

//go:generate stringer -type=Actions

var KiT_Actions = kit.Enums.AddEnum(ActionsN, kit.NotBitFlag, nil)

// General note on US / Drive indexes:
// The env does _not_ represent any built-in drives or USs (curiosity, effort, urgency)
// 0 = start of the sim-specific USs and Drives

// Env implements an N-armed maze ("bandit")
// with each Arm having a distinctive CS stimulus visible at the start
// (could be one of multiple possibilities) and (some probability of)
// a US outcome at the end of the maze (could be either positive
// or negative, with (variable) magnitude and probability.
type Env struct {

	// name of environment -- Train or Test
	Nm string `desc:"name of environment -- Train or Test"`

	// configuration parameters
	Config Config `desc:"configuration parameters"`

	// current drive strength for each of Config.NDrives in normalized 0-1 units of each drive: 0 = first sim drive, not curiosity
	Drives []float32 `desc:"current drive strength for each of Config.NDrives in normalized 0-1 units of each drive: 0 = first sim drive, not curiosity"`

	// parameters associated with each US.  The first NDrives are positive USs, and beyond that are negative USs
	USs []*USParams `desc:"parameters associated with each US.  The first NDrives are positive USs, and beyond that are negative USs"`

	// state of each arm: dist, effort, US, CS
	Arms []*Arm `desc:"state of each arm: dist, effort, US, CS"`

	// arm-wise location: either facing (Pos=0) or in (Pos > 0)
	Arm int `inactive:"+" desc:"arm-wise location: either facing (Pos=0) or in (Pos > 0)"`

	// current position in the Arm: 0 = at start looking in, otherwise at given distance into the arm
	Pos int `inactive:"+" desc:"current position in the Arm: 0 = at start looking in, otherwise at given distance into the arm"`

	// current target drive, in paradigms where that is used (e.g., Approach)
	TrgDrive int `inactive:"+" desc:"current target drive, in paradigms where that is used (e.g., Approach)"`

	// Current US being consumed -- is -1 unless being consumed
	USConsumed int `inactive:"+" desc:"Current US being consumed -- is -1 unless being consumed"`

	// reward or punishment value generated by the current US being consumed
	USValue float32 `inactive:"+" desc:"reward or punishment value generated by the current US being consumed"`

	// last action taken
	LastAct int `inactive:"+" desc:"last action taken"`

	// effort on current trial
	Effort float32 `inactive:"+" desc:"effort on current trial"`

	// last CS -- previous trial
	LastUS int `inactive:"+" desc:"last CS -- previous trial"`

	// true if looking at correct CS for first time
	ShouldGate bool `inactive:"+" desc:"true if looking at correct CS for first time"`

	// just gated on this trial -- set by sim-- used for instinct
	JustGated bool `inactive:"+" desc:"just gated on this trial -- set by sim-- used for instinct"`

	// has gated at some point during sequence -- set by sim -- used for instinct
	HasGated bool `inactive:"+" desc:"has gated at some point during sequence -- set by sim -- used for instinct"`

	// named states -- e.g., USs, CSs, etc
	States map[string]*etensor.Float32 `desc:"named states -- e.g., USs, CSs, etc"`

	// [view: -] random number generator for the env -- all random calls must use this
	Rand erand.SysRand `view:"-" desc:"random number generator for the env -- all random calls must use this"`

	// random seed
	RndSeed int64 `inactive:"+" desc:"random seed"`
}

const noUS = -1

func (ev *Env) Name() string {
	return ev.Nm
}

func (ev *Env) Desc() string {
	return "Maze"
}

// Defaults sets default params
func (ev *Env) Defaults() {
	ev.Config.Defaults()
	econfig.SetFromDefaults(&ev.Config)
	ev.Config.Update()
}

// ConfigEnv configures the environment.
// additional parameterization via specific configs
// is applied after this step, which initializes
// everything according to basic Ns
func (ev *Env) ConfigEnv() {
	if ev.Rand.Rand == nil {
		ev.Rand.NewRand(ev.RndSeed)
	} else {
		ev.Rand.Seed(ev.RndSeed)
	}

	switch ev.Config.Paradigm {
	case Approach:
		ev.ConfigApproach()
	}

	ev.Config.Update()

	ev.Drives = make([]float32, ev.Config.NDrives)
	ev.USs = make([]*USParams, ev.Config.NUSs)
	ev.Arms = make([]*Arm, ev.Config.NDrives)

	// defaults
	usperm := ev.Rand.Perm(ev.NUSs, -1)
	for i, arm := range ev.Arms {
		arm.Dist = 4
		arm.Effort.Set(1, 1)
		arm.US = usperm[i%ev.NUSs]
	}

	// defaults
	for i, us := range ev.USs {
		if i < ev.NDrives {
			us.Negative = false
		} else {
			us.Negative = true
		}
		us.Mag.Set(1, 1)
		us.Prob = 1
		us.CSProbs = make([]float32, ev.NCSs)
		us.CSProbs[i%ev.NCSs] = 1 // uniform allocation
	}

	ev.States = make(map[string]*etensor.Float32)
	ev.States["CS"] = etensor.NewFloat32([]int{ev.NYReps, ev.Config.NCSs}, nil, nil)
	ev.States["Action"] = etensor.NewFloat32([]int{ev.NYReps, NActions}, nil, nil)

	ev.NewStart()
}

func (ev *Env) Validate() error {
	return nil
}

func (ev *Env) Init(run int) {
	ev.Config()
}

func (ev *Env) Counter(scale env.TimeScales) (cur, prv int, changed bool) {
	return 0, 0, false
}

func (ev *Env) State(el string) etensor.Tensor {
	return ev.States[el]
}

// MinMaxRand returns a random number in the range between Min and Max
func MinMaxRand(mm minmax.F32, rand erand.SysRand) {
	return mm.Min + rand.Float32(-1)*mm.Range()
}

// InactiveVal returns a new random inactive value from Config.Params.Inactive
// param range.
func (ev *Env) InactiveVal() float32 {
	return MinMaxRand(ev.Config.Params.Inactive, ev.Rand)
}

// ForwardEffort returns a new random Effort value from Arm Effort range
func (ev *Env) ForwardEffort(arm *Arm) float32 {
	return MinMaxRand(arm.Effort, ev.Rand)
}

// TurnEffort returns a new random Effort value from Config.Params.TurnEffort
// param range.
func (ev *Env) TurnEffort() float32 {
	return MinMaxRand(ev.Config.Params.TurnEffort, ev.Rand)
}

// ConsumeEffort returns a new random Effort value from Config.Params.ConsumeEffort
// param range.
func (ev *Env) ConsumeEffort() float32 {
	return MinMaxRand(ev.Config.Params.ConsumeEffort, ev.Rand)
}

// ChooseCSs selects new CSs for each Arm as function of US CSProbs
// This must be called
func (ev *Env) ChooseCSs() {
	for i, arm := range ev.Arms {
		us := ev.USs[arm.US]
		arm.CS = erand.PChoose32(us.CSProbs, -1, ev.Rand) // choose by dist
	}
}

// NewStart starts a new approach run
func (ev *Env) NewStart() {
	if ev.Config.Params.RandomStart {
		ev.Arm = ev.Rand.Intn(ev.Arms, -1)
	}
	ev.Pos = 0
	ev.JustGated = false
	ev.HasGated = false
	ev.ChooseCSs()

	switch ev.Config.Paradigm {
	case Approach:
		ev.StartApproach()
	}

	ev.USConsumed = -1
	ev.USValue = 0
	ev.RenderState()
}

// RenderLocalist renders one localist state
func (ev *Env) RenderLocalist(name string, val int) {
	st := ev.States[name]
	st.SetZeros()
	if val >= st.Dim(1) {
		return
	}
	for y := 0; y < ev.NYReps; y++ {
		st.Set([]int{y, val}, 1.0)
	}
}

// RenderLocalist4D renders one localist state in 4D
func (ev *Env) RenderLocalist4D(name string, val int) {
	st := ev.States[name]
	st.SetZeros()
	for y := 0; y < ev.NYReps; y++ {
		st.Set([]int{0, val, y, 0}, 1.0)
	}
}

// RenderState renders the current state
func (ev *Env) RenderState() {
	ev.RenderLocalist("CS", ev.CS)
}

// RenderAction renders the action
func (ev *Env) RenderAction(act Actions) {
	ev.RenderLocalist("Action", int(act))
}

// Step does one step.  it is up to the driving sim to decide when to call NewStart
func (ev *Env) Step() bool {
	ev.TakeAct(ev.LastAct)
	switch ev.Config.Paradigm {
	case Approach:
		ev.StepApproach()
	}
	ev.RenderState()
	return true
}

func (ev *Env) DecodeAct(vt *etensor.Float32) (int, string) {
	mxi := ev.DecodeLocalist(vt)
	return mxi, ev.Acts[mxi]
}

func (ev *Env) DecodeLocalist(vt *etensor.Float32) int {
	dx := vt.Dim(1)
	var max float32
	var mxi int
	for i := 0; i < dx; i++ {
		var sum float32
		for j := 0; j < ev.NYReps; j++ {
			sum += vt.Value([]int{j, i})
		}
		if sum > max {
			max = sum
			mxi = i
		}
	}
	return mxi
}

// Action records the LastAct and renders it, but does not
// update the state accordingly.
func (ev *Env) Action(action string, nop etensor.Tensor) {
	act := None
	act.FromString(action)
	ev.LastAct = act
	ev.RenderAction(act) // plus phase input is action
	// note: action not taken via TakeAct until start of trial in Step()
}

func (ev *Env) TakeAct(act Actions) {
	narms := ev.Config.NArms
	arm := ev.Arms[ev.Arm]
	switch act {
	case Forward:
		ev.Effort = ev.ForwardEffort(arm) // pay effort regardless
		npos := ev.Pos + 1
		if npos <= arm.Length {
			ev.Pos = npos
		} else {
			// todo: bump into wall?
		}
	case Left:
		ev.Effort = ev.TurnEffort() // pay effort regardless
		if ev.Config.Params.OpenArms || ev.Pos == 0 {
			ev.Arm--
		}
		if ev.Arm < 0 {
			ev.Arm += narms
		}
	case Right:
		ev.Effort = ev.TurnEffort() // pay effort regardless
		if ev.Config.Params.OpenArms || ev.Pos == 0 {
			ev.Arm++
		}
		if ev.Arm >= narms {
			ev.Arm += narms
		}
	case Consume:
		ev.Effort = ev.ConsumeEffort()
		if ev.Pos == arm.Length {
			ev.ConsumeUS(arm)
		}
	}
}

// ConsumeUS implements the consume action at current position in given arm
func (ev *Env) ConsumeUS(arm *Arm) {
	us := ev.USs[arm.US]
	mag := MinMaxRand(us.Mag, ev.Rand)
	got := erand.BoolP32(ev.Rand)
	if got {
		ev.USConsumed = arm.US
		ev.USValue = ev.Drives[arm.US] * mag
	} else {
		ev.USConsumed = -1
		ev.USValue = 0
	}
}

// USForPos returns the US at given position
func (ev *Env) USForPos() int {
	uss := ev.States["USs"]
	return int(uss.Values[ev.Pos])
}

// PosHasDriveUS returns true if the current USForPos corresponds
// to the current Drive -- i.e., are we looking at the right thing?a
func (ev *Env) PosHasDriveUS() bool {
	return ev.Drive == ev.USForPos()
}

// InstinctAct returns an "instinctive" action that implements a basic policy
func (ev *Env) InstinctAct(justGated, hasGated bool) int {
	ev.JustGated = justGated
	ev.HasGated = hasGated
	ev.ShouldGate = ((hasGated && ev.US != noUS) || // To clear the goal after US
		(!hasGated && ev.PosHasDriveUS())) // looking at correct, haven't yet gated

	if ev.Dist == 0 {
		return ev.ActMap["Consume"]
	}
	if ev.HasGated {
		return ev.ActMap["Forward"]
	}
	lt := ev.ActMap["Left"]
	rt := ev.ActMap["Right"]
	if ev.LastAct == lt || ev.LastAct == rt {
		return ev.LastAct
	}
	if ev.AlwaysLeft || erand.BoolP(.5, -1, &ev.Rand) {
		return lt
	}
	return rt
}
