// Code generated by "core generate -add-types"; DO NOT EDIT.

package main

import (
	"cogentcore.org/core/gti"
)

var _ = gti.AddType(&gti.Type{Name: "main.EnvConfig", IDName: "env-config", Doc: "EnvConfig has config params for environment\nnote: only adding fields for key Env params that matter for both Network and Env\nother params are set via the Env map data mechanism.", Fields: []gti.Field{{Name: "Env", Doc: "env parameters -- can set any field/subfield on Env struct, using standard TOML formatting"}, {Name: "UnitsPer", Doc: "number of units per localist output unit"}, {Name: "FullSong", Doc: "train the full song -- else 30 notes"}, {Name: "PlayTarg", Doc: "during testing, play the target note instead of the actual network output"}, {Name: "TestClamp", Doc: "drive inputs from the training sequence during testing -- otherwise use network's own output"}}})

var _ = gti.AddType(&gti.Type{Name: "main.ParamConfig", IDName: "param-config", Doc: "ParamConfig has config parameters related to sim params", Fields: []gti.Field{{Name: "Network", Doc: "network parameters"}, {Name: "Hid2", Doc: "use a second hidden layer that predicts the first -- is not beneficia"}, {Name: "Sheet", Doc: "Extra Param Sheet name(s) to use (space separated if multiple) -- must be valid name as listed in compiled-in params or loaded params"}, {Name: "Tag", Doc: "extra tag to add to file names and logs saved from this run"}, {Name: "Note", Doc: "user note -- describe the run params etc -- like a git commit message for the run"}, {Name: "File", Doc: "Name of the JSON file to input saved parameters from."}, {Name: "SaveAll", Doc: "Save a snapshot of all current param and config settings in a directory named params_<datestamp> (or _good if Good is true), then quit -- useful for comparing to later changes and seeing multiple views of current params"}, {Name: "Good", Doc: "for SaveAll, save to params_good for a known good params state.  This can be done prior to making a new release after all tests are passing -- add results to git to provide a full diff record of all params over time."}}})

var _ = gti.AddType(&gti.Type{Name: "main.RunConfig", IDName: "run-config", Doc: "RunConfig has config parameters related to running the sim", Fields: []gti.Field{{Name: "GPU", Doc: "use the GPU for computation -- generally faster even for small models if NData ~16"}, {Name: "NData", Doc: "number of data-parallel items to process in parallel per trial -- works (and is significantly faster) for both CPU and GPU.  Results in an effective mini-batch of learning.  is noisy above 4 for 30 note case."}, {Name: "NThreads", Doc: "number of parallel threads for CPU computation -- 0 = use default"}, {Name: "Run", Doc: "starting run number -- determines the random seed -- runs counts from there -- can do all runs in parallel by launching separate jobs with each run, runs = 1"}, {Name: "NRuns", Doc: "total number of runs to do when running Train"}, {Name: "NEpochs", Doc: "total number of epochs per run"}, {Name: "NTrials", Doc: "total number of trials per epoch.  Should be an even multiple of NData."}, {Name: "PCAInterval", Doc: "how frequently (in epochs) to compute PCA on hidden representations to measure variance?"}, {Name: "TestInterval", Doc: "how often to run through all the test patterns, in terms of training epochs -- can use 0 or -1 for no testing"}}})

var _ = gti.AddType(&gti.Type{Name: "main.LogConfig", IDName: "log-config", Doc: "LogConfig has config parameters related to logging data", Fields: []gti.Field{{Name: "SaveWts", Doc: "if true, save final weights after each run"}, {Name: "Epoch", Doc: "if true, save train epoch log to file, as .epc.tsv typically"}, {Name: "Run", Doc: "if true, save run log to file, as .run.tsv typically"}, {Name: "Trial", Doc: "if true, save train trial log to file, as .trl.tsv typically. May be large."}, {Name: "TestEpoch", Doc: "if true, save testing epoch log to file, as .tst_epc.tsv typically.  In general it is better to copy testing items over to the training epoch log and record there."}, {Name: "TestTrial", Doc: "if true, save testing trial log to file, as .tst_trl.tsv typically. May be large."}, {Name: "NetData", Doc: "if true, save network activation etc data from testing trials, for later viewing in netview"}}})

var _ = gti.AddType(&gti.Type{Name: "main.Config", IDName: "config", Doc: "Config is a standard Sim config -- use as a starting point.", Fields: []gti.Field{{Name: "Includes", Doc: "specify include files here, and after configuration, it contains list of include files added"}, {Name: "GUI", Doc: "open the GUI -- does not automatically run -- if false, then runs automatically and quits"}, {Name: "Debug", Doc: "log debugging information"}, {Name: "Env", Doc: "environment configuration options"}, {Name: "Params", Doc: "parameter related configuration options"}, {Name: "Run", Doc: "sim running related configuration options"}, {Name: "Log", Doc: "data logging related configuration options"}}})

var _ = gti.AddType(&gti.Type{Name: "main.Sim", IDName: "sim", Doc: "Sim encapsulates the entire simulation model, and we define all the\nfunctionality as methods on this struct.  This structure keeps all relevant\nstate information organized and available without having to pass everything around\nas arguments to methods, and provides the core GUI interface (note the view tags\nfor the fields which provide hints to how things should be displayed).", Fields: []gti.Field{{Name: "Config", Doc: "simulation configuration parameters -- set by .toml config file and / or args"}, {Name: "Net", Doc: "the network -- click to view / edit parameters for layers, prjns, etc"}, {Name: "Params", Doc: "all parameter management"}, {Name: "Loops", Doc: "contains looper control loops for running sim"}, {Name: "Stats", Doc: "contains computed statistic values"}, {Name: "Logs", Doc: "Contains all the logs and information about the logs.'"}, {Name: "Envs", Doc: "Environments"}, {Name: "Context", Doc: "axon timing parameters and state"}, {Name: "ViewUpdt", Doc: "netview update parameters"}, {Name: "GUI", Doc: "manages all the gui elements"}, {Name: "RndSeeds", Doc: "a list of random seeds to use for each run"}}})

var _ = gti.AddType(&gti.Type{Name: "main.MusicEnv", IDName: "music-env", Doc: "MusicEnv reads in a midi SMF file and presents it as a sequence of notes.\nSongs with one note at a time per track are currently supported.\nRenders note to a tensor with localist note coding with duplicate units for spiking.", Fields: []gti.Field{{Name: "Nm", Doc: "name of this environment"}, {Name: "Debug", Doc: "emit debugging messages about the music file"}, {Name: "WrapNotes", Doc: "use only 1 octave of 12 notes for everything -- keeps it consistent"}, {Name: "TicksPer", Doc: "number of time ticks per row in table -- note transitions that are faster than this will be lost"}, {Name: "Track", Doc: "which track to process"}, {Name: "Play", Doc: "play output as it steps"}, {Name: "MaxSteps", Doc: "limit song length to given number of steps, if > 0"}, {Name: "DiOffset", Doc: "time offset for data parallel = Song.Rows / (NData+1)"}, {Name: "UnitsPer", Doc: "number of units per localist note value"}, {Name: "NoteRange", Doc: "range of notes in given track"}, {Name: "NNotes", Doc: "number of notes"}, {Name: "Song", Doc: "the song encoded into 200 msec increments, with columns as tracks"}, {Name: "Time", Doc: "current time step"}, {Name: "Note", Doc: "current note, rendered as a 4D tensor with shape:"}, {Name: "NoteIdx", Doc: "current note index"}, {Name: "Player", Doc: "the function for playing midi"}, {Name: "LastNotePlayed", Doc: "for playing notes"}}})
