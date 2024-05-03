// Code generated by "core generate -add-types"; DO NOT EDIT.

package main

import (
	"cogentcore.org/core/types"
)

var _ = types.AddType(&types.Type{Name: "main.EnvConfig", IDName: "env-config", Doc: "EnvConfig has config params for environment\nnote: only adding fields for key Env params that matter for both Network and Env\nother params are set via the Env map data mechanism.", Directives: []types.Directive{{Tool: "types", Directive: "add"}}, Fields: []types.Field{{Name: "Env", Doc: "env parameters -- can set any field/subfield on Env struct, using standard TOML formatting"}, {Name: "NOutPer", Doc: "number of units per localist output unit"}}})

var _ = types.AddType(&types.Type{Name: "main.ParamConfig", IDName: "param-config", Doc: "ParamConfig has config parameters related to sim params", Directives: []types.Directive{{Tool: "types", Directive: "add"}}, Fields: []types.Field{{Name: "Network", Doc: "network parameters"}, {Name: "Sheet", Doc: "Extra Param Sheet name(s) to use (space separated if multiple) -- must be valid name as listed in compiled-in params or loaded params"}, {Name: "Tag", Doc: "extra tag to add to file names and logs saved from this run"}, {Name: "Note", Doc: "user note -- describe the run params etc -- like a git commit message for the run"}, {Name: "File", Doc: "Name of the JSON file to input saved parameters from."}, {Name: "SaveAll", Doc: "Save a snapshot of all current param and config settings in a directory named params_<datestamp> (or _good if Good is true), then quit -- useful for comparing to later changes and seeing multiple views of current params"}, {Name: "Good", Doc: "for SaveAll, save to params_good for a known good params state.  This can be done prior to making a new release after all tests are passing -- add results to git to provide a full diff record of all params over time."}, {Name: "V1V4Path"}}})

var _ = types.AddType(&types.Type{Name: "main.RunConfig", IDName: "run-config", Doc: "RunConfig has config parameters related to running the sim", Directives: []types.Directive{{Tool: "types", Directive: "add"}}, Fields: []types.Field{{Name: "GPU", Doc: "use the GPU for computation -- generally faster even for small models if NData ~16"}, {Name: "NData", Doc: "number of data-parallel items to process in parallel per trial -- works (and is significantly faster) for both CPU and GPU.  Results in an effective mini-batch of learning."}, {Name: "NThreads", Doc: "number of parallel threads for CPU computation -- 0 = use default"}, {Name: "Run", Doc: "starting run number -- determines the random seed -- runs counts from there -- can do all runs in parallel by launching separate jobs with each run, runs = 1"}, {Name: "NRuns", Doc: "total number of runs to do when running Train"}, {Name: "NEpochs", Doc: "total number of epochs per run"}, {Name: "NTrials", Doc: "total number of trials per epoch.  Should be an even multiple of NData."}, {Name: "PCAInterval", Doc: "how frequently (in epochs) to compute PCA on hidden representations to measure variance?"}, {Name: "TestInterval", Doc: "how often to run through all the test patterns, in terms of training epochs -- can use 0 or -1 for no testing"}}})

var _ = types.AddType(&types.Type{Name: "main.LogConfig", IDName: "log-config", Doc: "LogConfig has config parameters related to logging data", Directives: []types.Directive{{Tool: "types", Directive: "add"}}, Fields: []types.Field{{Name: "SaveWts", Doc: "if true, save final weights after each run"}, {Name: "Epoch", Doc: "if true, save train epoch log to file, as .epc.tsv typically"}, {Name: "Run", Doc: "if true, save run log to file, as .run.tsv typically"}, {Name: "Trial", Doc: "if true, save train trial log to file, as .trl.tsv typically. May be large."}, {Name: "TestEpoch", Doc: "if true, save testing epoch log to file, as .tst_epc.tsv typically.  In general it is better to copy testing items over to the training epoch log and record there."}, {Name: "TestTrial", Doc: "if true, save testing trial log to file, as .tst_trl.tsv typically. May be large."}, {Name: "NetData", Doc: "if true, save network activation etc data from testing trials, for later viewing in netview"}}})

var _ = types.AddType(&types.Type{Name: "main.Config", IDName: "config", Doc: "Config is a standard Sim config -- use as a starting point.", Directives: []types.Directive{{Tool: "types", Directive: "add"}}, Fields: []types.Field{{Name: "Includes", Doc: "specify include files here, and after configuration, it contains list of include files added"}, {Name: "GUI", Doc: "open the GUI -- does not automatically run -- if false, then runs automatically and quits"}, {Name: "Debug", Doc: "log debugging information"}, {Name: "Env", Doc: "environment configuration options"}, {Name: "Params", Doc: "parameter related configuration options"}, {Name: "Run", Doc: "sim running related configuration options"}, {Name: "Log", Doc: "data logging related configuration options"}}})

var _ = types.AddType(&types.Type{Name: "main.LEDEnv", IDName: "led-env", Doc: "LEDEnv generates images of old-school \"LED\" style \"letters\" composed of a set of horizontal\nand vertical elements.  All possible such combinations of 3 out of 6 line segments are created.\nRenders using SVG.", Fields: []types.Field{{Name: "Nm", Doc: "name of this environment"}, {Name: "Dsc", Doc: "description of this environment"}, {Name: "Draw", Doc: "draws LEDs onto image"}, {Name: "Vis", Doc: "visual processing params"}, {Name: "NOutPer", Doc: "number of output units per LED item -- spiking benefits from replication"}, {Name: "MinLED", Doc: "minimum LED number to draw (0-19)"}, {Name: "MaxLED", Doc: "maximum LED number to draw (0-19)"}, {Name: "CurLED", Doc: "current LED number that was drawn"}, {Name: "PrvLED", Doc: "previous LED number that was drawn"}, {Name: "XFormRand", Doc: "random transform parameters"}, {Name: "XForm", Doc: "current -- prev transforms"}, {Name: "Run", Doc: "current run of model as provided during Init"}, {Name: "Epoch", Doc: "number of times through Seq.Max number of sequences"}, {Name: "Trial", Doc: "trial is the step counter within epoch"}, {Name: "OrigImg", Doc: "original image prior to random transforms"}, {Name: "Output", Doc: "CurLED one-hot output tensor"}}})

var _ = types.AddType(&types.Type{Name: "main.LEDraw", IDName: "le-draw", Doc: "LEDraw renders old-school \"LED\" style \"letters\" composed of a set of horizontal\nand vertical elements.  All possible such combinations of 3 out of 6 line segments are created.\nRenders using SVG.", Directives: []types.Directive{{Tool: "types", Directive: "add"}}, Fields: []types.Field{{Name: "Width", Doc: "line width of LEDraw as percent of display size"}, {Name: "Size", Doc: "size of overall LED as proportion of overall image size"}, {Name: "LineColor", Doc: "color name for drawing lines"}, {Name: "BgColor", Doc: "color name for background"}, {Name: "ImgSize", Doc: "size of image to render"}, {Name: "Image", Doc: "rendered image"}, {Name: "Paint", Doc: "painting context object"}}})

var _ = types.AddType(&types.Type{Name: "main.LEDSegs", IDName: "led-segs", Doc: "LEDSegs are the led segments"})

var _ = types.AddType(&types.Type{Name: "main.Sim", IDName: "sim", Doc: "Sim encapsulates the entire simulation model, and we define all the\nfunctionality as methods on this struct.  This structure keeps all relevant\nstate information organized and available without having to pass everything around\nas arguments to methods, and provides the core GUI interface (note the view tags\nfor the fields which provide hints to how things should be displayed).", Fields: []types.Field{{Name: "Config", Doc: "simulation configuration parameters -- set by .toml config file and / or args"}, {Name: "Net", Doc: "the network -- click to view / edit parameters for layers, paths, etc"}, {Name: "Params", Doc: "all parameter management"}, {Name: "Loops", Doc: "contains looper control loops for running sim"}, {Name: "Stats", Doc: "contains computed statistic values"}, {Name: "Logs", Doc: "Contains all the logs and information about the logs.'"}, {Name: "Envs", Doc: "Environments"}, {Name: "Context", Doc: "axon timing parameters and state"}, {Name: "ViewUpdate", Doc: "netview update parameters"}, {Name: "GUI", Doc: "manages all the gui elements"}, {Name: "RandSeeds", Doc: "a list of random seeds to use for each run"}}})

var _ = types.AddType(&types.Type{Name: "main.Vis", IDName: "vis", Doc: "Vis encapsulates specific visual processing pipeline for V1 filtering", Directives: []types.Directive{{Tool: "types", Directive: "add"}}, Fields: []types.Field{{Name: "V1sGabor", Doc: "V1 simple gabor filter parameters"}, {Name: "V1sGeom", Doc: "geometry of input, output for V1 simple-cell processing"}, {Name: "V1sNeighInhib", Doc: "neighborhood inhibition for V1s -- each unit gets inhibition from same feature in nearest orthogonal neighbors -- reduces redundancy of feature code"}, {Name: "V1sKWTA", Doc: "kwta parameters for V1s"}, {Name: "ImgSize", Doc: "target image size to use -- images will be rescaled to this size"}, {Name: "V1sGaborTsr", Doc: "V1 simple gabor filter tensor"}, {Name: "ImgTsr", Doc: "input image as tensor"}, {Name: "Img", Doc: "current input image"}, {Name: "V1sTsr", Doc: "V1 simple gabor filter output tensor"}, {Name: "V1sExtGiTsr", Doc: "V1 simple extra Gi from neighbor inhibition tensor"}, {Name: "V1sKwtaTsr", Doc: "V1 simple gabor filter output, kwta output tensor"}, {Name: "V1sPoolTsr", Doc: "V1 simple gabor filter output, max-pooled 2x2 of V1sKwta tensor"}, {Name: "V1sUnPoolTsr", Doc: "V1 simple gabor filter output, un-max-pooled 2x2 of V1sPool tensor"}, {Name: "V1sAngOnlyTsr", Doc: "V1 simple gabor filter output, angle-only features tensor"}, {Name: "V1sAngPoolTsr", Doc: "V1 simple gabor filter output, max-pooled 2x2 of AngOnly tensor"}, {Name: "V1cLenSumTsr", Doc: "V1 complex length sum filter output tensor"}, {Name: "V1cEndStopTsr", Doc: "V1 complex end stop filter output tensor"}, {Name: "V1AllTsr", Doc: "Combined V1 output tensor with V1s simple as first two rows, then length sum, then end stops = 5 rows total"}, {Name: "V1sInhibs", Doc: "inhibition values for V1s KWTA"}}})
