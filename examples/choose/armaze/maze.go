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

//go:generate core generate -add-types

import (
	"github.com/emer/axon/v2/axon"
	"github.com/emer/emergent/v2/econfig"
	"github.com/emer/emergent/v2/env"
	"github.com/emer/emergent/v2/erand"
	"github.com/emer/etable/v2/etensor"
	"github.com/emer/etable/v2/minmax"
)

// Actions is a list of mutually exclusive states
// for tracing the behavior and internal state of Emery
type Actions int32 //enums:enum

const (
	Forward Actions = iota
	Left
	Right
	Consume
	None
)

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
	Nm string

	// our data parallel index
	Di int `edit:"-"`

	// configuration parameters
	Config Config

	// current drive strength for each of Config.NDrives in normalized 0-1 units of each drive: 0 = first sim drive, not curiosity
	Drives []float32

	// arm-wise location: either facing (Pos=0) or in (Pos > 0)
	Arm int `edit:"-"`

	// current position in the Arm: 0 = at start looking in, otherwise at given distance into the arm
	Pos int `edit:"-"`

	// distance from US
	Dist int `edit:"-"`

	// current integer time step since last NewStart
	Tick int `edit:"-"`

	// current target drive, in paradigms where that is used (e.g., Approach)
	TrgDrive int `edit:"-"`

	// Current US being consumed -- is -1 unless being consumed
	USConsumed int `edit:"-"`

	// reward or punishment value generated by the current US being consumed -- just the Magnitude of the US -- does NOT include any modulation by Drive
	USValue float32 `edit:"-"`

	// just finished consuming a US -- ready to start doing something new
	JustConsumed bool `edit:"-"`

	// arm(s) with maximum Drive * Mag * Prob US outcomes
	ArmsMaxValue []int `edit:"-"`

	// maximum value for ArmsMaxValue arms
	MaxValue float32 `edit:"-"`

	// arm(s) with maximum Value outcome discounted by Effort
	ArmsMaxUtil []int `edit:"-"`

	// maximum utility for ArmsMaxUtil arms
	MaxUtil float32 `edit:"-"`

	// arm(s) with negative US outcomes
	ArmsNeg []int `edit:"-"`

	// last action taken
	LastAct Actions `edit:"-"`

	// effort on current trial
	Effort float32 `edit:"-"`

	// last CS seen
	LastCS int `edit:"-"`

	// last US -- previous trial
	LastUS int `edit:"-"`

	// true if looking at correct CS for first time
	ShouldGate bool `edit:"-"`

	// just gated on this trial -- set by sim-- used for instinct
	JustGated bool `edit:"-"`

	// has gated at some point during sequence -- set by sim -- used for instinct
	HasGated bool `edit:"-"`

	// named states -- e.g., USs, CSs, etc
	States map[string]*etensor.Float32

	// maximum length of any arm
	MaxLength int `edit:"-"`

	// random number generator for the env -- all random calls must use this
	Rand erand.SysRand `view:"-"`

	// random seed
	RndSeed int64 `edit:"-"`
}

const noUS = -1

func (ev *Env) Name() string {
	return ev.Nm
}

func (ev *Env) Desc() string {
	return "N-Arm Maze Environment"
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
// takes the data parallel index di
func (ev *Env) ConfigEnv(di int) {
	ev.Di = di
	cfg := &ev.Config

	if ev.Rand.Rand == nil {
		ev.Rand.NewRand(ev.RndSeed)
	} else {
		ev.Rand.Seed(ev.RndSeed)
	}

	switch cfg.Paradigm {
	case Approach:
		ev.ConfigApproach()
	}

	cfg.Update()

	ev.Drives = make([]float32, cfg.NDrives)
	cfg.USs = make([]*USParams, cfg.NUSs)
	cfg.Arms = make([]*Arm, cfg.NArms)

	// log.Printf("drives: %d, USs: %d, CSs: %d", cfg.NDrives, cfg.NUSs, cfg.NCSs)
	// log.Printf("max arm length: %d", cfg.MaxArmLength)

	// defaults
	for i := range cfg.Arms {
		// TODO: if we permute CSs do we also want to keep the USs aligned?
		length := cfg.ArmLengths.Min
		lrng := cfg.ArmLengths.Range()
		if lrng > 0 {
			length += ev.Rand.Intn(lrng, -1)
		}
		arm := &Arm{Length: length, CS: i % cfg.NCSs, US: i % cfg.NUSs}
		cfg.Arms[i] = arm
		arm.Effort.Set(1, 1)
	}

	// defaults
	for i := range cfg.USs {
		us := &USParams{Prob: 1}
		cfg.USs[i] = us
		if i < cfg.NDrives {
			us.Negative = false
		} else {
			us.Negative = true
		}
		us.Mag.Set(1, 1)
	}

	if cfg.GroupMinMax {
		ev.ConfigGroupMinMax()
	}

	ev.UpdateMaxLength()
}

func (ev *Env) ConfigGroupMinMax() {
	cfg := &ev.Config
	// defaults
	// narms := cfg.NArms
	// nalts := narms / cfg.NUSs
	for i, arm := range cfg.Arms {
		ci := i / cfg.NUSs
		ui := i % cfg.NUSs
		arm.CS = i
		arm.US = ui
		if ci%2 == 0 {
			arm.Length = cfg.ArmLengths.Max
		} else {
			arm.Length = cfg.ArmLengths.Min
		}
	}
}

func (ev *Env) Validate() error {
	return nil
}

// Init does updating preparing to run -- params could have changed since initial config
// so updates everything except broad overall config stuff.
func (ev *Env) Init(run int) {
	cfg := &ev.Config

	ev.UpdateMaxLength()

	ev.States = make(map[string]*etensor.Float32)
	ev.States["CS"] = etensor.NewFloat32([]int{cfg.Params.NYReps, cfg.NCSs}, nil, nil)
	ev.States["Pos"] = etensor.NewFloat32([]int{cfg.Params.NYReps, ev.MaxLength + 1}, nil, nil)
	ev.States["Dist"] = etensor.NewFloat32([]int{cfg.Params.NYReps, ev.MaxLength + 1}, nil, nil)
	ev.States["Arm"] = etensor.NewFloat32([]int{cfg.Params.NYReps, ev.Config.NArms}, nil, nil)
	ev.States["Action"] = etensor.NewFloat32([]int{cfg.Params.NYReps, int(ActionsN)}, nil, nil)

	ev.NewStart()
	ev.JustConsumed = true // will trigger a new start again on Step
}

func (ev *Env) Counter(scale env.TimeScales) (cur, prv int, changed bool) {
	return 0, 0, false
}

func (ev *Env) State(el string) etensor.Tensor {
	return ev.States[el]
}

// NewStart starts a new approach run
func (ev *Env) NewStart() {
	arm := ev.Config.Arms[ev.Arm]
	// choose a new CS that maps to the same US
	arm.CS = (arm.CS + ev.Config.NCSs) % ev.Config.NCSs

	if ev.Config.Params.RandomStart {
		ev.Arm = ev.Rand.Intn(len(ev.Config.Arms), -1)
	}
	ev.Pos = 0
	ev.Dist = arm.Length - ev.Pos
	ev.Tick = 0
	ev.JustGated = false
	ev.HasGated = false
	ev.USConsumed = -1
	ev.USValue = 0
	ev.JustConsumed = false

	switch ev.Config.Paradigm {
	case Approach:
		ev.StartApproach()
	}
	ev.RenderState()
}

func (ev *Env) ExValueUtil(pv *axon.Rubicon, ctx *axon.Context) {
	maxval := float32(0)
	maxutil := float32(0)
	ev.ArmsNeg = nil
	usPos := make([]float32, pv.NPosUSs)
	cost := make([]float32, pv.NCosts)
	for i, arm := range ev.Config.Arms {
		us := ev.Config.USs[arm.US]
		if us.Negative {
			ev.ArmsNeg = append(ev.ArmsNeg, i)
			continue
		}
		val := ev.Drives[arm.US] * us.Mag.Midpoint() * us.Prob
		arm.ExValue = val
		for j := range usPos { // reset
			usPos[j] = 0
		}
		for j := range cost { // reset
			cost[j] = 0
		}
		usPos[arm.US+1] = val
		_, pvPos := pv.PVposEstFromUSs(ctx, uint32(ev.Di), usPos)
		exTime := float32(arm.Length) + 1 // time
		cost[0] = exTime
		cost[1] = exTime * arm.Effort.Midpoint()
		_, pvNeg := pv.PVcostEstFromCosts(cost)
		burst, dip, da, rew := pv.DAFromPVs(pvPos, pvNeg, 0, 0)
		_, _, _ = burst, dip, rew
		arm.ExPVpos = pvPos
		arm.ExPVneg = pvNeg
		arm.ExUtil = da
		if val > maxval {
			maxval = val
		}
		if da > maxutil {
			maxutil = da
		}
	}
	ev.MaxValue = maxval
	ev.MaxUtil = maxutil
	ev.ArmsMaxValue = nil
	ev.ArmsMaxUtil = nil
	for i, arm := range ev.Config.Arms {
		if arm.ExValue == maxval {
			ev.ArmsMaxValue = append(ev.ArmsMaxValue, i)
		}
		if arm.ExUtil == maxutil {
			ev.ArmsMaxUtil = append(ev.ArmsMaxUtil, i)
		}
	}
}

// ArmIsMaxValue returns true if the given arm is (one of) the arms with the best
// current expected outcome value
func (ev *Env) ArmIsMaxValue(arm int) bool {
	for _, ai := range ev.ArmsMaxValue {
		if arm == ai {
			return true
		}
	}
	return false
}

// ArmIsMaxUtil returns true if the given arm is (one of) the arms with the best
// current expected outcome utility
func (ev *Env) ArmIsMaxUtil(arm int) bool {
	for _, ai := range ev.ArmsMaxUtil {
		if arm == ai {
			return true
		}
	}
	return false
}

// ArmIsNegative returns true if the given arm is (one of) the arms with
// negative outcomes
func (ev *Env) ArmIsNegative(arm int) bool {
	for _, ai := range ev.ArmsNeg {
		if arm == ai {
			return true
		}
	}
	return false
}

// Step does one step.  it is up to the driving sim to decide when to call NewStart
func (ev *Env) Step() bool {
	ev.LastCS = ev.CurCS()
	if ev.JustConsumed { // from last time, not this time.
		ev.NewStart()
	} else {
		ev.Tick++
	}
	ev.TakeAct(ev.LastAct)
	switch ev.Config.Paradigm {
	case Approach:
		ev.StepApproach()
	}
	ev.RenderState()
	return true
}

//////////////////////////////////////////////////
//   Render

// RenderLocalist renders one localist state
func (ev *Env) RenderLocalist(name string, val int) {
	st := ev.States[name]
	st.SetZeros()
	if val >= st.Dim(1) {
		return
	}
	for y := 0; y < ev.Config.Params.NYReps; y++ {
		st.Set([]int{y, val}, 1.0)
	}
}

// RenderLocalist4D renders one localist state in 4D
func (ev *Env) RenderLocalist4D(name string, val int) {
	st := ev.States[name]
	st.SetZeros()
	for y := 0; y < ev.Config.Params.NYReps; y++ {
		st.Set([]int{0, val, y, 0}, 1.0)
	}
}

// RenderState renders the current state
func (ev *Env) RenderState() {
	ev.RenderLocalist("CS", ev.CurCS())
	ev.RenderLocalist("Pos", ev.Pos)
	ev.RenderLocalist("Dist", ev.Dist)
	ev.RenderLocalist("Arm", ev.Arm)
}

// RenderAction renders the action
func (ev *Env) RenderAction(act Actions) {
	ev.RenderLocalist("Action", int(act))
}

//////////////////////////////////////////////////
//   Action

func (ev *Env) DecodeAct(vt *etensor.Float32) Actions {
	mxi := ev.DecodeLocalist(vt)
	return Actions(mxi)
}

func (ev *Env) DecodeLocalist(vt *etensor.Float32) int {
	dx := vt.Dim(1)
	var max float32
	var mxi int
	for i := 0; i < dx; i++ {
		var sum float32
		for j := 0; j < ev.Config.Params.NYReps; j++ {
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
	act.SetString(action)
	ev.LastAct = act
	ev.RenderAction(act) // plus phase input is action
	// note: action not taken via TakeAct until start of trial in Step()
}

func (ev *Env) TakeAct(act Actions) {
	narms := ev.Config.NArms
	arm := ev.Config.Arms[ev.Arm]
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
			if ev.USConsumed < 0 {
				ev.ConsumeUS(arm)
			}
		}
	}
	// always update Dist
	arm = ev.Config.Arms[ev.Arm]
	ev.Dist = arm.Length - ev.Pos
}

// ConsumeUS implements the consume action at current position in given arm
func (ev *Env) ConsumeUS(arm *Arm) {
	us := ev.Config.USs[arm.US]
	mag := MinMaxRand(us.Mag, ev.Rand)
	got := erand.BoolP32(us.Prob, -1, &ev.Rand)
	if got {
		ev.USConsumed = arm.US
		ev.USValue = mag
		ev.JustConsumed = true
	} else {
		ev.USConsumed = -1
		ev.USValue = 0
	}
}

// InstinctAct returns an "instinctive" action that implements a basic policy
func (ev *Env) InstinctAct(justGated, hasGated bool) Actions {
	ev.JustGated = justGated
	ev.HasGated = hasGated
	ev.ShouldGate = ((hasGated && ev.USConsumed >= 0) || // To clear the goal after US
		(!hasGated && ev.ArmIsMaxUtil(ev.Arm))) // looking at correct, haven't yet gated

	arm := ev.CurArm()
	if ev.Pos >= arm.Length {
		return Consume
	}
	if ev.HasGated {
		return Forward
	}
	if ev.LastAct == Left || ev.LastAct == Right {
		return ev.LastAct
	}
	if ev.Config.Params.AlwaysLeft || erand.BoolP(.5, -1, &ev.Rand) {
		return Left
	}
	return Right
}

//////////////////////////////////////////////////
//   Utils

// CurArm returns current Arm
func (ev *Env) CurArm() *Arm {
	return ev.Config.Arms[ev.Arm]
}

// CurCS returns current CS from current Arm
func (ev *Env) CurCS() int {
	return ev.CurArm().CS
}

// MinMaxRand returns a random number in the range between Min and Max
func MinMaxRand(mm minmax.F32, rand erand.SysRand) float32 {
	return mm.Min + rand.Float32(-1)*mm.Range()
}

// InactiveVal returns a new random inactive value from Config.Params.Inactive
// param range.
func (ev *Env) InactiveValue() float32 {
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

func (ev *Env) UpdateMaxLength() {
	ev.MaxLength = 0
	for _, arm := range ev.Config.Arms {
		if arm.Length > ev.MaxLength {
			ev.MaxLength = arm.Length
		}
	}
}
