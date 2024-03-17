// Code generated by "core generate -add-types"; DO NOT EDIT.

package fsfffb

import (
	"cogentcore.org/core/gti"
)

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/fsfffb.GiParams", IDName: "gi-params", Doc: "GiParams parameterizes feedforward (FF) and feedback (FB) inhibition (FFFB)\nbased on incoming spikes (FF) and outgoing spikes (FB)\nacross Fast (PV+) and Slow (SST+) timescales.\nFF -> PV -> FS fast spikes, FB -> SST -> SS slow spikes (slow to get going)", Directives: []gti.Directive{{Tool: "gosl", Directive: "start", Args: []string{"fsfffb"}}}, Fields: []gti.Field{{Name: "On", Doc: "enable this level of inhibition"}, {Name: "Gi", Doc: "overall inhibition gain -- this is main parameter to adjust to change overall activation levels -- it scales both the the FS and SS factors uniformly"}, {Name: "FB", Doc: "amount of FB spikes included in FF for driving FS -- for small networks, 0.5 or 1 works best; larger networks and more demanding inhibition requires higher levels."}, {Name: "FSTau", Doc: "fast spiking (PV+) intgration time constant in cycles (msec) -- tau is roughly how long it takes for value to change significantly -- 1.4x the half-life."}, {Name: "SS", Doc: "multiplier on SS slow-spiking (SST+) in contributing to the overall Gi inhibition -- FS contributes at a factor of 1"}, {Name: "SSfTau", Doc: "slow-spiking (SST+) facilitation decay time constant in cycles (msec) -- facilication factor SSf determines impact of FB spikes as a function of spike input-- tau is roughly how long it takes for value to change significantly -- 1.4x the half-life."}, {Name: "SSiTau", Doc: "slow-spiking (SST+) intgration time constant in cycles (msec) cascaded on top of FSTau -- tau is roughly how long it takes for value to change significantly -- 1.4x the half-life."}, {Name: "FS0", Doc: "fast spiking zero point -- below this level, no FS inhibition is computed, and this value is subtracted from the FSi"}, {Name: "FFAvgTau", Doc: "time constant for updating a running average of the feedforward inhibition over a longer time scale, for computing FFPrv"}, {Name: "FFPrv", Doc: "proportion of previous average feed-forward inhibition (FFAvgPrv) to add, resulting in an accentuated temporal-derivative dynamic where neurons respond most strongly to increases in excitation that exceeds inhibition from last time."}, {Name: "ClampExtMin", Doc: "minimum GeExt value required to drive external clamping dynamics (if clamp is set), where only GeExt drives inhibition.  If GeExt is below this value, then the usual FS-FFFB drivers are used."}, {Name: "FSDt", Doc: "rate = 1 / tau"}, {Name: "SSfDt", Doc: "rate = 1 / tau"}, {Name: "SSiDt", Doc: "rate = 1 / tau"}, {Name: "FFAvgDt", Doc: "rate = 1 / tau"}, {Name: "pad"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/fsfffb.Inhib", IDName: "inhib", Doc: "Inhib contains state values for computed FFFB inhibition", Directives: []gti.Directive{{Tool: "gosl", Directive: "start", Args: []string{"fsfffb"}}}, Fields: []gti.Field{{Name: "FFsRaw", Doc: "all feedforward incoming spikes into neurons in this pool -- raw aggregation"}, {Name: "FBsRaw", Doc: "all feedback outgoing spikes generated from neurons in this pool -- raw aggregation"}, {Name: "GeExtRaw", Doc: "all extra GeExt conductances added to neurons"}, {Name: "FFs", Doc: "all feedforward incoming spikes into neurons in this pool, normalized by pool size"}, {Name: "FBs", Doc: "all feedback outgoing spikes generated from neurons in this pool, normalized by pool size"}, {Name: "GeExts", Doc: "all extra GeExt conductances added to neurons, normalized by pool size"}, {Name: "Clamped", Doc: "if true, this layer is hard-clamped and should use GeExts exclusively for PV"}, {Name: "FSi", Doc: "fast spiking PV+ fast integration of FFs feedforward spikes"}, {Name: "SSi", Doc: "slow spiking SST+ integration of FBs feedback spikes"}, {Name: "SSf", Doc: "slow spiking facilitation factor, representing facilitating effects of recent activity"}, {Name: "FSGi", Doc: "overall fast-spiking inhibitory conductance"}, {Name: "SSGi", Doc: "overall slow-spiking inhibitory conductance"}, {Name: "Gi", Doc: "overall inhibitory conductance = FSGi + SSGi"}, {Name: "GiOrig", Doc: "original value of the inhibition (before pool or other effects)"}, {Name: "LayGi", Doc: "for pools, this is the layer-level inhibition that is MAX'd with the pool-level inhibition to produce the net inhibition"}, {Name: "FFAvg", Doc: "longer time scale running average FF drive -- used for FFAvgPrv"}, {Name: "FFAvgPrv", Doc: "previous theta cycle FFAvg value -- for FFPrv factor -- updated in Decay function that is called at start of new ThetaCycle"}, {Name: "FFsRawInt", Doc: "int32 atomic add compatible integration of FFsRaw"}, {Name: "FBsRawInt", Doc: "int32 atomic add compatible integration of FBsRaw"}, {Name: "GeExtRawInt", Doc: "int32 atomic add compatible integration of GeExtRaw"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/fsfffb.Inhibs", IDName: "inhibs", Doc: "Inhibs is a slice of Inhib records", Directives: []gti.Directive{{Tool: "gosl", Directive: "end", Args: []string{"fsfffb"}}, {Tool: "gosl", Directive: "hlsl", Args: []string{"fsfffb"}}, {Tool: "gosl", Directive: "end", Args: []string{"fsfffb"}}}})
