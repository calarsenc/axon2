// Code generated by "goki generate ./..."; DO NOT EDIT.

package main

import (
	"goki.dev/gti"
)

var _ = gti.AddType(&gti.Type{Name: "main.OnOff", IDName: "on-off", Doc: "OnOff represents stimulus On / Off timing", Fields: []gti.Field{{Name: "Act", Doc: "is this stimulus active -- use it?"}, {Name: "On", Doc: "when stimulus turns on"}, {Name: "Off", Doc: "when stimulu turns off"}, {Name: "P", Doc: "probability of being active on any given trial"}, {Name: "OnVar", Doc: "variability in onset timing (max number of trials before/after On that it could start)"}, {Name: "OffVar", Doc: "variability in offset timing (max number of trials before/after Off that it could end)"}, {Name: "CurAct", Doc: "current active status based on P probability"}, {Name: "CurOn", Doc: "current on / off values using Var variability"}, {Name: "CurOff", Doc: "current on / off values using Var variability"}}})

var _ = gti.AddType(&gti.Type{Name: "main.CondEnv", IDName: "cond-env", Doc: "CondEnv simulates an n-armed bandit, where each of n inputs is associated with\na specific probability of reward.", Fields: []gti.Field{{Name: "Nm", Doc: "name of this environment"}, {Name: "Dsc", Doc: "description of this environment"}, {Name: "TotTime", Doc: "total time for trial"}, {Name: "CSA", Doc: "Conditioned stimulus A (e.g., Tone)"}, {Name: "CSB", Doc: "Conditioned stimulus B (e.g., Light)"}, {Name: "CSC", Doc: "Conditioned stimulus C"}, {Name: "US", Doc: "Unconditioned stimulus -- reward"}, {Name: "RewVal", Doc: "value for reward"}, {Name: "NoRewVal", Doc: "value for non-reward"}, {Name: "Input", Doc: "one-hot input representation of current option"}, {Name: "Reward", Doc: "single reward value"}, {Name: "HasRew", Doc: "true if a US reward value was set"}, {Name: "Run", Doc: "current run of model as provided during Init"}, {Name: "Epoch", Doc: "number of times through Seq.Max number of sequences"}, {Name: "Trial", Doc: "one trial is a pass through all TotTime Events"}, {Name: "Event", Doc: "event is one time step within Trial -- e.g., CS turning on, etc"}}})

var _ = gti.AddType(&gti.Type{Name: "main.EnvConfig", IDName: "env-config", Doc: "EnvConfig has config params for environment\nnote: only adding fields for key Env params that matter for both Network and Env\nother params are set via the Env map data mechanism.", Fields: []gti.Field{{Name: "Env", Doc: "env parameters -- can set any field/subfield on Env struct, using standard TOML formatting"}}})

var _ = gti.AddType(&gti.Type{Name: "main.ParamConfig", IDName: "param-config", Doc: "ParamConfig has config parameters related to sim params", Fields: []gti.Field{{Name: "Network", Doc: "network parameters"}, {Name: "Sheet", Doc: "Extra Param Sheet name(s) to use (space separated if multiple) -- must be valid name as listed in compiled-in params or loaded params"}, {Name: "Tag", Doc: "extra tag to add to file names and logs saved from this run"}, {Name: "Note", Doc: "user note -- describe the run params etc -- like a git commit message for the run"}, {Name: "File", Doc: "Name of the JSON file to input saved parameters from."}, {Name: "SaveAll", Doc: "Save a snapshot of all current param and config settings in a directory named params_<datestamp> (or _good if Good is true), then quit -- useful for comparing to later changes and seeing multiple views of current params"}, {Name: "Good", Doc: "for SaveAll, save to params_good for a known good params state.  This can be done prior to making a new release after all tests are passing -- add results to git to provide a full diff record of all params over time."}}})

var _ = gti.AddType(&gti.Type{Name: "main.RunConfig", IDName: "run-config", Doc: "RunConfig has config parameters related to running the sim", Fields: []gti.Field{{Name: "GPU", Doc: "use the GPU for computation -- only for testing in this model -- not faster"}, {Name: "NThreads", Doc: "number of parallel threads for CPU computation -- 0 = use default"}, {Name: "Run", Doc: "starting run number -- determines the random seed -- runs counts from there -- can do all runs in parallel by launching separate jobs with each run, runs = 1"}, {Name: "NRuns", Doc: "total number of runs to do when running Train"}, {Name: "NEpochs", Doc: "total number of epochs per run"}, {Name: "NTrials", Doc: "total number of trials per epoch -- should be number of ticks in env."}}})

var _ = gti.AddType(&gti.Type{Name: "main.LogConfig", IDName: "log-config", Doc: "LogConfig has config parameters related to logging data", Fields: []gti.Field{{Name: "AggStats", Doc: "] stats to aggregate at higher levels"}, {Name: "SaveWts", Doc: "if true, save final weights after each run"}, {Name: "Epoch", Doc: "if true, save train epoch log to file, as .epc.tsv typically"}, {Name: "Run", Doc: "if true, save run log to file, as .run.tsv typically"}, {Name: "Trial", Doc: "if true, save train trial log to file, as .trl.tsv typically. May be large."}, {Name: "NetData", Doc: "if true, save network activation etc data from testing trials, for later viewing in netview"}}})

var _ = gti.AddType(&gti.Type{Name: "main.Config", IDName: "config", Doc: "Config is a standard Sim config -- use as a starting point.", Fields: []gti.Field{{Name: "Includes", Doc: "specify include files here, and after configuration, it contains list of include files added"}, {Name: "RW", Doc: "if true, use Rescorla-Wagner -- set in code or rebuild network"}, {Name: "GUI", Doc: "open the GUI -- does not automatically run -- if false, then runs automatically and quits"}, {Name: "Debug", Doc: "log debugging information"}, {Name: "Env", Doc: "environment configuration options"}, {Name: "Params", Doc: "parameter related configuration options"}, {Name: "Run", Doc: "sim running related configuration options"}, {Name: "Log", Doc: "data logging related configuration options"}}})

var _ = gti.AddType(&gti.Type{Name: "main.Sim", IDName: "sim", Doc: "Sim encapsulates the entire simulation model, and we define all the\nfunctionality as methods on this struct.  This structure keeps all relevant\nstate information organized and available without having to pass everything around\nas arguments to methods, and provides the core GUI interface (note the view tags\nfor the fields which provide hints to how things should be displayed).", Methods: []gti.Method{{Name: "VmRun", Doc: "VmRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "VGRun", Doc: "VGRun runs the V-G equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "SGRun", Doc: "SGRun runs the spike-g equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "VmRun", Doc: "VmRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "Run", Doc: "Run runs the equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "CaRun", Doc: "CaRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "CamRun", Doc: "CamRun plots the equation as a function of Ca", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "VmRun", Doc: "VmRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}}, Fields: []gti.Field{{Name: "Config", Doc: "simulation configuration parameters -- set by .toml config file and / or args"}, {Name: "Net", Doc: "the network -- click to view / edit parameters for layers, prjns, etc"}, {Name: "Params", Doc: "all parameter management"}, {Name: "Loops", Doc: "contains looper control loops for running sim"}, {Name: "Stats", Doc: "contains computed statistic values"}, {Name: "Logs", Doc: "Contains all the logs and information about the logs.'"}, {Name: "Envs", Doc: "Environments"}, {Name: "Context", Doc: "axon timing parameters and state"}, {Name: "ViewUpdt", Doc: "netview update parameters"}, {Name: "GUI", Doc: "manages all the gui elements"}, {Name: "RndSeeds", Doc: "a list of random seeds to use for each run"}}})
