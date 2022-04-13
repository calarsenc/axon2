// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ra25 runs a simple random-associator four-layer axon network
// that uses the standard supervised learning paradigm to learn
// mappings between 25 random input / output patterns
// defined over 5x5 input / output layers (i.e., 25 units)
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/emer/axon/axon"
	"github.com/emer/emergent/egui"
	"github.com/emer/emergent/elog"
	"github.com/emer/emergent/emer"
	"github.com/emer/emergent/envlp"
	"github.com/emer/emergent/estats"
	"github.com/emer/emergent/etime"
	"github.com/emer/emergent/looper"
	"github.com/emer/emergent/patgen"
	"github.com/emer/emergent/prjn"
	"github.com/emer/etable/agg"
	"github.com/emer/etable/etable"
	"github.com/emer/etable/etensor"
	_ "github.com/emer/etable/etview" // include to get gui views
	"github.com/emer/etable/split"
	"github.com/goki/gi/gi"
	"github.com/goki/gi/gimain"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
	"github.com/goki/mat32"
)

func main() {
	TheSim.New()
	TheSim.Config()
	if len(os.Args) > 1 {
		TheSim.CmdArgs() // simple assumption is that any args = no gui -- could add explicit arg if you want
	} else {
		gimain.Main(func() { // this starts gui -- requires valid OpenGL display connection (e.g., X11)
			guirun()
		})
	}
}

func guirun() {
	TheSim.Init()
	win := TheSim.ConfigGui()
	win.StartEventLoop()
}

// LogPrec is precision for saving float values in logs
const LogPrec = 4

// see params.go for params

// Sim encapsulates the entire simulation model, and we define all the
// functionality as methods on this struct.  This structure keeps all relevant
// state information organized and available without having to pass everything around
// as arguments to methods, and provides the core GUI interface (note the view tags
// for the fields which provide hints to how things should be displayed).
type Sim struct {
	Net          *axon.Network    `view:"no-inline" desc:"the network -- click to view / edit parameters for layers, prjns, etc"`
	Params       emer.Params      `view:"inline" desc:"all parameter management"`
	Tag          string           `desc:"extra tag string to add to any file names output from sim (e.g., weights files, log files, params for run)"`
	Loops        looper.Set       `desc:"contains looper control loops for running sim"`
	Stats        estats.Stats     `desc:"contains computed statistic values"`
	Logs         elog.Logs        `desc:"Contains all the logs and information about the logs.'"`
	ITICycles    int              `desc:"number of cycles between trials"`
	StartRun     int              `desc:"starting run number -- typically 0 but can be set in command args for parallel runs on a cluster"`
	MaxRuns      int              `desc:"maximum number of model runs to perform (starting from StartRun)"`
	MaxEpcs      int              `desc:"maximum number of epochs to run per model run"`
	NZeroStop    int              `desc:"if a positive number, training will stop after this many epochs with zero UnitErr"`
	Pats         *etable.Table    `view:"no-inline" desc:"the training patterns to use"`
	Envs         envlp.Envs       `desc:"Environments"`
	TrainEnv     envlp.FixedTable `desc:"Training environment -- contains everything about iterating over input / output patterns over training"`
	TestEnv      envlp.FixedTable `desc:"Testing environment -- manages iterating over testing"`
	Time         axon.Time        `desc:"axon timing parameters and state"`
	ViewOn       bool             `desc:"whether to update the network view while running"`
	TrainUpdt    axon.TimeScales  `desc:"at what time scale to update the display during training?  Anything longer than Epoch updates at Epoch in this model"`
	TestUpdt     axon.TimeScales  `desc:"at what time scale to update the display during testing?  Anything longer than Epoch updates at Epoch in this model"`
	TestInterval int              `desc:"how often to run through all the test patterns, in terms of training epochs -- can use 0 or -1 for no testing"`
	PCAInterval  int              `desc:"how frequently (in epochs) to compute PCA on hidden representations to measure variance?"`

	GUI         egui.GUI `view:"-" desc:"manages all the gui elements"`
	SaveWts     bool     `view:"-" desc:"for command-line run only, auto-save final weights after each run"`
	NoGui       bool     `view:"-" desc:"if true, runing in no GUI mode"`
	NeedsNewRun bool     `view:"-" desc:"flag to initialize NewRun if last one finished"`
	RndSeeds    []int64  `view:"-" desc:"a list of random seeds to use for each run"`
}

// this registers this Sim Type and gives it properties that e.g.,
// prompt for filename for save methods.
var KiT_Sim = kit.Types.AddType(&Sim{}, SimProps)

// TheSim is the overall state for this simulation
var TheSim Sim

// New creates new blank elements and initializes defaults
func (ss *Sim) New() {
	ss.Net = &axon.Network{}
	ss.Params.Params = ParamSetsMin // ParamSetsDefs
	ss.Params.AddNetwork(ss.Net)
	ss.Params.AddSim(ss)
	ss.Params.AddNetSize()
	ss.Stats.Init()
	ss.Pats = &etable.Table{}
	ss.RndSeeds = make([]int64, 100) // make enough for plenty of runs
	for i := 0; i < 100; i++ {
		ss.RndSeeds[i] = int64(i) + 1 // exclude 0
	}
	ss.ITICycles = 0
	ss.ViewOn = true
	ss.TrainUpdt = axon.AlphaCycle
	ss.TestUpdt = axon.Cycle
	ss.TestInterval = 500
	ss.PCAInterval = 5
	ss.Time.Defaults()
}

////////////////////////////////////////////////////////////////////////////////////////////
// 		Configs

// Config configures all the elements using the standard functions
func (ss *Sim) Config() {
	//ss.ConfigPats()
	ss.OpenPats()
	ss.ConfigEnv()
	ss.ConfigNet(ss.Net)
	ss.ConfigLogs()
	ss.ConfigLoops()
}

func (ss *Sim) ConfigEnv() {
	if ss.MaxRuns == 0 { // allow user override
		ss.MaxRuns = 5
	}
	if ss.MaxEpcs == 0 { // allow user override
		ss.MaxEpcs = 100
		ss.NZeroStop = 5
	}

	ss.Envs.Init()
	ss.TrainEnv.Nm = "TrainEnv"
	ss.TrainEnv.Dsc = "training params and state"
	ss.TrainEnv.Config(etable.NewIdxView(ss.Pats), etime.Train.String())
	ss.TrainEnv.Counter(etime.Run).Max = ss.MaxRuns
	ss.TrainEnv.Counter(etime.Epoch).Max = ss.MaxEpcs
	ss.TrainEnv.Validate()

	ss.TestEnv.Nm = "TestEnv"
	ss.TestEnv.Dsc = "testing params and state"
	ss.TestEnv.Config(etable.NewIdxView(ss.Pats), etime.Test.String())
	ss.TestEnv.Sequential = true
	ss.TestEnv.Counter(etime.Epoch).Max = 1
	ss.TestEnv.Validate()

	// note: to create a train / test split of pats, do this:
	// all := etable.NewIdxView(ss.Pats)
	// splits, _ := split.Permuted(all, []float64{.8, .2}, []string{"Train", "Test"})
	// ss.TrainEnv.Table = splits.Splits[0]
	// ss.TestEnv.Table = splits.Splits[1]

	ss.TrainEnv.Init()
	ss.TestEnv.Init()
	ss.Envs.Add(&ss.TrainEnv, &ss.TestEnv)
}

// Env returns the relevant environment based on Time Mode
func (ss *Sim) Env() envlp.Env {
	return ss.Envs[ss.Time.Mode]
}

func (ss *Sim) ConfigNet(net *axon.Network) {
	ss.Params.AddLayers([]string{"Hidden1", "Hidden2"}, "Hidden")
	ss.Params.SetObject("NetSize")

	net.InitName(net, "RA25")
	inp := net.AddLayer2D("Input", 5, 5, emer.Input)
	hid1 := net.AddLayer2D("Hidden1", ss.Params.LayY("Hidden1", 10), ss.Params.LayX("Hidden1", 10), emer.Hidden)
	hid2 := net.AddLayer2D("Hidden2", ss.Params.LayY("Hidden2", 10), ss.Params.LayX("Hidden2", 10), emer.Hidden)
	out := net.AddLayer2D("Output", 5, 5, emer.Target)

	// use this to position layers relative to each other
	// default is Above, YAlign = Front, XAlign = Center
	// hid2.SetRelPos(relpos.Rel{Rel: relpos.RightOf, Other: "Hidden1", YAlign: relpos.Front, Space: 2})

	// note: see emergent/prjn module for all the options on how to connect
	// NewFull returns a new prjn.Full connectivity pattern
	full := prjn.NewFull()

	net.ConnectLayers(inp, hid1, full, emer.Forward)
	net.BidirConnectLayers(hid1, hid2, full)
	net.BidirConnectLayers(hid2, out, full)

	// net.LateralConnectLayerPrjn(hid1, full, &axon.HebbPrjn{}).SetType(emer.Inhib)

	// note: can set these to do parallel threaded computation across multiple cpus
	// not worth it for this small of a model, but definitely helps for larger ones
	// if Thread {
	// 	hid2.SetThread(1)
	// 	out.SetThread(1)
	// }

	// note: if you wanted to change a layer type from e.g., Target to Compare, do this:
	// out.SetType(emer.Compare)
	// that would mean that the output layer doesn't reflect target values in plus phase
	// and thus removes error-driven learning -- but stats are still computed.

	net.Defaults()
	ss.Params.SetObject("Network")
	err := net.Build()
	if err != nil {
		log.Println(err)
		return
	}
	net.InitWts()
}

////////////////////////////////////////////////////////////////////////////////
// 	    Init, utils

// Init restarts the run, and initializes everything, including network weights
// and resets the epoch log table
func (ss *Sim) Init() {
	ss.InitRndSeed()
	ss.ConfigEnv() // re-config env just in case a different set of patterns was
	// selected or patterns have been modified etc
	ss.GUI.StopNow = false
	// ss.GUI.StopNow = true -- prints messages for params as set
	ss.Params.SetAll()
	// fmt.Println(ss.Params.NetHypers.JSONString())
	ss.NewRun()
	ss.GUI.UpdateNetView()
}

// InitRndSeed initializes the random seed based on current training run number
func (ss *Sim) InitRndSeed() {
	run := ss.TrainEnv.Counter(etime.Run).Cur
	rand.Seed(ss.RndSeeds[run])
}

// NewRndSeed gets a new set of random seeds based on current time -- otherwise uses
// the same random seeds for every run
func (ss *Sim) NewRndSeed() {
	rs := time.Now().UnixNano()
	for i := 0; i < 100; i++ {
		ss.RndSeeds[i] = rs + int64(i)
	}
}

func (ss *Sim) UpdateNetView() {
	if !ss.ViewOn {
		return
	}
	viewUpdt := ss.TrainUpdt
	if ss.Time.Testing {
		viewUpdt = ss.TestUpdt
	}
	switch viewUpdt {
	case axon.Cycle:
		ss.GUI.UpdateNetView()
	case axon.FastSpike:
		if ss.Time.Cycle%10 == 0 {
			ss.GUI.UpdateNetView()
		}
	case axon.GammaCycle:
		if ss.Time.Cycle%25 == 0 {
			ss.GUI.UpdateNetView()
		}
	case axon.AlphaCycle:
		if ss.Time.Cycle%100 == 0 {
			ss.GUI.UpdateNetView()
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// 	    Running the Network, starting bottom-up..

func (ss *Sim) ConfigLoops() {
	trn := looper.NewStackEnv(&ss.TrainEnv)
	tst := looper.NewStackEnv(&ss.TestEnv)
	ss.Loops.AddStack(trn)
	ss.Loops.AddStack(tst)
	axon.ConfigLoopsStd(&ss.Loops, ss.Net, &ss.Time, 150, 50)
	axon.AddCycle0(&ss.Loops, &ss.Time, "Sim:NewRun", func() {
		if ss.NeedsNewRun {
			ss.NewRun()
		}
	})
	axon.AddCycle0(&ss.Loops, &ss.Time, "Sim:ApplyInputs", ss.ApplyInputs)
	// note: AddLoopCycle adds in reverse order of where things end up!
	if !ss.NoGui { // todo: cmdline
		axon.AddLoopCycle(&ss.Loops, "GUI:UpdateNetView", ss.UpdateNetView)
		axon.AddLoopCycle(&ss.Loops, "GUI:RasterRec", ss.RasterRec)
	}
	tst.Loop(etime.Cycle).Main.InsertAfter("Axon:Cycle:Run", "Log:Test:Cycle", func() {
		ss.Log(etime.Test, etime.Cycle)
	})
	axon.AddLoopCycle(&ss.Loops, "Sim:SaveState", func() {
		if ss.Time.Phase == 0 {
			switch ss.Time.Cycle { // save states at beta-frequency -- not used computationally
			case 75:
				ss.Net.ActSt1(&ss.Time)
			case 100:
				ss.Net.ActSt2(&ss.Time)
			}
		}
	})
	axon.AddLoopCycle(&ss.Loops, "Sim:StatCounters", ss.StatCounters) // add last so comes first!

	axon.AddPhaseMain(&ss.Loops, "Sim:TrialStats", func() {
		if ss.Time.Phase == 1 {
			ss.TrialStats()
		}
	})
	if !ss.NoGui {
		// after dwt updated, grab it
		trn.Loop(etime.Phase).End.Add("GUI:UpdateNetView", ss.UpdateNetView)
		tst.Loop(etime.Phase).End.Add("GUI:UpdatePlot", func() {
			ss.GUI.UpdatePlot(etime.Test, etime.Cycle) // make sure always updated at end
		})
	}

	// prepend = before counter is incremented
	trn.Loop(etime.Trial).Main.Prepend("Log:Train:Trial", func() {
		ss.Log(etime.Train, etime.Trial)
	})
	trn.Loop(etime.Epoch).Main.Prepend("Log:Train:Epoch", func() {
		epc := ss.TrainEnv.Counter(etime.Epoch).Cur
		if (ss.TestInterval > 0) && (epc%ss.TestInterval == 0) { // note: epc is *next* so won't trigger first time
			ss.TestAll()
		}
		ss.Log(etime.Train, etime.Epoch)
	})

	trn.Loop(etime.Epoch).Stop.Add("Epoch:NZeroStop", func() bool { // early stopping
		return ss.NZeroStop > 0 && ss.Stats.Int("NZero") >= ss.NZeroStop
	})

	trn.Loop(etime.Run).Main.Prepend("Log:Train:Run", func() {
		ss.Log(etime.Train, etime.Run)
		if ss.SaveWts {
			fnm := ss.WeightsFileName()
			fmt.Printf("Saving Weights to: %s\n", fnm)
			ss.Net.SaveWtsJSON(gi.FileName(fnm))
		}
		ss.NeedsNewRun = true // next step will trigger new init
	})

	tst.Loop(etime.Trial).Main.Add("Log:Test:Trial", func() {
		ss.Log(etime.Test, etime.Trial)
		ss.GUI.NetDataRecord()
	})
	tst.Loop(etime.Epoch).Main.Add("Log:Test:Epoch", func() {
		ss.Log(etime.Test, etime.Epoch)
	})

	fmt.Println(trn.DocString())
	fmt.Println(tst.DocString())
}

/*
// ThetaCyc runs one theta cycle (200 msec) of processing.
// External inputs must have already been applied prior to calling,
// using ApplyExt method on relevant layers (see TrainTrial, TestTrial).
// If train is true, then learning DWt or WtFmDWt calls are made.
// Handles netview updating within scope, and calls TrainStats()
func (ss *Sim) ThetaCyc(train bool) {
	// ss.Win.PollEvents() // this can be used instead of running in a separate goroutine
	viewUpdt := ss.TrainUpdt
	if !train {
		viewUpdt = ss.TestUpdt
	}

	// update prior weight changes at start, so any DWt values remain visible at end
	// you might want to do this less frequently to achieve a mini-batch update
	// in which case, move it out to the TrainTrial method where the relevant
	// counters are being dealt with.
	if train {
		ss.Net.WtFmDWt(&ss.Time)
	}

	minusCyc := 150 // 150
	plusCyc := 50   // 50

	ss.Net.NewState()
	// ss.Time.NewState(train)
	for cyc := 0; cyc < minusCyc; cyc++ { // do the minus phase
		ss.Net.Cycle(&ss.Time)
		// ss.StatCounters(train)
		if !train {
			ss.Log(etime.Test, etime.Cycle)
		}
		if ss.GUI.Active {
			ss.RasterRec()
		}
		ss.Time.CycleInc()
		switch ss.Time.Cycle { // save states at beta-frequency -- not used computationally
		case 75:
			ss.Net.ActSt1(&ss.Time)
		case 100:
			ss.Net.ActSt2(&ss.Time)
		}

		if cyc == minusCyc-1 { // do before view update
			ss.Net.MinusPhase(&ss.Time)
		}
		// if ss.ViewOn {
		// 	ss.UpdateViewTime(train, viewUpdt)
		// }
	}
	ss.Time.NewPhase(true)
	// ss.StatCounters(train)
	if viewUpdt == axon.Phase {
		ss.GUI.UpdateNetView()
	}
	for cyc := 0; cyc < plusCyc; cyc++ { // do the plus phase
		ss.Net.Cycle(&ss.Time)
		// ss.StatCounters(train)
		if !train {
			ss.Log(etime.Test, etime.Cycle)
		}
		if ss.GUI.Active {
			ss.RasterRec()
		}
		ss.Time.CycleInc()

		if cyc == plusCyc-1 { // do before view update
			ss.Net.PlusPhase(&ss.Time)
		}
		// if ss.ViewOn {
		// 	ss.UpdateViewTime(train, viewUpdt)
		// }
	}
	ss.TrialStats()
	// ss.StatCounters(train)

	if viewUpdt == axon.Phase {
		ss.GUI.UpdateNetView()
	}
	ss.Net.InitExt()
	for cyc := 0; cyc < ss.ITICycles; cyc++ { // do the plus phase
		ss.Net.Cycle(&ss.Time)
		// ss.StatCounters(train)
		if !train {
			ss.Log(etime.Test, etime.Cycle)
		}
		if ss.GUI.Active {
			ss.RasterRec()
		}
		ss.Time.CycleInc()
		if ss.ViewOn {
			ss.UpdateViewTime(train, viewUpdt)
		}
	}

	if train {
		ss.Net.DWt(&ss.Time)
	}

	if ss.ViewOn && viewUpdt <= axon.ThetaCycle {
		ss.GUI.UpdateNetView()
	}

	if !train {
		ss.GUI.UpdatePlot(etime.Test, etime.Cycle) // make sure always updated at end
	}
}
*/

// ApplyInputs applies input patterns from given environment.
// It is good practice to have this be a separate method with appropriate
// args so that it can be used for various different contexts
// (training, testing, etc).
func (ss *Sim) ApplyInputs() {
	ev := ss.Env()
	// ss.Net.InitExt() // clear any existing inputs -- not strictly necessary if always
	// going to the same layers, but good practice and cheap anyway

	lays := []string{"Input", "Output"}
	for _, lnm := range lays {
		ly := ss.Net.LayerByName(lnm).(axon.AxonLayer).AsAxon()
		pats := ev.State(ly.Nm)
		if pats != nil {
			ly.ApplyExt(pats)
		}
	}
}

/*
// TrainTrial runs one trial of training using TrainEnv
func (ss *Sim) TrainTrial() {
	if ss.NeedsNewRun {
		ss.NewRun()
	}

	ss.TrainEnv.Step() // the Env encapsulates and manages all counter state

	// Key to query counters FIRST because current state is in NEXT epoch
	// if epoch counter has changed
	epc, _, chg := ss.TrainEnv.Counter(env.Epoch)
	if chg {
		ss.Log(etime.Train, etime.Epoch)
		if ss.ViewOn && ss.TrainUpdt > axon.AlphaCycle {
			ss.GUI.UpdateNetView()
		}
		if (ss.TestInterval > 0) && (epc%ss.TestInterval == 0) { // note: epc is *next* so won't trigger first time
			ss.TestAll()
		}
		if epc >= ss.MaxEpcs || (ss.NZeroStop > 0 && ss.Stats.Int("NZero") >= ss.NZeroStop) {
			// done with training..
			ss.RunEnd()
			if ss.TrainEnv.Run.Incr() { // we are done!
				ss.GUI.StopNow = true
				return
			} else {
				ss.NeedsNewRun = true
				return
			}
		}
	}

	ss.ApplyInputs(&ss.TrainEnv)
	ss.ThetaCyc(true)
	ss.Log(etime.Train, etime.Trial)
}

*/

/*
// RunEnd is called at the end of a run -- save weights, record final log, etc here
func (ss *Sim) RunEnd() {
	ss.Log(etime.Train, etime.Run)
	if ss.SaveWts {
		fnm := ss.WeightsFileName()
		fmt.Printf("Saving Weights to: %s\n", fnm)
		ss.Net.SaveWtsJSON(gi.FileName(fnm))
	}
}
*/

// NewRun intializes a new run of the model, using the TrainEnv.Run counter
// for the new run value
func (ss *Sim) NewRun() {
	ss.InitRndSeed()
	ss.TrainEnv.Init()
	ss.TestEnv.Init()
	ss.Time.Reset()
	ss.Net.InitWts()
	ss.InitStats()
	ss.StatCounters()
	ss.Logs.ResetLog(etime.Train, etime.Epoch)
	ss.Logs.ResetLog(etime.Test, etime.Epoch)
	ss.NeedsNewRun = false
}

/*
// TrainEpoch runs training trials for remainder of this epoch
func (ss *Sim) TrainEpoch() {
	ss.GUI.StopNow = false
	curEpc := ss.TrainEnv.Epoch.Cur
	for {
		ss.TrainTrial()
		if ss.GUI.StopNow || ss.TrainEnv.Epoch.Cur != curEpc {
			break
		}
	}
	ss.Stopped()
}

// TrainRun runs training trials for remainder of run
func (ss *Sim) TrainRun() {
	ss.GUI.StopNow = false
	curRun := ss.TrainEnv.Run.Cur
	for {
		ss.TrainTrial()
		if ss.GUI.StopNow || ss.TrainEnv.Run.Cur != curRun {
			break
		}
	}
	ss.Stopped()
}

// Train runs the full training from this point onward
func (ss *Sim) Train() {
	ss.GUI.StopNow = false
	for {
		ss.TrainTrial()
		if ss.GUI.StopNow {
			break
		}
	}
	ss.Stopped()
}
*/

// Stop tells the sim to stop running
func (ss *Sim) Stop() {
	ss.GUI.StopNow = true
	ss.Loops.StopFlag = true
}

// Stopped is called when a run method stops running -- updates the IsRunning flag and toolbar
func (ss *Sim) Stopped() {
	ss.GUI.Stopped()
}

// SaveWeights saves the network weights -- when called with giv.CallMethod
// it will auto-prompt for filename
func (ss *Sim) SaveWeights(filename gi.FileName) {
	ss.Net.SaveWtsJSON(filename)
}

////////////////////////////////////////////////////////////////////////////////////////////
// Testing

/*
// TestTrial runs one trial of testing -- always sequentially presented inputs
func (ss *Sim) TestTrial(returnOnChg bool) {
	ss.TestEnv.Step()

	// Query counters FIRST
	_, _, chg := ss.TestEnv.Counter(env.Epoch)
	if chg {
		if ss.ViewOn && ss.TestUpdt > axon.AlphaCycle {
			ss.GUI.UpdateNetView()
		}
		ss.Log(etime.Test, etime.Epoch)
		if returnOnChg {
			return
		}
	}

	ss.ApplyInputs(&ss.TestEnv)
	ss.ThetaCyc(false) // !train
	ss.Log(etime.Test, etime.Trial)
	ss.GUI.NetDataRecord()
}
*/

// todo: not clear how to do this

/*
// TestItem tests given item which is at given index in test item list
func (ss *Sim) TestItem(idx int) {
	cur := ss.TestEnv.Trial.Cur
	ss.TestEnv.Trial.Cur = idx
	ss.TestEnv.SetTrialName()
	ss.ApplyInputs()
	ss.ThetaCyc(false) // !train
	ss.TestEnv.Trial.Cur = cur
}
*/

// TestAll runs through the full set of testing items
func (ss *Sim) TestAll() {
	ss.TestEnv.Init()
	tst := ss.Loops.Stack(etime.Test)
	tst.Init()
	tst.Run()
}

/*
// RunTestAll runs through the full set of testing items, has stop running = false at end -- for gui
func (ss *Sim) RunTestAll() {
	ss.GUI.StopNow = false
	ss.TestAll()
	ss.Stopped()
}
*/

/////////////////////////////////////////////////////////////////////////
//   Pats

func (ss *Sim) ConfigPats() {
	dt := ss.Pats
	dt.SetMetaData("name", "TrainPats")
	dt.SetMetaData("desc", "Training patterns")
	sch := etable.Schema{
		{"Name", etensor.STRING, nil, nil},
		{"Input", etensor.FLOAT32, []int{5, 5}, []string{"Y", "X"}},
		{"Output", etensor.FLOAT32, []int{5, 5}, []string{"Y", "X"}},
	}
	dt.SetFromSchema(sch, 25)

	patgen.PermutedBinaryRows(dt.Cols[1], 6, 1, 0)
	patgen.PermutedBinaryRows(dt.Cols[2], 6, 1, 0)
	dt.SaveCSV("random_5x5_25_gen.tsv", etable.Tab, etable.Headers)
}

func (ss *Sim) OpenPats() {
	dt := ss.Pats
	dt.SetMetaData("name", "TrainPats")
	dt.SetMetaData("desc", "Training patterns")
	err := dt.OpenCSV("random_5x5_25.tsv", etable.Tab)
	if err != nil {
		log.Println(err)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////
// 		Stats

// InitStats initializes all the statistics.
// called at start of new run
func (ss *Sim) InitStats() {
	// clear rest just to make Sim look initialized
	ss.Stats.SetFloat("TrlErr", 0.0)
	ss.Stats.SetFloat("TrlUnitErr", 0.0)
	ss.Stats.SetFloat("TrlCosDiff", 0.0)
	ss.Stats.SetInt("FirstZero", -1) // critical to reset to -1
	ss.Stats.SetInt("NZero", 0)
}

// StatCounters saves current counters to Stats, so they are available for logging etc
// Also saves a string rep of them to the GUI, if the GUI is active
func (ss *Sim) StatCounters() {
	ev := ss.Env()
	if ev == nil {
		return
	}
	ev.CtrsToStats(&ss.Stats)
	ss.Stats.SetInt("Cycle", ss.Time.Cycle)
	ss.GUI.NetViewText = ss.Stats.Print([]string{"Run", "Epoch", "Trial", "TrialName", "Cycle", "TrlUnitErr", "TrlErr", "TrlCosDiff"})
}

// TrialStats computes the trial-level statistics.
// Aggregation is done directly from log data.
func (ss *Sim) TrialStats() {
	out := ss.Net.LayerByName("Output").(axon.AxonLayer).AsAxon()

	ss.Stats.SetFloat("TrlCosDiff", float64(out.CosDiff.Cos))
	ss.Stats.SetFloat("TrlUnitErr", out.PctUnitErr())

	if ss.Stats.Float("TrlUnitErr") > 0 {
		ss.Stats.SetFloat("TrlErr", 1)
	} else {
		ss.Stats.SetFloat("TrlErr", 0)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////
// 		Logging

func (ss *Sim) ConfigLogs() {
	ss.ConfigLogItems()
	ss.Logs.CreateTables()
	ss.Logs.SetContext(&ss.Stats, ss.Net)
	// don't plot certain combinations we don't use
	ss.Logs.NoPlot(etime.Train, etime.Cycle)
	ss.Logs.NoPlot(etime.Test, etime.Run)
	// note: Analyze not plotted by default
	ss.Logs.SetMeta(etime.Train, etime.Run, "LegendCol", "Params")
	ss.Stats.ConfigRasters(ss.Net, 200+ss.ITICycles, ss.Net.LayersByClass())
}

// Log is the main logging function, handles special things for different scopes
func (ss *Sim) Log(mode etime.Modes, time etime.Times) {
	dt := ss.Logs.Table(mode, time)
	row := dt.Rows
	switch {
	case mode == etime.Test && time == etime.Epoch:
		ss.LogTestErrors()
	case mode == etime.Train && time == etime.Epoch:
		epc := ss.TrainEnv.Counter(etime.Epoch).Cur
		if (ss.PCAInterval > 0) && ((epc-1)%ss.PCAInterval == 0) { // -1 so runs on first epc
			ss.PCAStats()
		}
	case time == etime.Cycle:
		row = ss.Stats.Int("Cycle")
	case time == etime.Trial:
		row = ss.Stats.Int("Trial")
	}

	ss.Logs.LogRow(mode, time, row) // also logs to file, etc
	if time == etime.Cycle {
		ss.GUI.UpdateCyclePlot(etime.Test, ss.Time.Cycle)
	} else {
		ss.GUI.UpdatePlot(mode, time)
	}

	// post-logging special statistics
	switch {
	case mode == etime.Train && time == etime.Run:
		ss.LogRunStats()
	case mode == etime.Train && time == etime.Trial:
		epc := ss.TrainEnv.Counter(etime.Epoch).Cur
		if (ss.PCAInterval > 0) && (epc%ss.PCAInterval == 0) {
			ss.Log(etime.Analyze, etime.Trial)
		}
	}
}

// LogTestErrors records all errors made across TestTrials, at Test Epoch scope
func (ss *Sim) LogTestErrors() {
	sk := etime.Scope(etime.Test, etime.Trial)
	lt := ss.Logs.TableDetailsScope(sk)
	ix, _ := lt.NamedIdxView("TestErrors")
	ix.Filter(func(et *etable.Table, row int) bool {
		return et.CellFloat("Err", row) > 0 // include error trials
	})
	ss.Logs.MiscTables["TestErrors"] = ix.NewTable()

	allsp := split.All(ix)
	split.Agg(allsp, "SSE", agg.AggSum)
	// note: can add other stats to compute
	ss.Logs.MiscTables["TestErrorStats"] = allsp.AggsToTable(etable.AddAggName)
}

// LogRunStats records stats across all runs, at Train Run scope
func (ss *Sim) LogRunStats() {
	sk := etime.Scope(etime.Train, etime.Run)
	lt := ss.Logs.TableDetailsScope(sk)
	ix, _ := lt.NamedIdxView("RunStats")

	spl := split.GroupBy(ix, []string{"Params"})
	split.Desc(spl, "FirstZero")
	split.Desc(spl, "PctCor")
	ss.Logs.MiscTables["RunStats"] = spl.AggsToTable(etable.AddAggName)
}

// PCAStats computes PCA statistics on recorded hidden activation patterns
// from Analyze, Trial log data
func (ss *Sim) PCAStats() {
	ss.Stats.PCAStats(ss.Logs.IdxView(etime.Analyze, etime.Trial), "ActM", ss.Net.LayersByClass("Hidden", "Target"))
	ss.Logs.ResetLog(etime.Analyze, etime.Trial)
}

// RasterRec updates spike raster record for current Time.Cycle
func (ss *Sim) RasterRec() {
	ss.Stats.RasterRec(ss.Net, ss.Time.Cycle, "Spike")
}

// RunName returns a name for this run that combines Tag and Params -- add this to
// any file names that are saved.
func (ss *Sim) RunName() string {
	rn := ""
	if ss.Tag != "" {
		rn += ss.Tag + "_"
	}
	rn += ss.Params.Name()
	if ss.StartRun > 0 {
		rn += fmt.Sprintf("_%03d", ss.StartRun)
	}
	return rn
}

// RunEpochName returns a string with the run and epoch numbers with leading zeros, suitable
// for using in weights file names.  Uses 3, 5 digits for each.
func (ss *Sim) RunEpochName(run, epc int) string {
	return fmt.Sprintf("%03d_%05d", run, epc)
}

// WeightsFileName returns default current weights file name
func (ss *Sim) WeightsFileName() string {
	return ss.Net.Nm + "_" + ss.RunName() + "_" + ss.RunEpochName(ss.TrainEnv.Counter(etime.Run).Cur, ss.TrainEnv.Counter(etime.Epoch).Cur) + ".wts"
}

// LogFileName returns default log file name
func (ss *Sim) LogFileName(lognm string) string {
	return ss.Net.Nm + "_" + ss.RunName() + "_" + lognm + ".tsv"
}

////////////////////////////////////////////////////////////////////////////////////////////
// 		Gui

// ConfigGui configures the GoGi gui interface for this simulation,
func (ss *Sim) ConfigGui() *gi.Window {
	title := "Leabra Random Associator"
	ss.GUI.MakeWindow(ss, "ra25", title, `This demonstrates a basic Leabra model. See <a href="https://github.com/emer/emergent">emergent on GitHub</a>.</p>`)
	ss.GUI.CycleUpdateInterval = 10
	ss.GUI.NetView.Params.MaxRecs = 300
	ss.GUI.NetView.SetNet(ss.Net)

	ss.GUI.NetView.Scene().Camera.Pose.Pos.Set(0, 1, 2.75) // more "head on" than default which is more "top down"
	ss.GUI.NetView.Scene().Camera.LookAt(mat32.Vec3{0, 0, 0}, mat32.Vec3{0, 1, 0})
	ss.GUI.AddPlots(title, &ss.Logs)

	stb := ss.GUI.TabView.AddNewTab(gi.KiT_Layout, "Spike Rasters").(*gi.Layout)
	stb.Lay = gi.LayoutVert
	stb.SetStretchMax()
	for _, lnm := range ss.Stats.Rasters {
		sr := ss.Stats.F32Tensor("Raster_" + lnm)
		ss.GUI.ConfigRasterGrid(stb, lnm, sr)
	}

	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "Init", Icon: "update",
		Tooltip: "Initialize everything including network weights, and start over.  Also applies current params.",
		Active:  egui.ActiveStopped,
		Func: func() {
			ss.Init()
			ss.GUI.UpdateWindow()
		},
	})
	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "Train",
		Icon:    "run",
		Tooltip: "Starts the network training, picking up from wherever it may have left off.  If not stopped, training will complete the specified number of Runs through the full number of Epochs of training, with testing automatically occuring at the specified interval.",
		Active:  egui.ActiveStopped,
		Func: func() {
			if !ss.GUI.IsRunning {
				ss.GUI.IsRunning = true
				ss.GUI.ToolBar.UpdateActions()
				// go ss.Train()
			}
		},
	})
	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "Stop",
		Icon:    "stop",
		Tooltip: "Interrupts running.  Hitting Train again will pick back up where it left off.",
		Active:  egui.ActiveRunning,
		Func: func() {
			ss.Stop()
		},
	})
	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "Step Trial",
		Icon:    "step-fwd",
		Tooltip: "Advances one training trial at a time.",
		Active:  egui.ActiveStopped,
		Func: func() {
			if !ss.GUI.IsRunning {
				ss.GUI.IsRunning = true
				// ss.TrainTrial()
				ss.GUI.IsRunning = false
				ss.GUI.UpdateWindow()
			}
		},
	})
	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "Step Epoch",
		Icon:    "fast-fwd",
		Tooltip: "Advances one epoch (complete set of training patterns) at a time.",
		Active:  egui.ActiveStopped,
		Func: func() {
			if !ss.GUI.IsRunning {
				ss.GUI.IsRunning = true
				ss.GUI.ToolBar.UpdateActions()
				// go ss.TrainEpoch()
			}
		},
	})
	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "Step Run",
		Icon:    "fast-fwd",
		Tooltip: "Advances one full training Run at a time.",
		Active:  egui.ActiveStopped,
		Func: func() {
			if !ss.GUI.IsRunning {
				ss.GUI.IsRunning = true
				ss.GUI.ToolBar.UpdateActions()
				// go ss.TrainRun()
			}
		},
	})

	////////////////////////////////////////////////
	ss.GUI.ToolBar.AddSeparator("test")
	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "Test Trial",
		Icon:    "fast-fwd",
		Tooltip: "Runs the next testing trial.",
		Active:  egui.ActiveStopped,
		Func: func() {
			if !ss.GUI.IsRunning {
				ss.GUI.IsRunning = true
				// ss.TestTrial(false) // don't return on change -- wrap
				ss.GUI.IsRunning = false
				ss.GUI.UpdateWindow()
			}
		},
	})
	/*
		ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "Test Item",
			Icon:    "step-fwd",
			Tooltip: "Prompts for a specific input pattern name to run, and runs it in testing mode.",
			Active:  egui.ActiveStopped,
			Func: func() {

				gi.StringPromptDialog(ss.GUI.ViewPort, "", "Test Item",
					gi.DlgOpts{Title: "Test Item", Prompt: "Enter the Name of a given input pattern to test (case insensitive, contains given string."},
					ss.GUI.Win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
						dlg := send.(*gi.Dialog)
						if sig == int64(gi.DialogAccepted) {
							val := gi.StringPromptDialogValue(dlg)
							idxs := []int{0} //TODO: //ss.TestEnv.Table.RowsByString("Name", val, etable.Contains, etable.IgnoreCase)
							if len(idxs) == 0 {
								gi.PromptDialog(nil, gi.DlgOpts{Title: "Name Not Found", Prompt: "No patterns found containing: " + val}, gi.AddOk, gi.NoCancel, nil, nil)
							} else {
								if !ss.GUI.IsRunning {
									ss.GUI.IsRunning = true
									fmt.Printf("testing index: %d\n", idxs[0])
									ss.TestItem(idxs[0])
									ss.GUI.IsRunning = false
									ss.GUI.UpdateWindow()
								}
							}
						}
					})
			},
		})
	*/
	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "Test All",
		Icon:    "step-fwd",
		Tooltip: "Prompts for a specific input pattern name to run, and runs it in testing mode.",
		Active:  egui.ActiveStopped,
		Func: func() {
			if !ss.GUI.IsRunning {
				ss.GUI.IsRunning = true
				ss.GUI.ToolBar.UpdateActions()
				// go ss.RunTestAll()
			}
		},
	})

	////////////////////////////////////////////////
	ss.GUI.ToolBar.AddSeparator("log")
	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "Reset RunLog",
		Icon:    "reset",
		Tooltip: "Reset the accumulated log of all Runs, which are tagged with the ParamSet used",
		Active:  egui.ActiveAlways,
		Func: func() {
			ss.Logs.ResetLog(etime.Train, etime.Run)
			ss.GUI.UpdatePlot(etime.Train, etime.Run)
		},
	})
	////////////////////////////////////////////////
	ss.GUI.ToolBar.AddSeparator("misc")
	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "New Seed",
		Icon:    "new",
		Tooltip: "Generate a new initial random seed to get different results.  By default, Init re-establishes the same initial seed every time.",
		Active:  egui.ActiveAlways,
		Func: func() {
			ss.NewRndSeed()
		},
	})
	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "README",
		Icon:    "file-markdown",
		Tooltip: "Opens your browser on the README file that contains instructions for how to run this model.",
		Active:  egui.ActiveAlways,
		Func: func() {
			gi.OpenURL("https://github.com/emer/axon/blob/master/examples/ra25/README.md")
		},
	})
	ss.GUI.FinalizeGUI(false)
	return ss.GUI.Win
}

// These props register Save methods so they can be used
var SimProps = ki.Props{
	"CallMethods": ki.PropSlice{
		{"SaveWeights", ki.Props{
			"desc": "save network weights to file",
			"icon": "file-save",
			"Args": ki.PropSlice{
				{"File Name", ki.Props{
					"ext": ".wts,.wts.gz",
				}},
			},
		}},
	},
}

func (ss *Sim) CmdArgs() {
	ss.NoGui = true
	var nogui bool
	var saveEpcLog bool
	var saveRunLog bool
	var saveNetData bool
	var note string
	flag.StringVar(&ss.Params.ExtraSets, "params", "", "ParamSet name to use -- must be valid name as listed in compiled-in params or loaded params")
	flag.StringVar(&ss.Tag, "tag", "", "extra tag to add to file names saved from this run")
	flag.StringVar(&note, "note", "", "user note -- describe the run params etc")
	flag.IntVar(&ss.StartRun, "run", 0, "starting run number -- determines the random seed -- runs counts from there -- can do all runs in parallel by launching separate jobs with each run, runs = 1")
	flag.IntVar(&ss.MaxRuns, "runs", 10, "number of runs to do (note that MaxEpcs is in paramset)")
	flag.BoolVar(&ss.Params.SetMsg, "setparams", false, "if true, print a record of each parameter that is set")
	flag.BoolVar(&ss.SaveWts, "wts", false, "if true, save final weights after each run")
	flag.BoolVar(&saveEpcLog, "epclog", true, "if true, save train epoch log to file")
	flag.BoolVar(&saveRunLog, "runlog", true, "if true, save run epoch log to file")
	flag.BoolVar(&saveNetData, "netdata", false, "if true, save network activation etc data from testing trials, for later viewing in netview")
	flag.BoolVar(&nogui, "nogui", true, "if not passing any other args and want to run nogui, use nogui")
	flag.Parse()
	ss.Init()

	if note != "" {
		fmt.Printf("note: %s\n", note)
	}
	if ss.Params.ExtraSets != "" {
		fmt.Printf("Using ParamSet: %s\n", ss.Params.ExtraSets)
	}

	if saveEpcLog {
		fnm := ss.LogFileName("epc")
		ss.Logs.SetLogFile(etime.Train, etime.Epoch, fnm)
	}
	if saveRunLog {
		fnm := ss.LogFileName("run")
		ss.Logs.SetLogFile(etime.Train, etime.Run, fnm)
	}
	if saveNetData {
		fmt.Printf("Saving NetView data from testing\n")
		ss.GUI.InitNetData(ss.Net, 200)
	}
	if ss.SaveWts {
		fmt.Printf("Saving final weights per run\n")
	}
	fmt.Printf("Running %d Runs starting at %d\n", ss.MaxRuns, ss.StartRun)
	rc := ss.TrainEnv.Counter(etime.Run)
	rc.Set(ss.StartRun)
	rc.Max = ss.StartRun + ss.MaxRuns
	ss.NewRun()
	ss.Loops.Run(etime.Train)

	ss.Logs.CloseLogFiles()

	if saveNetData {
		ss.GUI.SaveNetData(ss.RunName())
	}
}
