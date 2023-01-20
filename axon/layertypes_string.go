// Code generated by "stringer -type=LayerTypes"; DO NOT EDIT.

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
	_ = x[SuperLayer-0]
	_ = x[InputLayer-1]
	_ = x[TargetLayer-2]
	_ = x[CompareLayer-3]
	_ = x[CTLayer-4]
	_ = x[PulvinarLayer-5]
	_ = x[TRNLayer-6]
	_ = x[RewLayer-7]
	_ = x[RSalienceLayer-8]
	_ = x[RWPredLayer-9]
	_ = x[RWDaLayer-10]
	_ = x[TDPredLayer-11]
	_ = x[TDIntegLayer-12]
	_ = x[TDDaLayer-13]
	_ = x[LayerTypesN-14]
}

const _LayerTypes_name = "SuperLayerInputLayerTargetLayerCompareLayerCTLayerPulvinarLayerTRNLayerRewLayerRSalienceLayerRWPredLayerRWDaLayerTDPredLayerTDIntegLayerTDDaLayerLayerTypesN"

var _LayerTypes_index = [...]uint8{0, 10, 20, 31, 43, 50, 63, 71, 79, 93, 104, 113, 124, 136, 145, 156}

func (i LayerTypes) String() string {
	if i < 0 || i >= LayerTypes(len(_LayerTypes_index)-1) {
		return "LayerTypes(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _LayerTypes_name[_LayerTypes_index[i]:_LayerTypes_index[i+1]]
}

func (i *LayerTypes) FromString(s string) error {
	for j := 0; j < len(_LayerTypes_index)-1; j++ {
		if s == _LayerTypes_name[_LayerTypes_index[j]:_LayerTypes_index[j+1]] {
			*i = LayerTypes(j)
			return nil
		}
	}
	return errors.New("String: " + s + " is not a valid option for type: LayerTypes")
}
