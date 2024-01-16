// Code generated by "goki generate ./..."; DO NOT EDIT.

package armaze

import (
	"goki.dev/gti"
)

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/examples/boa/armaze.Arm", IDName: "arm", Doc: "Arm represents the properties of a given arm of the N-maze.\nArms have characteristic distance and effort factors for getting\ndown the arm, and typically have a distinctive CS visible at the start\nand a US at the end, which has US-specific parameters for\nactually delivering reward or punishment.", Fields: []gti.Field{{Name: "Length", Doc: "length of arm: distance from CS start to US end for this arm"}, {Name: "Effort", Doc: "range of different effort levels per step (uniformly randomly sampled per step) for going down this arm"}, {Name: "US", Doc: "index of US present at the end of this arm -- -1 if none"}, {Name: "CS", Doc: "index of CS visible at the start of this arm, -1 if none"}, {Name: "ExValue", Doc: "current expected value = US.Prob * US.Mag * Drives-- computed at start of new approach"}, {Name: "ExPVpos", Doc: "current expected PVpos value = normalized ExValue -- computed at start of new approach"}, {Name: "ExPVneg", Doc: "current expected PVneg value = normalized time and effort costs"}, {Name: "ExUtil", Doc: "current expected utility = effort discounted version of ExPVpos -- computed at start of new approach"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/examples/boa/armaze.USParams", IDName: "us-params", Doc: "USParams has parameters for different USs", Fields: []gti.Field{{Name: "Negative", Doc: "if true is a negative valence US -- these are after the first NDrives USs"}, {Name: "Mag", Doc: "range of different magnitudes (uniformly sampled)"}, {Name: "Prob", Doc: "probability of delivering the US"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/examples/boa/armaze.Params", IDName: "params", Doc: "Params are misc environment parameters", Fields: []gti.Field{{Name: "TurnEffort", Doc: "effort for turning"}, {Name: "ConsumeEffort", Doc: "effort for consuming US"}, {Name: "AlwaysLeft", Doc: "always turn left -- zoolander style -- reduces degrees of freedom in evaluating behavior"}, {Name: "PermuteCSs", Doc: "permute the order of CSs prior to applying them to arms -- having this off makes it easier to visually determine match between Drive and arm approach, and shouldn't make any difference to behavior (model doesn't know about this ordering)."}, {Name: "RandomStart", Doc: "after running down an Arm, a new random starting location is selected (otherwise same arm as last run)"}, {Name: "OpenArms", Doc: "if true, allow movement between arms just by going Left or Right -- otherwise once past the start, no switching is allowed"}, {Name: "Inactive", Doc: "strength of inactive inputs (e.g., Drives in Approach paradigm)"}, {Name: "NYReps", Doc: "number of Y-axis repetitions of localist stimuli -- for redundancy in spiking nets"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/examples/boa/armaze.Config", IDName: "config", Doc: "Config has environment configuration", Fields: []gti.Field{{Name: "Paradigm", Doc: "experimental paradigm that governs the configuration and updating of environment state over time and the appropriate evaluation criteria."}, {Name: "Debug", Doc: "for debugging, print out key steps including a trace of the action generation logic"}, {Name: "NDrives", Doc: "number of different drive-like body states (hunger, thirst, etc), that are satisfied by a corresponding positive US outcome -- this does not include the first curiosity drive"}, {Name: "NNegUSs", Doc: "number of negative US outcomes -- these are added after NDrives positive USs to total US list"}, {Name: "NUSs", Doc: "total number of USs = NDrives + NNegUSs"}, {Name: "NArms", Doc: "number of different arms"}, {Name: "MaxArmLength", Doc: "maximum arm length (distance)"}, {Name: "NCSs", Doc: "number of different CSs -- typically at least a unique CS per US -- relationship is determined in the US params"}, {Name: "USs", Doc: "parameters associated with each US.  The first NDrives are positive USs, and beyond that are negative USs"}, {Name: "Arms", Doc: "state of each arm: dist, effort, US, CS"}, {Name: "Params", Doc: "misc params"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/examples/boa/armaze.Geom", IDName: "geom", Doc: "Geom is overall geometry of the space", Fields: []gti.Field{{Name: "ArmWidth", Doc: "width of arm -- emery rodent is 1 unit wide"}, {Name: "ArmSpace", Doc: "total space between arms, ends up being divided on either side"}, {Name: "LengthScale", Doc: "multiplier per unit arm length -- keep square with width"}, {Name: "Thick", Doc: "thickness of walls, floor"}, {Name: "Height", Doc: "height of walls"}, {Name: "ArmWidthTot", Doc: "width + space"}, {Name: "Depth", Doc: "computed total depth, starts at 0 goes deep"}, {Name: "Width", Doc: "computed total width"}, {Name: "HalfWidth", Doc: "half width for centering on 0 X"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/examples/boa/armaze.GUI", IDName: "gui", Doc: "GUI renders multiple views of the flat world env", Methods: []gti.Method{{Name: "Left", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "Right", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "Forward", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "Consume", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}}, Fields: []gti.Field{{Name: "Disp", Doc: "update display -- turn off to make it faster"}, {Name: "Env", Doc: "the env being visualized"}, {Name: "EnvName", Doc: "name of current env -- number is NData index"}, {Name: "SceneView", Doc: "3D visualization of the Scene"}, {Name: "Scene2D", Doc: "2D visualization of the Scene"}, {Name: "MatColors", Doc: "list of material colors"}, {Name: "StateColors", Doc: "internal state colors"}, {Name: "WallSize", Doc: "thickness (X) and height (Y) of walls"}, {Name: "State", Doc: "current internal / behavioral state"}, {Name: "Trace", Doc: "trace record of recent activity"}, {Name: "StructView", Doc: "view of the gui obj"}, {Name: "WorldTabs", Doc: "ArmMaze TabView"}, {Name: "IsRunning", Doc: "ArmMaze is running"}, {Name: "DepthVals", Doc: "current depth map"}, {Name: "Camera", Doc: "offscreen render camera settings"}, {Name: "DepthMap", Doc: "color map to use for rendering depth map"}, {Name: "EyeRFullImg", Doc: "first-person right-eye full field view"}, {Name: "EyeRFovImg", Doc: "first-person right-eye fovea view"}, {Name: "DepthImg", Doc: "depth map bitmap view"}, {Name: "USposPlot", Doc: "plot of positive valence drives, active OFC US state, and reward"}, {Name: "USposData", Doc: "data for USPlot"}, {Name: "USnegPlot", Doc: "plot of negative valence active OFC US state, and outcomes"}, {Name: "USnegData", Doc: "data for USPlot"}, {Name: "Geom", Doc: "geometry of world"}, {Name: "World", Doc: "world"}, {Name: "View3D", Doc: "3D view of world"}, {Name: "Emery", Doc: "emer group"}, {Name: "Arms", Doc: "arms group"}, {Name: "Stims", Doc: "stims group"}, {Name: "EyeR", Doc: "Right eye of emery"}, {Name: "Contacts", Doc: "contacts from last step, for body"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/examples/boa/armaze.Actions", IDName: "actions", Doc: "Actions is a list of mutually exclusive states\nfor tracing the behavior and internal state of Emery", Directives: []gti.Directive{{Tool: "enums", Directive: "enum"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/examples/boa/armaze.Env", IDName: "env", Doc: "Env implements an N-armed maze (\"bandit\")\nwith each Arm having a distinctive CS stimulus visible at the start\n(could be one of multiple possibilities) and (some probability of)\na US outcome at the end of the maze (could be either positive\nor negative, with (variable) magnitude and probability.", Fields: []gti.Field{{Name: "Nm", Doc: "name of environment -- Train or Test"}, {Name: "Di", Doc: "our data parallel index"}, {Name: "Config", Doc: "configuration parameters"}, {Name: "Drives", Doc: "current drive strength for each of Config.NDrives in normalized 0-1 units of each drive: 0 = first sim drive, not curiosity"}, {Name: "Arm", Doc: "arm-wise location: either facing (Pos=0) or in (Pos > 0)"}, {Name: "Pos", Doc: "current position in the Arm: 0 = at start looking in, otherwise at given distance into the arm"}, {Name: "Tick", Doc: "current integer time step since last NewStart"}, {Name: "TrgDrive", Doc: "current target drive, in paradigms where that is used (e.g., Approach)"}, {Name: "USConsumed", Doc: "Current US being consumed -- is -1 unless being consumed"}, {Name: "USValue", Doc: "reward or punishment value generated by the current US being consumed -- just the Magnitude of the US -- does NOT include any modulation by Drive"}, {Name: "JustConsumed", Doc: "just finished consuming a US -- ready to start doing something new"}, {Name: "ArmsMaxValue", Doc: "arm(s) with maximum Drive * Mag * Prob US outcomes"}, {Name: "MaxValue", Doc: "maximum value for ArmsMaxValue arms"}, {Name: "ArmsMaxUtil", Doc: "arm(s) with maximum Value outcome discounted by Effort"}, {Name: "MaxUtil", Doc: "maximum utility for ArmsMaxUtil arms"}, {Name: "ArmsNeg", Doc: "arm(s) with negative US outcomes"}, {Name: "LastAct", Doc: "last action taken"}, {Name: "Effort", Doc: "effort on current trial"}, {Name: "LastCS", Doc: "last CS seen"}, {Name: "LastUS", Doc: "last US -- previous trial"}, {Name: "ShouldGate", Doc: "true if looking at correct CS for first time"}, {Name: "JustGated", Doc: "just gated on this trial -- set by sim-- used for instinct"}, {Name: "HasGated", Doc: "has gated at some point during sequence -- set by sim -- used for instinct"}, {Name: "States", Doc: "named states -- e.g., USs, CSs, etc"}, {Name: "MaxLength", Doc: "maximum length of any arm"}, {Name: "Rand", Doc: "random number generator for the env -- all random calls must use this"}, {Name: "RndSeed", Doc: "random seed"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/examples/boa/armaze.Paradigms", IDName: "paradigms", Doc: "Paradigms is a list of experimental paradigms that\ngovern the configuration and updating of environment\nstate over time and the appropriate evaluation criteria.", Directives: []gti.Directive{{Tool: "enums", Directive: "enum"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/examples/boa/armaze.TraceStates", IDName: "trace-states", Doc: "TraceStates is a list of mutually exclusive states\nfor tracing the behavior and internal state of Emery", Directives: []gti.Directive{{Tool: "enums", Directive: "enum"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/examples/boa/armaze.TraceRec", IDName: "trace-rec", Doc: "TraceRec holds record of info for tracing behavior, state", Fields: []gti.Field{{Name: "Time", Doc: "absolute time"}, {Name: "Trial", Doc: "trial counter"}, {Name: "Arm", Doc: "current arm"}, {Name: "Pos", Doc: "position in arm"}, {Name: "State", Doc: "behavioral / internal state summary"}, {Name: "Drives", Doc: "NDrives current drive state level"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/axon/v2/examples/boa/armaze.StateTrace", IDName: "state-trace", Doc: "StateTrace holds trace records"})
