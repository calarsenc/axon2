// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"
	"reflect"

	"cogentcore.org/core/math32"
	"cogentcore.org/core/math32/minmax"
	"cogentcore.org/core/tensor"
	"cogentcore.org/core/tensor/stats/stats"
	"github.com/emer/emergent/v2/decoder"
	"github.com/emer/emergent/v2/elog"
	"github.com/emer/emergent/v2/etime"
)

const (
	NBins        = 20
	CyclesPerBin = 10
	NOutputs     = 3
	NInputs      = NBins + 2 // per neuron
)

// KinaseNeuron has Neuron state
type KinaseNeuron struct {
	// Neuron spiking (0,1)
	Spike float32

	// Neuron probability of spiking
	SpikeP float32

	// CaSyn is spike-driven calcium trace for synapse-level Ca-driven learning: exponential integration of SpikeG * Spike at SynTau time constant (typically 30).  Synapses integrate send.CaSyn * recv.CaSyn across M, P, D time integrals for the synaptic trace driving credit assignment in learning. Time constant reflects binding time of Glu to NMDA and Ca buffering postsynaptically, and determines time window where pre * post spiking must overlap to drive learning.
	CaSyn float32

	// regression variables
	StartCaSyn float32

	TotalSpikes float32

	// binned count of spikes, for regression learning
	BinnedSpikes [NBins]float32
}

func (kn *KinaseNeuron) Init() {
	kn.Spike = 0
	kn.SpikeP = 1
	kn.CaSyn = 0
	kn.StartTrial()
}

func (kn *KinaseNeuron) StartTrial() {
	kn.StartCaSyn = kn.CaSyn
	kn.TotalSpikes = 0
	for i := range kn.BinnedSpikes {
		kn.BinnedSpikes[i] = 0
	}
}

// Cycle does one cycle of neuron updating, with given exponential spike interval
// based on target spiking firing rate.
func (kn *KinaseNeuron) Cycle(expInt float32, params *ParamConfig, cyc int) {
	kn.Spike = 0
	bin := cyc / CyclesPerBin
	if expInt > 0 {
		kn.SpikeP *= rand.Float32()
		if kn.SpikeP <= expInt {
			kn.Spike = 1
			kn.SpikeP = 1
			kn.TotalSpikes += 1
			kn.BinnedSpikes[bin] += 1
		}
	}
	kn.CaSyn += params.SynDt * (params.SpikeG*kn.Spike - kn.CaSyn)
}

func (kn *KinaseNeuron) SetInput(inputs []float32, off int) {
	inputs[off] = kn.StartCaSyn
	inputs[off+1] = kn.TotalSpikes
	for i, s := range kn.BinnedSpikes {
		inputs[off+2+i] = s
	}
}

// KinaseSynapse has Synapse state
type KinaseSynapse struct {
	// CaM is first stage running average (mean) Ca calcium level (like CaM = calmodulin), feeds into CaP
	CaM float32

	// CaP is shorter timescale integrated CaM value, representing the plus, LTP direction of weight change and capturing the function of CaMKII in the Kinase learning rule
	CaP float32

	// CaD is longer timescale integrated CaP value, representing the minus, LTD direction of weight change and capturing the function of DAPK1 in the Kinase learning rule
	CaD float32

	// DWt is the CaP - CaD
	DWt float32
}

func (ks *KinaseSynapse) Init() {
	ks.CaM = 0
	ks.CaP = 0
	ks.CaD = 0
	ks.DWt = 0
}

// KinaseState is basic Kinase equation state
type KinaseState struct {

	// if true, training decoder
	Train bool

	// SSE for decoder
	SSE float32

	// Condition counter
	Condition int

	// Condition description
	Cond string

	// Trial counter
	Trial int

	// Cycle counter
	Cycle int

	// phase-based firing rates
	MinusHz, PlusHz float32

	// ErrDWt is the target error dwt: PlusHz - MinusHz
	ErrDWt float32

	// Sending neuron
	Send KinaseNeuron

	// Receiving neuron
	Recv KinaseNeuron

	// Standard synapse values
	StdSyn KinaseSynapse

	// Linearion synapse values
	LinearSyn KinaseSynapse
}

func (ks *KinaseState) Init() {
	ks.Send.Init()
	ks.Recv.Init()
	ks.StdSyn.Init()
	ks.LinearSyn.Init()
}

func (kn *KinaseState) StartTrial() {
	kn.Send.StartTrial()
	kn.Recv.StartTrial()
}

func (ss *Sim) ConfigKinase() {
	ss.Linear.Init(NOutputs, NInputs*2, 0, decoder.IdentityFunc)
	ss.Linear.LRate = ss.Config.Params.LRate
}

// Sweep runs a sweep through minus-plus ranges
func (ss *Sim) Sweep() {
	ss.Kinase.Train = false
	// hz := []float32{25, 50, 100}
	// nhz := len(hz)

	nhz := 100 / 5
	hz := make([]float32, nhz)
	i := 0
	for h := float32(5); h <= 100; h += 5 {
		hz[i] = h
		i++
	}

	cond := 0
	for mi := 0; mi < nhz; mi++ {
		minusHz := hz[mi]
		for pi := 0; pi < nhz; pi++ {
			plusHz := hz[pi]
			condStr := fmt.Sprintf("%03d -> %03d", minusHz, plusHz)
			ss.Kinase.Condition = cond
			ss.Kinase.Cond = condStr
			ss.RunImpl(minusHz, plusHz, ss.Config.Run.NTrials)
			cond++
		}
	}
	// note: can get this by setting x axis
	// ss.Plot("DWtPlot").Update()
	// ss.Plot("DWtVarPlot").Update()
}

// Run runs for given parameters
func (ss *Sim) Run() {
	cr := &ss.Config.Run
	ss.RunImpl(cr.MinusHz, cr.PlusHz, cr.NTrials)
}

// RunImpl runs NTrials, recording to RunLog and TrialLog
func (ss *Sim) RunImpl(minusHz, plusHz float32, ntrials int) {
	ss.Kinase.Train = false
	ss.Kinase.Init()
	for trl := 0; trl < ntrials; trl++ {
		ss.Kinase.Trial = trl
		ss.TrialImpl(minusHz, plusHz)
	}
	ss.Logs.LogRow(etime.Test, etime.Condition, ss.Kinase.Condition)
	ss.GUI.UpdatePlot(etime.Test, etime.Condition)
}

func (ss *Sim) Trial() {
	cr := &ss.Config.Run
	ss.Kinase.Init()
	ss.TrialImpl(cr.MinusHz, cr.PlusHz)
}

// TrialImpl runs one trial for given parameters
func (ss *Sim) TrialImpl(minusHz, plusHz float32) {
	cfg := &ss.Config
	ks := &ss.Kinase
	ks.MinusHz = minusHz
	ks.PlusHz = plusHz
	ks.Cycle = 0
	ks.ErrDWt = (plusHz - minusHz) / 100

	minusCycles := cfg.NCycles - cfg.PlusCycles

	ks.StartTrial()
	for phs := 0; phs < 2; phs++ {
		var maxcyc int
		var rhz float32
		switch phs {
		case 0:
			rhz = minusHz
			maxcyc = minusCycles
		case 1:
			rhz = plusHz
			maxcyc = cfg.PlusCycles
		}
		shz := rhz + cfg.Run.SendDiffHz
		if shz < 0 {
			shz = 0
		}

		var Sint, Rint float32
		if rhz > 5 {
			Rint = math32.Exp(-1000.0 / float32(rhz))
		}
		if shz > 5 {
			Sint = math32.Exp(-1000.0 / float32(shz))
		}
		for t := 0; t < maxcyc; t++ {
			ks.Send.Cycle(Sint, &cfg.Params, ks.Cycle)
			ks.Recv.Cycle(Rint, &cfg.Params, ks.Cycle)

			ca := ks.Send.CaSyn * ks.Recv.CaSyn
			ss.CaParams.FromCa(ca, &ks.StdSyn.CaM, &ks.StdSyn.CaP, &ks.StdSyn.CaD)
			if !ks.Train {
				ss.Logs.LogRow(etime.Test, etime.Cycle, ks.Cycle)
			}
			ks.Cycle++
		}
	}
	ks.StdSyn.DWt = ks.StdSyn.CaP - ks.StdSyn.CaD

	ks.Send.SetInput(ss.Linear.Inputs, 0)
	ks.Recv.SetInput(ss.Linear.Inputs, NInputs)
	ss.Linear.Forward()
	out := make([]float32, NOutputs)
	ss.Linear.Output(&out)
	ks.LinearSyn.CaM = out[0]
	ks.LinearSyn.CaP = out[1]
	ks.LinearSyn.CaD = out[2]
	ks.LinearSyn.DWt = ks.LinearSyn.CaP - ks.LinearSyn.CaD

	if ks.Train {
		targ := [NOutputs]float32{ks.StdSyn.CaM, ks.StdSyn.CaP, ks.StdSyn.CaD}
		sse, _ := ss.Linear.Train(targ[:])
		ks.SSE = sse
		ss.Logs.LogRow(etime.Train, etime.Cycle, 0)
		ss.GUI.UpdatePlot(etime.Train, etime.Cycle)
		ss.Logs.LogRow(etime.Train, etime.Trial, ks.Trial)
		ss.GUI.UpdatePlot(etime.Train, etime.Trial)
	} else {
		ss.GUI.UpdatePlot(etime.Test, etime.Cycle)
		ss.Logs.LogRow(etime.Test, etime.Trial, ks.Trial)
		ss.GUI.UpdatePlot(etime.Test, etime.Trial)
	}
}

// Train trains the linear decoder
func (ss *Sim) Train() {
	ss.Kinase.Train = true
	ss.Kinase.Init()
	for epc := 0; epc < ss.Config.Run.NEpochs; epc++ {
		ss.Kinase.Condition = epc
		for trl := 0; trl < ss.Config.Run.NTrials; trl++ {
			ss.Kinase.Trial = trl
			minusHz := 100 * rand.Float32()
			plusHz := 100 * rand.Float32()
			ss.TrialImpl(minusHz, plusHz)
		}
		ss.Logs.LogRow(etime.Train, etime.Condition, ss.Kinase.Condition)
		ss.GUI.UpdatePlot(etime.Train, etime.Condition)
	}
	tensor.SaveCSV(&ss.Linear.Weights, "trained.wts", '\t')
}

func (ss *Sim) ConfigKinaseLogItems() {
	lg := &ss.Logs
	ks := &ss.Kinase
	val := reflect.ValueOf(ks).Elem()
	parName := ""
	times := []etime.Times{etime.Condition, etime.Trial, etime.Cycle}
	tn := len(times)
	WalkFields(val,
		func(parent reflect.Value, field reflect.StructField, value reflect.Value) bool {
			if field.Name == "BinnedSpikes" {
				return false
			}
			return true
		},
		func(parent reflect.Value, field reflect.StructField, value reflect.Value) {
			fkind := field.Type.Kind()
			fname := field.Name
			if val.Interface() == parent.Interface() { // top-level
				if fkind == reflect.Struct {
					parName = fname
					return
				}
			} else {
				fname = parName + "." + fname
			}
			itm := lg.AddItem(&elog.Item{
				Name:   fname,
				Type:   fkind,
				FixMax: false,
				Range:  minmax.F32{Max: 1},
				Write: elog.WriteMap{
					etime.Scope(etime.AllModes, etime.Cycle): func(ctx *elog.Context) {
						fany := value.Interface()
						switch fkind {
						case reflect.Float32:
							ctx.SetFloat32(fany.(float32))
						case reflect.Int:
							ctx.SetFloat32(float32(fany.(int)))
						case reflect.String:
							ctx.SetString(fany.(string))
						}
					},
				}})
			for ti := 0; ti < tn-1; ti++ {
				if fkind == reflect.Float32 {
					itm.Write[etime.Scope(etime.AllModes, times[ti])] = func(ctx *elog.Context) {
						ctx.SetAgg(ctx.Mode, times[ti+1], stats.Mean)
					}
				} else {
					itm.Write[etime.Scope(etime.Train, times[ti])] = func(ctx *elog.Context) {
						fany := value.Interface()
						switch fkind {
						case reflect.Int:
							ctx.SetFloat32(float32(fany.(int)))
						case reflect.String:
							ctx.SetString(fany.(string))
						}
					}
					itm.Write[etime.Scope(etime.Test, times[ti])] = func(ctx *elog.Context) {
						fany := value.Interface()
						switch fkind {
						case reflect.Int:
							ctx.SetFloat32(float32(fany.(int)))
						case reflect.String:
							ctx.SetString(fany.(string))
						}
					}
				}
			}
		})
}

func WalkFields(parent reflect.Value, should func(parent reflect.Value, field reflect.StructField, value reflect.Value) bool, walk func(parent reflect.Value, field reflect.StructField, value reflect.Value)) {
	typ := parent.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		value := parent.Field(i)
		if !should(parent, field, value) {
			continue
		}
		if field.Type.Kind() == reflect.Struct {
			walk(parent, field, value)
			WalkFields(value, should, walk)
		} else {
			walk(parent, field, value)
		}
	}
}
