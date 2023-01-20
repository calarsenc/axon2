// Code generated by "stringer -type=PrjnTypes"; DO NOT EDIT.

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
	_ = x[ForwardPrj-0]
	_ = x[BackPrjn-1]
	_ = x[LateralPrjn-2]
	_ = x[InhibPrjn-3]
	_ = x[CTCtxtPrjn-4]
	_ = x[RWPrjn-5]
	_ = x[TDRewPredPrjn-6]
	_ = x[PrjnTypesN-7]
}

const _PrjnTypes_name = "ForwardPrjBackPrjnLateralPrjnInhibPrjnCTCtxtPrjnRWPrjnTDRewPredPrjnPrjnTypesN"

var _PrjnTypes_index = [...]uint8{0, 10, 18, 29, 38, 48, 54, 67, 77}

func (i PrjnTypes) String() string {
	if i < 0 || i >= PrjnTypes(len(_PrjnTypes_index)-1) {
		return "PrjnTypes(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PrjnTypes_name[_PrjnTypes_index[i]:_PrjnTypes_index[i+1]]
}

func (i *PrjnTypes) FromString(s string) error {
	for j := 0; j < len(_PrjnTypes_index)-1; j++ {
		if s == _PrjnTypes_name[_PrjnTypes_index[j]:_PrjnTypes_index[j+1]] {
			*i = PrjnTypes(j)
			return nil
		}
	}
	return errors.New("String: " + s + " is not a valid option for type: PrjnTypes")
}
