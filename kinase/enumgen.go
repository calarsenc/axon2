// Code generated by "core generate -add-types"; DO NOT EDIT.

package kinase

import (
	"cogentcore.org/core/enums"
)

var _RulesValues = []Rules{0, 1, 2, 3}

// RulesN is the highest valid value for type Rules, plus one.
const RulesN Rules = 4

var _RulesValueMap = map[string]Rules{`SynSpkCont`: 0, `SynNMDACont`: 1, `SynSpkTheta`: 2, `NeurSpkTheta`: 3}

var _RulesDescMap = map[Rules]string{0: `SynSpkCont implements synaptic-level Ca signals at an abstract level, purely driven by spikes, not NMDA channel Ca, as a product of sender and recv CaSyn values that capture the decaying Ca trace from spiking, qualitatively as in the NMDA dynamics. These spike-driven Ca signals are integrated in a cascaded manner via CaM, then CaP (reflecting CaMKII) and finally CaD (reflecting DAPK1). It uses continuous learning based on temporary DWt (TDWt) values based on the TWindow around spikes, which convert into DWt after a pause in synaptic activity (no arbitrary ThetaCycle boundaries). There is an option to compare with SynSpkTheta by only doing DWt updates at the theta cycle level, in which case the key difference is the use of TDWt, which can remove some variability associated with the arbitrary timing of the end of trials.`, 1: `SynNMDACont is the same as SynSpkCont with NMDA-driven calcium signals computed according to the very close approximation to the Urakubo et al (2008) allosteric NMDA dynamics, then integrated at P vs. D time scales. This is the most biologically realistic yet computationally tractable verseion of the Kinase learning algorithm.`, 2: `SynSpkTheta abstracts the SynSpkCont algorithm by only computing the DWt change at the end of the ThetaCycle, instead of continuous updating. This allows an optimized implementation that is roughly 1/3 slower than the fastest NeurSpkTheta version, while still capturing much of the learning dynamics by virtue of synaptic-level integration.`, 3: `NeurSpkTheta uses neuron-level spike-driven calcium signals integrated at P vs. D time scales -- this is the original Leabra and Axon XCAL / CHL learning rule. It exhibits strong sensitivity to final spikes and thus high levels of variance.`}

var _RulesMap = map[Rules]string{0: `SynSpkCont`, 1: `SynNMDACont`, 2: `SynSpkTheta`, 3: `NeurSpkTheta`}

// String returns the string representation of this Rules value.
func (i Rules) String() string { return enums.String(i, _RulesMap) }

// SetString sets the Rules value from its string representation,
// and returns an error if the string is invalid.
func (i *Rules) SetString(s string) error { return enums.SetString(i, s, _RulesValueMap, "Rules") }

// Int64 returns the Rules value as an int64.
func (i Rules) Int64() int64 { return int64(i) }

// SetInt64 sets the Rules value from an int64.
func (i *Rules) SetInt64(in int64) { *i = Rules(in) }

// Desc returns the description of the Rules value.
func (i Rules) Desc() string { return enums.Desc(i, _RulesDescMap) }

// RulesValues returns all possible values for the type Rules.
func RulesValues() []Rules { return _RulesValues }

// Values returns all possible values for the type Rules.
func (i Rules) Values() []enums.Enum { return enums.Values(_RulesValues) }

// MarshalText implements the [encoding.TextMarshaler] interface.
func (i Rules) MarshalText() ([]byte, error) { return []byte(i.String()), nil }

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (i *Rules) UnmarshalText(text []byte) error { return enums.UnmarshalText(i, text, "Rules") }
