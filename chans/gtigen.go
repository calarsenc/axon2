// Code generated by "goki generate -add-types"; DO NOT EDIT.

package chans

import (
	"goki.dev/gti"
	"goki.dev/ordmap"
)

var _ = gti.AddType(&gti.Type{
	Name:       "github.com/emer/axon/v2/chans.AKParams",
	ShortName:  "chans.AKParams",
	IDName:     "ak-params",
	Doc:        "AKParams control an A-type K channel, which is voltage gated with maximal\nactivation around -37 mV.  It has two state variables, M (v-gated opening)\nand H (v-gated closing), which integrate with fast and slow time constants,\nrespectively.  H relatively quickly hits an asymptotic level of inactivation\nfor sustained activity patterns.\nIt is particularly important for counteracting the excitatory effects of\nvoltage gated calcium channels which can otherwise drive runaway excitatory currents.\nSee AKsParams for a much simpler version that works fine when full AP-like spikes are\nnot simulated, as in our standard axon models.",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Gbar", &gti.Field{Name: "Gbar", Type: "float32", LocalType: "float32", Doc: "strength of AK current", Directives: gti.Directives{}, Tag: "def:\"1,0.1,0.01\""}},
		{"Beta", &gti.Field{Name: "Beta", Type: "float32", LocalType: "float32", Doc: "multiplier for the beta term; 0.01446 for distal, 0.02039 for proximal dendrites", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"0.01446,02039\""}},
		{"Dm", &gti.Field{Name: "Dm", Type: "float32", LocalType: "float32", Doc: "Dm factor: 0.5 for distal, 0.25 for proximal", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"0.5,0.25\""}},
		{"Koff", &gti.Field{Name: "Koff", Type: "float32", LocalType: "float32", Doc: "offset for K, 1.8 for distal, 1.5 for proximal", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"1.8,1.5\""}},
		{"Voff", &gti.Field{Name: "Voff", Type: "float32", LocalType: "float32", Doc: "voltage offset for alpha and beta functions: 1 for distal, 11 for proximal", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"1,11\""}},
		{"Hf", &gti.Field{Name: "Hf", Type: "float32", LocalType: "float32", Doc: "h multiplier factor, 0.1133 for distal, 0.1112 for proximal", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"0.1133,0.1112\""}},
		{"pad", &gti.Field{Name: "pad", Type: "float32", LocalType: "float32", Doc: "", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/v2/chans.AKsParams",
	ShortName: "chans.AKsParams",
	IDName:    "a-ks-params",
	Doc:       "AKsParams provides a highly simplified stateless A-type K channel\nthat only has the voltage-gated activation (M) dynamic with a cutoff\nthat ends up capturing a close approximation to the much more complex AK function.\nThis is voltage gated with maximal activation around -37 mV.\nIt is particularly important for counteracting the excitatory effects of\nvoltage gated calcium channels which can otherwise drive runaway excitatory currents.",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gosl", Directive: "start", Args: []string{"chans"}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Gbar", &gti.Field{Name: "Gbar", Type: "float32", LocalType: "float32", Doc: "strength of AK current", Directives: gti.Directives{}, Tag: "def:\"2,0.1,0.01\""}},
		{"Hf", &gti.Field{Name: "Hf", Type: "float32", LocalType: "float32", Doc: "H factor as a constant multiplier on overall M factor result -- rescales M to level consistent with H being present at full strength", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"0.076\""}},
		{"Mf", &gti.Field{Name: "Mf", Type: "float32", LocalType: "float32", Doc: "multiplier for M -- determines slope of function", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"0.075\""}},
		{"Voff", &gti.Field{Name: "Voff", Type: "float32", LocalType: "float32", Doc: "voltage offset in biological units for M function", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"2\""}},
		{"Vmax", &gti.Field{Name: "Vmax", Type: "float32", LocalType: "float32", Doc: "", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:-37\" desc:\"voltage level of maximum channel opening -- stays flat above that\""}},
		{"pad", &gti.Field{Name: "pad", Type: "int32", LocalType: "int32", Doc: "", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/v2/chans.Chans",
	ShortName: "chans.Chans",
	IDName:    "chans",
	Doc:       "Chans are ion channels used in computing point-neuron activation function",
	Directives: gti.Directives{
		&gti.Directive{Tool: "go", Directive: "generate", Args: []string{"goki", "generate", "-add-types"}},
		&gti.Directive{Tool: "gosl", Directive: "hlsl", Args: []string{"chans"}},
		&gti.Directive{Tool: "gosl", Directive: "end", Args: []string{"chans"}},
		&gti.Directive{Tool: "gosl", Directive: "start", Args: []string{"chans"}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"E", &gti.Field{Name: "E", Type: "float32", LocalType: "float32", Doc: "excitatory sodium (Na) AMPA channels activated by synaptic glutamate", Directives: gti.Directives{}, Tag: ""}},
		{"L", &gti.Field{Name: "L", Type: "float32", LocalType: "float32", Doc: "constant leak (potassium, K+) channels -- determines resting potential (typically higher than resting potential of K)", Directives: gti.Directives{}, Tag: ""}},
		{"I", &gti.Field{Name: "I", Type: "float32", LocalType: "float32", Doc: "inhibitory chloride (Cl-) channels activated by synaptic GABA", Directives: gti.Directives{}, Tag: ""}},
		{"K", &gti.Field{Name: "K", Type: "float32", LocalType: "float32", Doc: "gated / active potassium channels -- typically hyperpolarizing relative to leak / rest", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/v2/chans.GABABParams",
	ShortName: "chans.GABABParams",
	IDName:    "gabab-params",
	Doc:       "GABABParams control the GABAB dynamics in PFC Maint neurons,\nbased on Brunel & Wang (2001) parameters.",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gosl", Directive: "start", Args: []string{"chans"}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Gbar", &gti.Field{Name: "Gbar", Type: "float32", LocalType: "float32", Doc: "overall strength multiplier of GABA-B current. The 0.015 default is a high value that works well in smaller networks -- larger networks may benefit from lower levels (e.g., 0.012).", Directives: gti.Directives{}, Tag: "def:\"0,0.012,0.015\""}},
		{"RiseTau", &gti.Field{Name: "RiseTau", Type: "float32", LocalType: "float32", Doc: "rise time for bi-exponential time dynamics of GABA-B", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"45\""}},
		{"DecayTau", &gti.Field{Name: "DecayTau", Type: "float32", LocalType: "float32", Doc: "decay time for bi-exponential time dynamics of GABA-B", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"50\""}},
		{"Gbase", &gti.Field{Name: "Gbase", Type: "float32", LocalType: "float32", Doc: "baseline level of GABA-B channels open independent of inhibitory input (is added to spiking-produced conductance)", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"0.2\""}},
		{"GiSpike", &gti.Field{Name: "GiSpike", Type: "float32", LocalType: "float32", Doc: "multiplier for converting Gi to equivalent GABA spikes", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"10\""}},
		{"MaxTime", &gti.Field{Name: "MaxTime", Type: "float32", LocalType: "float32", Doc: "time offset when peak conductance occurs, in msec, computed from RiseTau and DecayTau", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" inactive:\"+\""}},
		{"TauFact", &gti.Field{Name: "TauFact", Type: "float32", LocalType: "float32", Doc: "time constant factor used in integration: (Decay / Rise) ^ (Rise / (Decay - Rise))", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"RiseDt", &gti.Field{Name: "RiseDt", Type: "float32", LocalType: "float32", Doc: "1/Tau", Directives: gti.Directives{}, Tag: "view:\"-\" inactive:\"+\""}},
		{"DecayDt", &gti.Field{Name: "DecayDt", Type: "float32", LocalType: "float32", Doc: "1/Tau", Directives: gti.Directives{}, Tag: "view:\"-\" inactive:\"+\""}},
		{"pad", &gti.Field{Name: "pad", Type: "float32", LocalType: "float32", Doc: "", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/v2/chans.KNaParams",
	ShortName: "chans.KNaParams",
	IDName:    "k-na-params",
	Doc:       "KNaParams implements sodium (Na) gated potassium (K) currents\nthat drive adaptation (accommodation) in neural firing.\nAs neurons spike, driving an influx of Na, this activates\nthe K channels, which, like leak channels, pull the membrane\npotential back down toward rest (or even below).",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gosl", Directive: "start", Args: []string{"chans"}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"On", &gti.Field{Name: "On", Type: "goki.dev/gosl/v2/slbool.Bool", LocalType: "slbool.Bool", Doc: "if On, use this component of K-Na adaptation", Directives: gti.Directives{}, Tag: ""}},
		{"Rise", &gti.Field{Name: "Rise", Type: "float32", LocalType: "float32", Doc: "Rise rate of fast time-scale adaptation as function of Na concentration due to spiking -- directly multiplies -- 1/rise = tau for rise rate", Directives: gti.Directives{}, Tag: "viewif:\"On\""}},
		{"Max", &gti.Field{Name: "Max", Type: "float32", LocalType: "float32", Doc: "Maximum potential conductance of fast K channels -- divide nA biological value by 10 for the normalized units here", Directives: gti.Directives{}, Tag: "viewif:\"On\""}},
		{"Tau", &gti.Field{Name: "Tau", Type: "float32", LocalType: "float32", Doc: "time constant in cycles for decay of adaptation, which should be milliseconds typically (tau is roughly how long it takes for value to change significantly -- 1.4x the half-life)", Directives: gti.Directives{}, Tag: "viewif:\"On\""}},
		{"Dt", &gti.Field{Name: "Dt", Type: "float32", LocalType: "float32", Doc: "1/Tau rate constant", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"pad", &gti.Field{Name: "pad", Type: "int32", LocalType: "int32", Doc: "", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "github.com/emer/axon/v2/chans.KNaMedSlow",
	ShortName:  "chans.KNaMedSlow",
	IDName:     "k-na-med-slow",
	Doc:        "KNaMedSlow describes sodium-gated potassium channel adaptation mechanism.\nEvidence supports 2 different time constants:\nSlick (medium) and Slack (slow)",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"On", &gti.Field{Name: "On", Type: "goki.dev/gosl/v2/slbool.Bool", LocalType: "slbool.Bool", Doc: "if On, apply K-Na adaptation", Directives: gti.Directives{}, Tag: ""}},
		{"TrialSlow", &gti.Field{Name: "TrialSlow", Type: "goki.dev/gosl/v2/slbool.Bool", LocalType: "slbool.Bool", Doc: "engages an optional version of Slow that discretely turns on at start of new trial (NewState): nrn.GknaSlow += Slow.Max * nrn.SpkPrv -- achieves a strong form of adaptation", Directives: gti.Directives{}, Tag: ""}},
		{"pad", &gti.Field{Name: "pad", Type: "int32", LocalType: "int32", Doc: "", Directives: gti.Directives{}, Tag: ""}},
		{"Med", &gti.Field{Name: "Med", Type: "github.com/emer/axon/v2/chans.KNaParams", LocalType: "KNaParams", Doc: "medium time-scale adaptation", Directives: gti.Directives{}, Tag: "viewif:\"On\" view:\"inline\""}},
		{"Slow", &gti.Field{Name: "Slow", Type: "github.com/emer/axon/v2/chans.KNaParams", LocalType: "KNaParams", Doc: "slow time-scale adaptation", Directives: gti.Directives{}, Tag: "viewif:\"On\" view:\"inline\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/v2/chans.MahpParams",
	ShortName: "chans.MahpParams",
	IDName:    "mahp-params",
	Doc:       "MahpParams implements an M-type medium afterhyperpolarizing (mAHP) channel,\nwhere m also stands for muscarinic due to the ACh inactivation of this channel.\nIt has a slow activation and deactivation time constant, and opens at a lowish\nmembrane potential.\nThere is one gating variable n updated over time with a tau that is also voltage dependent.\nThe infinite-time value of n is voltage dependent according to a logistic function\nof the membrane potential, centered at Voff with slope Vslope.",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gosl", Directive: "start", Args: []string{"chans"}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Gbar", &gti.Field{Name: "Gbar", Type: "float32", LocalType: "float32", Doc: "strength of mAHP current", Directives: gti.Directives{}, Tag: ""}},
		{"Voff", &gti.Field{Name: "Voff", Type: "float32", LocalType: "float32", Doc: "voltage offset (threshold) in biological units for infinite time N gating function -- where the gate is at 50% strength", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"-30\""}},
		{"Vslope", &gti.Field{Name: "Vslope", Type: "float32", LocalType: "float32", Doc: "slope of the arget (infinite time) gating function", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"9\""}},
		{"TauMax", &gti.Field{Name: "TauMax", Type: "float32", LocalType: "float32", Doc: "maximum slow rate time constant in msec for activation / deactivation.  The effective Tau is much slower -- 1/20th in original temp, and 1/60th in standard 37 C temp", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"1000\""}},
		{"Tadj", &gti.Field{Name: "Tadj", Type: "float32", LocalType: "float32", Doc: "temperature adjustment factor: assume temp = 37 C, whereas original units were at 23 C", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" view:\"-\" inactive:\"+\""}},
		{"DtMax", &gti.Field{Name: "DtMax", Type: "float32", LocalType: "float32", Doc: "1/Tau", Directives: gti.Directives{}, Tag: "view:\"-\" inactive:\"+\""}},
		{"pad", &gti.Field{Name: "pad", Type: "int32", LocalType: "int32", Doc: "", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/v2/chans.NMDAParams",
	ShortName: "chans.NMDAParams",
	IDName:    "nmda-params",
	Doc:       "NMDAParams control the NMDA dynamics, based on Jahr & Stevens (1990) equations\nwhich are widely used in models, from Brunel & Wang (2001) to Sanders et al. (2013).\nThe overall conductance is a function of a voltage-dependent postsynaptic factor based\non Mg ion blockage, and presynaptic Glu-based opening, which in a simple model just\nincrements",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gosl", Directive: "start", Args: []string{"chans"}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Gbar", &gti.Field{Name: "Gbar", Type: "float32", LocalType: "float32", Doc: "overall multiplier for strength of NMDA current -- multiplies GnmdaSyn to get net conductance.", Directives: gti.Directives{}, Tag: "def:\"0,0.006,0.007\""}},
		{"Tau", &gti.Field{Name: "Tau", Type: "float32", LocalType: "float32", Doc: "decay time constant for NMDA channel activation  -- rise time is 2 msec and not worth extra effort for biexponential.  30 fits the Urakubo et al (2008) model with ITau = 100, but 100 works better in practice is small networks so far.", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"30,50,100,200,300\""}},
		{"ITau", &gti.Field{Name: "ITau", Type: "float32", LocalType: "float32", Doc: "decay time constant for NMDA channel inhibition, which captures the Urakubo et al (2008) allosteric dynamics (100 fits their model well) -- set to 1 to eliminate that mechanism.", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"1,100\""}},
		{"MgC", &gti.Field{Name: "MgC", Type: "float32", LocalType: "float32", Doc: "magnesium ion concentration: Brunel & Wang (2001) and Sanders et al (2013) use 1 mM, based on Jahr & Stevens (1990). Urakubo et al (2008) use 1.5 mM. 1.4 with Voff = 5 works best so far in large models, 1.2, Voff = 0 best in smaller nets.", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"1:1.5\""}},
		{"Voff", &gti.Field{Name: "Voff", Type: "float32", LocalType: "float32", Doc: "offset in membrane potential in biological units for voltage-dependent functions.  5 corresponds to the -65 mV rest, -45 threshold of the Urakubo et al (2008) model. 5 was used before in a buggy version of NMDA equation -- 0 is new default.", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"0\""}},
		{"Dt", &gti.Field{Name: "Dt", Type: "float32", LocalType: "float32", Doc: "rate = 1 / tau", Directives: gti.Directives{}, Tag: "view:\"-\" json:\"-\" xml:\"-\""}},
		{"IDt", &gti.Field{Name: "IDt", Type: "float32", LocalType: "float32", Doc: "rate = 1 / tau", Directives: gti.Directives{}, Tag: "view:\"-\" json:\"-\" xml:\"-\""}},
		{"MgFact", &gti.Field{Name: "MgFact", Type: "float32", LocalType: "float32", Doc: "MgFact = MgC / 3.57", Directives: gti.Directives{}, Tag: "view:\"-\" json:\"-\" xml:\"-\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/v2/chans.SahpParams",
	ShortName: "chans.SahpParams",
	IDName:    "sahp-params",
	Doc:       "SahpParams implements a slow afterhyperpolarizing (sAHP) channel,\nIt has a slowly accumulating calcium value, aggregated at the\ntheta cycle level, that then drives the logistic gating function,\nso that it only activates after a significant accumulation.\nAfter which point it decays.\nFor the theta-cycle updating, the normal m-type tau is all within\nthe scope of a single theta cycle, so we just omit the time integration\nof the n gating value, but tau is computed in any case.",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gosl", Directive: "start", Args: []string{"chans"}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Gbar", &gti.Field{Name: "Gbar", Type: "float32", LocalType: "float32", Doc: "strength of sAHP current", Directives: gti.Directives{}, Tag: "def:\"0.05,0.1\""}},
		{"CaTau", &gti.Field{Name: "CaTau", Type: "float32", LocalType: "float32", Doc: "time constant for integrating Ca across theta cycles", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"5,10\""}},
		{"Off", &gti.Field{Name: "Off", Type: "float32", LocalType: "float32", Doc: "integrated Ca offset (threshold) for infinite time N gating function -- where the gate is at 50% strength", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"0.8\""}},
		{"Slope", &gti.Field{Name: "Slope", Type: "float32", LocalType: "float32", Doc: "slope of the infinite time logistic gating function", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"0.02\""}},
		{"TauMax", &gti.Field{Name: "TauMax", Type: "float32", LocalType: "float32", Doc: "maximum slow rate time constant in msec for activation / deactivation.  The effective Tau is much slower -- 1/20th in original temp, and 1/60th in standard 37 C temp", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"1\""}},
		{"CaDt", &gti.Field{Name: "CaDt", Type: "float32", LocalType: "float32", Doc: "1/Tau", Directives: gti.Directives{}, Tag: "view:\"-\" inactive:\"+\""}},
		{"DtMax", &gti.Field{Name: "DtMax", Type: "float32", LocalType: "float32", Doc: "1/Tau", Directives: gti.Directives{}, Tag: "view:\"-\" inactive:\"+\""}},
		{"pad", &gti.Field{Name: "pad", Type: "int32", LocalType: "int32", Doc: "", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/v2/chans.SKCaParams",
	ShortName: "chans.SKCaParams",
	IDName:    "sk-ca-params",
	Doc:       "SKCaParams describes the small-conductance calcium-activated potassium channel,\nactivated by intracellular stores in a way that drives pauses in firing,\nand can require inactivity to recharge the Ca available for release.\nThese intracellular stores can release quickly, have a slow decay once released,\nand the stores can take a while to rebuild, leading to rapidly triggered,\nlong-lasting pauses that don't recur until stores have rebuilt, which is the\nobserved pattern of firing of STNp pausing neurons.\nCaIn = intracellular stores available for release; CaR = released amount from stores\nCaM = K channel conductance gating factor driven by CaR binding,\ncomputed using the Hill equations described in Fujita et al (2012), Gunay et al (2008)\n(also Muddapu & Chakravarthy, 2021): X^h / (X^h + C50^h) where h ~= 4 (hard coded)",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gosl", Directive: "start", Args: []string{"chans"}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Gbar", &gti.Field{Name: "Gbar", Type: "float32", LocalType: "float32", Doc: "overall strength of sKCa current -- inactive if 0", Directives: gti.Directives{}, Tag: "def:\"0,2,3\""}},
		{"C50", &gti.Field{Name: "C50", Type: "float32", LocalType: "float32", Doc: "50% Ca concentration baseline value in Hill equation -- set this to level that activates at reasonable levels of SKCaR", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"0.4,0.5\""}},
		{"ActTau", &gti.Field{Name: "ActTau", Type: "float32", LocalType: "float32", Doc: "K channel gating factor activation time constant -- roughly 5-15 msec in literature", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"15\""}},
		{"DeTau", &gti.Field{Name: "DeTau", Type: "float32", LocalType: "float32", Doc: "K channel gating factor deactivation time constant -- roughly 30-50 msec in literature", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"30\""}},
		{"KCaR", &gti.Field{Name: "KCaR", Type: "float32", LocalType: "float32", Doc: "proportion of CaIn intracellular stores that are released per spike, going into CaR", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"0.4,0.8\""}},
		{"CaRDecayTau", &gti.Field{Name: "CaRDecayTau", Type: "float32", LocalType: "float32", Doc: "SKCaR released calcium decay time constant", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"150,200\""}},
		{"CaInThr", &gti.Field{Name: "CaInThr", Type: "float32", LocalType: "float32", Doc: "level of time-integrated spiking activity (CaSpkD) below which CaIn intracelluar stores are replenished -- a low threshold can be used to require minimal activity to recharge -- set to a high value (e.g., 10) for constant recharge.", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"0.01\""}},
		{"CaInTau", &gti.Field{Name: "CaInTau", Type: "float32", LocalType: "float32", Doc: "time constant in msec for storing CaIn when activity is below CaInThr", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"50\""}},
		{"ActDt", &gti.Field{Name: "ActDt", Type: "float32", LocalType: "float32", Doc: "rate = 1 / tau", Directives: gti.Directives{}, Tag: "view:\"-\" json:\"-\" xml:\"-\""}},
		{"DeDt", &gti.Field{Name: "DeDt", Type: "float32", LocalType: "float32", Doc: "rate = 1 / tau", Directives: gti.Directives{}, Tag: "view:\"-\" json:\"-\" xml:\"-\""}},
		{"CaRDecayDt", &gti.Field{Name: "CaRDecayDt", Type: "float32", LocalType: "float32", Doc: "rate = 1 / tau", Directives: gti.Directives{}, Tag: "view:\"-\" json:\"-\" xml:\"-\""}},
		{"CaInDt", &gti.Field{Name: "CaInDt", Type: "float32", LocalType: "float32", Doc: "rate = 1 / tau", Directives: gti.Directives{}, Tag: "view:\"-\" json:\"-\" xml:\"-\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/v2/chans.VGCCParams",
	ShortName: "chans.VGCCParams",
	IDName:    "vgcc-params",
	Doc:       "VGCCParams control the standard L-type Ca channel\nAll functions based on Urakubo et al (2008).\nSource code available at http://kurodalab.bs.s.u-tokyo.ac.jp/info/STDP/Urakubo2008.tar.gz.\nIn particular look at the file MODEL/Poirazi_cell/CaL.g.",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gosl", Directive: "start", Args: []string{"chans"}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Gbar", &gti.Field{Name: "Gbar", Type: "float32", LocalType: "float32", Doc: "strength of VGCC current -- 0.12 value is from Urakubo et al (2008) model -- best fits actual model behavior using axon equations (1.5 nominal in that model), 0.02 works better in practice for not getting stuck in high plateau firing", Directives: gti.Directives{}, Tag: "def:\"0.02,0.12\""}},
		{"Ca", &gti.Field{Name: "Ca", Type: "float32", LocalType: "float32", Doc: "calcium from conductance factor -- important for learning contribution of VGCC", Directives: gti.Directives{}, Tag: "viewif:\"Gbar>0\" def:\"25\""}},
		{"pad", &gti.Field{Name: "pad", Type: "int32", LocalType: "int32", Doc: "", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})
