// Copyright (c) 2024, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kinase

import (
	"fmt"
	"math/rand"

	"cogentcore.org/core/tensor/table"
	"github.com/chewxy/math32"
)

type LinearState int

const (
	StartCaSyn LinearState = iota
	StartCaM
	StartCaP
	StartCaD

	FinalCaSyn
	FinalCaM
	FinalCaP
	FinalCaD

	TotalSpikes

	NLinearState
)

// Linear performs a linear regression to approximate the synaptic Ca
// integration between send and recv neurons.
type Linear struct {
	// Kinase Neuron params
	Neuron NeurCaParams

	// Kinase Synapse params
	Synapse SynCaParams

	// gain on S*R product for BinnedSums
	BinProd float32 `default:"10"`

	// total number of cycles (1 MSec) to run
	NCycles int `min:"10" default:"200"`

	// number of plus cycles
	PlusCycles int `default:"50"`

	// CyclesPerBin specifies the bin size for accumulating spikes
	CyclesPerBin int `default:"25"`

	// NumBins = NCycles / CyclesPerBin
	NumBins int `edit:"-"`

	// MaxHz is the maximum firing rate to sample in minus, plus phases
	MaxHz int `default:"120"`

	// StepHz is the step size for sampling Hz
	StepHz int `default:"10"`

	// NTrials is number of trials per Hz case
	NTrials int `default:"100"`

	// Total Trials is number of trials for all data
	TotalTrials int `edit:"-"`

	// Sending neuron
	Send Neuron

	// Receiving neuron
	Recv Neuron

	// Standard synapse values
	StdSyn Synapse

	// Linear synapse values
	LinearSyn Synapse

	// ErrDWt is the target error dwt: PlusHz - MinusHz
	ErrDWt float32

	// binned integration of send, recv spikes
	BinnedSums []float32

	// Data to fit the regression
	Data table.Table
}

func (ls *Linear) Defaults() {
	ls.Neuron.Defaults()
	ls.Synapse.Defaults()
	ls.BinProd = 10
	ls.NCycles = 200
	ls.PlusCycles = 50
	ls.CyclesPerBin = 25
	ls.MaxHz = 120
	ls.StepHz = 10
	ls.NTrials = 10
	ls.Update()
}

func (ls *Linear) Update() {
	ls.Neuron.Update()
	ls.Synapse.Update()
	ls.NumBins = ls.NCycles / ls.CyclesPerBin
	nhz := ls.MaxHz / ls.StepHz
	ls.TotalTrials = nhz * nhz * nhz * nhz * ls.NTrials
	ls.BinnedSums = make([]float32, ls.NumBins)
	ls.Send.BinnedSpikes = make([]float32, ls.NumBins)
	ls.Recv.BinnedSpikes = make([]float32, ls.NumBins)
}

func (ls *Linear) Init() {
	ls.Send.Init()
	ls.Recv.Init()
	ls.StdSyn.Init()
	ls.LinearSyn.Init()
	ls.InitTable()
}

func (ls *Linear) InitTable() {
	if ls.Data.NumColumns() > 0 {
		return
	}
	nneur := int(NLinearState*2) + ls.NumBins
	ls.Data.AddIntColumn("Trial")
	ls.Data.AddFloat64TensorColumn("Hz", []int{4}, "Send*Recv*Minus*Plus")
	ls.Data.AddFloat64TensorColumn("State", []int{nneur}, "States")
	ls.Data.AddFloat64TensorColumn("StdCa", []int{4}, "CaSyn,M,P,D")
	ls.Data.AddFloat64TensorColumn("PredCa", []int{4}, "CaSyn,M,P,D")
	ls.Data.AddFloat64TensorColumn("CaSSE", []int{4}, "CaSyn,M,P,D")
	ls.Data.AddFloat64Column("SSE") // total SSE
	ls.Data.SetNumRows(ls.TotalTrials)
}

func (ls *Linear) StartTrial() {
	ls.Send.StartTrial()
	ls.Recv.StartTrial()
}

// Neuron has Neuron state
type Neuron struct {
	// Neuron spiking (0,1)
	Spike float32

	// Neuron probability of spiking
	SpikeP float32

	// CaSyn is spike-driven calcium trace for synapse-level Ca-driven learning: exponential integration of SpikeG * Spike at SynTau time constant (typically 30).  Synapses integrate send.CaSyn * recv.CaSyn across M, P, D time integrals for the synaptic trace driving credit assignment in learning. Time constant reflects binding time of Glu to NMDA and Ca buffering postsynaptically, and determines time window where pre * post spiking must overlap to drive learning.
	CaSyn float32

	// neuron-level spike-driven Ca integration
	CaSpkM, CaSpkP, CaSpkD float32

	// regression variables
	StartCaSyn float32

	TotalSpikes float32

	// binned count of spikes, for regression learning
	BinnedSpikes []float32
}

func (kn *Neuron) Init() {
	kn.Spike = 0
	kn.SpikeP = 1
	kn.CaSyn = 0
	kn.CaSpkM = 0
	kn.CaSpkP = 0
	kn.CaSpkD = 0
	kn.StartTrial()
}

func (kn *Neuron) StartTrial() {
	kn.StartCaSyn = kn.CaSyn
	kn.TotalSpikes = 0
	for i := range kn.BinnedSpikes {
		kn.BinnedSpikes[i] = 0
	}
}

// Cycle does one cycle of neuron updating, with given exponential spike interval
// based on target spiking firing rate.
func (ls *Linear) Cycle(nr *Neuron, expInt float32, cyc int) {
	nr.Spike = 0
	bin := cyc / ls.CyclesPerBin
	if expInt > 0 {
		nr.SpikeP *= rand.Float32()
		if nr.SpikeP <= expInt {
			nr.Spike = 1
			nr.SpikeP = 1
			nr.TotalSpikes += 1
			nr.BinnedSpikes[bin] += 1
		}
	}
	ls.Neuron.CaFromSpike(nr.Spike, &nr.CaSyn, &nr.CaSpkM, &nr.CaSpkP, &nr.CaSpkD)
}

// Synapse has Synapse state
type Synapse struct {
	CaSyn float32

	// CaM is first stage running average (mean) Ca calcium level (like CaM = calmodulin), feeds into CaP
	CaM float32

	// CaP is shorter timescale integrated CaM value, representing the plus, LTP direction of weight change and capturing the function of CaMKII in the Kinase learning rule
	CaP float32

	// CaD is longer timescale integrated CaP value, representing the minus, LTD direction of weight change and capturing the function of DAPK1 in the Kinase learning rule
	CaD float32

	// DWt is the CaP - CaD
	DWt float32
}

func (ks *Synapse) Init() {
	ks.CaSyn = 0
	ks.CaM = 0
	ks.CaP = 0
	ks.CaD = 0
	ks.DWt = 0
}

// Run generates data
func (ls *Linear) Run() {
	nhz := ls.MaxHz / ls.StepHz
	hz := make([]float32, nhz)
	i := 0
	for h := float32(ls.StepHz); h <= float32(ls.MaxHz); h += float32(ls.StepHz) {
		hz[i] = h
		i++
	}
	row := 0
	for smi := 0; smi < nhz; smi++ {
		sendMinusHz := hz[smi]
		for spi := 0; spi < nhz; spi++ {
			sendPlusHz := hz[spi]
			for rmi := 0; rmi < nhz; rmi++ {
				recvMinusHz := hz[rmi]
				for rpi := 0; rpi < nhz; rpi++ {
					recvPlusHz := hz[rpi]

					for ti := 0; ti < ls.NTrials; ti++ {
						ls.Trial(sendMinusHz, sendPlusHz, recvMinusHz, recvPlusHz, ti, row)
						row++
					}
				}
			}
		}
	}
	fmt.Println("row:", row)
}

func (ls *Linear) SetNeurState(nr *Neuron, off, row int) {
	ls.Data.SetTensorFloat1D("State", row, off, float64(nr.CaSyn))
	ls.Data.SetTensorFloat1D("State", row, off+1, float64(nr.CaSpkM))
	ls.Data.SetTensorFloat1D("State", row, off+2, float64(nr.CaSpkP))
	ls.Data.SetTensorFloat1D("State", row, off+3, float64(nr.CaSpkD))
}

func (ls *Linear) SetSynState(sy *Synapse, row int) {
	ls.Data.SetTensorFloat1D("StdCa", row, 0, float64(sy.CaSyn))
	ls.Data.SetTensorFloat1D("StdCa", row, 1, float64(sy.CaM))
	ls.Data.SetTensorFloat1D("StdCa", row, 2, float64(sy.CaP))
	ls.Data.SetTensorFloat1D("StdCa", row, 3, float64(sy.CaD))
}

func (ls *Linear) SetBins(sn, rn *Neuron, off, row int) {
	for i, s := range sn.BinnedSpikes {
		r := rn.BinnedSpikes[i]
		bs := r + s + ls.BinProd*r*s
		ls.BinnedSums[i] = bs
		ls.Data.SetTensorFloat1D("State", row, off+i, float64(bs))
	}
}

// Trial runs one trial
func (ls *Linear) Trial(sendMinusHz, sendPlusHz, recvMinusHz, recvPlusHz float32, ti, row int) {
	// ls.ErrDWt = (plusHz - minusHz) / 100

	ls.Data.SetFloat("Trial", row, float64(ti))
	ls.Data.SetTensorFloat1D("Hz", row, 0, float64(sendMinusHz))
	ls.Data.SetTensorFloat1D("Hz", row, 1, float64(sendPlusHz))
	ls.Data.SetTensorFloat1D("Hz", row, 2, float64(recvMinusHz))
	ls.Data.SetTensorFloat1D("Hz", row, 3, float64(recvPlusHz))

	// capture starting
	ls.SetNeurState(&ls.Send, 0, row)
	ls.SetNeurState(&ls.Recv, int(NLinearState), row)

	minusCycles := ls.NCycles - ls.PlusCycles

	ls.StartTrial()
	cyc := 0
	for phs := 0; phs < 2; phs++ {
		var maxcyc int
		var rhz, shz float32
		switch phs {
		case 0:
			rhz = recvMinusHz
			shz = sendMinusHz
			maxcyc = minusCycles
		case 1:
			rhz = recvPlusHz
			shz = sendPlusHz
			maxcyc = ls.PlusCycles
		}
		Rint := math32.Exp(-1000.0 / float32(rhz))
		Sint := math32.Exp(-1000.0 / float32(shz))
		for t := 0; t < maxcyc; t++ {
			ls.Cycle(&ls.Send, Sint, cyc)
			ls.Cycle(&ls.Recv, Rint, cyc)

			ls.StdSyn.CaSyn = ls.Send.CaSyn * ls.Recv.CaSyn
			ls.Synapse.FromCa(ls.StdSyn.CaSyn, &ls.StdSyn.CaM, &ls.StdSyn.CaP, &ls.StdSyn.CaD)
			cyc++
		}
	}
	ls.StdSyn.DWt = ls.StdSyn.CaP - ls.StdSyn.CaD

	// capture final
	ls.SetNeurState(&ls.Send, int(FinalCaSyn), row)
	ls.SetNeurState(&ls.Recv, int(NLinearState+FinalCaSyn), row)
	ls.SetSynState(&ls.StdSyn, row)

	ls.SetBins(&ls.Send, &ls.Recv, int(NLinearState*2), row)
}
