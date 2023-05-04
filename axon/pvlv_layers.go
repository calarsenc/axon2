// Copyright (c) 2022, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package axon

import (
	"log"
	"strings"

	"github.com/goki/gosl/slbool"
	"github.com/goki/mat32"
)

//gosl: start pvlv_layers

// LDTParams compute reward salience as ACh global neuromodulatory signal
// as a function of the MAX activation of its inputs.
type LDTParams struct {
	RewThr      float32     `desc:"threshold per input source, on absolute value (magnitude), to count as a significant reward event, which then drives maximal ACh -- set to 0 to disable this nonlinear behavior"`
	Rew         slbool.Bool `desc:"use the global Context.NeuroMod.HasRew flag -- if there is some kind of external reward being given, then ACh goes to 1, else 0 for this component"`
	RewPred     slbool.Bool `desc:"use the global Context.NeuroMod.RewPred value"`
	MaintInhib  float32     `desc:"extent to which active maintenance (via Context.NeuroMod.NotMaint PTNotMaintLayer activity) inhibits ACh signals -- when goal engaged, distractability is lower."`
	NotMaintMax float32     `desc:"maximum NeuroMod.NotMaint activity for computing Maint as 1-NotMaint -- when NotMaint is >= NotMaintMax, then Maint = 0."`
	SrcLay1Idx  int32       `inactive:"+" desc:"idx of Layer to get max activity from -- set during Build from BuildConfig SrcLay1Name if present -- -1 if not used"`
	SrcLay2Idx  int32       `inactive:"+" desc:"idx of Layer to get max activity from -- set during Build from BuildConfig SrcLay2Name if present -- -1 if not used"`
	SrcLay3Idx  int32       `inactive:"+" desc:"idx of Layer to get max activity from -- set during Build from BuildConfig SrcLay3Name if present -- -1 if not used"`
	SrcLay4Idx  int32       `inactive:"+" desc:"idx of Layer to get max activity from -- set during Build from BuildConfig SrcLay4Name if present -- -1 if not used"`

	pad, pad1, pad2 int32
}

func (lp *LDTParams) Defaults() {
	lp.RewThr = 0.05
	lp.Rew.SetBool(true)
	lp.MaintInhib = 0.5
	lp.NotMaintMax = 0.4
}

func (lp *LDTParams) Update() {
}

// Thr applies threshold to given value
func (lp *LDTParams) Thr(val float32) float32 {
	if lp.RewThr <= 0 {
		return val
	}
	if mat32.Abs(val) < lp.RewThr {
		return 0
	}
	return 1
}

// MaintFmNotMaint returns a 0-1 value reflecting strength of active maintenance
// based on the activity of the PTNotMaintLayer as recorded in NeuroMod.NotMaint.
func (lp *LDTParams) MaintFmNotMaint(notMaint float32) float32 {
	if notMaint > lp.NotMaintMax {
		return 0
	}
	return (lp.NotMaintMax - notMaint) / lp.NotMaintMax // == 1 when notMaint = 0
}

// PVLVParams has parameters for readout of values as inputs to PVLV equations.
type PVLVParams struct {
	Thr  float32 `desc:"threshold on value prior to multiplying by Gain"`
	Gain float32 `desc:"multiplier applied after Thr threshold"`

	pad, pad1 float32
}

func (pp *PVLVParams) Defaults() {
	pp.Thr = 0.2
	pp.Gain = 4
}

func (pp *PVLVParams) Update() {

}

func (pp *PVLVParams) Val(val float32) float32 {
	vs := val - pp.Thr
	if vs < 0 {
		return 0
	}
	return pp.Gain * vs
}

// VSPatchParams parameters for VSPatch learning
type VSPatchParams struct {
	NoDALRate float32 `def:"0.1" desc:"learning rate when no positive dopamine is present (i.e., when not learning to predict a positive valence PV / US outcome.  if too high, extinguishes too quickly.  if too low, doesn't discriminate US vs. non-US trials as well."`
	NoDAThr   float32 `def:"0.01" desc:"threshold on DA level to engage the NoDALRate -- use a small positive number just in case"`

	pad, pad1 float32
}

func (pp *VSPatchParams) Defaults() {
	pp.NoDALRate = 0.1
	pp.NoDAThr = 0.01
}

func (pp *VSPatchParams) Update() {

}

// DALRate returns the learning rate modulation factor modlr based on dopamine level
func (pp *VSPatchParams) DALRate(da, modlr float32) float32 {
	if da <= pp.NoDAThr {
		if modlr < -pp.NoDALRate { // big dip: use it
			return modlr
		}
		return -pp.NoDALRate
	}
	return modlr
}

//gosl: end pvlv_layers

func (ly *Layer) BLADefaults() {
	isAcq := strings.Contains(ly.Nm, "Acq") || strings.Contains(ly.Nm, "Novel")

	lp := ly.Params
	lp.Act.Decay.Act = 0.2
	lp.Act.Decay.Glong = 0.6
	lp.Act.Dend.SSGi = 0
	lp.Inhib.Layer.On.SetBool(true)
	if isAcq {
		lp.Inhib.Layer.Gi = 2.2 // acq has more input
	} else {
		lp.Inhib.Layer.Gi = 1.8
		lp.Act.Gbar.L = 0.25
	}
	lp.Inhib.Pool.On.SetBool(true)
	lp.Inhib.Pool.Gi = 0.9
	lp.Inhib.ActAvg.Nominal = 0.025
	lp.Learn.RLRate.SigmoidMin = 1.0
	lp.Learn.TrgAvgAct.On.SetBool(false)
	lp.Learn.RLRate.Diff.SetBool(true)
	lp.Learn.RLRate.DiffThr = 0.01
	lp.CT.DecayTau = 0
	lp.CT.GeGain = 0.1 // 0.1 has effect, can go a bit lower if need to

	if isAcq {
		lp.Learn.NeuroMod.DALRateMod = 0.5
		lp.Learn.NeuroMod.BurstGain = 0.2
		lp.Learn.NeuroMod.DipGain = 0
	} else {
		lp.Learn.NeuroMod.BurstGain = 1
		lp.Learn.NeuroMod.DipGain = 1
	}
	lp.Learn.NeuroMod.AChLRateMod = 1
	lp.Learn.NeuroMod.AChDisInhib = 0 // needs to be always active

	for _, pj := range ly.RcvPrjns {
		slay := pj.Send
		if slay.LayerType() == BLALayer && !strings.Contains(slay.Nm, "Novel") { // inhibition from Ext
			pj.Params.SetFixedWts()
			pj.Params.PrjnScale.Abs = 2
		}
	}
}

// PVLVPostBuild is used for BLA, VSPatch, and PVLayer types to set NeuroMod params
func (ly *Layer) PVLVPostBuild() {
	dm, err := ly.BuildConfigByName("DAMod")
	if err == nil {
		err = ly.Params.Learn.NeuroMod.DAMod.FromString(dm)
		if err != nil {
			log.Println(err)
		}
	}
	vl, err := ly.BuildConfigByName("Valence")
	if err == nil {
		err = ly.Params.Learn.NeuroMod.Valence.FromString(vl)
		if err != nil {
			log.Println(err)
		}
	}
}

func (ly *Layer) CeMDefaults() {
	lp := ly.Params
	lp.Act.Decay.Act = 1
	lp.Act.Decay.Glong = 1
	lp.Act.Dend.SSGi = 0
	lp.Inhib.Layer.On.SetBool(true)
	lp.Inhib.Layer.Gi = 0.5
	lp.Inhib.Pool.On.SetBool(true)
	lp.Inhib.Pool.Gi = 0.3
	lp.Inhib.ActAvg.Nominal = 0.15
	lp.Learn.TrgAvgAct.On.SetBool(false)
	lp.Learn.RLRate.SigmoidMin = 1.0 // doesn't matter -- doesn't learn..

	for _, pj := range ly.RcvPrjns {
		pj.Params.SetFixedWts()
		pj.Params.PrjnScale.Abs = 1
	}
}

func (ly *Layer) LDTDefaults() {
	lp := ly.Params
	lp.Inhib.ActAvg.Nominal = 0.1
	lp.Inhib.Layer.On.SetBool(true)
	lp.Inhib.Layer.Gi = 1 // todo: explore
	lp.Inhib.Pool.On.SetBool(false)
	lp.Act.Decay.Act = 1
	lp.Act.Decay.Glong = 1
	lp.Act.Decay.LearnCa = 1 // uses CaSpkD as a readout!
	lp.Learn.TrgAvgAct.On.SetBool(false)
	lp.PVLV.Thr = 0.2
	lp.PVLV.Gain = 2

	for _, pj := range ly.RcvPrjns {
		pj.Params.SetFixedWts()
		pj.Params.PrjnScale.Abs = 1
	}
}

func (ly *LayerParams) VSPatchDefaults() {
	ly.Act.Decay.Act = 1
	ly.Act.Decay.Glong = 1
	ly.Act.Decay.LearnCa = 1 // uses CaSpkD as a readout!
	ly.Inhib.Pool.On.SetBool(true)
	ly.Inhib.Layer.On.SetBool(false)
	ly.Inhib.Layer.Gi = 0.5
	ly.Inhib.Layer.FB = 0
	ly.Inhib.Pool.FB = 0
	ly.Inhib.Pool.Gi = 0.5
	ly.Inhib.ActAvg.Nominal = 0.2
	ly.Learn.RLRate.Diff.SetBool(false)
	ly.Learn.RLRate.SigmoidMin = 1
	ly.Learn.TrgAvgAct.On.SetBool(false)

	// ms.Learn.NeuroMod.DAMod needs to be set via BuildConfig
	ly.Learn.NeuroMod.DALRateSign.SetBool(true)
	ly.Learn.NeuroMod.AChLRateMod = 0.8 // ACh now active for extinction, so this is ok
	ly.Learn.NeuroMod.AChDisInhib = 0   // essential: has to fire when expected but not present!
	ly.Learn.NeuroMod.BurstGain = 1
	ly.Learn.NeuroMod.DipGain = 0.1 // extinction -- reduce to slow
	ly.PVLV.Thr = 0.4
	ly.PVLV.Gain = 20
}

func (ly *LayerParams) DrivesDefaults() {
	ly.Inhib.ActAvg.Nominal = 0.01
	ly.Inhib.Layer.On.SetBool(false)
	ly.Inhib.Pool.On.SetBool(true)
	ly.Inhib.Pool.Gi = 0.5
	ly.Act.PopCode.On.SetBool(true)
	ly.Act.PopCode.MinAct = 0.2 // low activity for low drive -- also has special 0 case = nothing
	ly.Act.PopCode.MinSigma = 0.08
	ly.Act.PopCode.MaxSigma = 0.12
	ly.Act.Decay.Act = 1
	ly.Act.Decay.Glong = 1
	ly.Learn.TrgAvgAct.On.SetBool(false)
}

func (ly *LayerParams) EffortDefaults() {
	ly.Inhib.ActAvg.Nominal = 0.2
	ly.Inhib.Layer.On.SetBool(true)
	ly.Inhib.Layer.Gi = 0.5
	ly.Inhib.Pool.On.SetBool(false)
	ly.Act.PopCode.On.SetBool(true) // use only popcode
	ly.Act.Decay.Act = 1
	ly.Act.Decay.Glong = 1
	ly.Learn.TrgAvgAct.On.SetBool(false)
}

func (ly *LayerParams) UrgencyDefaults() {
	ly.Inhib.ActAvg.Nominal = 0.2
	ly.Inhib.Layer.On.SetBool(true)
	ly.Inhib.Layer.Gi = 0.5
	ly.Inhib.Pool.On.SetBool(false)
	ly.Act.PopCode.On.SetBool(true) // use only popcode
	ly.Act.PopCode.MinAct = 0
	ly.Act.Decay.Act = 1
	ly.Act.Decay.Glong = 1
	ly.Learn.TrgAvgAct.On.SetBool(false)
}

func (ly *LayerParams) USDefaults() {
	ly.Inhib.ActAvg.Nominal = 0.2
	ly.Inhib.Layer.On.SetBool(true)
	ly.Inhib.Layer.Gi = 0.5
	ly.Inhib.Pool.On.SetBool(false)
	ly.Act.Decay.Act = 1
	ly.Act.Decay.Glong = 1
	ly.Learn.TrgAvgAct.On.SetBool(false)
}

func (ly *LayerParams) PVDefaults() {
	ly.Inhib.ActAvg.Nominal = 0.2
	ly.Inhib.Layer.On.SetBool(true)
	ly.Inhib.Layer.Gi = 0.5
	ly.Inhib.Pool.On.SetBool(false)
	ly.Act.PopCode.On.SetBool(true)
	// note: may want to modulate rate code as well:
	// ly.Act.PopCode.MinAct = 0.2
	// ly.Act.PopCode.MinSigma = 0.08
	// ly.Act.PopCode.MaxSigma = 0.12
	ly.Act.Decay.Act = 1
	ly.Act.Decay.Glong = 1
	ly.Learn.TrgAvgAct.On.SetBool(false)
}
