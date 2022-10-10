// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/emer/emergent/params"

// ParamSets is the default set of parameters -- Base is always applied, and others can be optionally
// selected to apply on top of that
var ParamSets = params.Sets{
	{Name: "Base", Desc: "these are the best params", Sheets: params.Sheets{
		"Network": &params.Sheet{
			{Sel: "Layer", Desc: "generic layer params",
				Params: params.Params{
					"Layer.Inhib.ActAvg.Init": "0.1", // 0.05 needed to get hidden2 high to .1, 0.1 keeps it too low!
					"Layer.Inhib.Layer.Gi":    "1.0", // 1.0 > 1.1  trace
					"Layer.Act.Gbar.L":        "0.2", // std
					"Layer.Act.Decay.Act":     "0.0", // 0 == 0.2
					"Layer.Act.Decay.Glong":   "0.0",
					"Layer.Act.NMDA.MgC":      "1.4", // 1.4, 5 > 1.2, 0 ?
					"Layer.Act.NMDA.Voff":     "5",
					"Layer.Act.Mahp.Gbar":     "0.04", // 0.04 == 0.05+ > 0.02 -- reduces hidden activity
					"Layer.Act.Sahp.Gbar":     "0.1",  // 0.1 == 0.02 no real diff
					"Layer.Act.Sahp.Off":      "0.8",  //
					"Layer.Act.Sahp.Slope":    "0.02", //
					"Layer.Act.Sahp.CaTau":    "5",    // 5 > 10
				}},
			{Sel: ".Hidden", Desc: "fix avg act",
				Params: params.Params{
					"Layer.Inhib.ActAvg.Init": "0.1",
				}},
			{Sel: ".DepthIn", Desc: "",
				Params: params.Params{
					"Layer.Inhib.ActAvg.Init": "0.2", // was .13 -- Ge very high b/c of topo prjn
				}},
			{Sel: ".HeadDirIn", Desc: "",
				Params: params.Params{
					"Layer.Inhib.ActAvg.Init": "0.13", // 0.13 > 0.2 -- 0.13 is accurate but Ge is high..
				}},
			{Sel: ".CT", Desc: "CT NMDA gbar factor is key",
				Params: params.Params{
					"Layer.Inhib.ActAvg.Init": "0.12", // CT in general more active
					"Layer.Inhib.Layer.Gi":    "1.4",  // 1.4 > 1.6 with 1to1 super -> CT; 1.6 needed for full
					"Layer.Inhib.Pool.Gi":     "1.4",
					"Layer.CT.GeGain":         "1", //
					"Layer.CT.DecayTau":       "0", // 50 > 30 -- 30 ok but takes a bit to get going
					"Layer.Act.Decay.Act":     "0.0",
					"Layer.Act.Decay.Glong":   "0.0",
					"Layer.Act.GABAB.Gbar":    "0.15",
					"Layer.Act.NMDA.Gbar":     "0.15",  // .3 is min -- .25 fails, even with .35 in hidden2!
					"Layer.Act.NMDA.Tau":      "100",   // 300 >> 200, even with 300 in hidden2
					"Layer.Act.Noise.On":      "false", // todo?
					"Layer.Act.Noise.Ge":      "0.005",
					"Layer.Act.Noise.Gi":      "0.005",
				}},
			{Sel: "#DepthHid", Desc: "",
				Params: params.Params{
					"Layer.Inhib.Layer.Gi":    "1.2",  // 1.2 tiny bit > 1.4
					"Layer.Inhib.ActAvg.Init": "0.05", // 0.05 > 1. even though wrong -- @1.4 = ~0.05 actual, @1.2 = .1 actual
				}},
			{Sel: "#DepthHidCT", Desc: "",
				Params: params.Params{
					"Layer.Inhib.Layer.Gi":    "2.0",
					"Layer.Inhib.ActAvg.Init": "0.05", // 0.05 > .1 even though wrong
				}},
			{Sel: "#DepthHid2CT", Desc: "CT NMDA gbar factor is key",
				Params: params.Params{
					"Layer.CT.GeGain":         "0.8",  // 0.8, 50 small benefit
					"Layer.CT.DecayTau":       "50",   // 50 > 0
					"Layer.Inhib.ActAvg.Init": "0.12", // 2 even more active -- maybe try higher inhib
					"Layer.Inhib.Layer.Gi":    "1.4",  // todo
					"Layer.Act.GABAB.Gbar":    "0.3",
					"Layer.Act.NMDA.Gbar":     "0.3", // higher layer has more nmda..
					"Layer.Act.NMDA.Tau":      "300", // 300 > 200
					"Layer.Act.Sahp.CaTau":    "10",  // todo
				}},
			{Sel: "TRCLayer", Desc: "TRC = Pulvinar",
				Params: params.Params{
					"Layer.Inhib.Layer.Gi":          "0.9", // 0.9 > 1.0
					"Layer.TRC.DriveScale":          "0.1", // 0.1 > 0.05 > 0.15
					"Layer.TRC.FullDriveAct":        "0.6", // 0.6 def
					"Layer.Act.Decay.Act":           "0.0",
					"Layer.Act.Decay.Glong":         "0.0", // clear long
					"Layer.Act.Decay.AHP":           "0.0", // clear long
					"Layer.Act.GABAB.Gbar":          "0.2",
					"Layer.Act.NMDA.Gbar":           "0.1", // .1 was important
					"Layer.Learn.RLrate.SigmoidMin": "1",   // 1 > .05
				}},
			{Sel: "#DepthHidP", Desc: "distributed hidden-layer pulvinar",
				Params: params.Params{
					"Layer.Inhib.Layer.Gi": "0.9",  // 0.9 > 0.8 > 1
					"Layer.TRC.DriveScale": "0.05", // 0.05 > .1
					"Layer.Act.NMDA.Gbar":  "0.1",
				}},
			{Sel: "#Action", Desc: "",
				Params: params.Params{
					"Layer.Inhib.ActAvg.Init": "0.15",
				}},

			// Projections below
			{Sel: "Prjn", Desc: "std",
				Params: params.Params{
					"Prjn.Learn.Lrate.Base":   "0.005", // 0.005 > 0.002 > 0.01
					"Prjn.SWt.Adapt.Lrate":    "0.01",  // 0.01 == 0.0001 but 0.001 not as good..
					"Prjn.SWt.Adapt.DreamVar": "0.0",   // 0.01 is just tolerable
					"Prjn.SWt.Init.SPct":      "1.0",   // 1 works fine here -- .5 also ok
					"Prjn.Com.PFail":          "0.0",
					"Prjn.Learn.Trace.Tau":    "2", // 4 == 2 > 1
				}},
			{Sel: ".Back", Desc: "top-down back-projections MUST have lower relative weight scale, otherwise network hallucinates",
				Params: params.Params{
					"Prjn.PrjnScale.Rel": "0.2",
				}},
			{Sel: "#HeadDirHidCTToDepthHidCT", Desc: "ct top-down",
				Params: params.Params{
					"Prjn.PrjnScale.Rel": "0.2", // not much diff here
				}},
			{Sel: ".CTCtxt", Desc: "all CT context prjns",
				Params: params.Params{
					"Prjn.Learn.Lrate.Base": "0.002", // has almost no effect in 1to1
					"Prjn.Trace":            "false", // not as good with Trace here..
					"Prjn.Com.PFail":        "0.0",   // .2, .3 too high -- very slow learning
				}},
			{Sel: ".CTFmSuper", Desc: "1to1 > full",
				Params: params.Params{
					"Prjn.Learn.Learn":   "true", // learning > fixed 1to1
					"Prjn.SWt.Init.Mean": "0.5",  // if fixed, 0.8 > 0.5, var = 0
					"Prjn.SWt.Init.Var":  "0.25",
				}},
			{Sel: ".FmPulv", Desc: "defaults to .Back but generally weaker is better",
				Params: params.Params{
					"Prjn.PrjnScale.Rel": "0.1", // 0.1 == 0.15 > 0.05
				}},
			{Sel: ".CTSelfCtxt", Desc: "",
				Params: params.Params{
					"Prjn.PrjnScale.Rel": "0.5",  // 0.5 > 0.2 > 0.8
					"Prjn.Com.PFail":     "0.0",  // never useful for random gen
					"Prjn.SWt.Init.Sym":  "true", // true > false
				}},
			{Sel: ".CTSelfMaint", Desc: "",
				Params: params.Params{
					"Prjn.PrjnScale.Rel": "0.1", // 0.1  >= 0.05 > 0.2
					"Prjn.Com.PFail":     "0.0",
					"Prjn.SWt.Init.Sym":  "true", // no effect?  not sure why
				}},
			{Sel: "#ActionToDepthHidCT", Desc: "",
				Params: params.Params{
					"Prjn.PrjnScale.Rel": "0.2", // 0.5 is not better
				}},
			{Sel: "#DepthHid2CTToDepthP", Desc: "",
				Params: params.Params{
					"Prjn.PrjnScale.Rel": "0.1", // 0.1 == 0.15 > 0.05
				}},
		},
	}},
}
