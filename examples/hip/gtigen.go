// Code generated by "core generate ./..."; DO NOT EDIT.

package main

import (
	"cogentcore.org/core/gti"
)

var _ = gti.AddType(&gti.Type{Name: "main.ParamConfig", IDName: "param-config", Doc: "ParamConfig has config parameters related to sim params", Fields: []gti.Field{{Name: "Network", Doc: "network parameters"}, {Name: "Sheet", Doc: "Extra Param Sheet name(s) to use (space separated if multiple) -- must be valid name as listed in compiled-in params or loaded params"}, {Name: "Tag", Doc: "extra tag to add to file names and logs saved from this run"}, {Name: "Note", Doc: "user note -- describe the run params etc -- like a git commit message for the run"}, {Name: "File", Doc: "Name of the JSON file to input saved parameters from."}, {Name: "SaveAll", Doc: "Save a snapshot of all current param and config settings in a directory named params_<datestamp> (or _good if Good is true), then quit -- useful for comparing to later changes and seeing multiple views of current params"}, {Name: "Good", Doc: "for SaveAll, save to params_good for a known good params state.  This can be done prior to making a new release after all tests are passing -- add results to git to provide a full diff record of all params over time."}}})

var _ = gti.AddType(&gti.Type{Name: "main.RunConfig", IDName: "run-config", Doc: "RunConfig has config parameters related to running the sim", Fields: []gti.Field{{Name: "StopMem", Doc: "mem % correct level (proportion) above which training on current list stops (switch from AB to AC or stop on AC)"}, {Name: "GPU", Doc: "use the GPU for computation -- generally faster even for small models if NData ~16"}, {Name: "NThreads", Doc: "number of parallel threads for CPU computation -- 0 = use default"}, {Name: "Run", Doc: "starting run number -- determines the random seed -- runs counts from there -- can do all runs in parallel by launching separate jobs with each run, runs = 1"}, {Name: "Runs", Doc: "total number of runs to do when running Train"}, {Name: "Epochs", Doc: "total number of epochs per run"}, {Name: "NTrials", Doc: "total number of trials per epoch.  Should be an even multiple of NData."}, {Name: "NData", Doc: "number of data-parallel items to process in parallel per trial -- works (and is significantly faster) for both CPU and GPU.  Results in an effective mini-batch of learning."}, {Name: "TestInterval", Doc: "how often to run through all the test patterns, in terms of training epochs -- can use 0 or -1 for no testing"}}})

var _ = gti.AddType(&gti.Type{Name: "main.LogConfig", IDName: "log-config", Doc: "LogConfig has config parameters related to logging data", Fields: []gti.Field{{Name: "SaveWts", Doc: "if true, save final weights after each run"}, {Name: "Epoch", Doc: "if true, save train epoch log to file, as .epc.tsv typically"}, {Name: "Run", Doc: "if true, save run log to file, as .run.tsv typically"}, {Name: "Trial", Doc: "if true, save train trial log to file, as .trl.tsv typically. May be large."}, {Name: "TestEpoch", Doc: "if true, save testing epoch log to file, as .tst_epc.tsv typically.  In general it is better to copy testing items over to the training epoch log and record there."}, {Name: "TestTrial", Doc: "if true, save testing trial log to file, as .tst_trl.tsv typically. May be large."}, {Name: "NetData", Doc: "if true, save network activation etc data from testing trials, for later viewing in netview"}}})

var _ = gti.AddType(&gti.Type{Name: "main.PatConfig", IDName: "pat-config", Doc: "PatConfig have the pattern parameters", Fields: []gti.Field{{Name: "MinDiffPct", Doc: "minimum difference between item random patterns, as a proportion (0-1) of total active"}, {Name: "DriftCtxt", Doc: "use drifting context representations -- otherwise does bit flips from prototype"}, {Name: "CtxtFlipPct", Doc: "proportion (0-1) of active bits to flip for each context pattern, relative to a prototype, for non-drifting"}, {Name: "DriftPct", Doc: "percentage of active bits that drift, per step, for drifting context"}}})

var _ = gti.AddType(&gti.Type{Name: "main.ModConfig", IDName: "mod-config", Fields: []gti.Field{{Name: "InToEc2PCon", Doc: "percent connectivity from Input to EC2"}, {Name: "ECPctAct", Doc: "percent activation in EC pool, used in patgen for input generation\npercent activation in EC pool, used in patgen for input generation"}, {Name: "MemThr", Doc: "memory threshold"}}})

var _ = gti.AddType(&gti.Type{Name: "main.Config", IDName: "config", Doc: "Config is a standard Sim config -- use as a starting point.", Fields: []gti.Field{{Name: "Includes", Doc: "specify include files here, and after configuration, it contains list of include files added"}, {Name: "GUI", Doc: "open the GUI -- does not automatically run -- if false, then runs automatically and quits"}, {Name: "Debug", Doc: "log debugging information"}, {Name: "Mod", Doc: "misc model parameters"}, {Name: "Hip", Doc: "Hippocampus sizing parameters"}, {Name: "Pat", Doc: "parameters for the input patterns"}, {Name: "Params", Doc: "parameter related configuration options"}, {Name: "Run", Doc: "sim running related configuration options"}, {Name: "Log", Doc: "data logging related configuration options"}}})

var _ = gti.AddType(&gti.Type{Name: "main.Sim", IDName: "sim", Doc: "Sim encapsulates the entire simulation model, and we define all the\nfunctionality as methods on this struct.  This structure keeps all relevant\nstate information organized and available without having to pass everything around\nas arguments to methods, and provides the core GUI interface (note the view tags\nfor the fields which provide hints to how things should be displayed).", Methods: []gti.Method{{Name: "VmRun", Doc: "VmRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "VGRun", Doc: "VGRun runs the V-G equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "SGRun", Doc: "SGRun runs the spike-g equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "VmRun", Doc: "VmRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "Run", Doc: "Run runs the equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "CaRun", Doc: "CaRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "CamRun", Doc: "CamRun plots the equation as a function of Ca", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "VmRun", Doc: "VmRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}}, Fields: []gti.Field{{Name: "Config", Doc: "simulation configuration parameters -- set by .toml config file and / or args"}, {Name: "Net", Doc: "the network -- click to view / edit parameters for layers, prjns, etc"}, {Name: "Params", Doc: "all parameter management"}, {Name: "Loops", Doc: "contains looper control loops for running sim"}, {Name: "Stats", Doc: "contains computed statistic values"}, {Name: "Logs", Doc: "Contains all the logs and information about the logs.'"}, {Name: "PretrainMode", Doc: "if true, run in pretrain mode"}, {Name: "PoolVocab", Doc: "pool patterns vocabulary"}, {Name: "TrainAB", Doc: "AB training patterns to use"}, {Name: "TrainAC", Doc: "AC training patterns to use"}, {Name: "TestAB", Doc: "AB testing patterns to use"}, {Name: "TestAC", Doc: "AC testing patterns to use"}, {Name: "PreTrainLure", Doc: "Lure pretrain patterns to use"}, {Name: "TestLure", Doc: "Lure testing patterns to use"}, {Name: "TrainAll", Doc: "all training patterns -- for pretrain"}, {Name: "TestABAC", Doc: "TestAB + TestAC"}, {Name: "Envs", Doc: "Environments"}, {Name: "Context", Doc: "axon timing parameters and state"}, {Name: "ViewUpdt", Doc: "netview update parameters"}, {Name: "GUI", Doc: "manages all the gui elements"}, {Name: "RndSeeds", Doc: "a list of random seeds to use for each run"}}})
