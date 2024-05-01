// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// bench runs a benchmark model with 5 layers (3 hidden, Input, Output) all of the same
// size, for benchmarking different size networks.  These are not particularly realistic
// models for actual applications (e.g., large models tend to have much more topographic
// patterns of connectivity and larger layers with fewer connections), but they are
// easy to run..
package bench

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"

	"cogentcore.org/core/tensor"
	"cogentcore.org/core/tensor/table"
	"github.com/emer/axon/v2/axon"
	"github.com/emer/emergent/v2/erand"
	"github.com/emer/emergent/v2/etime"
	"github.com/emer/emergent/v2/params"
	"github.com/emer/emergent/v2/patgen"
	"github.com/emer/emergent/v2/prjn"
	"github.com/emer/emergent/v2/timer"
)

// note: with 2 hidden layers, this simple test case converges to perfect performance:
// ./bench -epochs 100 -pats 10 -units 100 -threads=1
// so these params below are reasonable for actually learning (eventually)

var ParamSets = params.Sets{
	"Base": {Desc: "these are the best params", Sheets: params.Sheets{
		"Network": &params.Sheet{
			{Sel: "Prjn", Desc: "",
				Params: params.Params{
					"Prjn.Learn.LRate.Base": "0.1", // 0.1 is default, 0.05 for TrSpk = .5
					"Prjn.SWts.Adapt.LRate": "0.1", // .1 >= .2,
					"Prjn.SWts.Init.SPct":   "0.5", // .5 >= 1 here -- 0.5 more reliable, 1.0 faster..
				}},
			{Sel: "Layer", Desc: "",
				Params: params.Params{
					"Layer.Inhib.ActAvg.Nominal": "0.08",
					"Layer.Inhib.Layer.Gi":       "1.05",
					"Layer.Acts.Gbar.L":          "0.2",
				}},
			{Sel: "#Input", Desc: "",
				Params: params.Params{
					"Layer.Inhib.Layer.Gi": "0.9", // 0.9 > 1.0
					"Layer.Acts.Clamp.Ge":  "1.5",
				}},
			{Sel: "#Output", Desc: "",
				Params: params.Params{
					"Layer.Inhib.Layer.Gi": "0.70",
					"Layer.Acts.Clamp.Ge":  "0.8",
				}},
			{Sel: ".BackPrjn", Desc: "top-down back-projections MUST have lower relative weight scale, otherwise network hallucinates",
				Params: params.Params{
					"Prjn.PrjnScale.Rel": "0.2",
				}},
		},
	}},
}

func ConfigNet(net *axon.Network, ctx *axon.Context, threads, units int, verbose bool) {
	net.InitName(net, "BenchNet")

	squn := int(math.Sqrt(float64(units)))
	shp := []int{squn, squn}

	inLay := net.AddLayer("Input", shp, axon.InputLayer)
	hid1Lay := net.AddLayer("Hidden1", shp, axon.SuperLayer)
	hid2Lay := net.AddLayer("Hidden2", shp, axon.SuperLayer)
	hid3Lay := net.AddLayer("Hidden3", shp, axon.SuperLayer)
	outLay := net.AddLayer("Output", shp, axon.TargetLayer)

	full := prjn.NewFull()

	net.ConnectLayers(inLay, hid1Lay, full, axon.ForwardPrjn)
	net.BidirConnectLayers(hid1Lay, hid2Lay, full)
	net.BidirConnectLayers(hid2Lay, hid3Lay, full)
	net.BidirConnectLayers(hid3Lay, outLay, full)

	net.RecFunTimes = verbose
	net.GPU.RecFunTimes = verbose

	// builds with default threads
	if err := net.Build(ctx); err != nil {
		panic(err)
	}
	net.Defaults()
	if _, err := net.ApplyParams(ParamSets["Base"].Sheets["Network"], false); err != nil {
		panic(err)
	}

	if threads == 0 {
		if verbose {
			fmt.Print("Threading: using default values\n")
		}
	} else {
		net.SetNThreads(threads)
	}

	net.InitWts(ctx)
}

func ConfigPats(dt *table.Table, pats, units int) {
	squn := int(math.Sqrt(float64(units)))
	shp := []int{squn, squn}
	// fmt.Printf("shape: %v\n", shp)

	dt.SetFromSchema(table.Schema{
		{"Name", tensor.STRING, nil, nil},
		{"Input", reflect.Float32, shp, []string{"Y", "X"}},
		{"Output", reflect.Float32, shp, []string{"Y", "X"}},
	}, pats)

	// note: actually can learn if activity is .15 instead of .25
	nOn := units / 8

	patgen.PermutedBinaryRows(dt.Columns[1], nOn, 1, 0)
	patgen.PermutedBinaryRows(dt.Columns[2], nOn, 1, 0)
}

func ConfigEpcLog(dt *table.Table) {
	dt.SetFromSchema(table.Schema{
		{"Epoch", reflect.Int, nil, nil},
		{"CorSim", reflect.Float32, nil, nil},
		{"AvgCorSim", reflect.Float32, nil, nil},
		{"SSE", reflect.Float32, nil, nil},
		{"CountErr", reflect.Float32, nil, nil},
		{"PctErr", reflect.Float32, nil, nil},
		{"PctCor", reflect.Float32, nil, nil},
		{"Hid1ActAvg", reflect.Float32, nil, nil},
		{"Hid2ActAvg", reflect.Float32, nil, nil},
		{"OutActAvg", reflect.Float32, nil, nil},
	}, 0)
}

func TrainNet(net *axon.Network, ctx *axon.Context, pats, epcLog *table.Table, epcs int, verbose, gpu bool) {
	net.InitWts(ctx)
	np := pats.NumRows()
	porder := rand.Perm(np) // randomly permuted order of ints

	if gpu {
		net.ConfigGPUnoGUI(ctx)
	}

	epcLog.SetNumRows(epcs)

	inLay := net.LayerByName("Input").(*axon.Layer)
	hid1Lay := net.LayerByName("Hidden1").(*axon.Layer)
	hid2Lay := net.LayerByName("Hidden2").(*axon.Layer)
	outLay := net.LayerByName("Output").(*axon.Layer)

	inPats := pats.ColumnByName("Input").(*tensor.Float32)
	outPats := pats.ColumnByName("Output").(*tensor.Float32)

	cycPerQtr := 50

	tmr := timer.Time{}
	tmr.Start()
	for epc := 0; epc < epcs; epc++ {
		erand.PermuteInts(porder)
		outCorSim := float32(0)
		cntErr := 0
		sse := 0.0
		for pi := 0; pi < np; pi++ {
			ppi := porder[pi]
			inp := inPats.SubSpace([]int{ppi})
			outp := outPats.SubSpace([]int{ppi})

			inLay.ApplyExt(ctx, 0, inp)
			outLay.ApplyExt(ctx, 0, outp)
			net.ApplyExts(ctx)

			net.NewState(ctx)
			ctx.NewState(etime.Train)
			for qtr := 0; qtr < 4; qtr++ {
				for cyc := 0; cyc < cycPerQtr; cyc++ {
					net.Cycle(ctx)
					ctx.CycleInc()
				}
				if qtr == 2 {
					net.MinusPhase(ctx)
					ctx.NewPhase(true)
					net.PlusPhaseStart(ctx)
				}
			}
			net.PlusPhase(ctx)
			net.DWt(ctx)
			net.WtFromDWt(ctx)
			outCorSim += outLay.Values[0].CorSim.Cor
			pSSE := outLay.PctUnitErr(ctx)[0]
			sse += pSSE
			if pSSE != 0 {
				cntErr++
			}
		}
		outCorSim /= float32(np)
		sse /= float64(np)
		pctErr := float64(cntErr) / float64(np)
		pctCor := 1 - pctErr

		t := tmr.Stop()
		tmr.Start()
		if verbose {
			fmt.Printf("epc: %v  \tCorSim: %v \tAvgCorSim: %v \tTime:%v\n", epc, outCorSim, outLay.Values[0].CorSim.Avg, t)
		}

		epcLog.SetFloat("Epoch", epc, float64(epc))
		epcLog.SetFloat("CorSim", epc, float64(outCorSim))
		epcLog.SetFloat("AvgCorSim", epc, float64(outLay.Values[0].CorSim.Avg))
		epcLog.SetFloat("SSE", epc, sse)
		epcLog.SetFloat("CountErr", epc, float64(cntErr))
		epcLog.SetFloat("PctErr", epc, pctErr)
		epcLog.SetFloat("PctCor", epc, pctCor)
		epcLog.SetFloat("Hid1ActAvg", epc, float64(hid1Lay.Values[0].ActAvg.ActMAvg))
		epcLog.SetFloat("Hid2ActAvg", epc, float64(hid2Lay.Values[0].ActAvg.ActMAvg))
		epcLog.SetFloat("OutActAvg", epc, float64(outLay.Values[0].ActAvg.ActMAvg))
	}
	tmr.Stop()
	if verbose {
		fmt.Printf("Took %6.4g secs for %v epochs, avg per epc: %6.4g\n", tmr.TotalSecs(), epcs, tmr.TotalSecs()/float64(epcs))
		net.TimerReport()
	} else {
		fmt.Printf("Total Secs: %6.3g\n", tmr.TotalSecs())
	}

	net.GPU.Destroy()
}
