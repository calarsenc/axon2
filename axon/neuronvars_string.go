// Code generated by "stringer -type=NeuronVars"; DO NOT EDIT.

package axon

import (
	"errors"
	"strconv"
)

var _ = errors.New("dummy error")

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Spike-0]
	_ = x[Spiked-1]
	_ = x[Act-2]
	_ = x[ActInt-3]
	_ = x[ActM-4]
	_ = x[ActP-5]
	_ = x[Ext-6]
	_ = x[Target-7]
	_ = x[Ge-8]
	_ = x[Gi-9]
	_ = x[Gk-10]
	_ = x[Inet-11]
	_ = x[Vm-12]
	_ = x[VmDend-13]
	_ = x[ISI-14]
	_ = x[ISIAvg-15]
	_ = x[CaSpkP-16]
	_ = x[CaSpkD-17]
	_ = x[CaSyn-18]
	_ = x[CaSpkM-19]
	_ = x[CaSpkPM-20]
	_ = x[CaLrn-21]
	_ = x[NrnCaM-22]
	_ = x[NrnCaP-23]
	_ = x[NrnCaD-24]
	_ = x[CaDiff-25]
	_ = x[Attn-26]
	_ = x[RLRate-27]
	_ = x[SpkMaxCa-28]
	_ = x[SpkMax-29]
	_ = x[SpkPrv-30]
	_ = x[SpkSt1-31]
	_ = x[SpkSt2-32]
	_ = x[GeNoiseP-33]
	_ = x[GeNoise-34]
	_ = x[GiNoiseP-35]
	_ = x[GiNoise-36]
	_ = x[GeExt-37]
	_ = x[GeRaw-38]
	_ = x[GeSyn-39]
	_ = x[GiRaw-40]
	_ = x[GiSyn-41]
	_ = x[GeInt-42]
	_ = x[GeIntMax-43]
	_ = x[GiInt-44]
	_ = x[GModRaw-45]
	_ = x[GModSyn-46]
	_ = x[GMaintRaw-47]
	_ = x[GMaintSyn-48]
	_ = x[SSGi-49]
	_ = x[SSGiDend-50]
	_ = x[Gak-51]
	_ = x[MahpN-52]
	_ = x[SahpCa-53]
	_ = x[SahpN-54]
	_ = x[GknaMed-55]
	_ = x[GknaSlow-56]
	_ = x[GnmdaSyn-57]
	_ = x[Gnmda-58]
	_ = x[GnmdaMaint-59]
	_ = x[GnmdaLrn-60]
	_ = x[NmdaCa-61]
	_ = x[GgabaB-62]
	_ = x[GABAB-63]
	_ = x[GABABx-64]
	_ = x[Gvgcc-65]
	_ = x[VgccM-66]
	_ = x[VgccH-67]
	_ = x[VgccCa-68]
	_ = x[VgccCaInt-69]
	_ = x[SKCaIn-70]
	_ = x[SKCaR-71]
	_ = x[SKCaM-72]
	_ = x[Gsk-73]
	_ = x[Burst-74]
	_ = x[BurstPrv-75]
	_ = x[CtxtGe-76]
	_ = x[CtxtGeRaw-77]
	_ = x[CtxtGeOrig-78]
	_ = x[NrnFlags-79]
	_ = x[NeuronVarsN-80]
}

const _NeuronVars_name = "SpikeSpikedActActIntActMActPExtTargetGeGiGkInetVmVmDendISIISIAvgCaSpkPCaSpkDCaSynCaSpkMCaSpkPMCaLrnNrnCaMNrnCaPNrnCaDCaDiffAttnRLRateSpkMaxCaSpkMaxSpkPrvSpkSt1SpkSt2GeNoisePGeNoiseGiNoisePGiNoiseGeExtGeRawGeSynGiRawGiSynGeIntGeIntMaxGiIntGModRawGModSynGMaintRawGMaintSynSSGiSSGiDendGakMahpNSahpCaSahpNGknaMedGknaSlowGnmdaSynGnmdaGnmdaMaintGnmdaLrnNmdaCaGgabaBGABABGABABxGvgccVgccMVgccHVgccCaVgccCaIntSKCaInSKCaRSKCaMGskBurstBurstPrvCtxtGeCtxtGeRawCtxtGeOrigNrnFlagsNeuronVarsN"

var _NeuronVars_index = [...]uint16{0, 5, 11, 14, 20, 24, 28, 31, 37, 39, 41, 43, 47, 49, 55, 58, 64, 70, 76, 81, 87, 94, 99, 105, 111, 117, 123, 127, 133, 141, 147, 153, 159, 165, 173, 180, 188, 195, 200, 205, 210, 215, 220, 225, 233, 238, 245, 252, 261, 270, 274, 282, 285, 290, 296, 301, 308, 316, 324, 329, 339, 347, 353, 359, 364, 370, 375, 380, 385, 391, 400, 406, 411, 416, 419, 424, 432, 438, 447, 457, 465, 476}

func (i NeuronVars) String() string {
	if i < 0 || i >= NeuronVars(len(_NeuronVars_index)-1) {
		return "NeuronVars(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NeuronVars_name[_NeuronVars_index[i]:_NeuronVars_index[i+1]]
}

func (i *NeuronVars) FromString(s string) error {
	for j := 0; j < len(_NeuronVars_index)-1; j++ {
		if s == _NeuronVars_name[_NeuronVars_index[j]:_NeuronVars_index[j+1]] {
			*i = NeuronVars(j)
			return nil
		}
	}
	return errors.New("String: " + s + " is not a valid option for type: NeuronVars")
}
