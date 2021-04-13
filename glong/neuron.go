// Copyright (c) 2020, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package glong

import (
	"fmt"
	"unsafe"

	"github.com/emer/axon/axon"
)

var (
	// NeuronVars are extra neuron variables for glong
	NeuronVars = []string{"AlphaMax", "VmEff", "Gnmda", "NMDA", "NMDASyn", "GgabaB", "GABAB", "GABABx"}

	// NeuronVarsAll is the glong collection of all neuron-level vars
	NeuronVarsAll []string

	NeuronVarsMap map[string]int

	// NeuronVarProps are integrated neuron var props including axon
	NeuronVarProps = map[string]string{
		"NMDA":   `auto-scale:"+"`,
		"GABAB":  `auto-scale:"+"`,
		"GABABx": `auto-scale:"+"`,
	}
)

func init() {
	ln := len(axon.NeuronVars)
	NeuronVarsAll = make([]string, len(NeuronVars)+ln)
	copy(NeuronVarsAll, axon.NeuronVars)
	copy(NeuronVarsAll[ln:], NeuronVars)

	NeuronVarsMap = make(map[string]int, len(NeuronVars))
	for i, v := range NeuronVars {
		NeuronVarsMap[v] = i
	}
	for v, p := range axon.NeuronVarProps {
		NeuronVarProps[v] = p
	}
}

// Neuron holds the extra neuron (unit) level variables for glong computation.
type Neuron struct {
	AlphaMax float32 `desc:"Maximum activation over Alpha cycle period"`
	VmEff    float32 `desc:"Effective membrane potential, including simulated backpropagating action potential contribution from activity level."`
	Gnmda    float32 `desc:"net NMDA conductance, after Vm gating and Gbar -- added directly to Ge as it has the same reversal potential."`
	NMDA     float32 `desc:"NMDA channel activation -- underlying time-integrated value with decay"`
	NMDASyn  float32 `desc:"synaptic NMDA activation directly from projection(s)"`
	GgabaB   float32 `desc:"net GABA-B conductance, after Vm gating and Gbar + Gbase -- set to Gk for GIRK, with .1 reversal potential."`
	GABAB    float32 `desc:"GABA-B / GIRK activation -- time-integrated value with rise and decay time constants"`
	GABABx   float32 `desc:"GABA-B / GIRK internal drive variable -- gets the raw activation and decays"`
}

func (nrn *Neuron) VarNames() []string {
	return NeuronVars
}

// NeuronVarIdxByName returns the index of the variable in the Neuron, or error
func NeuronVarIdxByName(varNm string) (int, error) {
	i, ok := NeuronVarsMap[varNm]
	if !ok {
		return 0, fmt.Errorf("Neuron VarByName: variable name: %v not valid", varNm)
	}
	return i, nil
}

// VarByIndex returns variable using index (0 = first variable in NeuronVars list)
func (nrn *Neuron) VarByIndex(idx int) float32 {
	fv := (*float32)(unsafe.Pointer(uintptr(unsafe.Pointer(nrn)) + uintptr(4*idx)))
	return *fv
}

// VarByName returns variable by name, or error
func (nrn *Neuron) VarByName(varNm string) (float32, error) {
	i, err := NeuronVarIdxByName(varNm)
	if err != nil {
		return 0, err
	}
	return nrn.VarByIndex(i), nil
}
