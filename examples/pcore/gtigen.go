// Code generated by "goki generate ./..."; DO NOT EDIT.

package main

import (
	"goki.dev/gti"
	"goki.dev/ordmap"
)

var _ = gti.AddType(&gti.Type{
	Name:       "main.EnvConfig",
	ShortName:  "main.EnvConfig",
	IDName:     "env-config",
	Doc:        "EnvConfig has config params for environment\nnote: only adding fields for key Env params that matter for both Network and Env\nother params are set via the Env map data mechanism.",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Env", &gti.Field{Name: "Env", Type: "map[string]any", LocalType: "map[string]any", Doc: "env parameters -- can set any field/subfield on Env struct, using standard TOML formatting", Directives: gti.Directives{}, Tag: ""}},
		{"ZeroTest", &gti.Field{Name: "ZeroTest", Type: "bool", LocalType: "bool", Doc: "test with no ACC activity at all -- params need to prevent gating in this situation too", Directives: gti.Directives{}, Tag: ""}},
		{"NPools", &gti.Field{Name: "NPools", Type: "int", LocalType: "int", Doc: "number of pools in BG / PFC -- if > 1 then does selection among options presented in parallel (not yet supported / tested) -- otherwise does go / no on a single optoin (default)", Directives: gti.Directives{}, Tag: "def:\"1\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "main.ParamConfig",
	ShortName:  "main.ParamConfig",
	IDName:     "param-config",
	Doc:        "ParamConfig has config parameters related to sim params",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Network", &gti.Field{Name: "Network", Type: "map[string]any", LocalType: "map[string]any", Doc: "network parameters", Directives: gti.Directives{}, Tag: ""}},
		{"Sheet", &gti.Field{Name: "Sheet", Type: "string", LocalType: "string", Doc: "Extra Param Sheet name(s) to use (space separated if multiple) -- must be valid name as listed in compiled-in params or loaded params", Directives: gti.Directives{}, Tag: ""}},
		{"Tag", &gti.Field{Name: "Tag", Type: "string", LocalType: "string", Doc: "extra tag to add to file names and logs saved from this run", Directives: gti.Directives{}, Tag: ""}},
		{"Note", &gti.Field{Name: "Note", Type: "string", LocalType: "string", Doc: "user note -- describe the run params etc -- like a git commit message for the run", Directives: gti.Directives{}, Tag: ""}},
		{"File", &gti.Field{Name: "File", Type: "string", LocalType: "string", Doc: "Name of the JSON file to input saved parameters from.", Directives: gti.Directives{}, Tag: "nest:\"+\""}},
		{"SaveAll", &gti.Field{Name: "SaveAll", Type: "bool", LocalType: "bool", Doc: "Save a snapshot of all current param and config settings in a directory named params_<datestamp> (or _good if Good is true), then quit -- useful for comparing to later changes and seeing multiple views of current params", Directives: gti.Directives{}, Tag: "nest:\"+\""}},
		{"Good", &gti.Field{Name: "Good", Type: "bool", LocalType: "bool", Doc: "for SaveAll, save to params_good for a known good params state.  This can be done prior to making a new release after all tests are passing -- add results to git to provide a full diff record of all params over time.", Directives: gti.Directives{}, Tag: "nest:\"+\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "main.RunConfig",
	ShortName:  "main.RunConfig",
	IDName:     "run-config",
	Doc:        "RunConfig has config parameters related to running the sim",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"GPU", &gti.Field{Name: "GPU", Type: "bool", LocalType: "bool", Doc: "use the GPU for computation -- generally faster even for small models if NData ~16", Directives: gti.Directives{}, Tag: "def:\"true\""}},
		{"NData", &gti.Field{Name: "NData", Type: "int", LocalType: "int", Doc: "number of data-parallel items to process in parallel per trial -- works (and is significantly faster) for both CPU and GPU.  Results in an effective mini-batch of learning.", Directives: gti.Directives{}, Tag: "def:\"16\" min:\"1\""}},
		{"NThreads", &gti.Field{Name: "NThreads", Type: "int", LocalType: "int", Doc: "number of parallel threads for CPU computation -- 0 = use default", Directives: gti.Directives{}, Tag: "def:\"0\""}},
		{"Run", &gti.Field{Name: "Run", Type: "int", LocalType: "int", Doc: "starting run number -- determines the random seed -- runs counts from there -- can do all runs in parallel by launching separate jobs with each run, runs = 1", Directives: gti.Directives{}, Tag: "def:\"0\""}},
		{"NRuns", &gti.Field{Name: "NRuns", Type: "int", LocalType: "int", Doc: "total number of runs to do when running Train", Directives: gti.Directives{}, Tag: "def:\"1\" min:\"1\""}},
		{"NEpochs", &gti.Field{Name: "NEpochs", Type: "int", LocalType: "int", Doc: "total number of epochs per run", Directives: gti.Directives{}, Tag: "def:\"30\""}},
		{"NTrials", &gti.Field{Name: "NTrials", Type: "int", LocalType: "int", Doc: "total number of trials per epoch.  Should be an even multiple of NData.", Directives: gti.Directives{}, Tag: "def:\"128\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "main.LogConfig",
	ShortName:  "main.LogConfig",
	IDName:     "log-config",
	Doc:        "LogConfig has config parameters related to logging data",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"SaveWts", &gti.Field{Name: "SaveWts", Type: "bool", LocalType: "bool", Doc: "if true, save final weights after each run", Directives: gti.Directives{}, Tag: ""}},
		{"Epoch", &gti.Field{Name: "Epoch", Type: "bool", LocalType: "bool", Doc: "if true, save train epoch log to file, as .epc.tsv typically", Directives: gti.Directives{}, Tag: "def:\"true\" nest:\"+\""}},
		{"Run", &gti.Field{Name: "Run", Type: "bool", LocalType: "bool", Doc: "if true, save run log to file, as .run.tsv typically", Directives: gti.Directives{}, Tag: "def:\"true\" nest:\"+\""}},
		{"Trial", &gti.Field{Name: "Trial", Type: "bool", LocalType: "bool", Doc: "if true, save train trial log to file, as .trl.tsv typically. May be large.", Directives: gti.Directives{}, Tag: "def:\"false\" nest:\"+\""}},
		{"TestEpoch", &gti.Field{Name: "TestEpoch", Type: "bool", LocalType: "bool", Doc: "if true, save testing epoch log to file, as .tst_epc.tsv typically.  In general it is better to copy testing items over to the training epoch log and record there.", Directives: gti.Directives{}, Tag: "def:\"false\" nest:\"+\""}},
		{"TestTrial", &gti.Field{Name: "TestTrial", Type: "bool", LocalType: "bool", Doc: "if true, save testing trial log to file, as .tst_trl.tsv typically. May be large.", Directives: gti.Directives{}, Tag: "def:\"false\" nest:\"+\""}},
		{"NetData", &gti.Field{Name: "NetData", Type: "bool", LocalType: "bool", Doc: "if true, save network activation etc data from testing trials, for later viewing in netview", Directives: gti.Directives{}, Tag: ""}},
		{"Testing", &gti.Field{Name: "Testing", Type: "bool", LocalType: "bool", Doc: "activates testing mode -- records detailed data for Go CI tests (not the same as running test mode on network, via Looper)", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "main.Config",
	ShortName:  "main.Config",
	IDName:     "config",
	Doc:        "Config is a standard Sim config -- use as a starting point.",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Includes", &gti.Field{Name: "Includes", Type: "[]string", LocalType: "[]string", Doc: "specify include files here, and after configuration, it contains list of include files added", Directives: gti.Directives{}, Tag: ""}},
		{"GUI", &gti.Field{Name: "GUI", Type: "bool", LocalType: "bool", Doc: "open the GUI -- does not automatically run -- if false, then runs automatically and quits", Directives: gti.Directives{}, Tag: "def:\"true\""}},
		{"Debug", &gti.Field{Name: "Debug", Type: "bool", LocalType: "bool", Doc: "log debugging information", Directives: gti.Directives{}, Tag: ""}},
		{"Env", &gti.Field{Name: "Env", Type: "github.com/emer/axon/examples/pcore.EnvConfig", LocalType: "EnvConfig", Doc: "environment configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
		{"Params", &gti.Field{Name: "Params", Type: "github.com/emer/axon/examples/pcore.ParamConfig", LocalType: "ParamConfig", Doc: "parameter related configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
		{"Run", &gti.Field{Name: "Run", Type: "github.com/emer/axon/examples/pcore.RunConfig", LocalType: "RunConfig", Doc: "sim running related configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
		{"Log", &gti.Field{Name: "Log", Type: "github.com/emer/axon/examples/pcore.LogConfig", LocalType: "LogConfig", Doc: "data logging related configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "main.GoNoEnv",
	ShortName:  "main.GoNoEnv",
	IDName:     "go-no-env",
	Doc:        "GoNoEnv implements simple Go vs. NoGo input patterns to test BG learning.",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Nm", &gti.Field{Name: "Nm", Type: "string", LocalType: "string", Doc: "name of environment -- Train or Test", Directives: gti.Directives{}, Tag: ""}},
		{"Mode", &gti.Field{Name: "Mode", Type: "github.com/emer/emergent/v2/etime.Modes", LocalType: "etime.Modes", Doc: "training or testing env?", Directives: gti.Directives{}, Tag: ""}},
		{"Trial", &gti.Field{Name: "Trial", Type: "github.com/emer/emergent/v2/env.Ctr", LocalType: "env.Ctr", Doc: "trial counter -- set by caller for testing", Directives: gti.Directives{}, Tag: ""}},
		{"ACCPos", &gti.Field{Name: "ACCPos", Type: "float32", LocalType: "float32", Doc: "activation of ACC positive valence -- drives go", Directives: gti.Directives{}, Tag: ""}},
		{"ACCNeg", &gti.Field{Name: "ACCNeg", Type: "float32", LocalType: "float32", Doc: "activation of ACC neg valence -- drives nogo", Directives: gti.Directives{}, Tag: ""}},
		{"PosNegThr", &gti.Field{Name: "PosNegThr", Type: "float32", LocalType: "float32", Doc: "threshold on diff between ACCPos - ACCNeg for counting as a Go trial", Directives: gti.Directives{}, Tag: ""}},
		{"ManualVals", &gti.Field{Name: "ManualVals", Type: "bool", LocalType: "bool", Doc: "ACCPos and Neg are set manually -- do not generate random vals for training or auto-increment ACCPos / Neg values during test", Directives: gti.Directives{}, Tag: ""}},
		{"TestInc", &gti.Field{Name: "TestInc", Type: "float32", LocalType: "float32", Doc: "increment in testing activation for test all", Directives: gti.Directives{}, Tag: ""}},
		{"TestReps", &gti.Field{Name: "TestReps", Type: "int", LocalType: "int", Doc: "number of repetitions per testing level", Directives: gti.Directives{}, Tag: ""}},
		{"NPools", &gti.Field{Name: "NPools", Type: "int", LocalType: "int", Doc: "number of pools for representing multiple different options to be evaluated in parallel, vs. 1 pool with a simple go nogo overall choice -- currently tested / configured for the 1 pool case", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"ACCPosInc", &gti.Field{Name: "ACCPosInc", Type: "float32", LocalType: "float32", Doc: "for case with multiple pools evaluated in parallel (not currently used), this is the across-pools multiplier in activation of ACC positive valence -- e.g., .9 daecrements subsequent units by 10%", Directives: gti.Directives{}, Tag: ""}},
		{"ACCNegInc", &gti.Field{Name: "ACCNegInc", Type: "float32", LocalType: "float32", Doc: "for case with multiple pools evaluated in parallel (not currently used), this is the across-pools multiplier in activation of ACC neg valence, e.g., 1.1 increments subsequent units by 10%", Directives: gti.Directives{}, Tag: ""}},
		{"NUnitsY", &gti.Field{Name: "NUnitsY", Type: "int", LocalType: "int", Doc: "number of units within each pool, Y", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"NUnitsX", &gti.Field{Name: "NUnitsX", Type: "int", LocalType: "int", Doc: "number of units within each pool, X", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"NUnits", &gti.Field{Name: "NUnits", Type: "int", LocalType: "int", Doc: "total number of units within each pool", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"InN", &gti.Field{Name: "InN", Type: "int", LocalType: "int", Doc: "number of different values for PFC to learn in input layer -- gives PFC network something to do", Directives: gti.Directives{}, Tag: ""}},
		{"PopCode", &gti.Field{Name: "PopCode", Type: "github.com/emer/emergent/v2/popcode.OneD", LocalType: "popcode.OneD", Doc: "pop code the values in ACCPos and Neg", Directives: gti.Directives{}, Tag: ""}},
		{"Rand", &gti.Field{Name: "Rand", Type: "github.com/emer/emergent/v2/erand.SysRand", LocalType: "erand.SysRand", Doc: "random number generator for the env -- all random calls must use this", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"RndSeed", &gti.Field{Name: "RndSeed", Type: "int64", LocalType: "int64", Doc: "random seed", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"States", &gti.Field{Name: "States", Type: "map[string]*goki.dev/etable/v2/etensor.Float32", LocalType: "map[string]*etensor.Float32", Doc: "named states: ACCPos, ACCNeg", Directives: gti.Directives{}, Tag: ""}},
		{"Should", &gti.Field{Name: "Should", Type: "bool", LocalType: "bool", Doc: "true if Pos - Neg > Thr", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"Gated", &gti.Field{Name: "Gated", Type: "bool", LocalType: "bool", Doc: "true if model gated on this trial", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"Match", &gti.Field{Name: "Match", Type: "bool", LocalType: "bool", Doc: "true if gated == should", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"Rew", &gti.Field{Name: "Rew", Type: "float32", LocalType: "float32", Doc: "reward based on match between Should vs. Gated", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"InCtr", &gti.Field{Name: "InCtr", Type: "int", LocalType: "int", Doc: "input counter -- gives PFC network something to do", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "main.Sim",
	ShortName:  "main.Sim",
	IDName:     "sim",
	Doc:        "Sim encapsulates the entire simulation model, and we define all the\nfunctionality as methods on this struct.  This structure keeps all relevant\nstate information organized and available without having to pass everything around\nas arguments to methods, and provides the core GUI interface (note the view tags\nfor the fields which provide hints to how things should be displayed).",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Config", &gti.Field{Name: "Config", Type: "github.com/emer/axon/examples/pcore.Config", LocalType: "Config", Doc: "simulation configuration parameters -- set by .toml config file and / or args", Directives: gti.Directives{}, Tag: ""}},
		{"Net", &gti.Field{Name: "Net", Type: "*invalid type", LocalType: "*axon.Network", Doc: "the network -- click to view / edit parameters for layers, prjns, etc", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"Params", &gti.Field{Name: "Params", Type: "github.com/emer/emergent/v2/emer.NetParams", LocalType: "emer.NetParams", Doc: "all parameter management", Directives: gti.Directives{}, Tag: "view:\"inline\""}},
		{"Loops", &gti.Field{Name: "Loops", Type: "*github.com/emer/emergent/v2/looper.Manager", LocalType: "*looper.Manager", Doc: "contains looper control loops for running sim", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"Stats", &gti.Field{Name: "Stats", Type: "github.com/emer/emergent/v2/estats.Stats", LocalType: "estats.Stats", Doc: "contains computed statistic values", Directives: gti.Directives{}, Tag: ""}},
		{"Logs", &gti.Field{Name: "Logs", Type: "github.com/emer/emergent/v2/elog.Logs", LocalType: "elog.Logs", Doc: "Contains all the logs and information about the logs.'", Directives: gti.Directives{}, Tag: ""}},
		{"Envs", &gti.Field{Name: "Envs", Type: "github.com/emer/emergent/v2/env.Envs", LocalType: "env.Envs", Doc: "Environments", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"Context", &gti.Field{Name: "Context", Type: "invalid type", LocalType: "axon.Context", Doc: "axon timing parameters and state", Directives: gti.Directives{}, Tag: ""}},
		{"ViewUpdt", &gti.Field{Name: "ViewUpdt", Type: "github.com/emer/emergent/v2/netview.ViewUpdt", LocalType: "netview.ViewUpdt", Doc: "netview update parameters", Directives: gti.Directives{}, Tag: "view:\"inline\""}},
		{"GUI", &gti.Field{Name: "GUI", Type: "github.com/emer/emergent/v2/egui.GUI", LocalType: "egui.GUI", Doc: "manages all the gui elements", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"RndSeeds", &gti.Field{Name: "RndSeeds", Type: "github.com/emer/emergent/v2/erand.Seeds", LocalType: "erand.Seeds", Doc: "a list of random seeds to use for each run", Directives: gti.Directives{}, Tag: "view:\"-\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})
