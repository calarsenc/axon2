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
		{"InputPct", &gti.Field{Name: "InputPct", Type: "float32", LocalType: "float32", Doc: "percent of active units in input layer (literally number of active units, because input has 100 units total)", Directives: gti.Directives{}, Tag: "def:\"15\" min:\"5\" max:\"50\" step:\"1\""}},
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
		{"NLayers", &gti.Field{Name: "NLayers", Type: "int", LocalType: "int", Doc: "number of hidden layers to add", Directives: gti.Directives{}, Tag: "def:\"2\" min:\"1\""}},
		{"HidSize", &gti.Field{Name: "HidSize", Type: "github.com/emer/emergent/v2/evec.Vec2i", LocalType: "evec.Vec2i", Doc: "size of hidden layers", Directives: gti.Directives{}, Tag: "def:\"{'X':10,'Y':10}\""}},
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
		{"Epoch", &gti.Field{Name: "Epoch", Type: "bool", LocalType: "bool", Doc: "if true, save train epoch log to file, as .epc.tsv typically", Directives: gti.Directives{}, Tag: "def:\"true\" nest:\"+\""}},
		{"Trial", &gti.Field{Name: "Trial", Type: "bool", LocalType: "bool", Doc: "if true, save train trial log to file, as .trl.tsv typically. May be large.", Directives: gti.Directives{}, Tag: "def:\"false\" nest:\"+\""}},
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
		{"Env", &gti.Field{Name: "Env", Type: "github.com/emer/axon/examples/inhib.EnvConfig", LocalType: "EnvConfig", Doc: "environment configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
		{"Params", &gti.Field{Name: "Params", Type: "github.com/emer/axon/examples/inhib.ParamConfig", LocalType: "ParamConfig", Doc: "parameter related configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
		{"Run", &gti.Field{Name: "Run", Type: "github.com/emer/axon/examples/inhib.RunConfig", LocalType: "RunConfig", Doc: "sim running related configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
		{"Log", &gti.Field{Name: "Log", Type: "github.com/emer/axon/examples/inhib.LogConfig", LocalType: "LogConfig", Doc: "data logging related configuration options", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
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
		{"Config", &gti.Field{Name: "Config", Type: "github.com/emer/axon/examples/inhib.Config", LocalType: "Config", Doc: "simulation configuration parameters -- set by .toml config file and / or args", Directives: gti.Directives{}, Tag: ""}},
		{"Net", &gti.Field{Name: "Net", Type: "*invalid type", LocalType: "*axon.Network", Doc: "the network -- click to view / edit parameters for layers, prjns, etc", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"Params", &gti.Field{Name: "Params", Type: "github.com/emer/emergent/v2/emer.NetParams", LocalType: "emer.NetParams", Doc: "all parameter management", Directives: gti.Directives{}, Tag: "view:\"inline\""}},
		{"Loops", &gti.Field{Name: "Loops", Type: "*github.com/emer/emergent/v2/looper.Manager", LocalType: "*looper.Manager", Doc: "contains looper control loops for running sim", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"Stats", &gti.Field{Name: "Stats", Type: "github.com/emer/emergent/v2/estats.Stats", LocalType: "estats.Stats", Doc: "contains computed statistic values", Directives: gti.Directives{}, Tag: ""}},
		{"Logs", &gti.Field{Name: "Logs", Type: "github.com/emer/emergent/v2/elog.Logs", LocalType: "elog.Logs", Doc: "Contains all the logs and information about the logs.'", Directives: gti.Directives{}, Tag: ""}},
		{"Pats", &gti.Field{Name: "Pats", Type: "*github.com/emer/etable/v2/etable.Table", LocalType: "*etable.Table", Doc: "the training patterns to use", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"Context", &gti.Field{Name: "Context", Type: "invalid type", LocalType: "axon.Context", Doc: "axon timing parameters and state", Directives: gti.Directives{}, Tag: ""}},
		{"ViewUpdt", &gti.Field{Name: "ViewUpdt", Type: "github.com/emer/emergent/v2/netview.ViewUpdt", LocalType: "netview.ViewUpdt", Doc: "netview update parameters", Directives: gti.Directives{}, Tag: "view:\"inline\""}},
		{"GUI", &gti.Field{Name: "GUI", Type: "github.com/emer/emergent/v2/egui.GUI", LocalType: "egui.GUI", Doc: "manages all the gui elements", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"RndSeeds", &gti.Field{Name: "RndSeeds", Type: "github.com/emer/emergent/v2/erand.Seeds", LocalType: "erand.Seeds", Doc: "a list of random seeds to use for each run", Directives: gti.Directives{}, Tag: "view:\"-\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})
