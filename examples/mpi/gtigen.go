// Code generated by "goki generate ./..."; DO NOT EDIT.

package main

import (
	"goki.dev/gti"
	"goki.dev/ordmap"
)

var _ = gti.AddType(&gti.Type{
	Name:       "main.ParamConfig",
	ShortName:  "main.ParamConfig",
	IDName:     "param-config",
	Doc:        "ParamConfig has config parameters related to sim params",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Network", &gti.Field{Name: "Network", Type: "map[string]any", LocalType: "map[string]any", Doc: "network parameters", Directives: gti.Directives{}, Tag: ""}},
		{"Hidden1Size", &gti.Field{Name: "Hidden1Size", Type: "github.com/emer/emergent/v2/evec.Vec2i", LocalType: "evec.Vec2i", Doc: "size of hidden layer -- can use emer.LaySize for 4D layers", Directives: gti.Directives{}, Tag: "def:\"{'X':10,'Y':10}\" nest:\"+\""}},
		{"Hidden2Size", &gti.Field{Name: "Hidden2Size", Type: "github.com/emer/emergent/v2/evec.Vec2i", LocalType: "evec.Vec2i", Doc: "size of hidden layer -- can use emer.LaySize for 4D layers", Directives: gti.Directives{}, Tag: "def:\"{'X':10,'Y':10}\" nest:\"+\""}},
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
		{"MPI", &gti.Field{Name: "MPI", Type: "bool", LocalType: "bool", Doc: "use MPI message passing interface for data parallel computation between nodes running identical copies of the same sim, sharing DWt changes", Directives: gti.Directives{}, Tag: ""}},
		{"GPU", &gti.Field{Name: "GPU", Type: "bool", LocalType: "bool", Doc: "use the GPU for computation -- generally faster even for small models if NData ~16", Directives: gti.Directives{}, Tag: "def:\"false\""}},
		{"NData", &gti.Field{Name: "NData", Type: "int", LocalType: "int", Doc: "number of data-parallel items to process in parallel per trial -- works (and is significantly faster) for both CPU and GPU.  Results in an effective mini-batch of learning.", Directives: gti.Directives{}, Tag: "def:\"16\" min:\"1\""}},
		{"NThreads", &gti.Field{Name: "NThreads", Type: "int", LocalType: "int", Doc: "number of parallel threads for CPU computation -- 0 = use default", Directives: gti.Directives{}, Tag: "def:\"0\""}},
		{"Run", &gti.Field{Name: "Run", Type: "int", LocalType: "int", Doc: "starting run number -- determines the random seed -- runs counts from there -- can do all runs in parallel by launching separate jobs with each run, runs = 1", Directives: gti.Directives{}, Tag: "def:\"0\""}},
		{"NRuns", &gti.Field{Name: "NRuns", Type: "int", LocalType: "int", Doc: "total number of runs to do when running Train", Directives: gti.Directives{}, Tag: "def:\"5\" min:\"1\""}},
		{"NEpochs", &gti.Field{Name: "NEpochs", Type: "int", LocalType: "int", Doc: "total number of epochs per run", Directives: gti.Directives{}, Tag: "def:\"100\""}},
		{"NZero", &gti.Field{Name: "NZero", Type: "int", LocalType: "int", Doc: "stop run after this number of perfect, zero-error epochs", Directives: gti.Directives{}, Tag: "def:\"2\""}},
		{"NTrials", &gti.Field{Name: "NTrials", Type: "int", LocalType: "int", Doc: "total number of trials per epoch.  Should be an even multiple of NData.", Directives: gti.Directives{}, Tag: "def:\"32\""}},
		{"TestInterval", &gti.Field{Name: "TestInterval", Type: "int", LocalType: "int", Doc: "how often to run through all the test patterns, in terms of training epochs -- can use 0 or -1 for no testing", Directives: gti.Directives{}, Tag: "def:\"5\""}},
		{"PCAInterval", &gti.Field{Name: "PCAInterval", Type: "int", LocalType: "int", Doc: "how frequently (in epochs) to compute PCA on hidden representations to measure variance?", Directives: gti.Directives{}, Tag: "def:\"5\""}},
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
		{"Params", &gti.Field{Name: "Params", Type: "github.com/emer/axon/examples/mpi.ParamConfig", LocalType: "ParamConfig", Doc: "parameter related configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
		{"Run", &gti.Field{Name: "Run", Type: "github.com/emer/axon/examples/mpi.RunConfig", LocalType: "RunConfig", Doc: "sim running related configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
		{"Log", &gti.Field{Name: "Log", Type: "github.com/emer/axon/examples/mpi.LogConfig", LocalType: "LogConfig", Doc: "data logging related configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
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
		{"Config", &gti.Field{Name: "Config", Type: "github.com/emer/axon/examples/mpi.Config", LocalType: "Config", Doc: "simulation configuration parameters -- set by .toml config file and / or args", Directives: gti.Directives{}, Tag: ""}},
		{"Net", &gti.Field{Name: "Net", Type: "*invalid type", LocalType: "*axon.Network", Doc: "the network -- click to view / edit parameters for layers, prjns, etc", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"Params", &gti.Field{Name: "Params", Type: "github.com/emer/emergent/v2/emer.NetParams", LocalType: "emer.NetParams", Doc: "network parameter management", Directives: gti.Directives{}, Tag: "view:\"inline\""}},
		{"Loops", &gti.Field{Name: "Loops", Type: "*github.com/emer/emergent/v2/looper.Manager", LocalType: "*looper.Manager", Doc: "contains looper control loops for running sim", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"Stats", &gti.Field{Name: "Stats", Type: "github.com/emer/emergent/v2/estats.Stats", LocalType: "estats.Stats", Doc: "contains computed statistic values", Directives: gti.Directives{}, Tag: ""}},
		{"Logs", &gti.Field{Name: "Logs", Type: "github.com/emer/emergent/v2/elog.Logs", LocalType: "elog.Logs", Doc: "Contains all the logs and information about the logs.'", Directives: gti.Directives{}, Tag: ""}},
		{"Pats", &gti.Field{Name: "Pats", Type: "*goki.dev/etable/v2/etable.Table", LocalType: "*etable.Table", Doc: "the training patterns to use", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"Envs", &gti.Field{Name: "Envs", Type: "github.com/emer/emergent/v2/env.Envs", LocalType: "env.Envs", Doc: "Environments", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"Context", &gti.Field{Name: "Context", Type: "invalid type", LocalType: "axon.Context", Doc: "axon timing parameters and state", Directives: gti.Directives{}, Tag: ""}},
		{"ViewUpdt", &gti.Field{Name: "ViewUpdt", Type: "github.com/emer/emergent/v2/netview.ViewUpdt", LocalType: "netview.ViewUpdt", Doc: "netview update parameters", Directives: gti.Directives{}, Tag: "view:\"inline\""}},
		{"GUI", &gti.Field{Name: "GUI", Type: "github.com/emer/emergent/v2/egui.GUI", LocalType: "egui.GUI", Doc: "manages all the gui elements", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"RndSeeds", &gti.Field{Name: "RndSeeds", Type: "github.com/emer/emergent/v2/erand.Seeds", LocalType: "erand.Seeds", Doc: "a list of random seeds to use for each run", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"Comm", &gti.Field{Name: "Comm", Type: "*github.com/emer/empi/v2/mpi.Comm", LocalType: "*mpi.Comm", Doc: "mpi communicator", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"AllDWts", &gti.Field{Name: "AllDWts", Type: "[]float32", LocalType: "[]float32", Doc: "buffer of all dwt weight changes -- for mpi sharing", Directives: gti.Directives{}, Tag: "view:\"-\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})
