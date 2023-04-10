// Copyright (c) 2022, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package axon

import (
	"github.com/emer/emergent/prjn"
	"github.com/emer/emergent/relpos"
)

// AddLDTLayer adds a LDTLayer
func (nt *Network) AddLDTLayer(prefix string) *Layer {
	ldt := nt.AddLayer2D(prefix+"LDT", 1, 1, LDTLayer)
	return ldt
}

// AddBLALayers adds two BLA layers, acquisition / extinction / D1 / D2,
// for positive or negative valence
func (nt *Network) AddBLALayers(prefix string, pos bool, nUs, nNeurY, nNeurX int, rel relpos.Relations, space float32) (acq, ext *Layer) {
	if pos {
		d1 := nt.AddLayer4D(prefix+"BLAPosAcqD1", 1, nUs, nNeurY, nNeurX, BLALayer)
		d1.SetBuildConfig("DAMod", "D1Mod")
		d1.SetBuildConfig("Valence", "Positive")
		d2 := nt.AddLayer4D(prefix+"BLAPosExtD2", 1, nUs, nNeurY, nNeurX, BLALayer)
		d2.SetBuildConfig("DAMod", "D2Mod")
		d2.SetBuildConfig("Valence", "Positive")
		acq = d1
		ext = d2
	} else {
		d1 := nt.AddLayer4D(prefix+"BLANegExtD1", 1, nUs, nNeurY, nNeurX, BLALayer)
		d1.SetBuildConfig("DAMod", "D1Mod")
		d1.SetBuildConfig("Valence", "Negative")
		d2 := nt.AddLayer4D(prefix+"BLANegAcqD2", 1, nUs, nNeurY, nNeurX, BLALayer)
		d2.SetBuildConfig("DAMod", "D2Mod")
		d2.SetBuildConfig("Valence", "Negative")
		acq = d2
		ext = d1
	}

	nt.ConnectLayers(ext, acq, prjn.NewPoolOneToOne(), InhibPrjn).SetClass("BLAExtToAcq")

	if rel == relpos.Behind {
		ext.PlaceBehind(acq, space)
	} else {
		ext.PlaceRightOf(acq, space)
	}
	acq.SetClass("BLA")
	ext.SetClass("BLA")
	return
}

// AddAmygdala adds a full amygdala complex including BLA,
// CeM, and LDT.  Inclusion of negative valence is optional with neg
// arg -- neg* layers are nil if not included.
func (nt *Network) AddAmygdala(prefix string, neg bool, nUs, nNeurY, nNeurX int, space float32) (blaPosAcq, blaPosExt, blaNegAcq, blaNegExt, cemPos, cemNeg *Layer) {
	blaPosAcq, blaPosExt = nt.AddBLALayers(prefix, true, nUs, nNeurY, nNeurX, relpos.Behind, space)
	if neg {
		blaNegAcq, blaNegExt = nt.AddBLALayers(prefix, false, nUs, nNeurY, nNeurX, relpos.Behind, space)
	}
	cemPos = nt.AddLayer4D(prefix+"CeMPos", 1, nUs, 1, nNeurX, CeMLayer)
	cemPos.SetBuildConfig("DAMod", "D1Mod") // not relevant but avoids warning
	cemPos.SetBuildConfig("Valence", "Positive")
	if neg {
		cemNeg = nt.AddLayer4D(prefix+"CeMNeg", 1, nUs, 1, nNeurX, CeMLayer)
		cemNeg.SetBuildConfig("DAMod", "D2Mod") // not relevant but avoids warning
		cemNeg.SetBuildConfig("Valence", "Negative")
	}

	p1to1 := prjn.NewPoolOneToOne()

	nt.ConnectLayers(blaPosAcq, cemPos, p1to1, ForwardPrjn).SetClass("BLAToCeM_Excite")
	nt.ConnectLayers(blaPosExt, cemPos, p1to1, InhibPrjn).SetClass("BLAToCeM_Inhib")

	if neg {
		nt.ConnectLayers(blaNegAcq, cemNeg, p1to1, ForwardPrjn).SetClass("BLAToCeM_Excite")
		nt.ConnectLayers(blaNegExt, cemNeg, p1to1, InhibPrjn).SetClass("BLAToCeM_Inhib")
	}

	cemPos.PlaceBehind(blaPosExt, space)
	if neg {
		blaNegAcq.PlaceBehind(blaPosExt, space)
		cemPos.PlaceBehind(blaNegExt, space)
		cemNeg.PlaceBehind(cemPos, space)
	}

	return
}

// ConnectToBLAAcq adds a BLAAcqPrjn from given sending layer to a BLA layer
func (nt *Network) ConnectToBLAAcq(send, recv *Layer, pat prjn.Pattern) *Prjn {
	return nt.ConnectLayers(send, recv, pat, BLAAcqPrjn)
}

// ConnectToBLAExt adds a BLAExtPrjn from given sending layer to a BLA layer
func (nt *Network) ConnectToBLAExt(send, recv *Layer, pat prjn.Pattern) *Prjn {
	return nt.ConnectLayers(send, recv, pat, BLAExtPrjn)
}

// AddUSLayers adds USpos and USneg layers for positive or negative valence
// unconditioned stimuli (USs).
// These track the ContextPVLV.USpos or USneg, for visualization purposes.
// Actual US inputs are set in DrivePVLV.
func (nt *Network) AddUSLayers(nUSpos, nUSneg, nYneur int, rel relpos.Relations, space float32) (usPos, usNeg *Layer) {
	usPos = nt.AddLayer4D("USpos", 1, nUSpos, nYneur, 1, USLayer)
	usPos.SetBuildConfig("DAMod", "D1Mod") // not relevant but avoids warning
	usPos.SetBuildConfig("Valence", "Positive")
	usNeg = nt.AddLayer4D("USneg", 1, nUSneg, nYneur, 1, USLayer)
	usNeg.SetBuildConfig("DAMod", "D2Mod") // not relevant but avoids warning
	usNeg.SetBuildConfig("Valence", "Negative")
	if rel == relpos.Behind {
		usNeg.PlaceBehind(usPos, space)
	} else {
		usNeg.PlaceRightOf(usPos, space)
	}
	return
}

// AddUSPulvLayers adds USpos and USneg layers for positive or negative valence
// unconditioned stimuli (USs).
// These track the ContextPVLV.USpos or USneg, for visualization purposes.
// Actual US inputs are set in DrivePVLV.
// Adds Pulvinar predictive layers for each.
func (nt *Network) AddUSPulvLayers(nUSpos, nUSneg, nYneur int, rel relpos.Relations, space float32) (usPos, usNeg, usPosP, usNegP *Layer) {
	usPos, usNeg = nt.AddUSLayers(nUSpos, nUSneg, nYneur, rel, space)
	usPosP = nt.AddPulvForLayer(usPos, space)
	usNegP = nt.AddPulvForLayer(usNeg, space)
	if rel == relpos.Behind {
		usNeg.PlaceBehind(usPosP, space)
	}
	usPosP.SetClass("USLayer")
	usNegP.SetClass("USLayer")
	return
}

// AddPVLayers adds PVpos and PVneg layers for positive or negative valence
// primary value representations, representing the total drive and effort weighted
// USpos outcome, or total USneg outcome.
// Uses a PopCode representation based on LayerParams.Act.PopCode, distributed over
// given numbers of neurons in the X and Y dimensions.
func (nt *Network) AddPVLayers(nNeurY, nNeurX int, rel relpos.Relations, space float32) (pvPos, pvNeg *Layer) {
	pvPos = nt.AddLayer2D("PVpos", nNeurY, nNeurX, PVLayer)
	pvPos.SetBuildConfig("DAMod", "D1Mod") // not relevant but avoids warning
	pvPos.SetBuildConfig("Valence", "Positive")
	pvNeg = nt.AddLayer2D("PVneg", nNeurY, nNeurX, PVLayer)
	pvNeg.SetBuildConfig("DAMod", "D2Mod") // not relevant but avoids warning
	pvNeg.SetBuildConfig("Valence", "Negative")
	if rel == relpos.Behind {
		pvNeg.PlaceBehind(pvPos, space)
	} else {
		pvNeg.PlaceRightOf(pvPos, space)
	}
	return
}

// AddPVLayers adds PVpos and PVneg layers for positive or negative valence
// primary value representations, representing the total drive and effort weighted
// USpos outcomes, or total USneg outcomes.
// Uses a PopCode representation based on LayerParams.Act.PopCode, distributed over
// given numbers of neurons in the X and Y dimensions.
// Adds Pulvinar predictive layers for each.
func (nt *Network) AddPVPulvLayers(nNeurY, nNeurX int, rel relpos.Relations, space float32) (pvPos, pvNeg, pvPosP, pvNegP *Layer) {
	pvPos, pvNeg = nt.AddPVLayers(nNeurX, nNeurY, rel, space)
	pvPosP = nt.AddPulvForLayer(pvPos, space)
	pvNegP = nt.AddPulvForLayer(pvNeg, space)
	if rel == relpos.Behind {
		pvNeg.PlaceBehind(pvPosP, space)
	}
	pvPosP.SetClass("PVLayer")
	pvNegP.SetClass("PVLayer")
	return
}

// AddVSPatchLayer adds VSPatch (Pos, D1)
func (nt *Network) AddVSPatchLayer(prefix string, nUs, nNeurY, nNeurX int) *Layer {
	d1 := nt.AddLayer4D(prefix+"VSPatch", 1, nUs, nNeurY, nNeurX, VSPatchLayer)
	d1.SetBuildConfig("DAMod", "D1Mod")
	d1.SetBuildConfig("Valence", "Positive")
	return d1

	// d2 = nt.AddLayer4D(prefix+"VSPatchPosD2", 1, nUs, nNeurY, nNeurX, VSPatchLayer)
	// d2.SetBuildConfig("DAMod", "D2Mod")
	// d2.SetBuildConfig("Valence", "Positive")
	// if rel == relpos.Behind {
	// 	d2.PlaceBehind(d1, space)
	// } else {
	// 	d2.PlaceRightOf(d1, space)
	// }
	//  else {
	// d2 = nt.AddLayer4D(prefix+"VSPatchNegD2", 1, nUs, nNeurY, nNeurX, VSPatchLayer)
	// d2.SetBuildConfig("DAMod", "D2Mod")
	// d2.SetBuildConfig("Valence", "Negative")
	// d1 = nt.AddLayer4D(prefix+"VSPatchNegD1", 1, nUs, nNeurY, nNeurX, VSPatchLayer)
	// d1.SetBuildConfig("DAMod", "D1Mod")
	// d1.SetBuildConfig("Valence", "Negative")
	// if rel == relpos.Behind {
	// 	d1.PlaceBehind(d2, space)
	// } else {
	// 	d1.PlaceRightOf(d2, space)
	// }
}

// ConnectToVSPatch adds a VSPatchPrjn from given sending layer to a VSPatch layer
func (nt *Network) ConnectToVSPatch(send, recv *Layer, pat prjn.Pattern) *Prjn {
	return nt.ConnectLayers(send, recv, pat, VSPatchPrjn)
}

// AddVTALHbLDTLayers adds VTA dopamine, LHb DA dipping, and LDT ACh layers
// which are driven by corresponding values in ContextPVLV
func (nt *Network) AddVTALHbLDTLayers(rel relpos.Relations, space float32) (vta, lhb, ldt *Layer) {
	vta = nt.AddLayer2D("VTA", 1, 1, VTALayer)
	lhb = nt.AddLayer2D("LHb", 1, 2, LHbLayer)
	ldt = nt.AddLDTLayer("")
	if rel == relpos.Behind {
		lhb.PlaceBehind(vta, space)
		ldt.PlaceBehind(lhb, space)
	} else {
		lhb.PlaceRightOf(vta, space)
		ldt.PlaceRightOf(lhb, space)
	}
	return
}

// AddSCLayer2D adds superior colliculcus 2D layer
// which computes stimulus onset via trial-delayed inhibition
// (Inhib.FFPrv) -- connect with fixed random input from sensory
// input layers.  Sets base name and class name to SC.
// Must set Inhib.FFPrv > 0 and Act.Decay.* = 0
func (nt *Network) AddSCLayer2D(prefix string, nNeurY, nNeurX int) *Layer {
	sc := nt.AddLayer2D(prefix+"SC", nNeurY, nNeurX, SuperLayer)
	sc.SetClass("SC")
	return sc
}

// AddSCLayer4D adds superior colliculcus 4D layer
// which computes stimulus onset via trial-delayed inhibition
// (Inhib.FFPrv) -- connect with fixed random input from sensory
// input layers.  Sets base name and class name to SC.
// Must set Inhib.FFPrv > 0 and Act.Decay.* = 0
func (nt *Network) AddSCLayer4D(prefix string, nPoolsY, nPoolsX, nNeurY, nNeurX int) *Layer {
	sc := nt.AddLayer4D(prefix+"SC", nPoolsY, nPoolsX, nNeurY, nNeurX, SuperLayer)
	sc.SetClass("SC")
	return sc
}

// ConnectToSC adds a ForwardPrjn from given sending layer to
// a SC layer, setting class as ToSC -- should set params
// as fixed random with more variance than usual.
func (nt *Network) ConnectToSC(send, recv *Layer, pat prjn.Pattern) *Prjn {
	pj := nt.ConnectLayers(send, recv, pat, ForwardPrjn)
	pj.SetClass("ToSC")
	return pj
}

// AddDrivesLayer adds DrivePVLV layer representing current drive activity,
// from ContextPVLV.Drive.Drives.
// Uses a PopCode representation based on LayerParams.Act.PopCode, distributed over
// given numbers of neurons in the X and Y dimensions, per drive pool.
func (nt *Network) AddDrivesLayer(ctx *Context, nNeurY, nNeurX int) *Layer {
	drv := nt.AddLayer4D("Drives", 1, int(ctx.PVLV.Drive.NActive), nNeurY, nNeurX, DrivesLayer)
	return drv
}

// AddDrivesPulvLayer adds DrivePVLV layer representing current drive activity,
// from ContextPVLV.Drive.Drives.
// Uses a PopCode representation based on LayerParams.Act.PopCode, distributed over
// given numbers of neurons in the X and Y dimensions, per drive pool.
// Adds Pulvinar predictive layers for Drives.
func (nt *Network) AddDrivesPulvLayer(ctx *Context, nNeurY, nNeurX int, space float32) (drv, drvP *Layer) {
	drv = nt.AddDrivesLayer(ctx, nNeurY, nNeurX)
	drvP = nt.AddPulvForLayer(drv, space)
	drvP.SetClass("DrivesLayer")
	return
}

// AddEffortLayer adds DrivePVLV layer representing current effort factor,
// from ContextPVLV.Effort.EffortDisc()
// Uses a PopCode representation based on LayerParams.Act.PopCode, distributed over
// given numbers of neurons in the X and Y dimensions, per drive pool.
func (nt *Network) AddEffortLayer(nNeurY, nNeurX int) *Layer {
	eff := nt.AddLayer2D("Effort", nNeurY, nNeurX, EffortLayer)
	return eff
}

// AddEffortPulvLayer adds DrivePVLV layer representing current effort factor,
// from ContextPVLV.Effort.EffortDisc()
// Uses a PopCode representation based on LayerParams.Act.PopCode, distributed over
// given numbers of neurons in the X and Y dimensions, per drive pool.
// Adds Pulvinar predictive layers for Effort.
func (nt *Network) AddEffortPulvLayer(nNeurY, nNeurX int, space float32) (eff, effP *Layer) {
	eff = nt.AddEffortLayer(nNeurY, nNeurX)
	effP = nt.AddPulvForLayer(eff, space)
	effP.SetClass("EffortLayer")
	return
}

// AddDrivePVLVPulvLayers adds PVLV layers for PV-related information visualizing
// the internal states of the ContextPVLV state, with Pulvinar prediction
// layers for training PFC layers.
// * drives = popcode representation of drive strength (no activity for 0)
// number of active drives comes from Context; popY, popX neurons per pool.
// * effort = popcode representation of effort discount factor, popY, popX neurons.
// * us = nYneur per US, represented as present or absent
// * pv = popcode representation of final primary value on positive and negative
// valences -- this is what the dopamine value ends up conding (pos - neg).
// Layers are organized in depth per type: USs in one column, PVs in the next,
// with Drives in the back.
func (nt *Network) AddDrivePVLVPulvLayers(ctx *Context, nUSneg, nYneur, popY, popX int, space float32) (drives, drivesP, effort, effortP, usPos, usNeg, usPosP, usNegP, pvPos, pvNeg, pvPosP, pvNegP *Layer) {
	rel := relpos.Behind
	nUSpos := int(ctx.PVLV.Drive.NActive)
	usPos, usNeg, usPosP, usNegP = nt.AddUSPulvLayers(nUSpos, nUSneg, nYneur, rel, space)
	pvPos, pvNeg, pvPosP, pvNegP = nt.AddPVPulvLayers(popY, popX, rel, space)
	drives, drivesP = nt.AddDrivesPulvLayer(ctx, popY, popX, space)
	effort, effortP = nt.AddEffortPulvLayer(popY, popX, space)

	pvPos.PlaceRightOf(usPos, space)
	drives.PlaceBehind(usNegP, space)
	effort.PlaceRightOf(drives, space)
	return
}
