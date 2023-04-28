// Copyright (c) 2020, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package axon

import (
	"github.com/emer/emergent/params"
	"github.com/goki/mat32"
)

//gosl: start deep_layers

// BurstParams determine how the 5IB Burst activation is computed from
// CaSpkP integrated spiking values in Super layers -- thresholded.
type BurstParams struct {
	ThrRel float32 `max:"1" def:"0.1" desc:"Relative component of threshold on superficial activation value, below which it does not drive Burst (and above which, Burst = CaSpkP).  This is the distance between the average and maximum activation values within layer (e.g., 0 = average, 1 = max).  Overall effective threshold is MAX of relative and absolute thresholds."`
	ThrAbs float32 `min:"0" max:"1" def:"0.1" desc:"Absolute component of threshold on superficial activation value, below which it does not drive Burst (and above which, Burst = CaSpkP).  Overall effective threshold is MAX of relative and absolute thresholds."`

	pad, pad1 float32
}

func (bp *BurstParams) Update() {
}

func (bp *BurstParams) Defaults() {
	bp.ThrRel = 0.1
	bp.ThrAbs = 0.1
}

// ThrFmAvgMax returns threshold from average and maximum values
func (bp *BurstParams) ThrFmAvgMax(avg, mx float32) float32 {
	thr := avg + bp.ThrRel*(mx-avg)
	thr = mat32.Max(thr, bp.ThrAbs)
	return thr
}

// CTParams control the CT corticothalamic neuron special behavior
type CTParams struct {
	GeGain   float32 `def:"0.8,1" desc:"gain factor for context excitatory input, which is constant as compared to the spiking input from other projections, so it must be downscaled accordingly.  This can make a difference and may need to be scaled up or down."`
	DecayTau float32 `def:"0,50" desc:"decay time constant for context Ge input -- if > 0, decays over time so intrinsic circuit dynamics have to take over.  For single-step copy-based cases, set to 0, while longer-time-scale dynamics should use 50"`
	DecayDt  float32 `view:"-" json:"-" xml:"-" desc:"1 / tau"`

	pad float32
}

func (cp *CTParams) Update() {
	if cp.DecayTau > 0 {
		cp.DecayDt = 1 / cp.DecayTau
	} else {
		cp.DecayDt = 0
	}
}

func (cp *CTParams) Defaults() {
	cp.GeGain = 0.8
	cp.DecayTau = 50
	cp.Update()
}

// PulvParams provides parameters for how the plus-phase (outcome)
// state of Pulvinar thalamic relay cell neurons is computed from
// the corresponding driver neuron Burst activation (or CaSpkP if not Super)
type PulvParams struct {
	DriveScale   float32 `def:"0.1" min:"0.0" desc:"multiplier on driver input strength, multiplies CaSpkP from driver layer to produce Ge excitatory input to Pulv unit."`
	FullDriveAct float32 `def:"0.6" min:"0.01" desc:"Level of Max driver layer CaSpkP at which the drivers fully drive the burst phase activation.  If there is weaker driver input, then (Max/FullDriveAct) proportion of the non-driver inputs remain and this critically prevents the network from learning to turn activation off, which is difficult and severely degrades learning."`
	DriveLayIdx  int32   `inactive:"+" desc:"index of layer that generates the driving activity into this one -- set via SetBuildConfig(DriveLayName) setting"`
	pad          float32
}

func (tp *PulvParams) Update() {
}

func (tp *PulvParams) Defaults() {
	tp.DriveScale = 0.1
	tp.FullDriveAct = 0.6
}

// DriveGe returns effective excitatory conductance
// to use for given driver input Burst activation
func (tp *PulvParams) DriveGe(act float32) float32 {
	return tp.DriveScale * act
}

// NonDrivePct returns the multiplier proportion of the non-driver based Ge to
// keep around, based on FullDriveAct and the max activity in driver layer.
func (tp *PulvParams) NonDrivePct(drvMax float32) float32 {
	return 1.0 - mat32.Min(1, drvMax/tp.FullDriveAct)
}

//gosl: end deep_layers

// note: Defaults not called on GPU

func (ly *LayerParams) CTDefaults() {
	ly.Act.Decay.Act = 0 // deep doesn't decay!
	ly.Act.Decay.Glong = 0
	ly.Act.Decay.AHP = 0
	ly.Act.Dend.SSGi = 0    // key: otherwise interferes with NMDA maint!
	ly.Inhib.Layer.Gi = 2.2 // higher inhib for more NMDA, recurrents.
	ly.Inhib.Pool.Gi = 2.2
	// these are for longer temporal integration:
	// ly.Act.NMDA.Gbar = 0.003
	// ly.Act.NMDA.Tau = 300
	// ly.Act.GABAB.Gbar = 0.008
}

// CTDefParamsFast sets fast time-integration parameters for CTLayer.
// This is what works best in the deep_move 1 trial history case,
// vs Medium and Long
func (ly *Layer) CTDefParamsFast() {
	ly.DefParams = params.Params{
		"Layer.CT.GeGain":       "1",
		"Layer.CT.DecayTau":     "0",
		"Layer.Inhib.Layer.Gi":  "2.0",
		"Layer.Inhib.Pool.Gi":   "2.0",
		"Layer.Act.GABAB.Gbar":  "0.006",
		"Layer.Act.NMDA.Gbar":   "0.004",
		"Layer.Act.NMDA.Tau":    "100",
		"Layer.Act.Decay.Act":   "0.0",
		"Layer.Act.Decay.Glong": "0.0",
		"Layer.Act.Sahp.Gbar":   "1.0",
	}
}

// CTDefParamsMedium sets medium time-integration parameters for CTLayer.
// This is what works best in the FSA test case, compared to Fast (deep_move)
// and Long (deep_music) time integration.
func (ly *Layer) CTDefParamsMedium() {
	ly.DefParams = params.Params{
		"Layer.CT.GeGain":       "0.8",
		"Layer.CT.DecayTau":     "50",
		"Layer.Inhib.Layer.Gi":  "2.2",
		"Layer.Inhib.Pool.Gi":   "2.2",
		"Layer.Act.GABAB.Gbar":  "0.009",
		"Layer.Act.NMDA.Gbar":   "0.008",
		"Layer.Act.NMDA.Tau":    "200",
		"Layer.Act.Decay.Act":   "0.0",
		"Layer.Act.Decay.Glong": "0.0",
		"Layer.Act.Sahp.Gbar":   "1.0",
	}
}

// CTDefParamsLong sets long time-integration parameters for CTLayer.
// This is what works best in the deep_music test case integrating over
// long time windows, compared to Medium and Fast.
func (ly *Layer) CTDefParamsLong() {
	ly.DefParams = params.Params{
		"Layer.CT.GeGain":       "1.0",
		"Layer.CT.DecayTau":     "50",
		"Layer.Inhib.Layer.Gi":  "2.8",
		"Layer.Inhib.Pool.Gi":   "2.8",
		"Layer.Act.GABAB.Gbar":  "0.01",
		"Layer.Act.NMDA.Gbar":   "0.01",
		"Layer.Act.NMDA.Tau":    "300",
		"Layer.Act.Decay.Act":   "0.0",
		"Layer.Act.Decay.Glong": "0.0",
		"Layer.Act.Dend.SSGi":   "0", // else kills nmda
		"Layer.Act.Sahp.Gbar":   "1.0",
	}
}

func (ly *Layer) PTMaintDefaults() {
	ly.Params.Act.Decay.Act = 0 // deep doesn't decay!
	ly.Params.Act.Decay.Glong = 0
	ly.Params.Act.Decay.AHP = 0
	ly.Params.Act.Decay.OnRew.SetBool(true)
	ly.Params.Act.Sahp.Gbar = 0.01 // not much pressure -- long maint
	ly.Params.Act.GABAB.Gbar = 0.01
	ly.Params.Act.Dend.ModGain = 1.5
	ly.Params.Inhib.ActAvg.Nominal = 0.3 // very active
	if ly.Is4D() {
		ly.Params.Inhib.ActAvg.Nominal = 0.05
	}
	ly.Params.Inhib.Layer.Gi = 1.8
	ly.Params.Inhib.Pool.Gi = 1.8
	ly.Params.Learn.TrgAvgAct.On.SetBool(false)

	for _, pj := range ly.RcvPrjns {
		slay := pj.Send
		if slay.LayerType() == BGThalLayer {
			pj.Params.Com.GType = ModulatoryG
		}
	}
}

func (ly *Layer) PTNotMaintDefaults() {
	ly.Params.Act.Decay.Act = 1
	ly.Params.Act.Decay.Glong = 1
	ly.Params.Act.Decay.OnRew.SetBool(true)
	ly.Params.Act.Init.GeBase = 1.2
	ly.Params.Learn.TrgAvgAct.On.SetBool(false)
	ly.Params.Inhib.ActAvg.Nominal = 0.2
	ly.Params.Inhib.Pool.On.SetBool(false)
	ly.Params.Inhib.Layer.On.SetBool(true)
	ly.Params.Inhib.Layer.Gi = 0.5
	ly.Params.CT.GeGain = 0.2
	ly.Params.CT.DecayTau = 0
	ly.Params.CT.Update()

	for _, pj := range ly.RcvPrjns {
		pj.Params.SetFixedWts()
	}
}

func (ly *LayerParams) PTPredDefaults() {
	ly.Act.Decay.Act = 0.12 // keep it dynamically changing
	ly.Act.Decay.Glong = 0.6
	ly.Act.Decay.AHP = 0
	ly.Act.Decay.OnRew.SetBool(true)
	ly.Act.Sahp.Gbar = 0.1    // more
	ly.Act.KNa.Slow.Max = 0.2 // todo: more?
	ly.Inhib.Layer.Gi = 0.8
	ly.Inhib.Pool.Gi = 0.8
	ly.CT.GeGain = 0.01
	ly.CT.DecayTau = 50

	// regular:
	// ly.Act.GABAB.Gbar = 0.006
	// ly.Act.NMDA.Gbar = 0.004
	// ly.Act.NMDA.Tau = 100
}

// called in Defaults for Pulvinar layer type
func (ly *LayerParams) PulvDefaults() {
	ly.Act.Decay.Act = 0
	ly.Act.Decay.Glong = 0
	ly.Act.Decay.AHP = 0
	ly.Learn.RLRate.SigmoidMin = 1.0 // 1.0 generally better but worth trying 0.05 too
}

// PulvPostBuild does post-Build config of Pulvinar based on BuildConfig options
func (ly *Layer) PulvPostBuild() {
	ly.Params.Pulv.DriveLayIdx = ly.BuildConfigFindLayer("DriveLayName", true)
}
