// Copyright (c) 2022, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package axon

import (
	"github.com/emer/emergent/erand"
	"github.com/goki/gosl/slbool"
	"github.com/goki/ki/bools"
	"github.com/goki/mat32"
)

//gosl: start pvlv

// DriveVals represents different internal drives,
// such as hunger, thirst, etc.  The first drive is
// typically reserved for novelty / curiosity.
// labels can be provided by specific environments.
type DriveVals struct {
	D0 float32
	D1 float32
	D2 float32
	D3 float32
	D4 float32
	D5 float32
	D6 float32
	D7 float32
}

func (ds *DriveVals) SetAll(val float32) {
	ds.D0 = val
	ds.D1 = val
	ds.D2 = val
	ds.D3 = val
	ds.D4 = val
	ds.D5 = val
	ds.D6 = val
	ds.D7 = val
}

func (ds *DriveVals) Zero() {
	ds.SetAll(0)
}

func (ds *DriveVals) Set(drv uint32, val float32) {
	switch drv {
	case 0:
		ds.D0 = val
	case 1:
		ds.D1 = val
	case 2:
		ds.D2 = val
	case 3:
		ds.D3 = val
	case 4:
		ds.D4 = val
	case 5:
		ds.D5 = val
	case 6:
		ds.D6 = val
	case 7:
		ds.D7 = val
	}
}

func (ds *DriveVals) Get(drv uint32) float32 {
	val := float32(0)
	switch drv {
	case 0:
		val = ds.D0
	case 1:
		val = ds.D1
	case 2:
		val = ds.D2
	case 3:
		val = ds.D3
	case 4:
		val = ds.D4
	case 5:
		val = ds.D5
	case 6:
		val = ds.D6
	case 7:
		val = ds.D7
	}
	return val
}

// Drives manages the drive parameters for updating drive state,
// and drive state.
type Drives struct {
	NActive  uint32  `max:"8" desc:"number of active drives -- first drive is novelty / curiosity drive -- total must be &lt;= 8"`
	NNegUSs  uint32  `min:"1" max:"8" desc:"number of active negative US states recognized -- the first is always reserved for the accumulated effort cost / dissapointment when an expected US is not achieved"`
	DriveMin float32 `desc:"minimum effective drive value -- this is an automatic baseline ensuring that a positive US results in at least some minimal level of reward.  Unlike Base values, this is not reflected in the activity of the drive values -- applies at the time of reward calculation as a minimum baseline."`

	pad int32

	Base  DriveVals `view:"inline" desc:"baseline levels for each drive -- what they naturally trend toward in the absence of any input.  Set inactive drives to 0 baseline, active ones typically elevated baseline (0-1 range)."`
	Tau   DriveVals `view:"inline" desc:"time constants in ThetaCycle (trial) units for natural update toward Base values -- 0 values means no natural update."`
	USDec DriveVals `view:"inline" desc:"decrement in drive value when Drive-US is consumed -- positive values are subtracted from current Drive value."`

	Drives DriveVals `inactive:"+" view:"inline" desc:"current drive state -- updated with optional homeostatic exponential return to baseline values"`

	Dt DriveVals `view:"-" desc:"1/Tau"`
}

func (dp *Drives) Defaults() {
	dp.NNegUSs = 1
	dp.DriveMin = 0.5
	dp.Update()
	dp.USDec.SetAll(1)
}

// ToBaseline sets all drives to their baseline levels
func (dp *Drives) ToBaseline() {
	dp.Drives = dp.Base
}

// ToZero sets all drives to 0
func (dp *Drives) ToZero() {
	dp.Drives = dp.Base
}

func (dp *Drives) Update() {
	for i := uint32(0); i < 8; i++ {
		tau := dp.Tau.Get(i)
		if tau <= 0 {
			dp.Dt.Set(i, 0)
		} else {
			dp.Dt.Set(i, 1.0/tau)
		}
	}
}

// Zero sets all NActive values of given drive vars to 0
func (dp *Drives) Zero(ctx *Context, di uint32, gvar GlobalVars) {
	for i := uint32(0); i < dp.NActive; i++ {
		SetGlobalDriveV(ctx, di, i, gvar, 0)
	}
}

// AddTo increments drive by given amount, subject to 0-1 range clamping.
// Returns new val.
func (dp *Drives) AddTo(ctx *Context, di uint32, drv uint32, delta float32) float32 {
	dv := GlobalDriveV(ctx, di, drv, GvDrives) + delta
	if dv > 1 {
		dv = 1
	} else if dv < 0 {
		dv = 0
	}
	SetGlobalDriveV(ctx, di, drv, GvDrives, dv)
	return dv
}

// SoftAdd increments drive by given amount, using soft-bounding to 0-1 extremes.
// if delta is positive, multiply by 1-val, else val.  Returns new val.
func (ds *Drives) SoftAdd(ctx *Context, di uint32, drv uint32, delta float32) float32 {
	dv := GlobalDriveV(ctx, di, drv, GvDrives)
	if delta > 0 {
		dv += (1 - dv) * delta
	} else {
		dv += dv * delta
	}
	if dv > 1 {
		dv = 1
	} else if dv < 0 {
		dv = 0
	}
	SetGlobalDriveV(ctx, di, drv, GvDrives, dv)
	return dv
}

// ExpStep updates drive with an exponential step with given dt value
// toward given baseline value.
func (ds *Drives) ExpStep(ctx *Context, di uint32, drv uint32, dt, base float32) float32 {
	dv := GlobalDriveV(ctx, di, drv, GvDrives)
	dv += dt * (base - dv)
	if dv > 1 {
		dv = 1
	} else if dv < 0 {
		dv = 0
	}
	SetGlobalDriveV(ctx, di, drv, GvDrives, dv)
	return dv
}

// ExpStepAll updates given drives with an exponential step using dt values
// toward baseline values.
func (dp *Drives) ExpStepAll(ctx *Context, di uint32) {
	for i := uint32(0); i < dp.NActive; i++ {
		dp.ExpStep(ctx, di, i, dp.Dt.Get(i), dp.Base.Get(i))
	}
}

// EffectiveDrive returns the Max of Drives at given index and DriveMin.
// note that index 0 is the novelty / curiosity drive.
func (dp *Drives) EffectiveDrive(ctx *Context, di uint32, i uint32) float32 {
	if i == 0 {
		return GlobalDriveV(ctx, di, uint32(0), GvDrives)
	}
	return mat32.Max(GlobalDriveV(ctx, di, i, GvDrives), dp.DriveMin)
}

///////////////////////////////////////////////////////////////////////////////
//  Effort

// Effort has effort and parameters for updating it
type Effort struct {
	Gain       float32 `desc:"gain factor for computing effort discount factor -- larger = quicker discounting"`
	CurMax     float32 `inactive:"-" desc:"current maximum raw effort level -- above this point, any current goal will be terminated during the GiveUp function, which also looks for accumulated disappointment.  See Max, MaxNovel, MaxPostDip for values depending on how the goal was triggered."`
	Max        float32 `desc:"default maximum raw effort level, when MaxNovel and MaxPostDip don't apply."`
	MaxNovel   float32 `desc:"maximum raw effort level when novelty / curiosity drive is engaged -- typically shorter than default Max"`
	MaxPostDip float32 `desc:"if the LowThr amount of VSPatch expectation is triggered, as accumulated in LHb.DipSum, then CurMax is set to the current Raw effort plus this increment, which is generally low -- once an expectation has been activated, don't wait around forever.."`
	MaxVar     float32 `desc:"variance in additional maximum effort level, applied whenever CurMax is updated"`
}

func (ef *Effort) Defaults() {
	ef.Gain = 0.1
	ef.Max = 100
	ef.MaxNovel = 8
	ef.MaxPostDip = 4
	ef.MaxVar = 2
}

func (ef *Effort) Update() {

}

// Reset resets the raw effort back to zero -- at start of new gating event
func (ef *Effort) Reset(ctx *Context, di uint32) {
	SetGlobalV(ctx, di, GvEffortRaw, 0)
	SetGlobalV(ctx, di, GvEffortCurMax, ef.Max)
	SetGlobalV(ctx, di, GvEffortDisc, 1)
}

// DiscFun is the effort discount function: 1 / (1 + ef.Gain * effort)
func (ef *Effort) DiscFun(effort float32) float32 {
	return 1.0 / (1.0 + ef.Gain*effort)
}

// DiscFmEffort computes Disc from Raw effort
func (ef *Effort) DiscFmEffort(ctx *Context, di uint32) float32 {
	disc := ef.DiscFun(GlobalV(ctx, di, GvEffortRaw))
	SetGlobalV(ctx, di, GvEffortDisc, disc)
	return disc
}

// AddEffort adds an increment of effort and updates the Disc discount factor
func (ef *Effort) AddEffort(ctx *Context, di uint32, inc float32) {
	AddGlobalV(ctx, di, GvEffortRaw, inc)
	ef.DiscFmEffort(ctx, di)
}

// GiveUp returns true if maximum effort has been exceeded
func (ef *Effort) GiveUp(ctx *Context, di uint32) bool {
	raw := GlobalV(ctx, di, GvEffortRaw)
	curMax := GlobalV(ctx, di, GvEffortCurMax)
	if curMax > 0 && raw > curMax {
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////
//  Urgency

// Urgency has urgency (increasing pressure to do something) and parameters for updating it.
// Raw urgency is incremented by same units as effort, but is only reset with a positive US.
// Could also make it a function of drives and bodily state factors
// e.g., desperate thirst, hunger.  Drive activations probably have limited range
// and renormalization, so urgency can be another dimension with more impact by directly biasing Go.
type Urgency struct {
	U50   float32 `desc:"value of raw urgency where the urgency activation level is 50%"`
	Power int32   `def:"4" desc:"exponent on the urge factor -- valid numbers are 1,2,4,6"`
	Thr   float32 `def:"0.2" desc:"threshold for urge -- cuts off small baseline values"`

	pad float32
}

func (ur *Urgency) Defaults() {
	ur.U50 = 20
	ur.Power = 4
	ur.Thr = 0.2
}

func (ur *Urgency) Update() {

}

// Reset resets the raw urgency back to zero -- at start of new gating event
func (ur *Urgency) Reset(ctx *Context, di uint32) {
	SetGlobalV(ctx, di, GvUrgencyRaw, 0)
	SetGlobalV(ctx, di, GvUrgency, 0)
}

// UrgeFun is the urgency function: urgency / (urgency + 1) where
// urgency = (Raw / U50)^Power
func (ur *Urgency) UrgeFun(urgency float32) float32 {
	urgency /= ur.U50
	switch ur.Power {
	case 2:
		urgency *= urgency
	case 4:
		urgency *= urgency * urgency * urgency
	case 6:
		urgency *= urgency * urgency * urgency * urgency * urgency
	}
	return urgency / (1.0 + urgency)
}

// UrgeFmUrgency computes Urge from Raw
func (ur *Urgency) UrgeFmUrgency(ctx *Context, di uint32) float32 {
	urge := ur.UrgeFun(GlobalV(ctx, di, GvUrgencyRaw))
	if urge < ur.Thr {
		urge = 0
	}
	SetGlobalV(ctx, di, GvUrgency, urge)
	return urge
}

// AddEffort adds an effort increment of urgency and updates the Urge factor
func (ur *Urgency) AddEffort(ctx *Context, di uint32, inc float32) {
	AddGlobalV(ctx, di, GvUrgencyRaw, inc)
	ur.UrgeFmUrgency(ctx, di)
}

///////////////////////////////////////////////////////////////////////////////
//  LHb & RMTg

// LHb has values for computing LHb & RMTg which drives dips / pauses in DA firing.
// Positive net LHb activity drives dips / pauses in VTA DA activity,
// e.g., when predicted pos > actual or actual neg > predicted.
// Negative net LHb activity drives bursts in VTA DA activity,
// e.g., when actual pos > predicted (redundant with LV / Amygdala)
// or "relief" burst when actual neg < predicted.
type LHb struct {
	PosGain   float32 `def:"1" desc:"gain multiplier on overall VSPatchPos - PosPV component"`
	NegGain   float32 `def:"1" desc:"gain multiplier on overall PVneg component"`
	GiveUpThr float32 `def:"0.2" desc:"threshold on summed LHbDip over trials for triggering a reset of goal engaged state"`
	DipLowThr float32 `def:"0.05" desc:"low threshold on summed LHbDip, used for triggering switch to a faster effort max timeout -- Effort.MaxPostDip"`

	pad float32
}

func (lh *LHb) Defaults() {
	lh.PosGain = 1
	lh.NegGain = 1
	lh.GiveUpThr = 0.2
	lh.DipLowThr = 0.05
}

func (lh *LHb) Update() {
}

func (lh *LHb) Reset(ctx *Context, di uint32) {
	SetGlobalV(ctx, di, GvLHbDip, 0)
	SetGlobalV(ctx, di, GvLHbBurst, 0)
	SetGlobalV(ctx, di, GvLHbDipSumCur, 0)
	SetGlobalV(ctx, di, GvLHbDipSum, 0)
	SetGlobalV(ctx, di, GvLHbGiveUp, 0)
}

// LHbFmPVVS computes the overall LHbDip and LHbBurst values from PV (primary value)
// and VSPatch inputs.
func (lh *LHb) LHbFmPVVS(ctx *Context, di uint32, pvPos, pvNeg, vsPatchPos float32) {
	pos := lh.PosGain * (vsPatchPos - pvPos)
	neg := lh.NegGain * pvNeg
	SetGlobalV(ctx, di, GvLHbPos, pos)
	SetGlobalV(ctx, di, GvLHbNeg, neg)
	netLHb := pos + neg

	if netLHb > 0 {
		SetGlobalV(ctx, di, GvLHbDip, netLHb)
		SetGlobalV(ctx, di, GvLHbBurst, 0)
	} else {
		SetGlobalV(ctx, di, GvLHbBurst, -netLHb)
		SetGlobalV(ctx, di, GvLHbDip, 0)
	}
}

// ShouldGiveUp increments DipSum and checks if should give up if above threshold
func (lh *LHb) ShouldGiveUp(ctx *Context, di uint32) bool {
	dip := GlobalV(ctx, di, GvLHbDip)
	AddGlobalV(ctx, di, GvLHbDipSumCur, dip)
	cur := GlobalV(ctx, di, GvLHbDipSumCur)
	SetGlobalV(ctx, di, GvLHbDipSum, cur)
	SetGlobalV(ctx, di, GvLHbGiveUp, 0)
	giveUp := false
	if cur > lh.GiveUpThr {
		giveUp = true
		SetGlobalV(ctx, di, GvLHbGiveUp, 1)
		SetGlobalV(ctx, di, GvLHbGiveUp, 1)
		SetGlobalV(ctx, di, GvLHbDipSumCur, 0)
	}
	return giveUp
}

///////////////////////////////////////////////////////////////////////////////
//  VTA

// VTAVals has values for all the inputs to the VTA.
// Used as gain factors and computed values.
type VTAVals struct {
	DA         float32 `desc:"overall dopamine value reflecting all of the different inputs"`
	USpos      float32 `desc:"total positive valence primary value = sum of USpos * Drive without effort discounting"`
	PVpos      float32 `desc:"total positive valence primary value = sum of USpos * Drive * (1-Effort.Disc) -- what actually drives DA bursting from actual USs received"`
	PVneg      float32 `desc:"total negative valence primary value = sum of USneg inputs"`
	CeMpos     float32 `desc:"positive valence central nucleus of the amygdala (CeM) LV (learned value) activity, reflecting |BLAPosAcqD1 - BLAPosExtD2|_+ positively rectified.  CeM sets Raw directly.  Note that a positive US onset even with no active Drive will be reflected here, enabling learning about unexpected outcomes."`
	CeMneg     float32 `desc:"negative valence central nucleus of the amygdala (CeM) LV (learned value) activity, reflecting |BLANegAcqD2 - BLANegExtD1|_+ positively rectified.  CeM sets Raw directly."`
	LHbDip     float32 `desc:"dip from LHb / RMTg -- net inhibitory drive on VTA DA firing = dips"`
	LHbBurst   float32 `desc:"burst from LHb / RMTg -- net excitatory drive on VTA DA firing = bursts"`
	VSPatchPos float32 `desc:"net shunting input from VSPatch (PosD1 -- PVi in original PVLV)"`

	pad, pad1, pad2 float32
}

func (vt *VTAVals) Set(usPos, pvPos, pvNeg, lhbDip, lhbBurst, vsPatchPos float32) {
	vt.USpos = usPos
	vt.PVpos = pvPos
	vt.PVneg = pvNeg
	vt.LHbDip = lhbDip
	vt.LHbBurst = lhbBurst
	vt.VSPatchPos = vsPatchPos
}

func (vt *VTAVals) SetAll(val float32) {
	vt.DA = val
	vt.USpos = val
	vt.PVpos = val
	vt.PVneg = val
	vt.CeMpos = val
	vt.CeMneg = val
	vt.LHbDip = val
	vt.LHbBurst = val
	vt.VSPatchPos = val
}

func (vt *VTAVals) Zero() {
	vt.SetAll(0)
}

// VTA has parameters and values for computing VTA DA dopamine,
// as a function of:
//   - PV (primary value) driving inputs reflecting US inputs,
//     which are modulated by Drives and discounted by Effort for positive.
//   - LV / Amygdala which drives bursting for unexpected CSs or USs via CeM.
//   - Shunting expectations of reward from VSPatchPosD1 - D2.
//   - Dipping / pausing inhibitory inputs from lateral habenula (LHb) reflecting
//     predicted positive outcome > actual, or actual negative > predicted.
//   - ACh from LDT (laterodorsal tegmentum) reflecting sensory / reward salience,
//     which disinhibits VTA activity.
type VTA struct {
	PVThr float32 `desc:"threshold for activity of PVpos or VSPatchPos to determine if a PV event (actual PV or omission thereof) is present"`

	pad, pad1, pad2 float32

	Gain VTAVals `view:"inline" desc:"gain multipliers on inputs from each input"`
	// Raw  VTAVals `view:"inline" inactive:"+" desc:"raw current values -- inputs to the computation"`
	// Vals VTAVals `view:"inline" inactive:"+" desc:"computed current values"`
	// Prev VTAVals `view:"inline" inactive:"+" desc:"previous computed  values -- to avoid a data race"`
}

func (vt *VTA) Defaults() {
	vt.PVThr = 0.05
	vt.Gain.SetAll(1)
}

func (vt *VTA) Update() {
}

func (vt *VTA) ZeroVals(ctx *Context, di uint32, vtaType GlobalVTAType) {
	for vv := GvVtaDA; vv <= GvVtaVSPatchPos; vv++ {
		SetGlobalVTA(ctx, di, vtaType, vv, 0)
	}
}

func (vt *VTA) Reset(ctx *Context, di uint32) {
	vt.ZeroVals(ctx, di, GvVtaRaw)
	vt.ZeroVals(ctx, di, GvVtaVals)
	vt.ZeroVals(ctx, di, GvVtaPrev)
}

// DAFmRaw computes the intermediate Vals and final DA value from
// Raw values that have been set prior to calling.
// ACh value from LDT is passed as a parameter.
func (vt *VTA) DAFmRaw(ctx *Context, di uint32, ach float32, hasRew bool) {
	SetGlobalVTA(ctx, di, GvVtaVals, GvVtaPVpos, vt.Gain.PVpos*GlobalVTA(ctx, di, GvVtaRaw, GvVtaPVpos))
	SetGlobalVTA(ctx, di, GvVtaVals, GvVtaCeMpos, vt.Gain.CeMpos*GlobalVTA(ctx, di, GvVtaRaw, GvVtaCeMpos))
	SetGlobalVTA(ctx, di, GvVtaVals, GvVtaCeMneg, vt.Gain.CeMneg*GlobalVTA(ctx, di, GvVtaRaw, GvVtaCeMneg))
	SetGlobalVTA(ctx, di, GvVtaVals, GvVtaLHbDip, vt.Gain.LHbDip*GlobalVTA(ctx, di, GvVtaRaw, GvVtaLHbDip))
	SetGlobalVTA(ctx, di, GvVtaVals, GvVtaLHbBurst, vt.Gain.LHbBurst*GlobalVTA(ctx, di, GvVtaRaw, GvVtaLHbBurst))
	SetGlobalVTA(ctx, di, GvVtaVals, GvVtaVSPatchPos, vt.Gain.VSPatchPos*GlobalVTA(ctx, di, GvVtaRaw, GvVtaVSPatchPos))

	if GlobalVTA(ctx, di, GvVtaVals, GvVtaVSPatchPos) < 0 {
		SetGlobalVTA(ctx, di, GvVtaVals, GvVtaVSPatchPos, 0)
	}
	pvDA := GlobalVTA(ctx, di, GvVtaVals, GvVtaPVpos) - GlobalVTA(ctx, di, GvVtaVals, GvVtaVSPatchPos) - GlobalVTA(ctx, di, GvVtaVals, GvVtaPVneg) - GlobalVTA(ctx, di, GvVtaVals, GvVtaLHbDip)
	csNet := GlobalVTA(ctx, di, GvVtaVals, GvVtaCeMpos) - GlobalVTA(ctx, di, GvVtaVals, GvVtaCeMneg)
	csDA := ach*mat32.Max(csNet, GlobalVTA(ctx, di, GvVtaVals, GvVtaLHbBurst)) - GlobalVTA(ctx, di, GvVtaVals, GvVtaLHbDip) // restore LHbDip contribution
	// note that ach is only on cs -- should be 1 for PV events anyway..
	netDA := float32(0)
	if hasRew {
		netDA = pvDA
	} else {
		netDA = csDA
	}
	SetGlobalVTA(ctx, di, GvVtaVals, GvVtaDA, vt.Gain.DA*netDA)
}

// func (vt *VSMatrix) Reset() {
// 	vt.JustGated.SetBool(false)
// 	vt.HasGated.SetBool(false)
// }
//
// // VSGated updates JustGated as function of VS gating
// // at end of the plus phase.
// func (vt *VSMatrix) VSGated(gated bool) {
// 	vt.JustGated.SetBool(gated)
// }

///////////////////////////////////////////////////////////////////////////////
//  PVLV

// PVLV represents the core brainstem-level (hypothalamus) bodily drives
// and resulting dopamine from US (unconditioned stimulus) inputs,
// as computed by the PVLV model of primary value (PV)
// and learned value (LV), describing the functions of the Amygala,
// Ventral Striatum, VTA and associated midbrain nuclei (LDT, LHb, RMTg)
// Core LHb (lateral habenula) and VTA (ventral tegmental area) dopamine
// are computed in equations using inputs from specialized network layers
// (LDTLayer driven by BLA, CeM layers, VSPatchLayer).
// Renders USLayer, PVLayer, DrivesLayer representations based on state updated here.
type PVLV struct {
	Drive   Drives    `desc:"parameters and state for built-in drives that form the core motivations of agent, controlled by lateral hypothalamus and associated body state monitoring such as glucose levels and thirst."`
	Effort  Effort    `view:"inline" desc:"effort parameters and state, tracking relative depletion of glucose levels and water levels as a function of time and exertion"`
	Urgency Urgency   `view:"inline" desc:"urgency (increasing pressure to do something) and parameters for updating it. Raw urgency is incremented by same units as effort, but is only reset with a positive US."`
	VTA     VTA       `desc:"parameters and values for computing VTA dopamine, as a function of PV primary values (via Pos / Neg US), LV learned values (Amygdala bursting from unexpected CSs, USs), shunting VSPatchPos expectations, and dipping / pausing inputs from LHb"`
	LHb     LHb       `view:"inline" desc:"lateral habenula (LHb) parameters and state, which drives dipping / pausing in dopamine when the predicted positive outcome > actual, or actual negative outcome > predicted.  Can also drive bursting for the converse, and via matrix phasic firing"`
	USpos   DriveVals `inactive:"+" view:"inline" desc:"current positive-valence drive-satisfying input(s) (unconditioned stimuli = US)"`
	USneg   DriveVals `inactive:"+" view:"inline" desc:"current negative-valence (aversive), non-drive-satisfying input(s) (unconditioned stimuli = US) -- does not have corresponding drive but uses DriveVals.  Number of active ones is Drive.NNegUSs -- the first is always reserved for the accumulated effort cost / dissapointment when an expected US is not achieved"`
	VSPatch DriveVals `inactive:"+" view:"inline" desc:"current positive-valence drive-satisfying reward predicting VSPatch (PosD1) values"`

	HasRewPrev   slbool.Bool `inactive:"+" desc:"HasRew state from the previous trial -- copied from HasRew in NewState -- used for updating Effort, Urgency at start of new trial"`
	HasPosUSPrev slbool.Bool `inactive:"+" desc:"HasPosUS state from the previous trial -- copied from HasPosUS in NewState -- used for updating Effort, Urgency at start of new trial"`
	pad, pad1    int32
}

func (pp *PVLV) Defaults() {
	pp.Drive.Defaults()
	pp.Effort.Defaults()
	pp.Urgency.Defaults()
	pp.VTA.Defaults()
	pp.LHb.Defaults()
	pp.USpos.Zero()
	pp.USneg.Zero()
	pp.VSPatch.Zero()
}

func (pp *PVLV) Update() {
	pp.Drive.Update()
	pp.Effort.Update()
	pp.Urgency.Update()
	pp.VTA.Update()
	pp.LHb.Update()
}

func (pp *PVLV) Reset(ctx *Context, di uint32) {
	pp.Drive.ToZero()
	pp.Effort.Reset(ctx, di)
	pp.Urgency.Reset(ctx, di)
	pp.LHb.Reset(ctx, di)
	pp.VTA.Reset(ctx, di)
	pp.USpos.Zero()
	pp.USneg.Zero()
	pp.VSPatch.Zero()
	SetGlobalV(ctx, di, GvVSMatrixJustGated, 0)
	SetGlobalV(ctx, di, GvVSMatrixHasGated, 0)
	SetGlobalV(ctx, di, GvHasRewPrev, 0)
	SetGlobalV(ctx, di, GvHasPosUS, 0)
	// pp.HasPosUSPrev.SetBool(false)
}

// NewState is called at start of new state (trial) of processing.
// hadRew indicates if there was a reward state the previous trial.
// It calls LHGiveUpFmSum to trigger a "give up" state on this trial
// if previous expectation of reward exceeds critical sum.
func (pp *PVLV) NewState(ctx *Context, di uint32, hadRew bool) {
	SetGlobalV(ctx, di, GvHasRewPrev, bools.ToFloat32(hadRew))
	pp.HasPosUSPrev.SetBool(pp.HasPosUS(ctx, di))

	if hadRew {
		SetGlobalV(ctx, di, GvVSMatrixHasGated, 0)
	} else if GlobalV(ctx, di, GvVSMatrixJustGated) > 0 {
		SetGlobalV(ctx, di, GvVSMatrixHasGated, 1)
	}
	SetGlobalV(ctx, di, GvVSMatrixJustGated, 0)
}

// InitUS initializes all the USs to zero
func (pp *PVLV) InitUS(ctx *Context, di uint32) {
	pp.Drive.Zero(ctx, di, GvUSpos)
	pp.Drive.Zero(ctx, di, GvUSneg)
}

// SetPosUS sets given positive US (associated with same-indexed Drive) to given value
func (pp *PVLV) SetPosUS(ctx *Context, di uint32, usn uint32, val float32) {
	SetGlobalDriveV(ctx, di, usn, GvUSpos, val)
}

// SetNegUS sets given negative US to given value
func (pp *PVLV) SetNegUS(ctx *Context, di uint32, usn uint32, val float32) {
	SetGlobalDriveV(ctx, di, usn, GvUSneg, val)
}

// InitDrives initializes all the Drives to zero
func (pp *PVLV) InitDrives(ctx *Context, di uint32) {
	pp.Drive.Zero(ctx, di, GvDrives)
}

// SetDrive sets given Drive to given value
func (pp *PVLV) SetDrive(ctx *Context, di uint32, dr uint32, val float32) {
	SetGlobalDriveV(ctx, di, dr, GvDrives, val)
}

// USStimVal returns stimulus value for US at given index
// and valence.  If US > 0.01, a full 1 US activation is returned.
func (pp *PVLV) USStimVal(ctx *Context, di uint32, usIdx uint32, valence ValenceTypes) float32 {
	us := float32(0)
	if valence == Positive {
		us = GlobalDriveV(ctx, di, usIdx, GvUSpos)
	} else {
		us = GlobalUSneg(ctx, di, usIdx)
	}
	if us > 0.01 { // threshold for presentation to net
		us = 1 // https://github.com/emer/axon/issues/194
	}
	return us
}

// PosPV returns the reward for current positive US state relative to current drives
func (pp *PVLV) PosPV(ctx *Context, di uint32) float32 {
	rew := float32(0)
	for i := uint32(0); i < pp.Drive.NActive; i++ {
		rew += GlobalDriveV(ctx, di, i, GvUSpos) * pp.Drive.EffectiveDrive(ctx, di, i)
	}
	return rew
}

// NegPV returns the reward for current negative US state -- just a sum of USneg
func (pp *PVLV) NegPV(ctx *Context, di uint32) float32 {
	rew := float32(0)
	for i := uint32(0); i < pp.Drive.NNegUSs; i++ {
		rew += GlobalDriveV(ctx, di, i, GvUSneg)
	}
	return rew
}

// VSPatchMax returns the max VSPatch value across drives
func (pp *PVLV) VSPatchMax(ctx *Context, di uint32) float32 {
	max := float32(0)
	for i := uint32(0); i < pp.Drive.NActive; i++ {
		vs := GlobalDriveV(ctx, di, i, GvVSPatch)
		if vs > max {
			max = vs
		}
	}
	return max
}

// HasPosUS returns true if there is at least one non-zero positive US
func (pp *PVLV) HasPosUS(ctx *Context, di uint32) bool {
	for i := uint32(0); i < pp.Drive.NActive; i++ {
		if GlobalDriveV(ctx, di, i, GvUSpos) > 0 {
			return true
		}
	}
	return false
}

// HasNegUS returns true if there is at least one non-zero negative US
func (pp *PVLV) HasNegUS(ctx *Context, di uint32) bool {
	for i := uint32(0); i < pp.Drive.NActive; i++ {
		if GlobalDriveV(ctx, di, i, GvUSpos) > 0 {
			return true
		}
	}
	return false
}

// NetPV returns VTA.Vals.PVpos - VTA.Vals.PVneg
func (pp *PVLV) NetPV(ctx *Context, di uint32) float32 {
	return GlobalVTA(ctx, di, GvVtaVals, GvVtaPVpos) - GlobalVTA(ctx, di, GvVtaVals, GvVtaPVneg)
}

// PosPVFmDriveEffort returns the net primary value ("reward") based on
// given US value and drive for that value (typically in 0-1 range),
// and total effort, from which the effort discount factor is computed an applied:
// usValue * drive * Effort.DiscFun(effort)
func (pp *PVLV) PosPVFmDriveEffort(usValue, drive, effort float32) float32 {
	return usValue * drive * pp.Effort.DiscFun(effort)
}

// DA computes the updated dopamine from all the current state,
// including ACh from LDT via Context.
// Call after setting USs, Effort, Drives, VSPatch vals etc.
// Resulting DA is in VTA.Vals.DA, and is returned
// (to be set to Context.NeuroMod.DA)
func (pp *PVLV) DA(ctx *Context, di uint32, ach float32, hasRew bool) float32 {
	usPos := pp.PosPV(ctx, di)
	pvNeg := pp.NegPV(ctx, di)
	giveUp := GlobalV(ctx, di, GvLHbGiveUp)
	effDisc := GlobalV(ctx, di, GvEffortDisc)
	if giveUp > 0 {
		pvNeg += 1.0 - effDisc // pay effort cost here..
	}
	pvPos := usPos * effDisc
	vsPatchPos := pp.VSPatchMax(ctx, di)
	pp.LHb.LHbFmPVVS(ctx, di, pvPos, pvNeg, vsPatchPos)

	SetGlobalVTA(ctx, di, GvVtaRaw, GvVtaUSpos, usPos)
	SetGlobalVTA(ctx, di, GvVtaRaw, GvVtaPVpos, pvPos)
	SetGlobalVTA(ctx, di, GvVtaRaw, GvVtaPVneg, pvNeg)
	SetGlobalVTA(ctx, di, GvVtaRaw, GvVtaLHbDip, GlobalV(ctx, di, GvLHbDip))
	SetGlobalVTA(ctx, di, GvVtaRaw, GvVtaLHbBurst, GlobalV(ctx, di, GvLHbBurst))
	SetGlobalVTA(ctx, di, GvVtaRaw, GvVtaVSPatchPos, vsPatchPos)

	pp.VTA.DAFmRaw(ctx, di, ach, hasRew)
	return GlobalVTA(ctx, di, GvVtaVals, GvVtaDA)
}

// DriveUpdt updates the drives based on the current USs,
// subtracting USDec * US from current Drive,
// and calling ExpStep with the Dt and Base params.
func (pp *PVLV) DriveUpdt(ctx *Context, di uint32) {
	pp.Drive.ExpStepAll(ctx, di)
	for i := uint32(0); i < pp.Drive.NActive; i++ {
		us := GlobalDriveV(ctx, di, i, GvUSpos)
		AddGlobalDriveV(ctx, di, i, GvDrives, -us*pp.Drive.USDec.Get(i))
	}
}

// UrgencyUpdt updates the urgency and urgency based on given effort increment,
// resetting instead if HasRewPrev and HasPosUSPrev is true indicating receipt
// of an actual positive US.
// Call this at the start of the trial, in ApplyPVLV method.
func (pp *PVLV) UrgencyUpdt(ctx *Context, di uint32, effort float32) {
	if (GlobalV(ctx, di, GvHasRewPrev) > 0) && (GlobalV(ctx, di, GvHasPosUS) > 0) {
		pp.Urgency.Reset(ctx, di)
	} else {
		pp.Urgency.AddEffort(ctx, di, effort)
	}
}

//gosl: end pvlv

// PlusVar returns value plus random variance
func (ef *Effort) PlusVar(rnd erand.Rand, max float32) float32 {
	if ef.MaxVar == 0 {
		return max
	}
	return max + ef.MaxVar*float32(rnd.NormFloat64(-1))
}

// ReStart restarts restarts the raw effort back to zero
// and sets the Max with random additional variance.
func (ef *Effort) ReStart(ctx *Context, di uint32, rnd erand.Rand) {
	SetGlobalV(ctx, di, GvEffortRaw, 0)
	SetGlobalV(ctx, di, GvEffortCurMax, ef.PlusVar(rnd, ef.Max))
	SetGlobalV(ctx, di, GvEffortDisc, 1)
}

// VSGated updates JustGated and HasGated as function of VS
// (ventral striatum / ventral pallidum) gating at end of the plus phase.
// Also resets effort and LHb.DipSumCur counters -- starting fresh at start
// of a new goal engaged state.
func (pp *PVLV) VSGated(ctx *Context, di uint32, rnd erand.Rand, gated, hasRew bool, poolIdx int) {
	if !hasRew && gated {
		pp.Effort.ReStart(ctx, di, rnd)
		SetGlobalV(ctx, di, GvLHbDipSumCur, 0)
		if poolIdx == 0 { // novelty / curiosity pool
			SetGlobalV(ctx, di, GvEffortCurMax, pp.Effort.MaxNovel)
		}
	}
	SetGlobalV(ctx, di, GvVSMatrixJustGated, bools.ToFloat32(gated))
}

// ShouldGiveUp tests whether it is time to give up on the current goal,
// based on sum of LHb Dip (missed expected rewards) and maximum effort.
func (pp *PVLV) ShouldGiveUp(ctx *Context, di uint32, rnd erand.Rand, hasRew bool) bool {
	SetGlobalV(ctx, di, GvLHbGiveUp, 0)
	if hasRew { // can't give up if got something now
		SetGlobalV(ctx, di, GvLHbDipSumCur, 0)
		return false
	}
	prevSum := GlobalV(ctx, di, GvLHbDipSumCur)
	giveUp := pp.LHb.ShouldGiveUp(ctx, di)
	if prevSum < pp.LHb.DipLowThr && GlobalV(ctx, di, GvLHbDipSumCur) >= pp.LHb.DipLowThr {
		SetGlobalV(ctx, di, GvEffortCurMax, pp.Effort.PlusVar(rnd, GlobalV(ctx, di, GvEffortRaw)+pp.Effort.MaxPostDip))
	}
	if pp.Effort.GiveUp(ctx, di) {
		SetGlobalV(ctx, di, GvLHbGiveUp, 1)
		giveUp = true
	}
	return giveUp
}

// EffortUpdt updates the effort based on given effort increment,
// resetting instead if HasRewPrev flag is true.
// Call this at the start of the trial, in ApplyPVLV method.
func (pp *PVLV) EffortUpdt(ctx *Context, di uint32, rnd erand.Rand, effort float32) {
	if GlobalV(ctx, di, GvHasRewPrev) > 0 {
		pp.Effort.ReStart(ctx, di, rnd)
	} else {
		pp.Effort.AddEffort(ctx, di, effort)
	}
}

// EffortUrgencyUpdt updates the Effort & Urgency based on
// given effort increment, resetting instead if HasRewPrev flag is true.
// Call this at the start of the trial, in ApplyPVLV method.
func (pp *PVLV) EffortUrgencyUpdt(ctx *Context, di uint32, rnd erand.Rand, effort float32) {
	pp.EffortUpdt(ctx, di, rnd, effort)
	pp.UrgencyUpdt(ctx, di, effort)
}
