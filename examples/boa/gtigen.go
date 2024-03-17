// Code generated by "core generate -add-types"; DO NOT EDIT.

package main

import (
	"cogentcore.org/core/gti"
)

var _ = gti.AddType(&gti.Type{Name: "main.Sim", IDName: "sim", Doc: "Sim encapsulates the entire simulation model, and we define all the\nfunctionality as methods on this struct.  This structure keeps all relevant\nstate information organized and available without having to pass everything around\nas arguments to methods, and provides the core GUI interface (note the view tags\nfor the fields which provide hints to how things should be displayed).", Fields: []gti.Field{{Name: "Config", Doc: "simulation configuration parameters -- set by .toml config file and / or args"}, {Name: "Net", Doc: "the network -- click to view / edit parameters for layers, prjns, etc"}, {Name: "StopOnSeq", Doc: "if true, stop running at end of a sequence (for NetView Di data parallel index)"}, {Name: "StopOnErr", Doc: "if true, stop running when an error programmed into the code occurs"}, {Name: "Params", Doc: "network parameter management"}, {Name: "Loops", Doc: "contains looper control loops for running sim"}, {Name: "Stats", Doc: "contains computed statistic values"}, {Name: "Logs", Doc: "Contains all the logs and information about the logs.'"}, {Name: "Envs", Doc: "Environments"}, {Name: "Context", Doc: "axon timing parameters and state"}, {Name: "ViewUpdt", Doc: "netview update parameters"}, {Name: "GUI", Doc: "manages all the gui elements"}, {Name: "EnvGUI", Doc: "gui for viewing env"}, {Name: "RndSeeds", Doc: "a list of random seeds to use for each run"}, {Name: "TestData", Doc: "testing data, from -test arg"}}})

var _ = gti.AddType(&gti.Type{Name: "main.EnvConfig", IDName: "env-config", Doc: "EnvConfig has config params for environment\nnote: only adding fields for key Env params that matter for both Network and Env\nother params are set via the Env map data mechanism.", Fields: []gti.Field{{Name: "Config", Doc: "name of config file that loads into Env.Config for setting environment parameters directly"}, {Name: "NDrives", Doc: "number of different drive-like body states (hunger, thirst, etc), that are satisfied by a corresponding US outcome"}, {Name: "PctCortexStEpc", Doc: "epoch when PctCortex starts increasing"}, {Name: "PctCortexNEpc", Doc: "number of epochs over which PctCortexMax is reached"}, {Name: "PctCortex", Doc: "proportion of behavioral approach sequences driven by the cortex vs. hard-coded reflexive subcortical"}, {Name: "SameSeed", Doc: "for testing, force each env to use same seed"}}})

var _ = gti.AddType(&gti.Type{Name: "main.ParamConfig", IDName: "param-config", Doc: "ParamConfig has config parameters related to sim params", Fields: []gti.Field{{Name: "PVLV", Doc: "PVLV parameters -- can set any field/subfield on Net.PVLV params, using standard TOML formatting"}, {Name: "Network", Doc: "network parameters"}, {Name: "Sheet", Doc: "Extra Param Sheet name(s) to use (space separated if multiple) -- must be valid name as listed in compiled-in params or loaded params"}, {Name: "Tag", Doc: "extra tag to add to file names and logs saved from this run"}, {Name: "Note", Doc: "user note -- describe the run params etc -- like a git commit message for the run"}, {Name: "File", Doc: "Name of the JSON file to input saved parameters from."}, {Name: "SaveAll", Doc: "Save a snapshot of all current param and config settings in a directory named params_<datestamp> (or _good if Good is true), then quit -- useful for comparing to later changes and seeing multiple views of current params"}, {Name: "Good", Doc: "for SaveAll, save to params_good for a known good params state.  This can be done prior to making a new release after all tests are passing -- add results to git to provide a full diff record of all params over time."}}})

var _ = gti.AddType(&gti.Type{Name: "main.RunConfig", IDName: "run-config", Doc: "RunConfig has config parameters related to running the sim", Fields: []gti.Field{{Name: "GPU", Doc: "use the GPU for computation -- generally faster even for small models if NData ~16"}, {Name: "NData", Doc: "number of data-parallel items to process in parallel per trial -- works (and is significantly faster) for both CPU and GPU.  Results in an effective mini-batch of learning."}, {Name: "NThreads", Doc: "number of parallel threads for CPU computation -- 0 = use default"}, {Name: "Run", Doc: "starting run number -- determines the random seed -- runs counts from there -- can do all runs in parallel by launching separate jobs with each run, runs = 1"}, {Name: "NRuns", Doc: "total number of runs to do when running Train"}, {Name: "NEpochs", Doc: "total number of epochs per run"}, {Name: "NTrials", Doc: "total number of trials per epoch.  Should be an even multiple of NData."}, {Name: "PCAInterval", Doc: "how frequently (in epochs) to compute PCA on hidden representations to measure variance?"}}})

var _ = gti.AddType(&gti.Type{Name: "main.LogConfig", IDName: "log-config", Doc: "LogConfig has config parameters related to logging data", Fields: []gti.Field{{Name: "SaveWts", Doc: "if true, save final weights after each run"}, {Name: "Epoch", Doc: "if true, save train epoch log to file, as .epc.tsv typically"}, {Name: "Run", Doc: "if true, save run log to file, as .run.tsv typically"}, {Name: "Trial", Doc: "if true, save train trial log to file, as .trl.tsv typically. May be large."}, {Name: "NetData", Doc: "if true, save network activation etc data from testing trials, for later viewing in netview"}, {Name: "Testing", Doc: "activates testing mode -- records detailed data for Go CI tests (not the same as running test mode on network, via Looper)"}}})

var _ = gti.AddType(&gti.Type{Name: "main.Config", IDName: "config", Doc: "Config is a standard Sim config -- use as a starting point.", Fields: []gti.Field{{Name: "Includes", Doc: "specify include files here, and after configuration, it contains list of include files added"}, {Name: "GUI", Doc: "open the GUI -- does not automatically run -- if false, then runs automatically and quits"}, {Name: "Debug", Doc: "log debugging information"}, {Name: "OpenWts", Doc: "if set, open given weights file at start of training"}, {Name: "Env", Doc: "environment configuration options"}, {Name: "Params", Doc: "parameter related configuration options"}, {Name: "Run", Doc: "sim running related configuration options"}, {Name: "Log", Doc: "data logging related configuration options"}}})
