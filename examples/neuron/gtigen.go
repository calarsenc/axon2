// Code generated by "core generate -add-types"; DO NOT EDIT.

package main

import (
	"cogentcore.org/core/gti"
)

var _ = gti.AddType(&gti.Type{Name: "main.ParamConfig", IDName: "param-config", Doc: "ParamConfig has config parameters related to sim params", Fields: []gti.Field{{Name: "Network", Doc: "network parameters"}, {Name: "Sheet", Doc: "Extra Param Sheet name(s) to use (space separated if multiple) -- must be valid name as listed in compiled-in params or loaded params"}, {Name: "Tag", Doc: "extra tag to add to file names and logs saved from this run"}, {Name: "Note", Doc: "user note -- describe the run params etc -- like a git commit message for the run"}, {Name: "File", Doc: "Name of the JSON file to input saved parameters from."}, {Name: "SaveAll", Doc: "Save a snapshot of all current param and config settings in a directory named params_<datestamp> (or _good if Good is true), then quit -- useful for comparing to later changes and seeing multiple views of current params"}, {Name: "Good", Doc: "for SaveAll, save to params_good for a known good params state.  This can be done prior to making a new release after all tests are passing -- add results to git to provide a full diff record of all params over time."}}})

var _ = gti.AddType(&gti.Type{Name: "main.RunConfig", IDName: "run-config", Doc: "RunConfig has config parameters related to running the sim", Fields: []gti.Field{{Name: "GPU", Doc: "use the GPU for computation -- only for testing in this model -- not faster"}, {Name: "NThreads", Doc: "number of parallel threads for CPU computation -- 0 = use default"}, {Name: "Run", Doc: "starting run number -- determines the random seed -- runs counts from there -- can do all runs in parallel by launching separate jobs with each run, runs = 1"}, {Name: "NRuns", Doc: "total number of runs to do when running Train"}, {Name: "NEpochs", Doc: "total number of epochs per run"}}})

var _ = gti.AddType(&gti.Type{Name: "main.LogConfig", IDName: "log-config", Doc: "LogConfig has config parameters related to logging data", Fields: []gti.Field{{Name: "SaveWts", Doc: "if true, save final weights after each run"}, {Name: "Cycle", Doc: "if true, save cycle log to file, as .cyc.tsv typically"}, {Name: "NetData", Doc: "if true, save network activation etc data from testing trials, for later viewing in netview"}}})

var _ = gti.AddType(&gti.Type{Name: "main.Config", IDName: "config", Doc: "Config is a standard Sim config -- use as a starting point.", Fields: []gti.Field{{Name: "GeClamp", Doc: "clamp constant Ge value -- otherwise drive discrete spiking input"}, {Name: "SpikeHz", Doc: "frequency of input spiking for !GeClamp mode"}, {Name: "Ge", Doc: "Raw synaptic excitatory conductance"}, {Name: "Gi", Doc: "Inhibitory conductance"}, {Name: "ErevE", Doc: "excitatory reversal (driving) potential -- determines where excitation pushes Vm up to"}, {Name: "ErevI", Doc: "leak reversal (driving) potential -- determines where excitation pulls Vm down to"}, {Name: "Noise", Doc: "the variance parameter for Gaussian noise added to unit activations on every cycle"}, {Name: "KNaAdapt", Doc: "apply sodium-gated potassium adaptation mechanisms that cause the neuron to reduce spiking over time"}, {Name: "MahpGbar", Doc: "strength of mAHP M-type channel -- used to be implemented by KNa but now using the more standard M-type channel mechanism"}, {Name: "NMDAGbar", Doc: "strength of NMDA current -- 0.006 default for posterior cortex"}, {Name: "GABABGbar", Doc: "strength of GABAB current -- 0.015 default for posterior cortex"}, {Name: "VGCCGbar", Doc: "strength of VGCC voltage gated calcium current -- only activated during spikes -- this is now an essential part of Ca-driven learning to reflect recv spiking in the Ca signal -- but if too strong leads to runaway excitatory bursting."}, {Name: "AKGbar", Doc: "strength of A-type potassium channel -- this is only active at high (depolarized) membrane potentials -- only during spikes -- useful to counteract VGCC's"}, {Name: "NCycles", Doc: "total number of cycles to run"}, {Name: "OnCycle", Doc: "when does excitatory input into neuron come on?"}, {Name: "OffCycle", Doc: "when does excitatory input into neuron go off?"}, {Name: "UpdtInterval", Doc: "how often to update display (in cycles)"}, {Name: "Includes", Doc: "specify include files here, and after configuration, it contains list of include files added"}, {Name: "GUI", Doc: "open the GUI -- does not automatically run -- if false, then runs automatically and quits"}, {Name: "Debug", Doc: "log debugging information"}, {Name: "Params", Doc: "parameter related configuration options"}, {Name: "Run", Doc: "sim running related configuration options"}, {Name: "Log", Doc: "data logging related configuration options"}}})

var _ = gti.AddType(&gti.Type{Name: "main.NeuronEx", IDName: "neuron-ex", Doc: "Extra state for neuron", Fields: []gti.Field{{Name: "InISI", Doc: "input ISI countdown for spiking mode -- counts up"}}})

var _ = gti.AddType(&gti.Type{Name: "main.Sim", IDName: "sim", Doc: "Sim encapsulates the entire simulation model, and we define all the\nfunctionality as methods on this struct.  This structure keeps all relevant\nstate information organized and available without having to pass everything around\nas arguments to methods, and provides the core GUI interface (note the view tags\nfor the fields which provide hints to how things should be displayed).", Fields: []gti.Field{{Name: "Config", Doc: "simulation configuration parameters -- set by .toml config file and / or args"}, {Name: "Net", Doc: "the network -- click to view / edit parameters for layers, prjns, etc"}, {Name: "NeuronEx", Doc: "extra neuron state for additional channels: VGCC, AK"}, {Name: "Context", Doc: "axon timing parameters and state"}, {Name: "Stats", Doc: "contains computed statistic values"}, {Name: "Logs", Doc: "logging"}, {Name: "Params", Doc: "all parameter management"}, {Name: "Cycle", Doc: "current cycle of updating"}, {Name: "ViewUpdt", Doc: "netview update parameters"}, {Name: "GUI", Doc: "manages all the gui elements"}, {Name: "ValMap", Doc: "map of values for detailed debugging / testing"}}})
