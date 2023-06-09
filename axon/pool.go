// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package axon

import (
	"github.com/emer/axon/fsfffb"
	"github.com/goki/gosl/slbool"
	"github.com/goki/mat32"
)

//gosl: hlsl pool
// #include "avgmaxi.hlsl"
//gosl: end pool

//gosl: start pool

// AvgMaxPhases contains the average and maximum values over a Pool of neurons,
// at different time scales within a standard ThetaCycle of updating.
// It is much more efficient on the GPU to just grab everything in one pass at
// the cycle level, and then take snapshots from there.
// All of the cycle level values are updated at the *start* of the cycle
// based on values from the prior cycle -- thus are 1 cycle behind in general.
type AvgMaxPhases struct {
	Cycle AvgMaxI32 `view:"inline" desc:"updated every cycle -- this is the source of all subsequent time scales"`
	Minus AvgMaxI32 `view:"inline" desc:"at the end of the minus phase"`
	Plus  AvgMaxI32 `view:"inline" desc:"at the end of the plus phase"`
	Prev  AvgMaxI32 `view:"inline" desc:"at the end of the previous plus phase"`
}

// CycleToMinus grabs current Cycle values into the Minus phase values
func (am *AvgMaxPhases) CycleToMinus() {
	am.Minus = am.Cycle
	am.Prev = am.Plus
}

// CycleToPlus grabs current Cycle values into the Plus phase values
func (am *AvgMaxPhases) CycleToPlus() {
	am.Plus = am.Cycle
}

// Calc does Calc on Cycle, which is then ready for aggregation again
func (am *AvgMaxPhases) Calc(refIdx int32) {
	am.Cycle.Calc(refIdx)
}

// Zero does a full reset on everything -- for InitActs
func (am *AvgMaxPhases) Zero() {
	am.Cycle.Zero()
	am.Minus.Zero()
	am.Plus.Zero()
	am.Prev.Zero()
}

// PoolAvgMax contains the average and maximum values over a Pool of neurons
// for different variables of interest, at Cycle, Minus and Plus phase timescales.
// All of the cycle level values are updated at the *start* of the cycle
// based on values from the prior cycle -- thus are 1 cycle behind in general.
type PoolAvgMax struct {
	CaSpkP   AvgMaxPhases `inactive:"+" view:"inline" desc:"avg and maximum CaSpkP (continuously updated at roughly 40 msec integration window timescale, ends up capturing potentiation, plus-phase signal) -- this is the primary variable to use for tracking overall pool activity"`
	CaSpkD   AvgMaxPhases `inactive:"+" view:"inline" desc:"avg and maximum CaSpkD longer-term depression / DAPK1 signal in layer"`
	SpkMax   AvgMaxPhases `inactive:"+" view:"inline" desc:"avg and maximum SpkMax value (based on CaSpkP) -- reflects peak activity at any point across the cycle"`
	Act      AvgMaxPhases `inactive:"+" view:"inline" desc:"avg and maximum Act firing rate value"`
	GeInt    AvgMaxPhases `inactive:"+" view:"inline" desc:"avg and maximum GeInt integrated running-average excitatory conductance value"`
	GeIntMax AvgMaxPhases `inactive:"+" view:"inline" desc:"avg and maximum GeIntMax integrated running-average excitatory conductance value"`
	GiInt    AvgMaxPhases `inactive:"+" view:"inline" desc:"avg and maximum GiInt integrated running-average inhibitory conductance value"`
}

// SetN sets the N for aggregation
func (am *PoolAvgMax) SetN(n int32) {
	am.CaSpkP.Cycle.N = n
	am.CaSpkD.Cycle.N = n
	am.SpkMax.Cycle.N = n
	am.Act.Cycle.N = n
	am.GeInt.Cycle.N = n
	am.GeIntMax.Cycle.N = n
	am.GiInt.Cycle.N = n
}

// CycleToMinus grabs current Cycle values into the Minus phase values
func (am *PoolAvgMax) CycleToMinus() {
	am.CaSpkP.CycleToMinus()
	am.CaSpkD.CycleToMinus()
	am.SpkMax.CycleToMinus()
	am.Act.CycleToMinus()
	am.GeInt.CycleToMinus()
	am.GeIntMax.CycleToMinus()
	am.GiInt.CycleToMinus()
}

// CycleToPlus grabs current Cycle values into the Plus phase values
func (am *PoolAvgMax) CycleToPlus() {
	am.CaSpkP.CycleToPlus()
	am.CaSpkD.CycleToPlus()
	am.SpkMax.CycleToPlus()
	am.Act.CycleToPlus()
	am.GeInt.CycleToPlus()
	am.GeIntMax.CycleToPlus()
	am.GiInt.CycleToPlus()
}

// Init does Init on Cycle vals-- for update start.
// always left init'd so generally unnecessary
func (am *PoolAvgMax) Init() {
	am.CaSpkP.Cycle.Init()
	am.CaSpkD.Cycle.Init()
	am.SpkMax.Cycle.Init()
	am.Act.Cycle.Init()
	am.GeInt.Cycle.Init()
	am.GeIntMax.Cycle.Init()
	am.GiInt.Cycle.Init()
}

// Zero does full reset on everything -- for InitActs
func (am *PoolAvgMax) Zero() {
	am.CaSpkP.Zero()
	am.CaSpkD.Zero()
	am.SpkMax.Zero()
	am.Act.Zero()
	am.GeInt.Zero()
	am.GeIntMax.Zero()
	am.GiInt.Zero()
}

// Calc does Calc on Cycle level, and re-inits
func (am *PoolAvgMax) Calc(refIdx int32) {
	am.CaSpkP.Calc(refIdx)
	am.CaSpkD.Calc(refIdx)
	am.SpkMax.Calc(refIdx)
	am.Act.Calc(refIdx)
	am.GeInt.Calc(refIdx)
	am.GeIntMax.Calc(refIdx)
	am.GiInt.Calc(refIdx)
}

//gosl: end pool

// note: the following is actually being used despite appearing to be
// commented out!  it is auto-uncommented when copied to hlsl
// MUST update whenever above UpdateVals code is updated.

//gosl: hlsl pool
/*
// // AtomicUpdatePoolAvgMax provides an atomic update using atomic ints
// // implemented by InterlockedAdd HLSL intrinsic.
// // This is a #define because it doesn't work on arg values --
// // must be directly operating on a RWStorageBuffer entity.
#define AtomicUpdatePoolAvgMax(am, ctx, ni, di) \
	AtomicUpdateAvgMaxI32(am.CaSpkP.Cycle, NrnV(ctx, ni, di, CaSpkP)); \
	AtomicUpdateAvgMaxI32(am.CaSpkD.Cycle, NrnV(ctx, ni, di, CaSpkD)); \
	AtomicUpdateAvgMaxI32(am.SpkMax.Cycle, NrnV(ctx, ni, di, SpkMax)); \
	AtomicUpdateAvgMaxI32(am.Act.Cycle, NrnV(ctx, ni, di, Act)); \
	AtomicUpdateAvgMaxI32(am.GeInt.Cycle, NrnV(ctx, ni, di, GeInt)); \
	AtomicUpdateAvgMaxI32(am.GeIntMax.Cycle, NrnV(ctx, ni, di, GeIntMax)); \
	AtomicUpdateAvgMaxI32(am.GiInt.Cycle, NrnV(ctx, ni, di, GiInt))
*/
//gosl: end pool

//gosl: start pool

// Pool contains computed values for FS-FFFB inhibition,
// and various other state values for layers
// and pools (unit groups) that can be subject to inhibition
type Pool struct {
	StIdx, EdIdx uint32      `inactive:"+" desc:"starting and ending (exlusive) layer-wise indexes for the list of neurons in this pool"`
	LayIdx       uint32      `view:"-" desc:"layer index in global layer list"`
	DataIdx      uint32      `view:"-" desc:"data parallel index (innermost index per layer)"`
	PoolIdx      uint32      `view:"-" desc:"pool index in global pool list: [Layer][Pool][Data]"`
	IsLayPool    slbool.Bool `inactive:"+" desc:"is this a layer-wide pool?  if not, it represents a sub-pool of units within a 4D layer"`
	Gated        slbool.Bool `inactive:"+" desc:"for special types where relevant (e.g., MatrixLayer, BGThalLayer), indicates if the pool was gated"`

	pad uint32

	Inhib fsfffb.Inhib `inactive:"+" desc:"fast-slow FFFB inhibition values"`
	// note: these last two have elements that are shared across data parallel -- not worth separating though?
	AvgMax PoolAvgMax `desc:"average and max values for relevant variables in this pool, at different time scales"`
	AvgDif AvgMaxI32  `inactive:"+" view:"inline" desc:"absolute value of AvgDif differences from actual neuron ActPct relative to TrgAvg"`
}

// Init is callled during InitActs
func (pl *Pool) Init() {
	pl.Inhib.Init()
	pl.AvgMax.Zero()
	pl.AvgMax.SetN(int32(pl.NNeurons()))
	pl.AvgDif.N = int32(pl.NNeurons())
	pl.AvgDif.Init()
	pl.Gated.SetBool(false)
}

// NNeurons returns the number of neurons in the pool: EdIdx - StIdx
func (pl *Pool) NNeurons() int {
	return int(pl.EdIdx - pl.StIdx)
}

//gosl: end pool

// AvgMaxUpdate updates the AvgMax values based on current neuron values
func (pl *Pool) AvgMaxUpdate(ctx *Context, ni, di uint32) {
	pl.AvgMax.CaSpkP.Cycle.UpdateVal(NrnV(ctx, ni, di, CaSpkP))
	pl.AvgMax.CaSpkD.Cycle.UpdateVal(NrnV(ctx, ni, di, CaSpkD))
	pl.AvgMax.SpkMax.Cycle.UpdateVal(NrnV(ctx, ni, di, SpkMax))
	pl.AvgMax.Act.Cycle.UpdateVal(mat32.Abs(NrnV(ctx, ni, di, Act))) // can be neg
	pl.AvgMax.GeInt.Cycle.UpdateVal(NrnV(ctx, ni, di, GeInt))
	pl.AvgMax.GeIntMax.Cycle.UpdateVal(NrnV(ctx, ni, di, GeIntMax))
	pl.AvgMax.GiInt.Cycle.UpdateVal(NrnV(ctx, ni, di, GiInt))
}

// TestVals returns a map of CaSpkD.Avg, which provides an
// integrated summary of pool activity for testing
func (pl *Pool) TestVals(layKey string, vals map[string]float32) {
	vals[layKey+" CaSpkD Avg"] = pl.AvgMax.CaSpkD.Cycle.Avg
}
