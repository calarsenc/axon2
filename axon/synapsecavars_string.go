// Code generated by "stringer -type=SynapseCaVars"; DO NOT EDIT.

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
	_ = x[CaM-0]
	_ = x[CaP-1]
	_ = x[CaD-2]
	_ = x[CaUpT-3]
	_ = x[Tr-4]
	_ = x[DTr-5]
	_ = x[DiDWt-6]
	_ = x[SynapseCaVarsN-7]
}

const _SynapseCaVars_name = "CaMCaPCaDCaUpTTrDTrDiDWtSynapseCaVarsN"

var _SynapseCaVars_index = [...]uint8{0, 3, 6, 9, 14, 16, 19, 24, 38}

func (i SynapseCaVars) String() string {
	if i < 0 || i >= SynapseCaVars(len(_SynapseCaVars_index)-1) {
		return "SynapseCaVars(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SynapseCaVars_name[_SynapseCaVars_index[i]:_SynapseCaVars_index[i+1]]
}

func (i *SynapseCaVars) FromString(s string) error {
	for j := 0; j < len(_SynapseCaVars_index)-1; j++ {
		if s == _SynapseCaVars_name[_SynapseCaVars_index[j]:_SynapseCaVars_index[j+1]] {
			*i = SynapseCaVars(j)
			return nil
		}
	}
	return errors.New("String: " + s + " is not a valid option for type: SynapseCaVars")
}

var _SynapseCaVars_descMap = map[SynapseCaVars]string{
	0: `CaM is first stage running average (mean) Ca calcium level (like CaM = calmodulin), feeds into CaP`,
	1: `CaP is shorter timescale integrated CaM value, representing the plus, LTP direction of weight change and capturing the function of CaMKII in the Kinase learning rule`,
	2: `CaD is longer timescale integrated CaP value, representing the minus, LTD direction of weight change and capturing the function of DAPK1 in the Kinase learning rule`,
	3: `CaUpT is time in CyclesTotal of last updating of Ca values at the synapse level, for optimized synaptic-level Ca integration -- converted to / from uint32`,
	4: `Tr is trace of synaptic activity over time -- used for credit assignment in learning. In MatrixPrjn this is a tag that is then updated later when US occurs.`,
	5: `DTr is delta (change in) Tr trace of synaptic activity over time`,
	6: `DiDWt is delta weight for each data parallel index (Di) -- this is directly computed from the Ca values (in cortical version) and then aggregated into the overall DWt (which may be further integrated across MPI nodes), which then drives changes in Wt values`,
	7: ``,
}

func (i SynapseCaVars) Desc() string {
	if str, ok := _SynapseCaVars_descMap[i]; ok {
		return str
	}
	return "SynapseCaVars(" + strconv.FormatInt(int64(i), 10) + ")"
}
