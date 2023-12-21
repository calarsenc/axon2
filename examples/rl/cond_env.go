// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"

	"github.com/emer/emergent/v2/env"
	"github.com/emer/emergent/v2/erand"
	"goki.dev/etable/v2/etensor"
)

// OnOff represents stimulus On / Off timing
type OnOff struct {

	// is this stimulus active -- use it?
	Act bool

	// when stimulus turns on
	On int

	// when stimulu turns off
	Off int

	// probability of being active on any given trial
	P float32

	// variability in onset timing (max number of trials before/after On that it could start)
	OnVar int

	// variability in offset timing (max number of trials before/after Off that it could end)
	OffVar int

	// current active status based on P probability
	CurAct bool `view:"-"`

	// current on / off values using Var variability
	CurOn, CurOff int `view:"-"`
}

func (oo *OnOff) Set(act bool, on, off int) {
	oo.Act = act
	oo.On = on
	oo.Off = off
	oo.P = 1 // default
}

// TrialUpdt updates Cur state at start of trial
func (oo *OnOff) TrialUpdt() {
	if !oo.Act {
		return
	}
	oo.CurAct = erand.BoolP32(oo.P, -1)
	oo.CurOn = oo.On - oo.OnVar + 2*rand.Intn(oo.OnVar+1)
	oo.CurOff = oo.Off - oo.OffVar + 2*rand.Intn(oo.OffVar+1)
}

// IsOn returns true if should be on according current time
func (oo *OnOff) IsOn(tm int) bool {
	return oo.Act && oo.CurAct && tm >= oo.CurOn && tm < oo.CurOff
}

// CondEnv simulates an n-armed bandit, where each of n inputs is associated with
// a specific probability of reward.
type CondEnv struct {

	// name of this environment
	Nm string

	// description of this environment
	Dsc string

	// total time for trial
	TotTime int

	// Conditioned stimulus A (e.g., Tone)
	CSA OnOff `view:"inline"`

	// Conditioned stimulus B (e.g., Light)
	CSB OnOff `view:"inline"`

	// Conditioned stimulus C
	CSC OnOff `view:"inline"`

	// Unconditioned stimulus -- reward
	US OnOff `view:"inline"`

	// value for reward
	RewVal float32

	// value for non-reward
	NoRewVal float32

	// one-hot input representation of current option
	Input etensor.Float64

	// single reward value
	Reward etensor.Float64

	// true if a US reward value was set
	HasRew bool

	// current run of model as provided during Init
	Run env.Ctr `view:"inline"`

	// number of times through Seq.Max number of sequences
	Epoch env.Ctr `view:"inline"`

	// one trial is a pass through all TotTime Events
	Trial env.Ctr `view:"inline"`

	// event is one time step within Trial -- e.g., CS turning on, etc
	Event env.Ctr `view:"inline"`
}

func (ev *CondEnv) Name() string { return ev.Nm }
func (ev *CondEnv) Desc() string { return ev.Dsc }

func (ev *CondEnv) Defaults() {
	ev.TotTime = 20
	ev.CSA.Set(true, 1, 6) // 10, 16
	ev.CSB.Set(false, 2, 10)
	ev.CSC.Set(false, 2, 5)
	ev.US.Set(true, 5, 6) // 15, 16
}

func (ev *CondEnv) Validate() error {
	if ev.TotTime == 0 {
		ev.Defaults()
	}
	return nil
}

func (ev *CondEnv) Counters() []env.TimeScales {
	return []env.TimeScales{env.Run, env.Epoch, env.Trial, env.Event}
}

func (ev *CondEnv) States() env.Elements {
	els := env.Elements{
		{"Input", []int{3, ev.TotTime}, []string{"3", "TotTime"}}, // CSC
		{"Reward", []int{1}, nil},
	}
	return els
}

func (ev *CondEnv) State(element string) etensor.Tensor {
	switch element {
	case "Input":
		return &ev.Input
	case "Reward":
		return &ev.Reward
	}
	return nil
}

func (ev *CondEnv) Actions() env.Elements {
	return nil
}

// String returns the current state as a string
func (ev *CondEnv) String() string {
	return fmt.Sprintf("S_%d_%g", ev.Event.Cur, ev.Reward.Values[0])
}

func (ev *CondEnv) Init(run int) {
	ev.Input.SetShape([]int{3, ev.TotTime}, nil, []string{"3", "TotTime"})
	ev.Reward.SetShape([]int{1}, nil, []string{"1"})
	ev.Run.Scale = env.Run
	ev.Epoch.Scale = env.Epoch
	ev.Trial.Scale = env.Trial
	ev.Event.Scale = env.Event
	ev.Run.Init()
	ev.Epoch.Init()
	ev.Trial.Init()
	ev.Event.Init()
	ev.Run.Cur = run
	ev.Event.Max = ev.TotTime
	ev.Event.Cur = -1 // init state -- key so that first Step() = 0
	ev.TrialUpdt()
}

// TrialUpdt updates all random vars at start of trial
func (ev *CondEnv) TrialUpdt() {
	ev.CSA.TrialUpdt()
	ev.CSB.TrialUpdt()
	ev.CSC.TrialUpdt()
	ev.US.TrialUpdt()
}

// SetInput sets the input state
func (ev *CondEnv) SetInput() {
	ev.Input.SetZeros()
	tm := ev.Event.Cur
	if ev.CSA.IsOn(tm) {
		ev.Input.Values[tm] = 1
	}
	if ev.CSB.IsOn(tm) {
		ev.Input.Values[ev.TotTime+tm] = 1
	}
	if ev.CSC.IsOn(tm) {
		ev.Input.Values[2*ev.TotTime+tm] = 1
	}
}

// SetReward sets reward for current option according to probability -- returns true if rewarded
func (ev *CondEnv) SetReward() bool {
	tm := ev.Event.Cur
	rw := ev.US.IsOn(tm)
	if rw {
		ev.HasRew = true
		ev.Reward.Values[0] = float64(ev.RewVal)
	} else {
		ev.HasRew = false
		ev.Reward.Values[0] = float64(ev.NoRewVal)
	}
	return rw
}

func (ev *CondEnv) Step() bool {
	ev.Epoch.Same() // good idea to just reset all non-inner-most counters at start
	ev.Trial.Same() // this ensures that they only report changed when actually changed

	incr := ev.Event.Incr()
	ev.SetInput()
	ev.SetReward()

	if incr {
		ev.TrialUpdt()
		if ev.Trial.Incr() {
			ev.Epoch.Incr()
		}
	}
	return true
}

func (ev *CondEnv) Action(element string, input etensor.Tensor) {
	// nop
}

func (ev *CondEnv) Counter(scale env.TimeScales) (cur, prv int, chg bool) {
	switch scale {
	case env.Run:
		return ev.Run.Query()
	case env.Epoch:
		return ev.Epoch.Query()
	case env.Trial:
		return ev.Trial.Query()
	case env.Event:
		return ev.Event.Query()
	}
	return -1, -1, false
}

// Compile-time check that implements Env interface
var _ env.Env = (*CondEnv)(nil)
