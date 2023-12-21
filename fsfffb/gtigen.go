// Code generated by "goki generate -add-types"; DO NOT EDIT.

package fsfffb

import (
	"goki.dev/gti"
	"goki.dev/ordmap"
)

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/fsfffb.GiParams",
	ShortName: "fsfffb.GiParams",
	IDName:    "gi-params",
	Doc:       "GiParams parameterizes feedforward (FF) and feedback (FB) inhibition (FFFB)\nbased on incoming spikes (FF) and outgoing spikes (FB)\nacross Fast (PV+) and Slow (SST+) timescales.\nFF -> PV -> FS fast spikes, FB -> SST -> SS slow spikes (slow to get going)",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gosl", Directive: "start", Args: []string{"fsfffb"}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"On", &gti.Field{Name: "On", Type: "goki.dev/gosl/v2/slbool.Bool", LocalType: "slbool.Bool", Doc: "enable this level of inhibition", Directives: gti.Directives{}, Tag: ""}},
		{"Gi", &gti.Field{Name: "Gi", Type: "float32", LocalType: "float32", Doc: "overall inhibition gain -- this is main parameter to adjust to change overall activation levels -- it scales both the the FS and SS factors uniformly", Directives: gti.Directives{}, Tag: "viewif:\"On\" min:\"0\" def:\"1,1.1,0.75,0.9\""}},
		{"FB", &gti.Field{Name: "FB", Type: "float32", LocalType: "float32", Doc: "amount of FB spikes included in FF for driving FS -- for small networks, 0.5 or 1 works best; larger networks and more demanding inhibition requires higher levels.", Directives: gti.Directives{}, Tag: "viewif:\"On\" min:\"0\" def:\"0.5,1,4\""}},
		{"FSTau", &gti.Field{Name: "FSTau", Type: "float32", LocalType: "float32", Doc: "fast spiking (PV+) intgration time constant in cycles (msec) -- tau is roughly how long it takes for value to change significantly -- 1.4x the half-life.", Directives: gti.Directives{}, Tag: "viewif:\"On\" min:\"0\" def:\"6\""}},
		{"SS", &gti.Field{Name: "SS", Type: "float32", LocalType: "float32", Doc: "multiplier on SS slow-spiking (SST+) in contributing to the overall Gi inhibition -- FS contributes at a factor of 1", Directives: gti.Directives{}, Tag: "viewif:\"On\" min:\"0\" def:\"30\""}},
		{"SSfTau", &gti.Field{Name: "SSfTau", Type: "float32", LocalType: "float32", Doc: "slow-spiking (SST+) facilitation decay time constant in cycles (msec) -- facilication factor SSf determines impact of FB spikes as a function of spike input-- tau is roughly how long it takes for value to change significantly -- 1.4x the half-life.", Directives: gti.Directives{}, Tag: "viewif:\"On\" min:\"0\" def:\"20\""}},
		{"SSiTau", &gti.Field{Name: "SSiTau", Type: "float32", LocalType: "float32", Doc: "slow-spiking (SST+) intgration time constant in cycles (msec) cascaded on top of FSTau -- tau is roughly how long it takes for value to change significantly -- 1.4x the half-life.", Directives: gti.Directives{}, Tag: "viewif:\"On\" min:\"0\" def:\"50\""}},
		{"FS0", &gti.Field{Name: "FS0", Type: "float32", LocalType: "float32", Doc: "fast spiking zero point -- below this level, no FS inhibition is computed, and this value is subtracted from the FSi", Directives: gti.Directives{}, Tag: "viewif:\"On\" def:\"0.1\""}},
		{"FFAvgTau", &gti.Field{Name: "FFAvgTau", Type: "float32", LocalType: "float32", Doc: "time constant for updating a running average of the feedforward inhibition over a longer time scale, for computing FFPrv", Directives: gti.Directives{}, Tag: "viewif:\"On\" def:\"50\""}},
		{"FFPrv", &gti.Field{Name: "FFPrv", Type: "float32", LocalType: "float32", Doc: "proportion of previous average feed-forward inhibition (FFAvgPrv) to add, resulting in an accentuated temporal-derivative dynamic where neurons respond most strongly to increases in excitation that exceeds inhibition from last time.", Directives: gti.Directives{}, Tag: "viewif:\"On\" def:\"0\""}},
		{"ClampExtMin", &gti.Field{Name: "ClampExtMin", Type: "float32", LocalType: "float32", Doc: "minimum GeExt value required to drive external clamping dynamics (if clamp is set), where only GeExt drives inhibition.  If GeExt is below this value, then the usual FS-FFFB drivers are used.", Directives: gti.Directives{}, Tag: "viewif:\"On\" def:\"0.05\""}},
		{"FSDt", &gti.Field{Name: "FSDt", Type: "float32", LocalType: "float32", Doc: "rate = 1 / tau", Directives: gti.Directives{}, Tag: "inactive:\"+\" view:\"-\" json:\"-\" xml:\"-\""}},
		{"SSfDt", &gti.Field{Name: "SSfDt", Type: "float32", LocalType: "float32", Doc: "rate = 1 / tau", Directives: gti.Directives{}, Tag: "inactive:\"+\" view:\"-\" json:\"-\" xml:\"-\""}},
		{"SSiDt", &gti.Field{Name: "SSiDt", Type: "float32", LocalType: "float32", Doc: "rate = 1 / tau", Directives: gti.Directives{}, Tag: "inactive:\"+\" view:\"-\" json:\"-\" xml:\"-\""}},
		{"FFAvgDt", &gti.Field{Name: "FFAvgDt", Type: "float32", LocalType: "float32", Doc: "rate = 1 / tau", Directives: gti.Directives{}, Tag: "inactive:\"+\" view:\"-\" json:\"-\" xml:\"-\""}},
		{"pad", &gti.Field{Name: "pad", Type: "float32", LocalType: "float32", Doc: "", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/fsfffb.Inhib",
	ShortName: "fsfffb.Inhib",
	IDName:    "inhib",
	Doc:       "Inhib contains state values for computed FFFB inhibition",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gosl", Directive: "start", Args: []string{"fsfffb"}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"FFsRaw", &gti.Field{Name: "FFsRaw", Type: "float32", LocalType: "float32", Doc: "all feedforward incoming spikes into neurons in this pool -- raw aggregation", Directives: gti.Directives{}, Tag: ""}},
		{"FBsRaw", &gti.Field{Name: "FBsRaw", Type: "float32", LocalType: "float32", Doc: "all feedback outgoing spikes generated from neurons in this pool -- raw aggregation", Directives: gti.Directives{}, Tag: ""}},
		{"GeExtRaw", &gti.Field{Name: "GeExtRaw", Type: "float32", LocalType: "float32", Doc: "all extra GeExt conductances added to neurons", Directives: gti.Directives{}, Tag: ""}},
		{"FFs", &gti.Field{Name: "FFs", Type: "float32", LocalType: "float32", Doc: "all feedforward incoming spikes into neurons in this pool, normalized by pool size", Directives: gti.Directives{}, Tag: ""}},
		{"FBs", &gti.Field{Name: "FBs", Type: "float32", LocalType: "float32", Doc: "all feedback outgoing spikes generated from neurons in this pool, normalized by pool size", Directives: gti.Directives{}, Tag: ""}},
		{"GeExts", &gti.Field{Name: "GeExts", Type: "float32", LocalType: "float32", Doc: "all extra GeExt conductances added to neurons, normalized by pool size", Directives: gti.Directives{}, Tag: ""}},
		{"Clamped", &gti.Field{Name: "Clamped", Type: "goki.dev/gosl/v2/slbool.Bool", LocalType: "slbool.Bool", Doc: "if true, this layer is hard-clamped and should use GeExts exclusively for PV", Directives: gti.Directives{}, Tag: ""}},
		{"FSi", &gti.Field{Name: "FSi", Type: "float32", LocalType: "float32", Doc: "fast spiking PV+ fast integration of FFs feedforward spikes", Directives: gti.Directives{}, Tag: ""}},
		{"SSi", &gti.Field{Name: "SSi", Type: "float32", LocalType: "float32", Doc: "slow spiking SST+ integration of FBs feedback spikes", Directives: gti.Directives{}, Tag: ""}},
		{"SSf", &gti.Field{Name: "SSf", Type: "float32", LocalType: "float32", Doc: "slow spiking facilitation factor, representing facilitating effects of recent activity", Directives: gti.Directives{}, Tag: ""}},
		{"FSGi", &gti.Field{Name: "FSGi", Type: "float32", LocalType: "float32", Doc: "overall fast-spiking inhibitory conductance", Directives: gti.Directives{}, Tag: ""}},
		{"SSGi", &gti.Field{Name: "SSGi", Type: "float32", LocalType: "float32", Doc: "overall slow-spiking inhibitory conductance", Directives: gti.Directives{}, Tag: ""}},
		{"Gi", &gti.Field{Name: "Gi", Type: "float32", LocalType: "float32", Doc: "overall inhibitory conductance = FSGi + SSGi", Directives: gti.Directives{}, Tag: ""}},
		{"GiOrig", &gti.Field{Name: "GiOrig", Type: "float32", LocalType: "float32", Doc: "original value of the inhibition (before pool or other effects)", Directives: gti.Directives{}, Tag: ""}},
		{"LayGi", &gti.Field{Name: "LayGi", Type: "float32", LocalType: "float32", Doc: "for pools, this is the layer-level inhibition that is MAX'd with the pool-level inhibition to produce the net inhibition", Directives: gti.Directives{}, Tag: ""}},
		{"FFAvg", &gti.Field{Name: "FFAvg", Type: "float32", LocalType: "float32", Doc: "longer time scale running average FF drive -- used for FFAvgPrv", Directives: gti.Directives{}, Tag: ""}},
		{"FFAvgPrv", &gti.Field{Name: "FFAvgPrv", Type: "float32", LocalType: "float32", Doc: "previous theta cycle FFAvg value -- for FFPrv factor -- updated in Decay function that is called at start of new ThetaCycle", Directives: gti.Directives{}, Tag: ""}},
		{"FFsRawInt", &gti.Field{Name: "FFsRawInt", Type: "int32", LocalType: "int32", Doc: "int32 atomic add compatible integration of FFsRaw", Directives: gti.Directives{}, Tag: ""}},
		{"FBsRawInt", &gti.Field{Name: "FBsRawInt", Type: "int32", LocalType: "int32", Doc: "int32 atomic add compatible integration of FBsRaw", Directives: gti.Directives{}, Tag: ""}},
		{"GeExtRawInt", &gti.Field{Name: "GeExtRawInt", Type: "int32", LocalType: "int32", Doc: "int32 atomic add compatible integration of GeExtRaw", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/fsfffb.Inhibs",
	ShortName: "fsfffb.Inhibs",
	IDName:    "inhibs",
	Doc:       "Inhibs is a slice of Inhib records",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gosl", Directive: "end", Args: []string{"fsfffb"}},
		&gti.Directive{Tool: "gosl", Directive: "hlsl", Args: []string{"fsfffb"}},
		&gti.Directive{Tool: "gosl", Directive: "end", Args: []string{"fsfffb"}},
	},

	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})
