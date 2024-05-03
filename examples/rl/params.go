// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/emer/emergent/v2/netparams"
	"github.com/emer/emergent/v2/params"
)

// ParamSets is the default set of parameters -- Base is always applied, and others can be optionally
// selected to apply on top of that
var ParamSets = netparams.Sets{
	"Base": {
		{Sel: "#Input", Desc: "input fixed act",
			Params: params.Params{
				"Layer.Acts.Decay.Act":       "1",
				"Layer.Acts.Decay.Glong":     "1",
				"Layer.Inhib.ActAvg.Nominal": "0.05",
			}},
		{Sel: "#Rew", Desc: "",
			Params: params.Params{
				"Layer.Inhib.Layer.Gi":       "0.2",
				"Layer.Inhib.ActAvg.Nominal": "1",
			}},
	},
	"RW": {
		{Sel: ".RWPredLayer", Desc: "",
			Params: params.Params{
				"Layer.Inhib.Layer.Gi":       "0.2",
				"Layer.Inhib.ActAvg.Nominal": "1",
				"Layer.Acts.Dt.GeTau":        "40",
			}},
		{Sel: ".RWPath", Desc: "RW pred",
			Params: params.Params{
				"Path.SWts.Init.Mean":      "0",
				"Path.SWts.Init.Var":       "0",
				"Path.SWts.Init.Sym":       "false",
				"Path.Learn.LRate.Base":    "0.1",
				"Path.RLPred.OppSignLRate": "1.0",
				"Path.RLPred.DaTol":        "0.0",
			}},
	},
	"TD": {
		{Sel: ".TDPredLayer", Desc: "",
			Params: params.Params{
				"Layer.Inhib.Layer.Gi":       "0.2",
				"Layer.Inhib.ActAvg.Nominal": "1",
				"Layer.Acts.Dt.GeTau":        "40",
			}},
		{Sel: ".TDIntegLayer", Desc: "",
			Params: params.Params{
				"Layer.Inhib.Layer.Gi":       "0.2",
				"Layer.Inhib.ActAvg.Nominal": "1",
				"Layer.TDInteg.Discount":     "0.9",
				"Layer.TDInteg.PredGain":     "1.0",
			}},
		{Sel: "#InputToRewPred", Desc: "input to rewpred",
			Params: params.Params{
				"Path.SWts.Init.Mean":      "0",
				"Path.SWts.Init.Var":       "0",
				"Path.SWts.Init.Sym":       "false",
				"Path.Learn.LRate.Base":    "0.1",
				"Path.RLPred.OppSignLRate": "1.0",
			}},
	},
}
