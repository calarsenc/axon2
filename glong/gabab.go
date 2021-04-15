// Copyright (c) 2020, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package glong

import (
	"github.com/chewxy/math32"
)

// GABABParams control the GABAB dynamics in PFC Maint neurons,
// based on Brunel & Wang (2001) parameters.
type GABABParams struct {
	RiseTau  float32 `def:"45" desc:"rise time for bi-exponential time dynamics of GABA-B"`
	DecayTau float32 `def:"50" desc:"decay time for bi-exponential time dynamics of GABA-B"`
	Gbar     float32 `def:"0.2" desc:"overall strength multiplier of GABA-B current"`
	Gbase    float32 `def:"0.2" desc:"baseline level of GABA-B channels open independent of inhibitory input (is added to spiking-produced conductance)"`
	Smult    float32 `def:"15" desc:"multiplier for converting Gi from FFFB to GABA spikes"`
	MaxTime  float32 `inactive:"+" desc:"time offset when peak conductance occurs, in msec, computed from RiseTau and DecayTau"`
	TauFact  float32 `view:"-" desc:"time constant factor used in integration: (Decay / Rise) ^ (Rise / (Decay - Rise))"`
}

func (gp *GABABParams) Defaults() {
	gp.RiseTau = 45
	gp.DecayTau = 50
	gp.Gbar = 0.2
	gp.Gbase = 0.2
	gp.Smult = 15
	gp.Update()
}

func (gp *GABABParams) Update() {
	gp.TauFact = math32.Pow(gp.DecayTau/gp.RiseTau, gp.RiseTau/(gp.DecayTau-gp.RiseTau))
	gp.MaxTime = ((gp.RiseTau * gp.DecayTau) / (gp.DecayTau - gp.RiseTau)) * math32.Log(gp.DecayTau/gp.RiseTau)
}

// GFmV returns the GABA-B conductance as a function of normalized membrane potential
func (gp *GABABParams) GFmV(v float32) float32 {
	vbio := v*100 - 100
	return 1 / (1 + math32.Exp(0.1*((vbio+90)+10)))
}

// GFmS returns the GABA-B conductance as a function of GABA spiking rate,
// based on normalized spiking factor (i.e., Gi from FFFB etc)
func (gp *GABABParams) GFmS(s float32) float32 {
	ss := s * gp.Smult // convert to spikes
	return 1 / (1 + math32.Exp(-(ss-7.1)/1.4))
}

// BiExp computes bi-exponential update, returns dG and dX deltas to add to g and x
func (gp *GABABParams) BiExp(g, x float32) (dG, dX float32) {
	dG = (gp.TauFact*x - g) / gp.RiseTau
	dX = -x / gp.DecayTau
	return
}

// GABAB returns the updated GABA-B / GIRK activation and underlying x value
// based on current values and gi inhibitory conductance (proxy for GABA spikes)
func (gp *GABABParams) GABAB(gabaB, gabaBx, gi float32) (g, x float32) {
	dG, dX := gp.BiExp(gabaB, gabaBx)
	x = gabaBx + gp.GFmS(gi) + dX // gets new input
	g = gabaB + dG
	return
}

// GgabaB returns the overall net GABAB / GIRK conductance including
// Gbar, Gbase, and voltage-gating
func (gp *GABABParams) GgabaB(gabaB, vm float32) float32 {
	return gp.Gbar * gp.GFmV(vm) * (gabaB + gp.Gbase)
}
