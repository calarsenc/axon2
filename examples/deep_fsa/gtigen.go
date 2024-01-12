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
		{"UnitsPer", &gti.Field{Name: "UnitsPer", Type: "int", LocalType: "int", Doc: "number of units per localist output unit -- 1 works better than 5 here", Directives: gti.Directives{}, Tag: "def:\"1\""}},
		{"InputNames", &gti.Field{Name: "InputNames", Type: "[]string", LocalType: "[]string", Doc: "] names of input letters", Directives: gti.Directives{}, Tag: "def:\"['B','T','S','X','V','P','E']\""}},
		{"InputNameMap", &gti.Field{Name: "InputNameMap", Type: "map[string]int", LocalType: "map[string]int", Doc: "map of input names -- initialized during Configenv", Directives: gti.Directives{}, Tag: ""}},
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
		{"NData", &gti.Field{Name: "NData", Type: "int", LocalType: "int", Doc: "number of data-parallel items to process in parallel per trial -- works (and is significantly faster) for both CPU and GPU.  Results in an effective mini-batch of learning.  16 learns just as well as 1 -- no diffs.", Directives: gti.Directives{}, Tag: "def:\"16\" min:\"1\""}},
		{"NThreads", &gti.Field{Name: "NThreads", Type: "int", LocalType: "int", Doc: "number of parallel threads for CPU computation -- 0 = use default", Directives: gti.Directives{}, Tag: "def:\"0\""}},
		{"Run", &gti.Field{Name: "Run", Type: "int", LocalType: "int", Doc: "starting run number -- determines the random seed -- runs counts from there -- can do all runs in parallel by launching separate jobs with each run, runs = 1", Directives: gti.Directives{}, Tag: "def:\"0\""}},
		{"NRuns", &gti.Field{Name: "NRuns", Type: "int", LocalType: "int", Doc: "total number of runs to do when running Train", Directives: gti.Directives{}, Tag: "def:\"5\" min:\"1\""}},
		{"NEpochs", &gti.Field{Name: "NEpochs", Type: "int", LocalType: "int", Doc: "total number of epochs per run", Directives: gti.Directives{}, Tag: "def:\"100\""}},
		{"NTrials", &gti.Field{Name: "NTrials", Type: "int", LocalType: "int", Doc: "total number of trials per epoch.  Should be an even multiple of NData.", Directives: gti.Directives{}, Tag: "def:\"196\""}},
		{"PCAInterval", &gti.Field{Name: "PCAInterval", Type: "int", LocalType: "int", Doc: "how frequently (in epochs) to compute PCA on hidden representations to measure variance?", Directives: gti.Directives{}, Tag: "def:\"5\""}},
		{"TestInterval", &gti.Field{Name: "TestInterval", Type: "int", LocalType: "int", Doc: "how often to run through all the test patterns, in terms of training epochs -- can use 0 or -1 for no testing", Directives: gti.Directives{}, Tag: "def:\"-1\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "main.LogConfig",
	ShortName: "main.LogConfig",
	IDName:    "log-config",
	Doc:       "LogConfig has config parameters related to logging data",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"SaveWts", &gti.Field{Name: "SaveWts", Type: "bool", LocalType: "bool", Doc: "if true, save final weights after each run", Directives: gti.Directives{}, Tag: ""}},
		{"Epoch", &gti.Field{Name: "Epoch", Type: "bool", LocalType: "bool", Doc: "if true, save train epoch log to file, as .epc.tsv typically", Directives: gti.Directives{}, Tag: "def:\"true\" nest:\"+\""}},
		{"Run", &gti.Field{Name: "Run", Type: "bool", LocalType: "bool", Doc: "if true, save run log to file, as .run.tsv typically", Directives: gti.Directives{}, Tag: "def:\"true\" nest:\"+\""}},
		{"Trial", &gti.Field{Name: "Trial", Type: "bool", LocalType: "bool", Doc: "if true, save train trial log to file, as .trl.tsv typically. May be large.", Directives: gti.Directives{}, Tag: "def:\"false\" nest:\"+\""}},
		{"TestEpoch", &gti.Field{Name: "TestEpoch", Type: "bool", LocalType: "bool", Doc: "if true, save testing epoch log to file, as .tst_epc.tsv typically.  In general it is better to copy testing items over to the training epoch log and record there.", Directives: gti.Directives{}, Tag: "def:\"false\" nest:\"+\""}},
		{"TestTrial", &gti.Field{Name: "TestTrial", Type: "bool", LocalType: "bool", Doc: "if true, save testing trial log to file, as .tst_trl.tsv typically. May be large.", Directives: gti.Directives{}, Tag: "def:\"false\" nest:\"+\""}},
		{"NetData", &gti.Field{Name: "NetData", Type: "bool", LocalType: "bool", Doc: "if true, save network activation etc data from testing trials, for later viewing in netview", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "main.Config",
	ShortName: "main.Config",
	IDName:    "config",
	Doc:       "Config is a standard Sim config -- use as a starting point.",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Includes", &gti.Field{Name: "Includes", Type: "[]string", LocalType: "[]string", Doc: "specify include files here, and after configuration, it contains list of include files added", Directives: gti.Directives{}, Tag: ""}},
		{"GUI", &gti.Field{Name: "GUI", Type: "bool", LocalType: "bool", Doc: "open the GUI -- does not automatically run -- if false, then runs automatically and quits", Directives: gti.Directives{}, Tag: "def:\"true\""}},
		{"Debug", &gti.Field{Name: "Debug", Type: "bool", LocalType: "bool", Doc: "log debugging information", Directives: gti.Directives{}, Tag: ""}},
		{"Env", &gti.Field{Name: "Env", Type: "github.com/emer/axon/examples/deep_fsa.EnvConfig", LocalType: "EnvConfig", Doc: "environment configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
		{"Params", &gti.Field{Name: "Params", Type: "github.com/emer/axon/examples/deep_fsa.ParamConfig", LocalType: "ParamConfig", Doc: "parameter related configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
		{"Run", &gti.Field{Name: "Run", Type: "github.com/emer/axon/examples/deep_fsa.RunConfig", LocalType: "RunConfig", Doc: "sim running related configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
		{"Log", &gti.Field{Name: "Log", Type: "github.com/emer/axon/examples/deep_fsa.LogConfig", LocalType: "LogConfig", Doc: "data logging related configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
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
		{"Config", &gti.Field{Name: "Config", Type: "github.com/emer/axon/examples/deep_fsa.Config", LocalType: "Config", Doc: "simulation configuration parameters -- set by .toml config file and / or args", Directives: gti.Directives{}, Tag: ""}},
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

var _ = gti.AddType(&gti.Type{
	Name:       "main.FSAEnv",
	ShortName:  "main.FSAEnv",
	IDName:     "fsa-env",
	Doc:        "FSAEnv generates states in a finite state automaton (FSA) which is a\nsimple form of grammar for creating non-deterministic but still\noverall structured sequences.",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Nm", &gti.Field{Name: "Nm", Type: "string", LocalType: "string", Doc: "name of this environment", Directives: gti.Directives{}, Tag: ""}},
		{"Dsc", &gti.Field{Name: "Dsc", Type: "string", LocalType: "string", Doc: "description of this environment", Directives: gti.Directives{}, Tag: ""}},
		{"TMat", &gti.Field{Name: "TMat", Type: "github.com/emer/etable/v2/etensor.Float64", LocalType: "etensor.Float64", Doc: "transition matrix, which is a square NxN tensor with outer dim being current state and inner dim having probability of transitioning to that state", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"Labels", &gti.Field{Name: "Labels", Type: "github.com/emer/etable/v2/etensor.String", LocalType: "etensor.String", Doc: "transition labels, one for each transition cell in TMat matrix", Directives: gti.Directives{}, Tag: ""}},
		{"AState", &gti.Field{Name: "AState", Type: "github.com/emer/emergent/v2/env.CurPrvInt", LocalType: "env.CurPrvInt", Doc: "automaton state within FSA that we're in", Directives: gti.Directives{}, Tag: ""}},
		{"NNext", &gti.Field{Name: "NNext", Type: "github.com/emer/etable/v2/etensor.Int", LocalType: "etensor.Int", Doc: "number of next states in current state output (scalar)", Directives: gti.Directives{}, Tag: ""}},
		{"NextStates", &gti.Field{Name: "NextStates", Type: "github.com/emer/etable/v2/etensor.Int", LocalType: "etensor.Int", Doc: "next states that have non-zero probability, with actual randomly chosen next state at start", Directives: gti.Directives{}, Tag: ""}},
		{"NextLabels", &gti.Field{Name: "NextLabels", Type: "github.com/emer/etable/v2/etensor.String", LocalType: "etensor.String", Doc: "transition labels for next states that have non-zero probability, with actual randomly chosen one for next state at start", Directives: gti.Directives{}, Tag: ""}},
		{"Run", &gti.Field{Name: "Run", Type: "github.com/emer/emergent/v2/env.Ctr", LocalType: "env.Ctr", Doc: "current run of model as provided during Init", Directives: gti.Directives{}, Tag: "view:\"inline\""}},
		{"Epoch", &gti.Field{Name: "Epoch", Type: "github.com/emer/emergent/v2/env.Ctr", LocalType: "env.Ctr", Doc: "number of times through Seq.Max number of sequences", Directives: gti.Directives{}, Tag: "view:\"inline\""}},
		{"Seq", &gti.Field{Name: "Seq", Type: "github.com/emer/emergent/v2/env.Ctr", LocalType: "env.Ctr", Doc: "sequence counter within epoch", Directives: gti.Directives{}, Tag: "view:\"inline\""}},
		{"Tick", &gti.Field{Name: "Tick", Type: "github.com/emer/emergent/v2/env.Ctr", LocalType: "env.Ctr", Doc: "tick counter within sequence", Directives: gti.Directives{}, Tag: "view:\"inline\""}},
		{"Trial", &gti.Field{Name: "Trial", Type: "github.com/emer/emergent/v2/env.Ctr", LocalType: "env.Ctr", Doc: "trial is the step counter within sequence - how many steps taken within current sequence -- it resets to 0 at start of each sequence", Directives: gti.Directives{}, Tag: "view:\"inline\""}},
		{"Rand", &gti.Field{Name: "Rand", Type: "github.com/emer/emergent/v2/erand.SysRand", LocalType: "erand.SysRand", Doc: "random number generator for the env -- all random calls must use this -- set seed here for weight initialization values", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"RndSeed", &gti.Field{Name: "RndSeed", Type: "int64", LocalType: "int64", Doc: "random seed", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})
