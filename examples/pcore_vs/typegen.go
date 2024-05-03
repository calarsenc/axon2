// Code generated by "core generate -add-types"; DO NOT EDIT.

package main

import (
	"cogentcore.org/core/types"
)

var _ = types.AddType(&types.Type{Name: "main.EnvConfig", IDName: "env-config", Doc: "EnvConfig has config params for environment\nnote: only adding fields for key Env params that matter for both Network and Env\nother params are set via the Env map data mechanism.", Fields: []types.Field{{Name: "Env", Doc: "env parameters -- can set any field/subfield on Env struct, using standard TOML formatting"}, {Name: "ZeroTest", Doc: "test with no ACC activity at all -- params need to prevent gating in this situation too"}}})

var _ = types.AddType(&types.Type{Name: "main.ParamConfig", IDName: "param-config", Doc: "ParamConfig has config parameters related to sim params", Fields: []types.Field{{Name: "Tweak", Doc: "If true, perform automated parameter tweaking for parameters marked Hypers Tweak = log,incr, or [vals]"}, {Name: "Baseline", Doc: "for Tweak, if true, first run a baseline with current default params"}, {Name: "DryRun", Doc: "for Tweak, if true, only print what would be done, don't run"}, {Name: "Network", Doc: "network parameters"}, {Name: "Sheet", Doc: "Extra Param Sheet name(s) to use (space separated if multiple) -- must be valid name as listed in compiled-in params or loaded params"}, {Name: "Tag", Doc: "extra tag to add to file names and logs saved from this run"}, {Name: "Note", Doc: "user note -- describe the run params etc -- like a git commit message for the run"}, {Name: "File", Doc: "Name of the JSON file to input saved parameters from."}, {Name: "SaveAll", Doc: "Save a snapshot of all current param and config settings in a directory named params_<datestamp> (or _good if Good is true), then quit -- useful for comparing to later changes and seeing multiple views of current params"}, {Name: "Good", Doc: "for SaveAll, save to params_good for a known good params state.  This can be done prior to making a new release after all tests are passing -- add results to git to provide a full diff record of all params over time."}}})

var _ = types.AddType(&types.Type{Name: "main.RunConfig", IDName: "run-config", Doc: "RunConfig has config parameters related to running the sim", Fields: []types.Field{{Name: "GPU", Doc: "use the GPU for computation -- generally faster even for small models if NData ~16"}, {Name: "NData", Doc: "number of data-parallel items to process in parallel per trial -- works (and is significantly faster) for both CPU and GPU.  Results in an effective mini-batch of learning."}, {Name: "NThreads", Doc: "number of parallel threads for CPU computation -- 0 = use default"}, {Name: "Run", Doc: "starting run number -- determines the random seed -- runs counts from there -- can do all runs in parallel by launching separate jobs with each run, runs = 1"}, {Name: "NRuns", Doc: "total number of runs to do when running Train"}, {Name: "NEpochs", Doc: "total number of epochs per run"}, {Name: "NTrials", Doc: "total number of trials per epoch.  Should be an even multiple of NData."}, {Name: "NCycles", Doc: "total number of cycles per trial.  needs more time to decide, also for actions."}}})

var _ = types.AddType(&types.Type{Name: "main.LogConfig", IDName: "log-config", Doc: "LogConfig has config parameters related to logging data", Fields: []types.Field{{Name: "SaveWts", Doc: "if true, save final weights after each run"}, {Name: "Epoch", Doc: "if true, save train epoch log to file, as .epc.tsv typically"}, {Name: "Run", Doc: "if true, save run log to file, as .run.tsv typically"}, {Name: "Trial", Doc: "if true, save train trial log to file, as .trl.tsv typically. May be large."}, {Name: "TestEpoch", Doc: "if true, save testing epoch log to file, as .tst_epc.tsv typically.  In general it is better to copy testing items over to the training epoch log and record there."}, {Name: "TestTrial", Doc: "if true, save testing trial log to file, as .tst_trl.tsv typically. May be large."}, {Name: "NetData", Doc: "if true, save network activation etc data from testing trials, for later viewing in netview"}, {Name: "Testing", Doc: "activates testing mode -- records detailed data for Go CI tests (not the same as running test mode on network, via Looper)"}}})

var _ = types.AddType(&types.Type{Name: "main.Config", IDName: "config", Doc: "Config is a standard Sim config -- use as a starting point.", Fields: []types.Field{{Name: "Includes", Doc: "specify include files here, and after configuration, it contains list of include files added"}, {Name: "GUI", Doc: "open the GUI -- does not automatically run -- if false, then runs automatically and quits"}, {Name: "Debug", Doc: "log debugging information"}, {Name: "Env", Doc: "environment configuration options"}, {Name: "Params", Doc: "parameter related configuration options"}, {Name: "Run", Doc: "sim running related configuration options"}, {Name: "Log", Doc: "data logging related configuration options"}}})

var _ = types.AddType(&types.Type{Name: "main.GoNoEnv", IDName: "go-no-env", Doc: "GoNoEnv implements simple Go vs. NoGo input patterns to test BG learning.", Fields: []types.Field{{Name: "Nm", Doc: "name of environment -- Train or Test"}, {Name: "Mode", Doc: "training or testing env?"}, {Name: "Trial", Doc: "trial counter -- set by caller for testing"}, {Name: "ManualValues", Doc: "if true, ACCPos and Neg are set manually for testing specific cases;\ndo not generate random vals for training or auto-increment ACCPos / Neg values during test"}, {Name: "ACCPos", Doc: "activation of ACC positive valence -- drives go"}, {Name: "ACCNeg", Doc: "activation of ACC neg valence -- drives nogo"}, {Name: "PosNegThr", Doc: "threshold on diff between ACCPos - ACCNeg for counting as a Go trial"}, {Name: "RewPredLRate", Doc: "learning rate for reward prediction"}, {Name: "RewPredMin", Doc: "minimum rewpred value"}, {Name: "GatedShould", Doc: "reward value for case where it gated and it should have:\nnominally 1 but can lead to over-learning, RPE would decrease over time"}, {Name: "NoGatedShould", Doc: "reward value for case where it did not gate and it should have:\nin real case, would not get anything for this, but 1 is a cheat to improve perf"}, {Name: "GatedShouldnt", Doc: "reward value for case where it gated and it should not have.  should be -1"}, {Name: "NoGatedShouldnt", Doc: "reward value for case where it did not gate and it should not have:\nshould be 0"}, {Name: "TestInc", Doc: "increment in testing activation for test all"}, {Name: "TestReps", Doc: "number of repetitions per testing level"}, {Name: "NUnitsY", Doc: "number of units, Y"}, {Name: "NUnitsX", Doc: "number of units, X"}, {Name: "NUnits", Doc: "total number of units"}, {Name: "PopCode", Doc: "pop code the values in ACCPos and Neg"}, {Name: "Rand", Doc: "random number generator for the env -- all random calls must use this"}, {Name: "RandSeed", Doc: "random seed"}, {Name: "States", Doc: "named states: ACCPos, ACCNeg"}, {Name: "Should", Doc: "true if Pos - Neg > Thr"}, {Name: "Gated", Doc: "true if model gated on this trial"}, {Name: "Match", Doc: "true if gated == should"}, {Name: "Rew", Doc: "reward based on match between Should vs. Gated"}, {Name: "RewPred", Doc: "reward prediction based on incremental learning: RewPredLRate * (Rew - RewPred)"}, {Name: "RPE", Doc: "reward prediction error: Rew - RewPred"}}})

var _ = types.AddType(&types.Type{Name: "main.Sim", IDName: "sim", Doc: "Sim encapsulates the entire simulation model, and we define all the\nfunctionality as methods on this struct.  This structure keeps all relevant\nstate information organized and available without having to pass everything around\nas arguments to methods, and provides the core GUI interface (note the view tags\nfor the fields which provide hints to how things should be displayed).", Fields: []types.Field{{Name: "Config", Doc: "simulation configuration parameters -- set by .toml config file and / or args"}, {Name: "Net", Doc: "the network -- click to view / edit parameters for layers, paths, etc"}, {Name: "Params", Doc: "all parameter management"}, {Name: "Loops", Doc: "contains looper control loops for running sim"}, {Name: "Stats", Doc: "contains computed statistic values"}, {Name: "Logs", Doc: "Contains all the logs and information about the logs.'"}, {Name: "Envs", Doc: "Environments"}, {Name: "Context", Doc: "axon timing parameters and state"}, {Name: "ViewUpdate", Doc: "netview update parameters"}, {Name: "GUI", Doc: "manages all the gui elements"}, {Name: "RandSeeds", Doc: "a list of random seeds to use for each run"}}})
