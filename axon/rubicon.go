// Copyright (c) 2022, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package axon

import (
	"fmt"
	"log/slog"

	"cogentcore.org/core/base/num"
	"cogentcore.org/core/base/randx"
	"cogentcore.org/core/math32"
	"github.com/emer/emergent/v2/popcode"
)

// DriveParams manages the drive parameters for computing and updating drive state.
// Most of the params are for optional case where drives are automatically
// updated based on US consumption (which satisfies drives) and time passing
// (which increases drives).
type DriveParams struct {

	// minimum effective drive value, which is an automatic baseline ensuring
	// that a positive US results in at least some minimal level of reward.
	// Unlike Base values, this is not reflected in the activity of the drive
	// values, and applies at the time of reward calculation as a minimum baseline.
	DriveMin float32

	// baseline levels for each drive, which is what they naturally trend toward
	// in the absence of any input.  Set inactive drives to 0 baseline,
	// active ones typically elevated baseline (0-1 range).
	Base []float32

	// time constants in ThetaCycle (trial) units for natural update toward
	// Base values. 0 values means no natural update (can be updated externally).
	Tau []float32

	// decrement in drive value when US is consumed, thus partially satisfying
	// the drive. Positive values are subtracted from current Drive value.
	Satisfaction []float32

	// 1/Tau
	Dt []float32 `display:"-"`
}

func (dp *DriveParams) Alloc(nDrives int) {
	if len(dp.Base) == nDrives {
		return
	}
	dp.Base = make([]float32, nDrives)
	dp.Tau = make([]float32, nDrives)
	dp.Dt = make([]float32, nDrives)
	dp.Satisfaction = make([]float32, nDrives)
}

func (dp *DriveParams) Defaults() {
	dp.DriveMin = 0.5
	for i := range dp.Satisfaction {
		dp.Satisfaction[i] = 0
	}
	dp.Update()
}

func (dp *DriveParams) Update() {
	for i, tau := range dp.Tau {
		if tau <= 0 {
			dp.Dt[i] = 0
		} else {
			dp.Dt[i] = 1 / tau
		}
	}
}

// VarToZero sets all values of given drive-sized variable to 0
func (dp *DriveParams) VarToZero(ctx *Context, di uint32, gvar GlobalVars) {
	for i := range dp.Base {
		SetGlbUSposV(ctx, di, gvar, uint32(i), 0)
	}
}

// ToZero sets all drives to 0
func (dp *DriveParams) ToZero(ctx *Context, di uint32) {
	dp.VarToZero(ctx, di, GvDrives)
}

// ToBaseline sets all drives to their baseline levels
func (dp *DriveParams) ToBaseline(ctx *Context, di uint32) {
	for i := range dp.Base {
		SetGlbUSposV(ctx, di, GvDrives, uint32(i), dp.Base[i])
	}
}

// AddTo increments drive by given amount, subject to 0-1 range clamping.
// Returns new val.
func (dp *DriveParams) AddTo(ctx *Context, di uint32, drv uint32, delta float32) float32 {
	dv := GlbUSposV(ctx, di, GvDrives, drv) + delta
	if dv > 1 {
		dv = 1
	} else if dv < 0 {
		dv = 0
	}
	SetGlbUSposV(ctx, di, GvDrives, drv, dv)
	return dv
}

// SoftAdd increments drive by given amount, using soft-bounding to 0-1 extremes.
// if delta is positive, multiply by 1-val, else val.  Returns new val.
func (dp *DriveParams) SoftAdd(ctx *Context, di uint32, drv uint32, delta float32) float32 {
	dv := GlbUSposV(ctx, di, GvDrives, drv)
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
	SetGlbUSposV(ctx, di, GvDrives, drv, dv)
	return dv
}

// ExpStep updates drive with an exponential step with given dt value
// toward given baseline value.
func (dp *DriveParams) ExpStep(ctx *Context, di uint32, drv uint32, dt, base float32) float32 {
	dv := GlbUSposV(ctx, di, GvDrives, drv)
	dv += dt * (base - dv)
	if dv > 1 {
		dv = 1
	} else if dv < 0 {
		dv = 0
	}
	SetGlbUSposV(ctx, di, GvDrives, drv, dv)
	return dv
}

// ExpStepAll updates given drives with an exponential step using dt values
// toward baseline values.
func (dp *DriveParams) ExpStepAll(ctx *Context, di uint32) {
	for i := range dp.Base {
		dp.ExpStep(ctx, di, uint32(i), dp.Dt[i], dp.Base[i])
	}
}

// EffectiveDrive returns the Max of Drives at given index and DriveMin.
// note that index 0 is the novelty / curiosity drive, which doesn't use DriveMin.
func (dp *DriveParams) EffectiveDrive(ctx *Context, di uint32, i uint32) float32 {
	if i == 0 {
		return GlbUSposV(ctx, di, GvDrives, uint32(0))
	}
	return math32.Max(GlbUSposV(ctx, di, GvDrives, i), dp.DriveMin)
}

///////////////////////////////////////////////////////////////////////////////
//  UrgencyParams

// UrgencyParams has urgency (increasing pressure to do something)
// and parameters for updating it.
// Raw urgency integrates effort when _not_ goal engaged
// while effort (negative US 0) integrates when a goal _is_ engaged.
type UrgencyParams struct {

	// value of raw urgency where the urgency activation level is 50%
	U50 float32

	// exponent on the urge factor -- valid numbers are 1,2,4,6
	Power int32 `default:"4"`

	// threshold for urge -- cuts off small baseline values
	Thr float32 `default:"0.2"`

	// gain factor for driving tonic DA levels as a function of urgency
	DAtonic float32 `default:"50"`
}

func (ur *UrgencyParams) Defaults() {
	ur.U50 = 10
	ur.Power = 4
	ur.Thr = 0.2
	ur.DAtonic = 50
}

func (ur *UrgencyParams) Update() {

}

// UrgeFun is the urgency function: urgency / (urgency + 1) where
// urgency = (Raw / U50)^Power
func (ur *UrgencyParams) UrgeFun(urgency float32) float32 {
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

// Reset resets the raw urgency back to zero -- at start of new gating event
func (ur *UrgencyParams) Reset(ctx *Context, di uint32) {
	SetGlbV(ctx, di, GvUrgencyRaw, 0)
	SetGlbV(ctx, di, GvUrgency, 0)
}

// Urge computes normalized Urge value from Raw, and sets DAtonic from that
func (ur *UrgencyParams) Urge(ctx *Context, di uint32) float32 {
	urge := ur.UrgeFun(GlbV(ctx, di, GvUrgencyRaw))
	if urge < ur.Thr {
		urge = 0
	}
	SetGlbV(ctx, di, GvUrgency, urge)
	SetGlbV(ctx, di, GvDAtonic, ur.DAtonic*urge) // simple equation for now
	return urge
}

// AddEffort adds an effort increment of urgency and updates the Urge factor
func (ur *UrgencyParams) AddEffort(ctx *Context, di uint32, inc float32) {
	AddGlbV(ctx, di, GvUrgencyRaw, inc)
	ur.Urge(ctx, di)
}

/////////////////////////////////////////////////////////
// USParams

// RubiconLNormFun is the normalizing function applied to the sum of all
// weighted raw values: 1 - (1 / (1 + usRaw.Sum()))
func RubiconNormFun(raw float32) float32 {
	return 1.0 - (1.0 / (1.0 + raw))
}

// USParams control how positive and negative USs and Costs are
// weighted and integrated to compute an overall PV primary value.
type USParams struct {

	// gain factor applied to sum of weighted, drive-scaled positive USs
	// to compute PVpos primary summary value.
	// This is multiplied prior to 1/(1+x) normalization.
	// Use this to adjust the overall scaling of PVpos reward within 0-1
	// normalized range (see also PVnegGain).
	// Each USpos is assumed to be in 0-1 range, with a default of 1.
	PVposGain float32 `default:"2"`

	// gain factor applied to sum of weighted negative USs and Costs
	// to compute PVneg primary summary value.
	// This is multiplied prior to 1/(1+x) normalization.
	// Use this to adjust overall scaling of PVneg within 0-1
	// normalized range (see also PVposGain).
	PVnegGain float32 `default:"1"`

	// Negative US gain factor for encoding each individual negative US,
	// within their own separate input pools, multiplied prior to 1/(1+x)
	// normalization of each term for activating the USneg pools.
	// These gains are _not_ applied in computing summary PVneg value
	// (see PVnegWts), and generally must be larger than the weights to leverage
	// the dynamic range within each US pool.
	USnegGains []float32

	// Cost gain factor for encoding the individual Time, Effort etc costs
	// within their own separate input pools, multiplied prior to 1/(1+x)
	// normalization of each term for activating the Cost pools.
	// These gains are _not_ applied in computing summary PVneg value
	// (see CostWts), and generally must be larger than the weights to use
	// the full dynamic range within each US pool.
	CostGains []float32

	// weight factor applied to each separate positive US on the way to computing
	// the overall PVpos summary value, to control the weighting of each US
	// relative to the others. Each pos US is also multiplied by its dynamic
	// Drive factor as well.
	// Use PVposGain to control the overall scaling of the PVpos value.
	PVposWts []float32

	// weight factor applied to each separate negative US on the way to computing
	// the overall PVneg summary value, to control the weighting of each US
	// relative to the others, and to the Costs.  These default to 1.
	PVnegWts []float32

	// weight factor applied to each separate Cost (Time, Effort, etc) on the
	// way to computing the overall PVneg summary value, to control the weighting
	// of each Cost relative to the others, and relative to the negative USs.
	// The first pool is Time, second is Effort, and these are typically weighted
	// lower (.02) than salient simulation-specific USs (1).
	PVcostWts []float32

	// computed estimated US values, based on OFCposPT and VSMatrix gating, in PVposEst
	USposEst []float32 `edit:"-"`
}

func (us *USParams) Alloc(nPos, nNeg, nCost int) {
	if len(us.PVposWts) != nPos {
		us.PVposWts = make([]float32, nPos)
		us.USposEst = make([]float32, nPos)
	}
	if len(us.PVnegWts) != nNeg {
		us.USnegGains = make([]float32, nNeg)
		us.PVnegWts = make([]float32, nNeg)
	}
	if len(us.PVcostWts) != nCost {
		us.CostGains = make([]float32, nCost)
		us.PVcostWts = make([]float32, nCost)
	}
}

func (us *USParams) Defaults() {
	us.PVposGain = 2
	us.PVnegGain = 1
	for i := range us.PVposWts {
		us.PVposWts[i] = 1
	}
	for i := range us.USnegGains {
		us.USnegGains[i] = 2
		us.PVnegWts[i] = 1
	}
	for i := range us.CostGains {
		us.CostGains[i] = 0.1
		us.PVcostWts[i] = 0.02
	}
}

func (us *USParams) Update() {
}

// USnegCostFromRaw sets normalized NegUS, Cost values from Raw values
func (us *USParams) USnegCostFromRaw(ctx *Context, di uint32) {
	for i, ng := range us.USnegGains {
		raw := GlbUSnegV(ctx, di, GvUSnegRaw, uint32(i))
		norm := RubiconNormFun(ng * raw)
		SetGlbUSnegV(ctx, di, GvUSneg, uint32(i), norm)
	}
	for i, ng := range us.CostGains {
		raw := GlbCostV(ctx, di, GvCostRaw, uint32(i))
		norm := RubiconNormFun(ng * raw)
		SetGlbCostV(ctx, di, GvCost, uint32(i), norm)
	}
}

// USnegToZero sets all values of USneg, USnegRaw to zero
func (us *USParams) USnegToZero(ctx *Context, di uint32) {
	for i := range us.USnegGains {
		SetGlbUSnegV(ctx, di, GvUSneg, uint32(i), 0)
		SetGlbUSnegV(ctx, di, GvUSnegRaw, uint32(i), 0)
	}
}

// CostToZero sets all values of Cost, CostRaw to zero
func (us *USParams) CostToZero(ctx *Context, di uint32) {
	for i := range us.CostGains {
		SetGlbCostV(ctx, di, GvCost, uint32(i), 0)
		SetGlbCostV(ctx, di, GvCostRaw, uint32(i), 0)
	}
}

// USposToZero sets all values of USpos to zero
func (us *USParams) USposToZero(ctx *Context, di uint32) {
	for i := range us.PVposWts {
		SetGlbUSposV(ctx, di, GvUSpos, uint32(i), 0)
	}
}

///////////////////////////////////////////////////////////////////
//  LHb & RMTg

// LHbParams has values for computing LHb & RMTg which drives dips / pauses in DA firing.
// LHb handles all US-related (PV = primary value) processing.
// Positive net LHb activity drives dips / pauses in VTA DA activity,
// e.g., when predicted pos > actual or actual neg > predicted.
// Negative net LHb activity drives bursts in VTA DA activity,
// e.g., when actual pos > predicted (redundant with LV / Amygdala)
// or "relief" burst when actual neg < predicted.
type LHbParams struct {

	// threshold on VSPatch prediction during a non-reward trial
	VSPatchNonRewThr float32 `default:"0.1"`

	// gain on the VSPatchD1 - D2 difference to drive the net VSPatch DA
	// prediction signal, which goes in VSPatchPos and RewPred global variables
	VSPatchGain float32 `default:"4"`

	// decay time constant for computing the temporal variance in VSPatch
	// values over time
	VSPatchVarTau float32 `default:"2"`

	// threshold factor that multiplies integrated pvNeg value
	// to establish a threshold for whether the integrated pvPos value
	// is good enough to drive overall net positive reward.
	// If pvPos wins, it is then multiplicatively discounted by pvNeg;
	// otherwise, pvNeg is discounted by pvPos.
	NegThr float32 `default:"1"`

	// gain multiplier on PVpos for purposes of generating bursts
	// (not for discounting negative dips).
	BurstGain float32 `default:"1"`

	// gain multiplier on PVneg for purposes of generating dips
	// (not for discounting positive bursts).
	DipGain float32 `default:"1"`

	// 1/tau
	VSPatchVarDt float32 `display:"-"`
}

func (lh *LHbParams) Defaults() {
	lh.VSPatchNonRewThr = 0.1
	lh.VSPatchGain = 4
	lh.VSPatchVarTau = 2
	lh.NegThr = 1
	lh.BurstGain = 1
	lh.DipGain = 1
}

func (lh *LHbParams) Update() {
	lh.VSPatchVarDt = 1 / lh.VSPatchVarTau
}

// Reset resets all LHb vars back to 0
func (lh *LHbParams) Reset(ctx *Context, di uint32) {
	SetGlbV(ctx, di, GvLHbDip, 0)
	SetGlbV(ctx, di, GvLHbBurst, 0)
	SetGlbV(ctx, di, GvLHbPVDA, 0)
	SetGlbV(ctx, di, GvVSPatchPosRPE, 0)
}

// DAFromPVs computes the overall PV DA in terms of LHb burst and dip
// activity from given pvPos, pvNeg, and vsPatchPos values.
// Also returns the net "reward" value as the discounted PV value,
// separate from the vsPatchPos prediction error factor.
func (lh *LHbParams) DAFromPVs(pvPos, pvNeg, vsPatchPos, vsPatchPosSum float32) (burst, dip, da, rew float32) {
	thr := lh.NegThr * pvNeg
	net := pvPos - thr // if > 0, net positive outcome; else net negative (not worth it)
	if net > 0 {       // worth it
		rew = lh.BurstGain * pvPos * (1 - pvNeg) // positive reward value: pos with mult neg discount factor
		rpe := rew - vsPatchPos                  // prediction error relative to pos reward value
		if rpe < 0 {
			dip = -rpe // positive dip = negative value
		} else {
			burst = rpe
		}
	} else { // not worth it: net negative but moderated (discounted) by strength of positive
		// also ensure that we at least get the accumulated expectation from vspatch
		neg := max(lh.DipGain*pvNeg*(1-pvPos), vsPatchPosSum)
		rew = -neg
		dip = neg // magnitude
	}
	da = burst - dip
	return
}

// DAforUS computes the overall LHb Dip or Burst (one is always 0),
// and PVDA ~= Burst - Dip, for case when there is a primary
// positive reward value or a give-up state has triggered.
// Returns the overall net reward magnitude, prior to VSPatch discounting.
func (lh *LHbParams) DAforUS(ctx *Context, di uint32, pvPos, pvNeg, vsPatchPos, vsPatchPosSum float32) float32 {
	burst, dip, da, rew := lh.DAFromPVs(pvPos, pvNeg, vsPatchPos, vsPatchPosSum)
	SetGlbV(ctx, di, GvLHbDip, dip)
	SetGlbV(ctx, di, GvLHbBurst, burst)
	SetGlbV(ctx, di, GvLHbPVDA, da)
	SetGlbV(ctx, di, GvVSPatchPosThr, vsPatchPos) // no thresholding for US
	SetGlbV(ctx, di, GvVSPatchPosRPE, da)
	return rew
}

// DAforNoUS computes the LHb response when there is _NOT_ a primary
// positive reward value or a give-up state.
// In this case, inhibition of VS via tonic ACh is assumed to prevent
// activity of PVneg (and there is no PVpos).
// Because the LHb only responds when it decides to GiveUp,
// there is no response in this case.
// DA is instead driven by CS-based computation, in rubicon_layers.go, VTAParams.VTADA
func (lh *LHbParams) DAforNoUS(ctx *Context, di uint32) float32 {
	SetGlbV(ctx, di, GvLHbDip, 0)
	SetGlbV(ctx, di, GvLHbBurst, 0)
	SetGlbV(ctx, di, GvLHbPVDA, 0)
	return 0
}

//////////////////////////////////////////////////////////
//  GiveUpParams

// GiveUpParams are parameters for computing when to give up,
// based on Utility, Timing and Progress factors.
type GiveUpParams struct {

	// threshold on GiveUp probability, below which no give up is triggered
	ProbThr float32 `default:"0.5"`

	// minimum GiveUpSum value, which is the denominator in the sigmoidal function.
	// This minimum prevents division by zero and any other degenerate values.
	MinGiveUpSum float32 `default:"0.1"`

	// the factor multiplying utility values: cost and expected positive outcome
	Utility float32 `default:"1"`

	// the factor multiplying timing values from VSPatch
	Timing float32 `default:"2"`

	// the factor multiplying progress values based on time-integrated progress
	// toward the goal
	Progress float32 `default:"1"`

	// minimum utility cost and reward estimate values -- when they are below
	// these levels (at the start) then utility is effectively neutral,
	// so the other factors take precedence.
	MinUtility float32 `default:"0.2"`

	// maximum VSPatchPosSum for normalizing the value for give-up weighing
	VSPatchSumMax float32 `default:"1"`

	// maximum VSPatchPosVar for normalizing the value for give-up weighing
	VSPatchVarMax float32 `default:"0.5"`

	// time constant for integrating the ProgressRate
	// values over time
	ProgressRateTau float32 `default:"2"`

	// 1/tau
	ProgressRateDt float32 `display:"-"`
}

func (gp *GiveUpParams) Defaults() {
	gp.ProbThr = 0.5
	gp.MinGiveUpSum = 0.1
	gp.Utility = 1
	gp.Timing = 2
	gp.Progress = 1
	gp.MinUtility = 0.2
	gp.VSPatchSumMax = 1
	gp.VSPatchVarMax = 0.5
	gp.ProgressRateTau = 2
}

func (gp *GiveUpParams) Update() {
	gp.ProgressRateDt = 1 / gp.ProgressRateTau
}

// SigmoidFun is the sigmoid function for computing give up probabilities
func SigmoidFun(cnSum, guSum float32) float32 {
	return 1 / (1 + (cnSum / guSum))
}

// Prob returns the probability and discrete bool give up for giving up based on
// given sums of continue and give up factors
func (gp *GiveUpParams) Prob(cnSum, guSum float32, rnd randx.Rand) (float32, bool) {
	prob := SigmoidFun(cnSum, guSum)
	giveUp := randx.BoolP32(prob, rnd)
	if prob <= gp.ProbThr {
		giveUp = false
	}
	return prob, giveUp
}

// Sums computes the summed weighting factors that drive continue and give up
// contributions to the probability function.
func (gp *GiveUpParams) Sums(ctx *Context, di uint32) (cnSum, guSum float32) {
	negSum := GlbV(ctx, di, GvPVnegSum)
	guU := gp.Utility * max(gp.MinUtility, negSum)
	cnU := gp.Utility * max(gp.MinUtility, GlbV(ctx, di, GvPVposEst)) // todo: var?
	SetGlbV(ctx, di, GvGiveUpUtility, guU)
	SetGlbV(ctx, di, GvContUtility, cnU)

	vspSum := min(GlbV(ctx, di, GvVSPatchPosSum), gp.VSPatchSumMax) / gp.VSPatchSumMax
	vspVar := min(GlbV(ctx, di, GvVSPatchPosVar), gp.VSPatchVarMax) / gp.VSPatchVarMax
	guT := gp.Timing * vspSum * (1 - vspVar)
	cnT := gp.Timing * (1 - vspSum) * vspVar
	SetGlbV(ctx, di, GvGiveUpTiming, guT)
	SetGlbV(ctx, di, GvContTiming, cnT)

	prog := GlbV(ctx, di, GvProgressRate)
	guP := -gp.Progress * prog
	cnP := gp.Progress * prog
	SetGlbV(ctx, di, GvGiveUpProgress, guP)
	SetGlbV(ctx, di, GvContProgress, cnP)

	guSum = guU + guT + guP
	cnSum = cnU + cnT + cnP
	SetGlbV(ctx, di, GvGiveUpSum, guSum)
	SetGlbV(ctx, di, GvContSum, cnSum)

	if guSum < gp.MinGiveUpSum {
		guSum = gp.MinGiveUpSum
	}

	return
}

//////////////////////////////////////////////////////////
//  Rubicon

// Rubicon implements core elements of the Rubicon goal-directed motivational
// model, representing the core brainstem-level (hypothalamus) bodily drives
// and resulting dopamine from US (unconditioned stimulus) inputs,
// subsuming the earlier Rubicon model of primary value (PV)
// and learned value (LV), describing the functions of the Amygala,
// Ventral Striatum, VTA and associated midbrain nuclei (LDT, LHb, RMTg).
// Core LHb (lateral habenula) and VTA (ventral tegmental area) dopamine
// are computed in equations using inputs from specialized network layers
// (LDTLayer driven by BLA, CeM layers, VSPatchLayer).
// The Drives, Effort, US and resulting LHb PV dopamine computation all happens at the
// at the start of each trial (NewState, Step).  The LV / CS dopamine is computed
// cycle-by-cycle by the VTA layer using parameters set by the VTA layer.
// Renders USLayer, PVLayer, DrivesLayer representations based on state updated here.
type Rubicon struct {

	// number of possible positive US states and corresponding drives.
	// The first is always reserved for novelty / curiosity.
	// Must be set programmatically via SetNUSs method,
	// which allocates corresponding parameters.
	NPosUSs uint32 `edit:"-"`

	// number of possible phasic negative US states (e.g., shock, impact etc).
	// Must be set programmatically via SetNUSs method, which allocates corresponding
	// parameters.
	NNegUSs uint32 `edit:"-"`

	// number of possible costs, typically including accumulated time and effort costs.
	// Must be set programmatically via SetNUSs method, which allocates corresponding
	// parameters.
	NCosts uint32 `edit:"-"`

	// parameters and state for built-in drives that form the core motivations
	// of the agent, controlled by lateral hypothalamus and associated
	// body state monitoring such as glucose levels and thirst.
	Drive DriveParams `display:"inline"`

	// urgency (increasing pressure to do something) and parameters for
	//  updating it. Raw urgency is incremented by same units as effort,
	// but is only reset with a positive US.
	Urgency UrgencyParams `display:"inline"`

	// controls how positive and negative USs are weighted and integrated to
	// compute an overall PV primary value.
	USs USParams `display:"inline"`

	// lateral habenula (LHb) parameters and state, which drives
	// dipping / pausing in dopamine when the predicted positive
	// outcome > actual, or actual negative outcome > predicted.
	// Can also drive bursting for the converse, and via matrix phasic firing.
	LHb LHbParams `display:"inline"`

	// parameters for giving up based on PV pos - neg difference
	GiveUp GiveUpParams `display:"inline"`

	// population code decoding parameters for estimates from layers
	ValDecode popcode.OneD `display:"inline"`

	decodeActs []float32
}

func (rp *Rubicon) Defaults() {
	if rp.LHb.VSPatchGain != 0 { // already done
		return
	}
	rp.Drive.Defaults()
	rp.Urgency.Defaults()
	rp.USs.Defaults()
	rp.LHb.Defaults()
	rp.GiveUp.Defaults()
	rp.ValDecode.Defaults()
	rp.Update()
}

func (rp *Rubicon) Update() {
	rp.Drive.Update()
	rp.Urgency.Update()
	rp.USs.Update()
	rp.LHb.Update()
	rp.GiveUp.Update()
}

// USposIndex adds 1 to the given _simulation specific_ positive US index
// to get the actual US / Drive index, where the first pool is reserved
// for curiosity / novelty.
func (rp *Rubicon) USposIndex(simUsIndex int) int {
	return simUsIndex + 1
}

// USnegIndex allows for the possibility of automatically managed
// negative USs, by adding those to the given _simulation specific_
// negative US index to get the actual US index.
func (rp *Rubicon) USnegIndex(simUsIndex int) int {
	return simUsIndex
}

// SetNUSs sets the number of _additional_ simulation-specific
// phasic positive and negative USs (primary value outcomes).
// This must be called _before_ network Build, which allocates global values
// that depend on these numbers.  Any change must also call network.BuildGlobals.
// 1 PosUS (curiosity / novelty) is managed automatically by the Rubicon code.
// Two costs (Time, Effort) are also automatically allocated and managed.
// The USs specified here need to be managed by the simulation via the SetUS method.
// Positive USs each have corresponding Drives.
func (rp *Rubicon) SetNUSs(ctx *Context, nPos, nNeg int) {
	nPos = rp.USposIndex(max(nPos, 1))
	nNeg = rp.USnegIndex(max(nNeg, 1)) // ensure at least 1
	rp.NPosUSs = uint32(nPos)
	rp.NNegUSs = uint32(nNeg)
	rp.NCosts = 2 // default
	ctx.NetIndexes.RubiconNPosUSs = rp.NPosUSs
	ctx.NetIndexes.RubiconNNegUSs = rp.NNegUSs
	ctx.NetIndexes.RubiconNCosts = rp.NCosts
	rp.Drive.Alloc(nPos)
	rp.USs.Alloc(nPos, nNeg, int(rp.NCosts))
}

// Reset resets all Rubicon state
func (rp *Rubicon) Reset(ctx *Context, di uint32) {
	rp.Drive.ToBaseline(ctx, di)
	rp.TimeEffortReset(ctx, di)
	rp.Urgency.Reset(ctx, di)
	rp.InitUS(ctx, di)
	rp.LHb.Reset(ctx, di)
	rp.Drive.VarToZero(ctx, di, GvVSPatchD1)
	rp.Drive.VarToZero(ctx, di, GvVSPatchD2)
	rp.ResetGoalState(ctx, di)
	SetGlbV(ctx, di, GvVtaDA, 0)
	SetGlbV(ctx, di, GvVSMatrixJustGated, 0)
	SetGlbV(ctx, di, GvVSMatrixHasGated, 0)
	SetGlbV(ctx, di, GvHadRew, 0)
	// pp.HasPosUSPrev.SetBool(false) // key to not reset!!
}

// InitUS initializes all the USs to zero
func (rp *Rubicon) InitUS(ctx *Context, di uint32) {
	rp.USs.USposToZero(ctx, di)
	rp.USs.USnegToZero(ctx, di)
	rp.USs.CostToZero(ctx, di)
	SetGlbV(ctx, di, GvHasRew, 0)
	SetGlbV(ctx, di, GvRew, 0)
}

// InitDrives initializes all the Drives to baseline values (default = 0)
func (rp *Rubicon) InitDrives(ctx *Context, di uint32) {
	rp.Drive.ToBaseline(ctx, di)
}

// AddTimeEffort adds a unit of time and an increment of effort
func (rp *Rubicon) AddTimeEffort(ctx *Context, di uint32, effort float32) {
	AddGlbV(ctx, di, GvTime, 1)
	tm := GlbV(ctx, di, GvTime)
	SetGlbCostV(ctx, di, GvCostRaw, 0, tm) // time is neg 0

	AddGlbV(ctx, di, GvEffort, effort)
	eff := GlbV(ctx, di, GvEffort)
	SetGlbCostV(ctx, di, GvCostRaw, 1, eff) // effort is neg 1
}

// EffortUrgencyUpdate updates the Effort or Urgency based on
// given effort increment.
// Effort is incremented when VSMatrixHasGated (i.e., goal engaged)
// and Urgency updates otherwise (when not goal engaged)
// Call this at the start of the trial, in ApplyRubicon method,
// after NewState.
func (rp *Rubicon) EffortUrgencyUpdate(ctx *Context, di uint32, effort float32) {
	if GlbV(ctx, di, GvVSMatrixHasGated) > 0 {
		rp.AddTimeEffort(ctx, di, effort)
	} else {
		rp.Urgency.AddEffort(ctx, di, effort)
	}
}

// TimeEffortReset resets the raw time and effort back to zero,
// at start of new gating event
func (rp *Rubicon) TimeEffortReset(ctx *Context, di uint32) {
	SetGlbV(ctx, di, GvTime, 0)
	SetGlbV(ctx, di, GvEffort, 0)
	SetGlbCostV(ctx, di, GvCostRaw, 0, 0) // effort is neg 0
	SetGlbCostV(ctx, di, GvCost, 0, 0)
}

// PVposFromDriveEffort returns the net primary value ("reward") based on
// given US value and drive for that value (typically in 0-1 range),
// and total effort, from which the effort discount factor is computed an applied:
// usValue * drive * Effort.DiscFun(effort).
// This is not called directly in the Rubicon code -- can be used to compute
// what the Rubicon code itself will compute -- see LHbPVDA
// todo: this is not very meaningful anymore
// func (pp *Rubicon) PVposFromDriveEffort(ctx *Context, usValue, drive, effort float32) float32 {
// 	return usValue * drive * (1 - RubiconNormFun(pp.USs.PVnegWts[0]*effort))
// }

// RubiconSetDrive sets given Drive to given value
func (rp *Rubicon) SetDrive(ctx *Context, di uint32, dr uint32, val float32) {
	SetGlbUSposV(ctx, di, GvDrives, dr, val)
}

// SetDrives is used when directly controlling drive levels externally.
// curiosity sets the strength for the curiosity drive
// and drives are strengths of the remaining sim-specified drives, in order.
// any drives not so specified are at the InitDrives baseline level.
func (rp *Rubicon) SetDrives(ctx *Context, di uint32, curiosity float32, drives ...float32) {
	rp.InitDrives(ctx, di)
	rp.SetDrive(ctx, di, 0, curiosity)
	for i, v := range drives {
		rp.SetDrive(ctx, di, uint32(1+i), v)
	}
}

// DriveUpdate is used when auto-updating drive levels based on US consumption,
// which partially satisfies (decrements) corresponding drive,
// and on time passing, where drives adapt to their overall baseline levels.
func (rp *Rubicon) DriveUpdate(ctx *Context, di uint32) {
	rp.Drive.ExpStepAll(ctx, di)
	nd := rp.NPosUSs
	for i := uint32(0); i < nd; i++ {
		us := GlbUSposV(ctx, di, GvUSpos, i)
		nwdrv := GlbUSposV(ctx, di, GvDrives, i) - us*rp.Drive.Satisfaction[i]
		if nwdrv < 0 {
			nwdrv = 0
		}
		SetGlbUSposV(ctx, di, GvDrives, i, nwdrv)
	}
}

// SetUS sets the given _simulation specific_ unconditioned
// stimulus (US) state for Rubicon algorithm.  usIndex = 0 is first US, etc.
// The US then drives activity of relevant Rubicon-rendered inputs, and dopamine,
// and sets the global HasRew flag, thus triggering a US learning event.
// Note that costs can be used to track negative USs that are not strong
// enough to trigger a US learning event.
func (rp *Rubicon) SetUS(ctx *Context, di uint32, valence ValenceTypes, usIndex int, magnitude float32) {
	SetGlbV(ctx, di, GvHasRew, 1)
	if valence == Positive {
		usIndex = rp.USposIndex(usIndex)
		SetGlbUSposV(ctx, di, GvUSpos, uint32(usIndex), magnitude)
	} else {
		usIndex = rp.USnegIndex(usIndex)
		SetGlbUSnegV(ctx, di, GvUSnegRaw, uint32(usIndex), magnitude)
		SetGlbV(ctx, di, GvNegUSOutcome, 1)
	}
}

// ResetGoalState resets all the goal-engaged global values.
// Critically, this is only called after goal accomplishment,
// not after goal gating -- prevents "shortcutting" by re-gating.
func (rp *Rubicon) ResetGoalState(ctx *Context, di uint32) {
	SetGlbV(ctx, di, GvVSMatrixHasGated, 0)
	rp.Urgency.Reset(ctx, di)
	rp.TimeEffortReset(ctx, di)
	rp.USs.USnegToZero(ctx, di) // all negs restart
	rp.USs.CostToZero(ctx, di)
	rp.ResetGiveUp(ctx, di)
	SetGlbV(ctx, di, GvVSPatchPos, 0)
	SetGlbV(ctx, di, GvVSPatchPosSum, 0)
	SetGlbV(ctx, di, GvVSPatchPosPrev, 0)
	SetGlbV(ctx, di, GvVSPatchPosVar, 0)
	SetGlbV(ctx, di, GvRewPred, 0)
	SetGlbV(ctx, di, GvGoalDistEst, 0)
	SetGlbV(ctx, di, GvGoalDistPrev, 0)
	SetGlbV(ctx, di, GvProgressRate, 0)
	nd := rp.NPosUSs
	for i := uint32(0); i < nd; i++ {
		SetGlbUSposV(ctx, di, GvOFCposPTMaint, i, 0)
		SetGlbUSposV(ctx, di, GvVSMatrixPoolGated, i, 0)
	}
}

// ResetGiveUp resets all the give-up related global values.
func (rp *Rubicon) ResetGiveUp(ctx *Context, di uint32) {
	SetGlbV(ctx, di, GvPVposEst, 0)
	SetGlbV(ctx, di, GvPVposVar, 0)
	SetGlbV(ctx, di, GvGiveUpProb, 0)
	SetGlbV(ctx, di, GvGiveUp, 0)
}

// NewState is called at very start of new state (trial) of processing.
// sets HadRew = HasRew from last trial -- used to then reset various things
// after reward.
func (rp *Rubicon) NewState(ctx *Context, di uint32, rnd randx.Rand) {
	hadRewF := GlbV(ctx, di, GvHasRew)
	hadRew := num.ToBool(hadRewF)
	SetGlbV(ctx, di, GvHadRew, hadRewF)
	SetGlbV(ctx, di, GvHadPosUS, GlbV(ctx, di, GvHasPosUS))
	SetGlbV(ctx, di, GvHadNegUSOutcome, GlbV(ctx, di, GvNegUSOutcome))
	SetGlbV(ctx, di, GvGaveUp, GlbV(ctx, di, GvGiveUp))

	SetGlbV(ctx, di, GvHasRew, 0)
	SetGlbV(ctx, di, GvNegUSOutcome, 0)

	rp.VSPatchNewState(ctx, di)

	if hadRew {
		rp.ResetGoalState(ctx, di)
	} else if GlbV(ctx, di, GvVSMatrixJustGated) > 0 {
		SetGlbV(ctx, di, GvVSMatrixHasGated, 1)
		rp.Urgency.Reset(ctx, di)
	}
	SetGlbV(ctx, di, GvVSMatrixJustGated, 0)
	rp.USs.USposToZero(ctx, di) // pos USs must be set fresh every time
}

// Step does one step (trial) after applying USs, Drives,
// and updating Effort.  It should be the final call in ApplyRubicon.
// Calls PVDA which does all US, PV, LHb, GiveUp updating.
func (rp *Rubicon) Step(ctx *Context, di uint32, rnd randx.Rand) {
	rp.PVDA(ctx, di, rnd)
}

// SetGoalMaintFromLayer sets the GoalMaint global state variable
// from the average activity (CaSpkD) of the given layer name.
// GoalMaint is normalized 0-1 based on the given max activity level,
// with anything out of range clamped to 0-1 range.
// Returns (and logs) an error if layer name not found.
func (rp *Rubicon) SetGoalMaintFromLayer(ctx *Context, di uint32, net *Network, layName string, maxAct float32) error {
	ly := net.AxonLayerByName(layName)
	if ly == nil {
		err := fmt.Errorf("SetGoalMaintFromLayer: layer named: %q not found", layName)
		slog.Error(err.Error())
		return err
	}
	act := ly.Pool(0, di).AvgMax.CaSpkD.Cycle.Avg
	gm := float32(0)
	if act > maxAct {
		gm = 1
	} else {
		gm = act / maxAct
	}
	SetGlbV(ctx, di, GvGoalMaint, gm)
	return nil
}

// DecodeFromLayer decodes value and variance from the average activity (CaSpkD)
// of the given layer name.  Use for decoding PVposEst and Var, and PVnegEst and Var
func (rp *Rubicon) DecodeFromLayer(ctx *Context, di uint32, net *Network, layName string) (val, vr float32, err error) {
	ly := net.AxonLayerByName(layName)
	if ly == nil {
		err = fmt.Errorf("DecodeFromLayer: layer named: %q not found", layName)
		slog.Error(err.Error())
		return
	}
	ly.UnitValues(&rp.decodeActs, "CaSpkD", int(di))
	val = ly.Params.Acts.PopCode.Decode(rp.decodeActs)
	vr = ly.Params.Acts.PopCode.Uncertainty(val, rp.decodeActs)
	return
}

// DecodePVEsts decodes estimated PV outcome values from PVposP and PVnegP
// prediction layers, saves in global PVposEst, Var and PVnegEst, Var
func (rp *Rubicon) DecodePVEsts(ctx *Context, di uint32, net *Network) {
	posEst, posVar, err := rp.DecodeFromLayer(ctx, di, net, "PVposP")
	if err == nil {
		SetGlbV(ctx, di, GvPVposEst, posEst)
		SetGlbV(ctx, di, GvPVposVar, posVar)
	}

	negEst, negVar, err := rp.DecodeFromLayer(ctx, di, net, "PVnegP")
	if err == nil {
		SetGlbV(ctx, di, GvPVnegEst, negEst)
		SetGlbV(ctx, di, GvPVnegVar, negVar)
	}
}

// SetGoalDistEst sets the current estimated distance to the goal,
// in trial step units, which should decrease to 0 at the goal.
// This should be set at the start of every trial.
// Also computes the ProgressRate.
func (rp *Rubicon) SetGoalDistEst(ctx *Context, di uint32, dist float32) {
	if GlbV(ctx, di, GvVSMatrixHasGated) == 0 {
		SetGlbV(ctx, di, GvGoalDistPrev, dist)
		SetGlbV(ctx, di, GvGoalDistEst, dist)
		SetGlbV(ctx, di, GvProgressRate, 0)
		return
	}
	prev := GlbV(ctx, di, GvGoalDistEst)
	SetGlbV(ctx, di, GvGoalDistPrev, prev)
	SetGlbV(ctx, di, GvGoalDistEst, dist)
	rate := GlbV(ctx, di, GvProgressRate)
	del := prev - dist
	rate += rp.GiveUp.ProgressRateDt * (del - rate)
	SetGlbV(ctx, di, GvProgressRate, rate)
}

//////////////////////////////////////////////////////////////////////////////////////
//    methods below used in computing Rubicon state, not generally called from sims

// HasPosUS returns true if there is at least one non-zero positive US
func (rp *Rubicon) HasPosUS(ctx *Context, di uint32) bool {
	nd := rp.NPosUSs
	for i := uint32(0); i < nd; i++ {
		if GlbUSposV(ctx, di, GvUSpos, i) > 0 {
			return true
		}
	}
	return false
}

// PVpos returns the summed weighted positive value of
// current positive US state, where each US is multiplied by
// its current drive and weighting factor (pvPosSum),
// and the normalized version of this sum (PVpos = overall positive PV)
// as 1 / (1 + (PVposGain * pvPosSum))
func (rp *Rubicon) PVpos(ctx *Context, di uint32) (pvPosSum, pvPos float32) {
	nd := rp.NPosUSs
	wts := rp.USs.PVposWts
	for i := uint32(0); i < nd; i++ {
		pvPosSum += wts[i] * GlbUSposV(ctx, di, GvUSpos, i) * rp.Drive.EffectiveDrive(ctx, di, i)
	}
	pvPos = RubiconNormFun(rp.USs.PVposGain * pvPosSum)
	return
}

// PVneg returns the summed weighted negative value of current
// costs and negative US state, where each US is multiplied
// by a weighting factor and summed (usNegSum)
// and the normalized version of this sum (PVneg = overall negative PV)
// as 1 / (1 + (PVnegGain * PVnegSum))
func (rp *Rubicon) PVneg(ctx *Context, di uint32) (pvNegSum, pvNeg float32) {
	nn := rp.NNegUSs
	wts := rp.USs.PVnegWts
	for i := uint32(0); i < nn; i++ {
		pvNegSum += wts[i] * GlbUSnegV(ctx, di, GvUSnegRaw, i)
	}
	nn = rp.NCosts
	wts = rp.USs.PVcostWts
	for i := uint32(0); i < nn; i++ {
		pvNegSum += wts[i] * GlbCostV(ctx, di, GvCostRaw, i)
	}
	pvNeg = RubiconNormFun(rp.USs.PVnegGain * pvNegSum)
	return
}

// PVsFromUSs updates the current PV summed, weighted, normalized values
// from the underlying US values.
func (rp *Rubicon) PVsFromUSs(ctx *Context, di uint32) {
	pvPosSum, pvPos := rp.PVpos(ctx, di)
	SetGlbV(ctx, di, GvPVposSum, pvPosSum)
	SetGlbV(ctx, di, GvPVpos, pvPos)
	SetGlbV(ctx, di, GvHasPosUS, num.FromBool[float32](rp.HasPosUS(ctx, di)))

	pvNegSum, pvNeg := rp.PVneg(ctx, di)
	SetGlbV(ctx, di, GvPVnegSum, pvNegSum)
	SetGlbV(ctx, di, GvPVneg, pvNeg)
}

// VSPatchNewState does VSPatch processing in NewState:
// updates global VSPatchPos and VSPatchPosSum, sets to RewPred.
// uses max across recorded VSPatch activity levels.
func (rp *Rubicon) VSPatchNewState(ctx *Context, di uint32) {
	prev := GlbV(ctx, di, GvVSPatchPos)
	SetGlbV(ctx, di, GvVSPatchPosPrev, prev)
	mx := float32(0)
	nd := rp.NPosUSs
	for i := uint32(0); i < nd; i++ {
		vsD1 := GlbUSposV(ctx, di, GvVSPatchD1, i)
		vsD2 := GlbUSposV(ctx, di, GvVSPatchD2, i)
		vs := rp.LHb.VSPatchGain * (vsD1 - vsD2)
		if vs > mx {
			mx = vs
		}
	}
	v := math32.Abs(mx - prev)
	pv := GlbV(ctx, di, GvVSPatchPosVar)
	if v > pv {
		pv = v
	} else {
		pv += rp.LHb.VSPatchVarDt * (v - pv) // decay -- negative
	}
	SetGlbV(ctx, di, GvVSPatchPosVar, pv)
	SetGlbV(ctx, di, GvVSPatchPos, mx)
	SetGlbV(ctx, di, GvRewPred, mx)
	thr := mx
	if mx < rp.LHb.VSPatchNonRewThr {
		thr = 0
	}
	SetGlbV(ctx, di, GvVSPatchPosRPE, -mx) // default for non-us cases
	SetGlbV(ctx, di, GvVSPatchPosThr, thr)
	AddGlbV(ctx, di, GvVSPatchPosSum, thr) // key: use thresholded!
}

// PVposEstFromUSs returns the estimated positive PV value
// based on drives and given US values.  This can be used
// to compute estimates to compare network performance.
func (rp *Rubicon) PVposEstFromUSs(ctx *Context, di uint32, uss []float32) (pvPosSum, pvPos float32) {
	nd := rp.NPosUSs
	if len(uss) < int(nd) {
		nd = uint32(len(uss))
	}
	wts := rp.USs.PVposWts
	for i := uint32(0); i < nd; i++ {
		pvPosSum += wts[i] * uss[i] * rp.Drive.EffectiveDrive(ctx, di, i)
	}
	pvPos = RubiconNormFun(rp.USs.PVposGain * pvPosSum)
	return
}

// PVposEstFromUSsDrives returns the estimated positive PV value
// based on given externally provided drives and US values.
// This can be used to compute estimates to compare network performance.
func (rp *Rubicon) PVposEstFromUSsDrives(uss, drives []float32) (pvPosSum, pvPos float32) {
	nd := rp.NPosUSs
	if len(uss) < int(nd) {
		nd = uint32(len(uss))
	}
	wts := rp.USs.PVposWts
	for i := uint32(0); i < nd; i++ {
		pvPosSum += wts[i] * uss[i] * drives[i]
	}
	pvPos = RubiconNormFun(rp.USs.PVposGain * pvPosSum)
	return
}

// PVnegEstFromUSs returns the estimated negative PV value
// based on given externally provided US values.
// This can be used to compute estimates to compare network performance.
func (rp *Rubicon) PVnegEstFromUSs(uss []float32) (pvNegSum, pvNeg float32) {
	nn := rp.NNegUSs
	wts := rp.USs.PVnegWts
	for i := uint32(0); i < nn; i++ {
		pvNegSum += wts[i] * uss[i]
	}
	pvNeg = RubiconNormFun(rp.USs.PVnegGain * pvNegSum)
	return
}

// PVcostEstFromUSs returns the estimated negative PV value
// based on given externally provided Cost values.
// This can be used to compute estimates to compare network performance.
func (rp *Rubicon) PVcostEstFromCosts(costs []float32) (pvCostSum, pvNeg float32) {
	nn := rp.NCosts
	wts := rp.USs.PVcostWts
	for i := uint32(0); i < nn; i++ {
		pvCostSum += wts[i] * costs[i]
	}
	pvNeg = RubiconNormFun(rp.USs.PVnegGain * pvCostSum)
	return
}

// DAFromPVs computes the overall PV DA in terms of LHb burst and dip
// activity from given pvPos, pvNeg, and vsPatchPos values.
// Also returns the net "reward" value as the discounted PV value,
// separate from the vsPatchPos prediction error factor.
func (rp *Rubicon) DAFromPVs(pvPos, pvNeg, vsPatchPos, vsPatchPosSum float32) (burst, dip, da, rew float32) {
	return rp.LHb.DAFromPVs(pvPos, pvNeg, vsPatchPos, vsPatchPosSum)
}

// GiveUpOnGoal determines whether to give up on current goal
// based on Utility, Timing, and Progress weight factors.
func (rp *Rubicon) GiveUpOnGoal(ctx *Context, di uint32, rnd randx.Rand) bool {
	cnSum, guSum := rp.GiveUp.Sums(ctx, di)
	prob, giveUp := rp.GiveUp.Prob(cnSum, guSum, rnd)
	SetGlbV(ctx, di, GvGiveUpProb, prob)
	SetGlbV(ctx, di, GvGiveUp, num.FromBool[float32](giveUp))
	return giveUp
}

// PVDA computes the PV (primary value) based dopamine
// based on current state information, at the start of a trial.
// PV DA is computed by the VS (ventral striatum) and the LHb / RMTg,
// and the resulting values are stored in global variables.
// Called after updating USs, Effort, Drives at start of trial step,
// in Step.
func (rp *Rubicon) PVDA(ctx *Context, di uint32, rnd randx.Rand) {
	rp.USs.USnegCostFromRaw(ctx, di)
	rp.PVsFromUSs(ctx, di)

	hasRew := (GlbV(ctx, di, GvHasRew) > 0)
	pvPos := GlbV(ctx, di, GvPVpos)
	pvNeg := GlbV(ctx, di, GvPVneg)
	vsPatchPos := GlbV(ctx, di, GvVSPatchPos)
	vsPatchPosSum := GlbV(ctx, di, GvVSPatchPosSum)

	if hasRew {
		rp.ResetGiveUp(ctx, di)
		rew := rp.LHb.DAforUS(ctx, di, pvPos, pvNeg, vsPatchPos, vsPatchPosSum) // only when actual pos rew
		SetGlbV(ctx, di, GvRew, rew)
		return
	}

	if GlbV(ctx, di, GvVSMatrixHasGated) > 0 {
		giveUp := rp.GiveUpOnGoal(ctx, di, rnd)
		if giveUp {
			SetGlbV(ctx, di, GvHasRew, 1)                                           // key for triggering reset
			rew := rp.LHb.DAforUS(ctx, di, pvPos, pvNeg, vsPatchPos, vsPatchPosSum) // only when actual rew
			SetGlbV(ctx, di, GvRew, rew)
			return
		}
	}

	// no US regular case
	rp.LHb.DAforNoUS(ctx, di)
	SetGlbV(ctx, di, GvRew, 0)
}
