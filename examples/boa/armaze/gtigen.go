// Code generated by "goki generate -add-types"; DO NOT EDIT.

package armaze

import (
	"goki.dev/gti"
	"goki.dev/ordmap"
)

var _ = gti.AddType(&gti.Type{
	Name:       "github.com/emer/axon/v2/examples/boa/armaze.Arm",
	ShortName:  "armaze.Arm",
	IDName:     "arm",
	Doc:        "Arm represents the properties of a given arm of the N-maze.\nArms have characteristic distance and effort factors for getting\ndown the arm, and typically have a distinctive CS visible at the start\nand a US at the end, which has US-specific parameters for\nactually delivering reward or punishment.",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Length", &gti.Field{Name: "Length", Type: "int", LocalType: "int", Doc: "length of arm: distance from CS start to US end for this arm", Directives: gti.Directives{}, Tag: ""}},
		{"Effort", &gti.Field{Name: "Effort", Type: "goki.dev/etable/v2/minmax.F32", LocalType: "minmax.F32", Doc: "range of different effort levels per step (uniformly randomly sampled per step) for going down this arm", Directives: gti.Directives{}, Tag: ""}},
		{"US", &gti.Field{Name: "US", Type: "int", LocalType: "int", Doc: "index of US present at the end of this arm -- -1 if none", Directives: gti.Directives{}, Tag: ""}},
		{"CS", &gti.Field{Name: "CS", Type: "int", LocalType: "int", Doc: "index of CS visible at the start of this arm, -1 if none", Directives: gti.Directives{}, Tag: ""}},
		{"ExValue", &gti.Field{Name: "ExValue", Type: "float32", LocalType: "float32", Doc: "current expected value = US.Prob * US.Mag * Drives-- computed at start of new approach", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"ExPVpos", &gti.Field{Name: "ExPVpos", Type: "float32", LocalType: "float32", Doc: "current expected PVpos value = normalized ExValue -- computed at start of new approach", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"ExPVneg", &gti.Field{Name: "ExPVneg", Type: "float32", LocalType: "float32", Doc: "current expected PVneg value = normalized time and effort costs", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"ExUtil", &gti.Field{Name: "ExUtil", Type: "float32", LocalType: "float32", Doc: "current expected utility = effort discounted version of ExPVpos -- computed at start of new approach", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "github.com/emer/axon/v2/examples/boa/armaze.USParams",
	ShortName:  "armaze.USParams",
	IDName:     "us-params",
	Doc:        "USParams has parameters for different USs",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Negative", &gti.Field{Name: "Negative", Type: "bool", LocalType: "bool", Doc: "if true is a negative valence US -- these are after the first NDrives USs", Directives: gti.Directives{}, Tag: ""}},
		{"Mag", &gti.Field{Name: "Mag", Type: "goki.dev/etable/v2/minmax.F32", LocalType: "minmax.F32", Doc: "range of different magnitudes (uniformly sampled)", Directives: gti.Directives{}, Tag: ""}},
		{"Prob", &gti.Field{Name: "Prob", Type: "float32", LocalType: "float32", Doc: "probability of delivering the US", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "github.com/emer/axon/v2/examples/boa/armaze.Params",
	ShortName:  "armaze.Params",
	IDName:     "params",
	Doc:        "Params are misc environment parameters",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"TurnEffort", &gti.Field{Name: "TurnEffort", Type: "goki.dev/etable/v2/minmax.F32", LocalType: "minmax.F32", Doc: "effort for turning", Directives: gti.Directives{}, Tag: "nest:\"+\" def:\"{'Min':0.5, 'Max':0.5}\""}},
		{"ConsumeEffort", &gti.Field{Name: "ConsumeEffort", Type: "goki.dev/etable/v2/minmax.F32", LocalType: "minmax.F32", Doc: "effort for consuming US", Directives: gti.Directives{}, Tag: "nest:\"+\" def:\"{'Min':0.5, 'Max':0.5}\""}},
		{"AlwaysLeft", &gti.Field{Name: "AlwaysLeft", Type: "bool", LocalType: "bool", Doc: "always turn left -- zoolander style -- reduces degrees of freedom in evaluating behavior", Directives: gti.Directives{}, Tag: "def:\"true\""}},
		{"PermuteCSs", &gti.Field{Name: "PermuteCSs", Type: "bool", LocalType: "bool", Doc: "permute the order of CSs prior to applying them to arms -- having this off makes it easier to visually determine match between Drive and arm approach, and shouldn't make any difference to behavior (model doesn't know about this ordering).", Directives: gti.Directives{}, Tag: "def:\"false\""}},
		{"RandomStart", &gti.Field{Name: "RandomStart", Type: "bool", LocalType: "bool", Doc: "after running down an Arm, a new random starting location is selected (otherwise same arm as last run)", Directives: gti.Directives{}, Tag: "def:\"true\""}},
		{"OpenArms", &gti.Field{Name: "OpenArms", Type: "bool", LocalType: "bool", Doc: "if true, allow movement between arms just by going Left or Right -- otherwise once past the start, no switching is allowed", Directives: gti.Directives{}, Tag: "def:\"true\""}},
		{"Inactive", &gti.Field{Name: "Inactive", Type: "goki.dev/etable/v2/minmax.F32", LocalType: "minmax.F32", Doc: "strength of inactive inputs (e.g., Drives in Approach paradigm)", Directives: gti.Directives{}, Tag: "nest:\"+\" def:\"{'Min':0, 'Max':0}\" view:\"inline\""}},
		{"NYReps", &gti.Field{Name: "NYReps", Type: "int", LocalType: "int", Doc: "number of Y-axis repetitions of localist stimuli -- for redundancy in spiking nets", Directives: gti.Directives{}, Tag: "def:\"4\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "github.com/emer/axon/v2/examples/boa/armaze.Config",
	ShortName:  "armaze.Config",
	IDName:     "config",
	Doc:        "Config has environment configuration",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Paradigm", &gti.Field{Name: "Paradigm", Type: "github.com/emer/axon/v2/examples/boa/armaze.Paradigms", LocalType: "Paradigms", Doc: "experimental paradigm that governs the configuration and updating of environment state over time and the appropriate evaluation criteria.", Directives: gti.Directives{}, Tag: ""}},
		{"Debug", &gti.Field{Name: "Debug", Type: "bool", LocalType: "bool", Doc: "for debugging, print out key steps including a trace of the action generation logic", Directives: gti.Directives{}, Tag: ""}},
		{"NDrives", &gti.Field{Name: "NDrives", Type: "int", LocalType: "int", Doc: "number of different drive-like body states (hunger, thirst, etc), that are satisfied by a corresponding positive US outcome -- this does not include the first curiosity drive", Directives: gti.Directives{}, Tag: ""}},
		{"NNegUSs", &gti.Field{Name: "NNegUSs", Type: "int", LocalType: "int", Doc: "number of negative US outcomes -- these are added after NDrives positive USs to total US list", Directives: gti.Directives{}, Tag: ""}},
		{"NUSs", &gti.Field{Name: "NUSs", Type: "int", LocalType: "int", Doc: "total number of USs = NDrives + NNegUSs", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"NArms", &gti.Field{Name: "NArms", Type: "int", LocalType: "int", Doc: "number of different arms", Directives: gti.Directives{}, Tag: ""}},
		{"MaxArmLength", &gti.Field{Name: "MaxArmLength", Type: "int", LocalType: "int", Doc: "maximum arm length (distance)", Directives: gti.Directives{}, Tag: ""}},
		{"NCSs", &gti.Field{Name: "NCSs", Type: "int", LocalType: "int", Doc: "number of different CSs -- typically at least a unique CS per US -- relationship is determined in the US params", Directives: gti.Directives{}, Tag: ""}},
		{"USs", &gti.Field{Name: "USs", Type: "[]*github.com/emer/axon/v2/examples/boa/armaze.USParams", LocalType: "[]*USParams", Doc: "parameters associated with each US.  The first NDrives are positive USs, and beyond that are negative USs", Directives: gti.Directives{}, Tag: ""}},
		{"Arms", &gti.Field{Name: "Arms", Type: "[]*github.com/emer/axon/v2/examples/boa/armaze.Arm", LocalType: "[]*Arm", Doc: "state of each arm: dist, effort, US, CS", Directives: gti.Directives{}, Tag: ""}},
		{"Params", &gti.Field{Name: "Params", Type: "github.com/emer/axon/v2/examples/boa/armaze.Params", LocalType: "Params", Doc: "misc params", Directives: gti.Directives{}, Tag: "view:\"add-fields\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "github.com/emer/axon/v2/examples/boa/armaze.Geom",
	ShortName:  "armaze.Geom",
	IDName:     "geom",
	Doc:        "Geom is overall geometry of the space",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"ArmWidth", &gti.Field{Name: "ArmWidth", Type: "float32", LocalType: "float32", Doc: "width of arm -- emery rodent is 1 unit wide", Directives: gti.Directives{}, Tag: "def:\"2\""}},
		{"ArmSpace", &gti.Field{Name: "ArmSpace", Type: "float32", LocalType: "float32", Doc: "total space between arms, ends up being divided on either side", Directives: gti.Directives{}, Tag: "def:\"1\""}},
		{"LengthScale", &gti.Field{Name: "LengthScale", Type: "float32", LocalType: "float32", Doc: "multiplier per unit arm length -- keep square with width", Directives: gti.Directives{}, Tag: "def:\"2\""}},
		{"Thick", &gti.Field{Name: "Thick", Type: "float32", LocalType: "float32", Doc: "thickness of walls, floor", Directives: gti.Directives{}, Tag: "def:\"0.1\""}},
		{"Height", &gti.Field{Name: "Height", Type: "float32", LocalType: "float32", Doc: "height of walls", Directives: gti.Directives{}, Tag: "def:\"0.2\""}},
		{"ArmWidthTot", &gti.Field{Name: "ArmWidthTot", Type: "float32", LocalType: "float32", Doc: "width + space", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"Depth", &gti.Field{Name: "Depth", Type: "float32", LocalType: "float32", Doc: "computed total depth, starts at 0 goes deep", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"Width", &gti.Field{Name: "Width", Type: "float32", LocalType: "float32", Doc: "computed total width", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"HalfWidth", &gti.Field{Name: "HalfWidth", Type: "float32", LocalType: "float32", Doc: "half width for centering on 0 X", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "github.com/emer/axon/v2/examples/boa/armaze.GUI",
	ShortName:  "armaze.GUI",
	IDName:     "gui",
	Doc:        "GUI renders multiple views of the flat world env",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Disp", &gti.Field{Name: "Disp", Type: "bool", LocalType: "bool", Doc: "update display -- turn off to make it faster", Directives: gti.Directives{}, Tag: ""}},
		{"Env", &gti.Field{Name: "Env", Type: "*github.com/emer/axon/v2/examples/boa/armaze.Env", LocalType: "*Env", Doc: "the env being visualized", Directives: gti.Directives{}, Tag: ""}},
		{"EnvName", &gti.Field{Name: "EnvName", Type: "string", LocalType: "string", Doc: "name of current env -- number is NData index", Directives: gti.Directives{}, Tag: ""}},
		{"SceneView", &gti.Field{Name: "SceneView", Type: "*goki.dev/gi/v2/xyzv.SceneView", LocalType: "*xyzv.SceneView", Doc: "3D visualization of the Scene", Directives: gti.Directives{}, Tag: ""}},
		{"Scene2D", &gti.Field{Name: "Scene2D", Type: "*goki.dev/gi/v2/gi.SVG", LocalType: "*gi.SVG", Doc: "2D visualization of the Scene", Directives: gti.Directives{}, Tag: ""}},
		{"MatColors", &gti.Field{Name: "MatColors", Type: "[]string", LocalType: "[]string", Doc: "list of material colors", Directives: gti.Directives{}, Tag: ""}},
		{"StateColors", &gti.Field{Name: "StateColors", Type: "map[string]string", LocalType: "map[string]string", Doc: "internal state colors", Directives: gti.Directives{}, Tag: ""}},
		{"WallSize", &gti.Field{Name: "WallSize", Type: "goki.dev/mat32/v2.Vec2", LocalType: "mat32.Vec2", Doc: "thickness (X) and height (Y) of walls", Directives: gti.Directives{}, Tag: ""}},
		{"State", &gti.Field{Name: "State", Type: "github.com/emer/axon/v2/examples/boa/armaze.TraceStates", LocalType: "TraceStates", Doc: "current internal / behavioral state", Directives: gti.Directives{}, Tag: ""}},
		{"Trace", &gti.Field{Name: "Trace", Type: "github.com/emer/axon/v2/examples/boa/armaze.StateTrace", LocalType: "StateTrace", Doc: "trace record of recent activity", Directives: gti.Directives{}, Tag: ""}},
		{"StructView", &gti.Field{Name: "StructView", Type: "*goki.dev/gi/v2/giv.StructView", LocalType: "*giv.StructView", Doc: "view of the gui obj", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"WorldTabs", &gti.Field{Name: "WorldTabs", Type: "*goki.dev/gi/v2/gi.Tabs", LocalType: "*gi.Tabs", Doc: "ArmMaze TabView", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"IsRunning", &gti.Field{Name: "IsRunning", Type: "bool", LocalType: "bool", Doc: "ArmMaze is running", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"DepthVals", &gti.Field{Name: "DepthVals", Type: "[]float32", LocalType: "[]float32", Doc: "current depth map", Directives: gti.Directives{}, Tag: ""}},
		{"Camera", &gti.Field{Name: "Camera", Type: "github.com/emer/eve/v2/evev.Camera", LocalType: "evev.Camera", Doc: "offscreen render camera settings", Directives: gti.Directives{}, Tag: ""}},
		{"DepthMap", &gti.Field{Name: "DepthMap", Type: "goki.dev/gi/v2/giv.ColorMapName", LocalType: "giv.ColorMapName", Doc: "color map to use for rendering depth map", Directives: gti.Directives{}, Tag: ""}},
		{"EyeRFullImg", &gti.Field{Name: "EyeRFullImg", Type: "*goki.dev/gi/v2/gi.Image", LocalType: "*gi.Image", Doc: "first-person right-eye full field view", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"EyeRFovImg", &gti.Field{Name: "EyeRFovImg", Type: "*goki.dev/gi/v2/gi.Image", LocalType: "*gi.Image", Doc: "first-person right-eye fovea view", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"DepthImg", &gti.Field{Name: "DepthImg", Type: "*goki.dev/gi/v2/gi.Image", LocalType: "*gi.Image", Doc: "depth map bitmap view", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"USposPlot", &gti.Field{Name: "USposPlot", Type: "*goki.dev/etable/v2/eplot.Plot2D", LocalType: "*eplot.Plot2D", Doc: "plot of positive valence drives, active OFC US state, and reward", Directives: gti.Directives{}, Tag: ""}},
		{"USposData", &gti.Field{Name: "USposData", Type: "*goki.dev/etable/v2/etable.Table", LocalType: "*etable.Table", Doc: "data for USPlot", Directives: gti.Directives{}, Tag: ""}},
		{"USnegPlot", &gti.Field{Name: "USnegPlot", Type: "*goki.dev/etable/v2/eplot.Plot2D", LocalType: "*eplot.Plot2D", Doc: "plot of negative valence active OFC US state, and outcomes", Directives: gti.Directives{}, Tag: ""}},
		{"USnegData", &gti.Field{Name: "USnegData", Type: "*goki.dev/etable/v2/etable.Table", LocalType: "*etable.Table", Doc: "data for USPlot", Directives: gti.Directives{}, Tag: ""}},
		{"Geom", &gti.Field{Name: "Geom", Type: "github.com/emer/axon/v2/examples/boa/armaze.Geom", LocalType: "Geom", Doc: "geometry of world", Directives: gti.Directives{}, Tag: ""}},
		{"World", &gti.Field{Name: "World", Type: "*github.com/emer/eve/v2/eve.Group", LocalType: "*eve.Group", Doc: "world", Directives: gti.Directives{}, Tag: ""}},
		{"View3D", &gti.Field{Name: "View3D", Type: "*github.com/emer/eve/v2/evev.View", LocalType: "*evev.View", Doc: "3D view of world", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"Emery", &gti.Field{Name: "Emery", Type: "*github.com/emer/eve/v2/eve.Group", LocalType: "*eve.Group", Doc: "emer group", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"Arms", &gti.Field{Name: "Arms", Type: "*github.com/emer/eve/v2/eve.Group", LocalType: "*eve.Group", Doc: "arms group", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"Stims", &gti.Field{Name: "Stims", Type: "*github.com/emer/eve/v2/eve.Group", LocalType: "*eve.Group", Doc: "stims group", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"EyeR", &gti.Field{Name: "EyeR", Type: "github.com/emer/eve/v2/eve.Body", LocalType: "eve.Body", Doc: "Right eye of emery", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"Contacts", &gti.Field{Name: "Contacts", Type: "github.com/emer/eve/v2/eve.Contacts", LocalType: "eve.Contacts", Doc: "contacts from last step, for body", Directives: gti.Directives{}, Tag: "view:\"-\""}},
	}),
	Embeds: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{
		{"Left", &gti.Method{Name: "Left", Doc: "", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"Right", &gti.Method{Name: "Right", Doc: "", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"Forward", &gti.Method{Name: "Forward", Doc: "", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"Consume", &gti.Method{Name: "Consume", Doc: "", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
	}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/v2/examples/boa/armaze.Actions",
	ShortName: "armaze.Actions",
	IDName:    "actions",
	Doc:       "Actions is a list of mutually exclusive states\nfor tracing the behavior and internal state of Emery",
	Directives: gti.Directives{
		&gti.Directive{Tool: "enums", Directive: "enum", Args: []string{}},
	},

	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "github.com/emer/axon/v2/examples/boa/armaze.Env",
	ShortName:  "armaze.Env",
	IDName:     "env",
	Doc:        "Env implements an N-armed maze (\"bandit\")\nwith each Arm having a distinctive CS stimulus visible at the start\n(could be one of multiple possibilities) and (some probability of)\na US outcome at the end of the maze (could be either positive\nor negative, with (variable) magnitude and probability.",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Nm", &gti.Field{Name: "Nm", Type: "string", LocalType: "string", Doc: "name of environment -- Train or Test", Directives: gti.Directives{}, Tag: ""}},
		{"Di", &gti.Field{Name: "Di", Type: "int", LocalType: "int", Doc: "our data parallel index", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"Config", &gti.Field{Name: "Config", Type: "github.com/emer/axon/v2/examples/boa/armaze.Config", LocalType: "Config", Doc: "configuration parameters", Directives: gti.Directives{}, Tag: ""}},
		{"Drives", &gti.Field{Name: "Drives", Type: "[]float32", LocalType: "[]float32", Doc: "current drive strength for each of Config.NDrives in normalized 0-1 units of each drive: 0 = first sim drive, not curiosity", Directives: gti.Directives{}, Tag: ""}},
		{"Arm", &gti.Field{Name: "Arm", Type: "int", LocalType: "int", Doc: "arm-wise location: either facing (Pos=0) or in (Pos > 0)", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"Pos", &gti.Field{Name: "Pos", Type: "int", LocalType: "int", Doc: "current position in the Arm: 0 = at start looking in, otherwise at given distance into the arm", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"Tick", &gti.Field{Name: "Tick", Type: "int", LocalType: "int", Doc: "current integer time step since last NewStart", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"TrgDrive", &gti.Field{Name: "TrgDrive", Type: "int", LocalType: "int", Doc: "current target drive, in paradigms where that is used (e.g., Approach)", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"USConsumed", &gti.Field{Name: "USConsumed", Type: "int", LocalType: "int", Doc: "Current US being consumed -- is -1 unless being consumed", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"USValue", &gti.Field{Name: "USValue", Type: "float32", LocalType: "float32", Doc: "reward or punishment value generated by the current US being consumed -- just the Magnitude of the US -- does NOT include any modulation by Drive", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"JustConsumed", &gti.Field{Name: "JustConsumed", Type: "bool", LocalType: "bool", Doc: "just finished consuming a US -- ready to start doing something new", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"ArmsMaxValue", &gti.Field{Name: "ArmsMaxValue", Type: "[]int", LocalType: "[]int", Doc: "arm(s) with maximum Drive * Mag * Prob US outcomes", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"MaxValue", &gti.Field{Name: "MaxValue", Type: "float32", LocalType: "float32", Doc: "maximum value for ArmsMaxValue arms", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"ArmsMaxUtil", &gti.Field{Name: "ArmsMaxUtil", Type: "[]int", LocalType: "[]int", Doc: "arm(s) with maximum Value outcome discounted by Effort", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"MaxUtil", &gti.Field{Name: "MaxUtil", Type: "float32", LocalType: "float32", Doc: "maximum utility for ArmsMaxUtil arms", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"ArmsNeg", &gti.Field{Name: "ArmsNeg", Type: "[]int", LocalType: "[]int", Doc: "arm(s) with negative US outcomes", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"LastAct", &gti.Field{Name: "LastAct", Type: "github.com/emer/axon/v2/examples/boa/armaze.Actions", LocalType: "Actions", Doc: "last action taken", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"Effort", &gti.Field{Name: "Effort", Type: "float32", LocalType: "float32", Doc: "effort on current trial", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"LastCS", &gti.Field{Name: "LastCS", Type: "int", LocalType: "int", Doc: "last CS seen", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"LastUS", &gti.Field{Name: "LastUS", Type: "int", LocalType: "int", Doc: "last US -- previous trial", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"ShouldGate", &gti.Field{Name: "ShouldGate", Type: "bool", LocalType: "bool", Doc: "true if looking at correct CS for first time", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"JustGated", &gti.Field{Name: "JustGated", Type: "bool", LocalType: "bool", Doc: "just gated on this trial -- set by sim-- used for instinct", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"HasGated", &gti.Field{Name: "HasGated", Type: "bool", LocalType: "bool", Doc: "has gated at some point during sequence -- set by sim -- used for instinct", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"States", &gti.Field{Name: "States", Type: "map[string]*goki.dev/etable/v2/etensor.Float32", LocalType: "map[string]*etensor.Float32", Doc: "named states -- e.g., USs, CSs, etc", Directives: gti.Directives{}, Tag: ""}},
		{"MaxLength", &gti.Field{Name: "MaxLength", Type: "int", LocalType: "int", Doc: "maximum length of any arm", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
		{"Rand", &gti.Field{Name: "Rand", Type: "github.com/emer/emergent/v2/erand.SysRand", LocalType: "erand.SysRand", Doc: "random number generator for the env -- all random calls must use this", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"RndSeed", &gti.Field{Name: "RndSeed", Type: "int64", LocalType: "int64", Doc: "random seed", Directives: gti.Directives{}, Tag: "inactive:\"+\""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/v2/examples/boa/armaze.Paradigms",
	ShortName: "armaze.Paradigms",
	IDName:    "paradigms",
	Doc:       "Paradigms is a list of experimental paradigms that\ngovern the configuration and updating of environment\nstate over time and the appropriate evaluation criteria.",
	Directives: gti.Directives{
		&gti.Directive{Tool: "enums", Directive: "enum", Args: []string{}},
	},

	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "github.com/emer/axon/v2/examples/boa/armaze.TraceStates",
	ShortName: "armaze.TraceStates",
	IDName:    "trace-states",
	Doc:       "TraceStates is a list of mutually exclusive states\nfor tracing the behavior and internal state of Emery",
	Directives: gti.Directives{
		&gti.Directive{Tool: "enums", Directive: "enum", Args: []string{}},
	},

	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "github.com/emer/axon/v2/examples/boa/armaze.TraceRec",
	ShortName:  "armaze.TraceRec",
	IDName:     "trace-rec",
	Doc:        "TraceRec holds record of info for tracing behavior, state",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Time", &gti.Field{Name: "Time", Type: "float32", LocalType: "float32", Doc: "absolute time", Directives: gti.Directives{}, Tag: ""}},
		{"Trial", &gti.Field{Name: "Trial", Type: "int", LocalType: "int", Doc: "trial counter", Directives: gti.Directives{}, Tag: ""}},
		{"Arm", &gti.Field{Name: "Arm", Type: "int", LocalType: "int", Doc: "current arm", Directives: gti.Directives{}, Tag: ""}},
		{"Pos", &gti.Field{Name: "Pos", Type: "int", LocalType: "int", Doc: "position in arm", Directives: gti.Directives{}, Tag: ""}},
		{"State", &gti.Field{Name: "State", Type: "github.com/emer/axon/v2/examples/boa/armaze.TraceStates", LocalType: "TraceStates", Doc: "behavioral / internal state summary", Directives: gti.Directives{}, Tag: ""}},
		{"Drives", &gti.Field{Name: "Drives", Type: "[]float32", LocalType: "[]float32", Doc: "NDrives current drive state level", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "github.com/emer/axon/v2/examples/boa/armaze.StateTrace",
	ShortName:  "armaze.StateTrace",
	IDName:     "state-trace",
	Doc:        "StateTrace holds trace records",
	Directives: gti.Directives{},

	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})
