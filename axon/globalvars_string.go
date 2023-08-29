// Code generated by "stringer -type=GlobalVars"; DO NOT EDIT.

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
	_ = x[GvRew-0]
	_ = x[GvHasRew-1]
	_ = x[GvRewPred-2]
	_ = x[GvPrevPred-3]
	_ = x[GvHadRew-4]
	_ = x[GvDA-5]
	_ = x[GvACh-6]
	_ = x[GvNE-7]
	_ = x[GvSer-8]
	_ = x[GvAChRaw-9]
	_ = x[GvNotMaint-10]
	_ = x[GvEffortRaw-11]
	_ = x[GvEffortCurMax-12]
	_ = x[GvUrgency-13]
	_ = x[GvUrgencyRaw-14]
	_ = x[GvVSMatrixJustGated-15]
	_ = x[GvCuriosityPoolGated-16]
	_ = x[GvVSMatrixHasGated-17]
	_ = x[GvHasPosUS-18]
	_ = x[GvHadPosUS-19]
	_ = x[GvLHbDip-20]
	_ = x[GvLHbBurst-21]
	_ = x[GvLHbPVDA-22]
	_ = x[GvLHbDipSumCur-23]
	_ = x[GvLHbDipSum-24]
	_ = x[GvLHbGiveUp-25]
	_ = x[GvLHbGaveUp-26]
	_ = x[GvLHbVSPatchPos-27]
	_ = x[GvLHbPVposSum-28]
	_ = x[GvLHbPVpos-29]
	_ = x[GvLHbPVnegSum-30]
	_ = x[GvLHbPVneg-31]
	_ = x[GvCeMpos-32]
	_ = x[GvCeMneg-33]
	_ = x[GvVtaDA-34]
	_ = x[GvUSneg-35]
	_ = x[GvUSnegRaw-36]
	_ = x[GvDrives-37]
	_ = x[GvUSpos-38]
	_ = x[GvVSPatch-39]
	_ = x[GlobalVarsN-40]
}

const _GlobalVars_name = "GvRewGvHasRewGvRewPredGvPrevPredGvHadRewGvDAGvAChGvNEGvSerGvAChRawGvNotMaintGvEffortRawGvEffortCurMaxGvUrgencyGvUrgencyRawGvVSMatrixJustGatedGvCuriosityPoolGatedGvVSMatrixHasGatedGvHasPosUSGvHadPosUSGvLHbDipGvLHbBurstGvLHbPVDAGvLHbDipSumCurGvLHbDipSumGvLHbGiveUpGvLHbGaveUpGvLHbVSPatchPosGvLHbPVposSumGvLHbPVposGvLHbPVnegSumGvLHbPVnegGvCeMposGvCeMnegGvVtaDAGvUSnegGvUSnegRawGvDrivesGvUSposGvVSPatchGlobalVarsN"

var _GlobalVars_index = [...]uint16{0, 5, 13, 22, 32, 40, 44, 49, 53, 58, 66, 76, 87, 101, 110, 122, 141, 161, 179, 189, 199, 207, 217, 226, 240, 251, 262, 273, 288, 301, 311, 324, 334, 342, 350, 357, 364, 374, 382, 389, 398, 409}

func (i GlobalVars) String() string {
	if i < 0 || i >= GlobalVars(len(_GlobalVars_index)-1) {
		return "GlobalVars(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _GlobalVars_name[_GlobalVars_index[i]:_GlobalVars_index[i+1]]
}

func (i *GlobalVars) FromString(s string) error {
	for j := 0; j < len(_GlobalVars_index)-1; j++ {
		if s == _GlobalVars_name[_GlobalVars_index[j]:_GlobalVars_index[j+1]] {
			*i = GlobalVars(j)
			return nil
		}
	}
	return errors.New("String: " + s + " is not a valid option for type: GlobalVars")
}

var _GlobalVars_descMap = map[GlobalVars]string{
	0:  `Rew is reward value -- this is set here in the Context struct, and the RL Rew layer grabs it from there -- must also set HasRew flag when rew is set -- otherwise is ignored.`,
	1:  `HasRew must be set to true when a reward is present -- otherwise Rew is ignored. Also set when PVLV BOA model gives up. This drives ACh release in the PVLV model.`,
	2:  `RewPred is reward prediction -- computed by a special reward prediction layer`,
	3:  `PrevPred is previous time step reward prediction -- e.g., for TDPredLayer`,
	4:  `HadRew is HasRew state from the previous trial -- copied from HasRew in NewState -- used for updating Effort, Urgency at start of new trial`,
	5:  `DA is dopamine -- represents reward prediction error, signaled as phasic increases or decreases in activity relative to a tonic baseline, which is represented by a value of 0. Released by the VTA -- ventral tegmental area, or SNc -- substantia nigra pars compacta.`,
	6:  `ACh is acetylcholine -- activated by salient events, particularly at the onset of a reward / punishment outcome (US), or onset of a conditioned stimulus (CS). Driven by BLA -&gt; PPtg that detects changes in BLA activity, via LDTLayer type`,
	7:  `NE is norepinepherine -- not yet in use`,
	8:  `Ser is serotonin -- not yet in use`,
	9:  `AChRaw is raw ACh value used in updating global ACh value by LDTLayer`,
	10: `NotMaint is activity of the PTNotMaintLayer -- drives top-down inhibition of LDT layer / ACh activity.`,
	11: `EffortRaw is raw effort -- increments linearly upward for each additional effort step This is also copied directly into NegUS[0] which tracks effort, but we maintain a separate effort value to make it clearer.`,
	12: `EffortCurMax is current maximum raw effort level -- above this point, any current goal will be terminated during the GiveUp function, which also looks for accumulated disappointment. See Max, MaxNovel, MaxPostDip for values depending on how the goal was triggered`,
	13: `Urgency is the overall urgency activity level (normalized 0-1), computed from logistic function of GvUrgencyRaw`,
	14: `UrgencyRaw is raw effort for urgency -- increments linearly upward from effort increments per step`,
	15: `VSMatrixJustGated is VSMatrix just gated (to engage goal maintenance in PFC areas), set at end of plus phase -- this excludes any gating happening at time of US`,
	16: `CuriosityPoolGated is true if VSMatrixJustGated and the first pool representing the curiosity / novelty drive gated -- this can change the giving up Effort.Max parameter.`,
	17: `VSMatrixHasGated is VSMatrix has gated since the last time HasRew was set (US outcome received or expected one failed to be received`,
	18: `HasPosUS has positive US on this trial`,
	19: `HadPosUS is state from the previous trial -- copied from HasPosUS in NewState -- used for updating Effort, Urgency at start of new trial`,
	20: `computed LHb activity level that drives dipping / pausing of DA firing, when VSPatch pos prediction &gt; actual PV reward drive or PVNeg &gt; PVPos`,
	21: `LHbBurst is computed LHb activity level that drives bursts of DA firing, when actual PV reward drive &gt; VSPatch pos prediction`,
	22: `LHbPVDA is GvLHbBurst - GvLHbDip -- the LHb contribution to DA, reflecting PV and VSPatch (PVi), but not the CS (LV) contributions`,
	23: `LHbDipSumCur is current sum of LHbDip over trials, which is reset when there is a PV value, an above-threshold PPTg value, or when it triggers reset`,
	24: `LHbDipSum is copy of DipSum that is not reset -- used for driving negative dopamine dips on GiveUp trials`,
	25: `LHbGiveUp is true if a reset was triggered from LHbDipSum &gt; Reset Thr`,
	26: `LHbGaveUp is copy of LHbGiveUp from previous trial`,
	27: `LHbVSPatchPos is net shunting input from VSPatch (PosD1, named PVi in original PVLV)`,
	28: `LHbPVposSum is total weighted positive valence primary value = sum of Weight * USpos * Drive`,
	29: `LHbPVpos is positive valence primary value (normalized USpos) = (1 - 1/(1+LHb.PosGain * USpos))`,
	30: `LHbPVnegSum is total weighted negative valence primary value = sum of Weight * USneg`,
	31: `LHbPVpos is positive valence primary value (normalized USpos) = (1 - 1/(1+LHb.NegGain * USpos))`,
	32: `CeMpos is positive valence central nucleus of the amygdala (CeM) LV (learned value) activity, reflecting |BLAPosAcqD1 - BLAPosExtD2|_+ positively rectified. CeM sets Raw directly. Note that a positive US onset even with no active Drive will be reflected here, enabling learning about unexpected outcomes`,
	33: `CeMneg is negative valence central nucleus of the amygdala (CeM) LV (learned value) activity, reflecting |BLANegAcqD2 - BLANegExtD1|_+ positively rectified. CeM sets Raw directly`,
	34: `VtaDA is overall dopamine value reflecting all of the different inputs`,
	35: `USneg are negative valence US outcomes -- normalized version of raw, NNegUSs of them`,
	36: `USnegRaw are raw, linearly incremented negative valence US outcomes, this value is also integrated together with all US vals for PVneg`,
	37: `Drives is current drive state -- updated with optional homeostatic exponential return to baseline values`,
	38: `USpos is current positive-valence drive-satisfying input(s) (unconditioned stimuli = US)`,
	39: `VSPatch is current reward predicting VSPatch (PosD1) values`,
	40: ``,
}

func (i GlobalVars) Desc() string {
	if str, ok := _GlobalVars_descMap[i]; ok {
		return str
	}
	return "GlobalVars(" + strconv.FormatInt(int64(i), 10) + ")"
}
