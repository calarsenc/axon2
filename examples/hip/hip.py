#!/usr/local/bin/pyleabra -i

# Copyright (c) 2019, The Emergent Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# hip project

from leabra import go, leabra, emer, relpos, eplot, env, agg, patgen, prjn, etable, efile, split, etensor, params, netview, rand, erand, gi, giv, pygiv, pyparams, mat32, hip, evec, simat, metric

import importlib as il  #il.reload(ra25) -- doesn't seem to work for reasons unknown
import io, sys, getopt
from datetime import datetime, timezone

# OuterLoopParams are the parameters to run for outer crossed factor testing
# var OuterLoopParams = []string{"SmallHip", "MedHip"} //, "BigHip"}
OuterLoopParams = go.Slice_string(["MedHip"]) #, "BigHip"}

# InnerLoopParams are the parameters to run for inner crossed factor testing
# var InnerLoopParams = []string{"List020", "List040", "List050", "List060", "List070", "List080"} // , "List100"}
InnerLoopParams = go.Slice_string(["List040", "List080", "List120", "List160", "List200"]) # , "List100"}

# import numpy as np
# import matplotlib
# matplotlib.use('SVG')
# import matplotlib.pyplot as plt
# plt.rcParams['svg.fonttype'] = 'none'  # essential for not rendering fonts as paths

# this will become Sim later.. 
TheSim = 1

# LogPrec is precision for saving float values in logs
LogPrec = 4

# note: we cannot use methods for callbacks from Go -- must be separate functions
# so below are all the callbacks from the GUI toolbar actions

def InitCB(recv, send, sig, data):
    TheSim.Init()
    TheSim.UpdateClassView()
    TheSim.vp.SetNeedsFullRender()

def TrainCB(recv, send, sig, data):
    if not TheSim.IsRunning:
        TheSim.IsRunning = True
        TheSim.ToolBar.UpdateActions()
        TheSim.Train()

def StopCB(recv, send, sig, data):
    TheSim.Stop()

def StepTrialCB(recv, send, sig, data):
    if not TheSim.IsRunning:
        TheSim.IsRunning = True
        TheSim.TrainTrial()
        TheSim.IsRunning = False
        TheSim.UpdateClassView()
        TheSim.vp.SetNeedsFullRender()

def StepEpochCB(recv, send, sig, data):
    if not TheSim.IsRunning:
        TheSim.IsRunning = True
        TheSim.ToolBar.UpdateActions()
        TheSim.TrainEpoch()

def StepRunCB(recv, send, sig, data):
    if not TheSim.IsRunning:
        TheSim.IsRunning = True
        TheSim.ToolBar.UpdateActions()
        TheSim.TrainRun()

def TestTrialCB(recv, send, sig, data):
    if not TheSim.IsRunning:
        TheSim.IsRunning = True
        TheSim.TestTrial(False)
        TheSim.IsRunning = False
        TheSim.UpdateClassView()
        TheSim.vp.SetNeedsFullRender()

def TestItemCB2(recv, send, sig, data):
    win = gi.Window(handle=recv)
    vp = win.WinViewport2D()
    dlg = gi.Dialog(handle=send)
    if sig != gi.DialogAccepted:
        return
    val = gi.StringPromptDialogValue(dlg)
    idxs = TheSim.TestEnv.Table.RowsByString("Name", val, True, True) # contains, ignoreCase
    if len(idxs) == 0:
        gi.PromptDialog(vp, gi.DlgOpts(Title="Name Not Found", Prompt="No patterns found containing: " + val), True, False, go.nil, go.nil)
    else:
        if not TheSim.IsRunning:
            TheSim.IsRunning = True
            print("testing index: %s" % idxs[0])
            TheSim.TestItem(idxs[0])
            TheSim.IsRunning = False
            vp.SetNeedsFullRender()

def TestItemCB(recv, send, sig, data):
    win = gi.Window(handle=recv)
    gi.StringPromptDialog(win.WinViewport2D(), "", "Test Item",
        gi.DlgOpts(Title="Test Item", Prompt="Enter the Name of a given input pattern to test (case insensitive, contains given string."), win, TestItemCB2)

def TestAllCB(recv, send, sig, data):
    if not TheSim.IsRunning:
        TheSim.IsRunning = True
        TheSim.ToolBar.UpdateActions()
        TheSim.RunTestAll()

def ResetRunLogCB(recv, send, sig, data):
    TheSim.RunLog.SetNumRows(0)
    TheSim.RunPlot.Update()

def NewRndSeedCB(recv, send, sig, data):
    TheSim.NewRndSeed()

def ReadmeCB(recv, send, sig, data):
    gi.TheApp.OpenURL("https://github.com/emer/leabra/blob/master/examples/hip/README.md")

def FilterSSE(et, row):
    return etable.Table(handle=et).CellFloat("SSE", row) > 0 # include error trials    

def AggIfGt0(idx, val):
    return val > 0

def AggIfEq0(idx, val):
    return val == 0

def UpdtFuncNotRunning(act):
    act.SetActiveStateUpdt(not TheSim.IsRunning)
    
def UpdtFuncRunning(act):
    act.SetActiveStateUpdt(TheSim.IsRunning)

    
#####################################################    
#     Sim

class HipParams(pygiv.ClassViewObj):
    """
    see def_params.go for the default params, and params.go for user-saved versions
    from the gui.
    """

    def __init__(self):
        super(HipParams, self).__init__()
        self.ECSize = evec.Vec2i()
        self.SetTags("ECSize", 'desc:"size of EC in terms of overall pools (outer dimension)"')
        self.ECPool = evec.Vec2i()
        self.SetTags("ECPool", 'desc:"size of one EC pool"')
        self.CA1Pool = evec.Vec2i()
        self.SetTags("CA1Pool", 'desc:"size of one CA1 pool"')
        self.CA3Size = evec.Vec2i()
        self.SetTags("CA3Size", 'desc:"size of CA3"')
        self.DGRatio = float()
        self.SetTags("DGRatio", 'desc:"size of DG / CA3"')
        self.DGSize = evec.Vec2i()
        self.SetTags("DGSize", 'inactive:"+" desc:"size of DG"')
        self.DGPCon = float()
        self.SetTags("DGPCon", 'desc:"percent connectivity into DG"')
        self.CA3PCon = float()
        self.SetTags("CA3PCon", 'desc:"percent connectivity into CA3"')
        self.MossyPCon = float()
        self.SetTags("MossyPCon", 'desc:"percent connectivity into CA3 from DG"')
        self.ECPctAct = float()
        self.SetTags("ECPctAct", 'desc:"percent activation in EC pool"')
        self.MossyDel = float()
        self.SetTags("MossyDel", 'desc:"delta in mossy effective strength between minus and plus phase"')
        self.MossyDelTest = float()
        self.SetTags("MossyDelTest", 'desc:"delta in mossy strength for testing (relative to base param)"')

    def Update(hp):
        hp.DGSize.X = int(float(hp.CA3Size.X) * hp.DGRatio)
        hp.DGSize.Y = int(float(hp.CA3Size.Y) * hp.DGRatio)

    def Defaults(hp):
        hp.ECSize.Set(2, 3)
        hp.ECPool.Set(7, 7)
        hp.CA1Pool.Set(10, 10)
        hp.CA3Size.Set(20, 20)
        hp.DGRatio = 1.5

        # ratio
        hp.DGPCon = 0.25 # .35 is sig worse, .2 learns faster but AB recall is worse
        hp.CA3PCon = 0.25
        hp.MossyPCon = 0.02 # .02 > .05 > .01 (for small net)
        hp.ECPctAct = 0.2

        hp.MossyDel = 4     # 4 > 2 -- best is 4 del on 4 rel baseline
        hp.MossyDelTest = 3 # for rel = 4: 3 > 2 > 0 > 4 -- 4 is very bad -- need a small amount..



class PatParams(pygiv.ClassViewObj):
    """
    PatParams have the pattern parameters
    """

    def __init__(self):
        super(PatParams, self).__init__()
        self.ListSize = int()
        self.SetTags("ListSize", 'desc:"number of A-B, A-C patterns each"')
        self.MinDiffPct = float()
        self.SetTags("MinDiffPct", 'desc:"minimum difference between item random patterns, as a proportion (0-1) of total active"')
        self.DriftCtxt = bool()
        self.SetTags("DriftCtxt", 'desc:"use drifting context representations -- otherwise does bit flips from prototype"')
        self.CtxtFlipPct = float()
        self.SetTags("CtxtFlipPct", 'desc:"proportion (0-1) of active bits to flip for each context pattern, relative to a prototype, for non-drifting"')
        self.DriftPct = float()
        self.SetTags("DriftPct", 'desc:"percentage of active bits that drift, per step, for drifting context"')

    def Defaults(pp):
        pp.ListSize = 20 # 10 is too small to see issues..
        pp.MinDiffPct = 0.5
        pp.CtxtFlipPct = .25
        pp.DriftPct = .2



class Sim(pygiv.ClassViewObj):
    """
    Sim encapsulates the entire simulation model, and we define all the
    functionality as methods on this struct.  This structure keeps all relevant
    state information organized and available without having to pass everything around
    as arguments to methods, and provides the core GUI interface (note the view tags
    for the fields which provide hints to how things should be displayed).
    """

    def __init__(self):
        super(Sim, self).__init__()
        self.Net = leabra.Network()
        self.SetTags("Net", 'view:"no-inline"')
        self.Hip = HipParams()
        self.SetTags("Hip", 'desc:"hippocampus sizing parameters"')
        self.Pat = PatParams()
        self.SetTags("Pat", 'desc:"parameters for the input patterns"')
        self.PoolVocab = patgen.Vocab()
        self.SetTags("PoolVocab", 'view:"no-inline" desc:"pool patterns vocabulary"')
        self.TrainAB = etable.Table()
        self.SetTags("TrainAB", 'view:"no-inline" desc:"AB training patterns to use"')
        self.TrainAC = etable.Table()
        self.SetTags("TrainAC", 'view:"no-inline" desc:"AC training patterns to use"')
        self.TestAB = etable.Table()
        self.SetTags("TestAB", 'view:"no-inline" desc:"AB testing patterns to use"')
        self.TestAC = etable.Table()
        self.SetTags("TestAC", 'view:"no-inline" desc:"AC testing patterns to use"')
        self.TestLure = etable.Table()
        self.SetTags("TestLure", 'view:"no-inline" desc:"Lure testing patterns to use"')
        self.TrainAll = etable.Table()
        self.SetTags("TrainAll", 'view:"no-inline" desc:"all training patterns -- for pretrain"')
        self.TrnTrlLog = etable.Table()
        self.SetTags("TrnTrlLog", 'view:"no-inline" desc:"training trial-level log data"')
        self.TrnEpcLog = etable.Table()
        self.SetTags("TrnEpcLog", 'view:"no-inline" desc:"training epoch-level log data"')
        self.TstEpcLog = etable.Table()
        self.SetTags("TstEpcLog", 'view:"no-inline" desc:"testing epoch-level log data"')
        self.TstTrlLog = etable.Table()
        self.SetTags("TstTrlLog", 'view:"no-inline" desc:"testing trial-level log data"')
        self.TstCycLog = etable.Table()
        self.SetTags("TstCycLog", 'view:"no-inline" desc:"testing cycle-level log data"')
        self.RunLog = etable.Table()
        self.SetTags("RunLog", 'view:"no-inline" desc:"summary log of each run"')
        self.RunStats = etable.Table()
        self.SetTags("RunStats", 'view:"no-inline" desc:"aggregate stats on all runs"')
        self.TstStats = etable.Table()
        self.SetTags("TstStats", 'view:"no-inline" desc:"testing stats"')
        self.SimMats = {}
        self.SetTags("SimMats", 'view:"no-inline" desc:"similarity matrix results for layers"')
        self.Params = params.Sets()
        self.SetTags("Params", 'view:"no-inline" desc:"full collection of param sets"')
        self.ParamSet = str()
        self.SetTags("ParamSet", 'desc:"which set of *additional* parameters to use -- always applies Base and optionaly this next if set"')
        self.Tag = str()
        self.SetTags("Tag", 'desc:"extra tag string to add to any file names output from sim (e.g., weights files, log files, params)"')
        self.MaxRuns = int(10)
        self.SetTags("MaxRuns", 'desc:"maximum number of model runs to perform"')
        self.MaxEpcs = int(30)
        self.SetTags("MaxEpcs", 'desc:"maximum number of epochs to run per model run"')
        self.PreTrainEpcs = int(5)
        self.SetTags("PreTrainEpcs", 'desc:"number of epochs to run for pretraining"')
        self.NZeroStop = int(1)
        self.SetTags("NZeroStop", 'desc:"if a positive number, training will stop after this many epochs with zero mem errors"')
        self.TrainEnv = env.FixedTable()
        self.SetTags("TrainEnv", 'desc:"Training environment -- contains everything about iterating over input / output patterns over training"')
        self.TestEnv = env.FixedTable()
        self.SetTags("TestEnv", 'desc:"Testing environment -- manages iterating over testing"')
        self.Time = leabra.Time()
        self.SetTags("Time", 'desc:"leabra timing parameters and state"')
        self.ViewOn = True
        self.SetTags("ViewOn", 'desc:"whether to update the network view while running"')
        self.TrainUpdt = leabra.TimeScales.AlphaCycle
        self.SetTags("TrainUpdt", 'desc:"at what time scale to update the display during training?  Anything longer than Epoch updates at Epoch in this model"')
        self.TestUpdt = leabra.TimeScales.AlphaCycle
        self.SetTags("TestUpdt", 'desc:"at what time scale to update the display during testing?  Anything longer than Epoch updates at Epoch in this model"')
        self.TestInterval = int(1)
        self.SetTags("TestInterval", 'desc:"how often to run through all the test patterns, in terms of training epochs -- can use 0 or -1 for no testing"')
        self.MemThr = float(0.34)
        self.SetTags("MemThr", 'desc:"threshold to use for memory test -- if error proportion is below this number, it is scored as a correct trial"')

        # statistics: note use float64 as that is best for etable.Table
        self.TestNm = str()
        self.SetTags("TestNm", 'inactive:"+" desc:"what set of patterns are we currently testing"')
        self.Mem = float()
        self.SetTags("Mem", 'inactive:"+" desc:"whether current trial\'s ECout met memory criterion"')
        self.TrgOnWasOffAll = float()
        self.SetTags("TrgOnWasOffAll", 'inactive:"+" desc:"current trial\'s proportion of bits where target = on but ECout was off ( < 0.5), for all bits"')
        self.TrgOnWasOffCmp = float()
        self.SetTags("TrgOnWasOffCmp", 'inactive:"+" desc:"current trial\'s proportion of bits where target = on but ECout was off ( < 0.5), for only completion bits that were not active in ECin"')
        self.TrgOffWasOn = float()
        self.SetTags("TrgOffWasOn", 'inactive:"+" desc:"current trial\'s proportion of bits where target = off but ECout was on ( > 0.5)"')
        self.TrlSSE = float()
        self.SetTags("TrlSSE", 'inactive:"+" desc:"current trial\'s sum squared error"')
        self.TrlAvgSSE = float()
        self.SetTags("TrlAvgSSE", 'inactive:"+" desc:"current trial\'s average sum squared error"')
        self.TrlCosDiff = float()
        self.SetTags("TrlCosDiff", 'inactive:"+" desc:"current trial\'s cosine difference"')

        self.EpcSSE = float()
        self.SetTags("EpcSSE", 'inactive:"+" desc:"last epoch\'s total sum squared error"')
        self.EpcAvgSSE = float()
        self.SetTags("EpcAvgSSE", 'inactive:"+" desc:"last epoch\'s average sum squared error (average over trials, and over units within layer)"')
        self.EpcPctErr = float()
        self.SetTags("EpcPctErr", 'inactive:"+" desc:"last epoch\'s percent of trials that had SSE > 0 (subject to .5 unit-wise tolerance)"')
        self.EpcPctCor = float()
        self.SetTags("EpcPctCor", 'inactive:"+" desc:"last epoch\'s percent of trials that had SSE == 0 (subject to .5 unit-wise tolerance)"')
        self.EpcCosDiff = float()
        self.SetTags("EpcCosDiff", 'inactive:"+" desc:"last epoch\'s average cosine difference for output layer (a normalized error measure, maximum of 1 when the minus phase exactly matches the plus)"')
        self.EpcPerTrlMSec = float()
        self.SetTags("EpcPerTrlMSec", 'inactive:"+" desc:"how long did the epoch take per trial in wall-clock milliseconds"')
        self.FirstZero = int()
        self.SetTags("FirstZero", 'inactive:"+" desc:"epoch at when Mem err first went to zero"')
        self.NZero = int()
        self.SetTags("NZero", 'inactive:"+" desc:"number of epochs in a row with zero Mem err"')

        # internal state - view:"-"
        self.SumSSE = float()
        self.SetTags("SumSSE", 'view:"-" inactive:"+" desc:"sum to increment as we go through epoch"')
        self.SumAvgSSE = float()
        self.SetTags("SumAvgSSE", 'view:"-" inactive:"+" desc:"sum to increment as we go through epoch"')
        self.SumCosDiff = float()
        self.SetTags("SumCosDiff", 'view:"-" inactive:"+" desc:"sum to increment as we go through epoch"')
        self.CntErr = int()
        self.SetTags("CntErr", 'view:"-" inactive:"+" desc:"sum of errs to increment as we go through epoch"')
        self.Win = 0
        self.SetTags("Win", 'view:"-" desc:"main GUI window"')
        self.NetView = 0
        self.SetTags("NetView", 'view:"-" desc:"the network viewer"')
        self.ToolBar = 0
        self.SetTags("ToolBar", 'view:"-" desc:"the master toolbar"')
        self.TrnTrlPlot = 0
        self.SetTags("TrnTrlPlot", 'view:"-" desc:"the training trial plot"')
        self.TrnEpcPlot = 0
        self.SetTags("TrnEpcPlot", 'view:"-" desc:"the training epoch plot"')
        self.TstEpcPlot = 0
        self.SetTags("TstEpcPlot", 'view:"-" desc:"the testing epoch plot"')
        self.TstTrlPlot = 0
        self.SetTags("TstTrlPlot", 'view:"-" desc:"the test-trial plot"')
        self.TstCycPlot = 0
        self.SetTags("TstCycPlot", 'view:"-" desc:"the test-cycle plot"')
        self.RunPlot = 0
        self.SetTags("RunPlot", 'view:"-" desc:"the run plot"')
        self.RunStatsPlot = 0
        self.SetTags("RunStatsPlot", 'view:"-" desc:"the run stats plot"')
        self.TrnEpcFile = 0
        self.SetTags("TrnEpcFile", 'view:"-" desc:"log file"')
        self.TrnEpcHdrs = False
        self.SetTags("TrnEpcHdrs", 'view:"-" desc:"headers written"')
        self.TstEpcFile = 0
        self.SetTags("TstEpcFile", 'view:"-" desc:"log file"')
        self.TstEpcHdrs = False
        self.SetTags("TstEpcHdrs", 'view:"-" desc:"headers written"')
        self.RunFile = 0
        self.SetTags("RunFile", 'view:"-" desc:"log file"')
        self.TmpVals = go.Slice_float32()
        self.SetTags("TmpVals", 'view:"-" desc:"temp slice for holding values -- prevent mem allocs"')
        self.LayStatNms = go.Slice_string(["ECin", "ECout", "DG", "CA3", "CA1"])
        self.SetTags("LayStatNms", 'view:"-" desc:"names of layers to collect more detailed stats on (avg act, etc)"')
        self.TstNms = go.Slice_string(["AB", "AC", "Lure"])
        self.SetTags("TstNms", 'view:"-" desc:"names of test tables"')
        self.SimMatStats = go.Slice_string(["Within", "Between"])
        self.SetTags("SimMatStats", 'view:"-" desc:"names of sim mat stats"')
        self.TstStatNms = go.Slice_string(["Mem", "TrgOnWasOff", "TrgOffWasOn"])
        self.SetTags("TstStatNms", 'view:"-" desc:"names of test stats"')
        self.ValsTsrs = {}
        self.SetTags("ValsTsrs", 'view:"-" desc:"for holding layer values"')
        self.SaveWts = False
        self.SetTags("SaveWts", 'view:"-" desc:"for command-line run only, auto-save final weights after each run"')
        self.PreTrainWts = ""
        self.SetTags("PreTrainWts", 'view:"-" desc:"name of pretrained wts file"')
        self.NoGui = False
        self.SetTags("NoGui", 'view:"-" desc:"if true, runing in no GUI mode"')
        self.LogSetParams = False
        self.SetTags("LogSetParams", 'view:"-" desc:"if true, print message for all params that are set"')
        self.IsRunning = False
        self.SetTags("IsRunning", 'view:"-" desc:"true if sim is running"')
        self.StopNow = False
        self.SetTags("StopNow", 'view:"-" desc:"flag to stop running"')
        self.NeedsNewRun = False
        self.SetTags("NeedsNewRun", 'view:"-" desc:"flag to initialize NewRun if last one finished"')
        self.RndSeed = int(2)
        self.SetTags("RndSeed", 'view:"-" desc:"the current random seed"')
        self.LastEpcTime = 0
        self.SetTags("LastEpcTime", 'view:"-" desc:"timer for last epoch"')
        self.vp  = 0
        self.SetTags("vp", 'view:"-" desc:"viewport"')

    def InitParams(ss):
        """
        Sets the default set of parameters -- Base is always applied, and others can be optionally
        selected to apply on top of that
        """
        ss.Params.OpenJSON("def.params")
        ss.Defaults()

    def Defaults(ss):
        ss.Hip.Defaults()
        ss.Pat.Defaults()
        ss.Time.CycPerQtr = 25 # note: key param - 25 seems like it is actually fine?
        ss.Update()

    def Update(ss):
        ss.Hip.Update()

    def Config(ss):
        """
        Config configures all the elements using the standard functions
        """
        ss.InitParams()
        ss.ConfigPats()
        ss.ConfigEnv()
        ss.ConfigNet(ss.Net)
        ss.ConfigTrnTrlLog(ss.TrnTrlLog)
        ss.ConfigTrnEpcLog(ss.TrnEpcLog)
        ss.ConfigTstEpcLog(ss.TstEpcLog)
        ss.ConfigTstTrlLog(ss.TstTrlLog)
        ss.ConfigTstCycLog(ss.TstCycLog)
        ss.ConfigRunLog(ss.RunLog)

    def ConfigEnv(ss):
        if ss.MaxRuns == 0: # allow user override
            ss.MaxRuns = 10
        if ss.MaxEpcs == 0: # allow user override
            ss.MaxEpcs = 30
            ss.NZeroStop = 1
            ss.PreTrainEpcs = 5 # seems sufficient?

        ss.TrainEnv.Nm = "TrainEnv"
        ss.TrainEnv.Dsc = "training params and state"
        ss.TrainEnv.Table = etable.NewIdxView(ss.TrainAB)
        ss.TrainEnv.Validate()
        ss.TrainEnv.Run.Max = ss.MaxRuns # note: we are not setting epoch max -- do that manually

        ss.TestEnv.Nm = "TestEnv"
        ss.TestEnv.Dsc = "testing params and state"
        ss.TestEnv.Table = etable.NewIdxView(ss.TestAB)
        ss.TestEnv.Sequential = True
        ss.TestEnv.Validate()

        ss.TrainEnv.Init(0)
        ss.TestEnv.Init(0)

    def SetEnv(ss, trainAC):
        """
        SetEnv select which set of patterns to train on: AB or AC
        """
        if trainAC:
            ss.TrainEnv.Table = etable.NewIdxView(ss.TrainAC)
        else:
            ss.TrainEnv.Table = etable.NewIdxView(ss.TrainAB)
        ss.TrainEnv.Init(0)

    def ConfigNet(ss, net):
        net.InitName(net, "Hip_bench")
        hp = ss.Hip
        inl = net.AddLayer4D("Input", hp.ECSize.Y, hp.ECSize.X, hp.ECPool.Y, hp.ECPool.X, emer.Input)
        ecin = net.AddLayer4D("ECin", hp.ECSize.Y, hp.ECSize.X, hp.ECPool.Y, hp.ECPool.X, emer.Hidden)
        ecout = net.AddLayer4D("ECout", hp.ECSize.Y, hp.ECSize.X, hp.ECPool.Y, hp.ECPool.X, emer.Target)
        ca1 = net.AddLayer4D("CA1", hp.ECSize.Y, hp.ECSize.X, hp.CA1Pool.Y, hp.CA1Pool.X, emer.Hidden)
        dg = net.AddLayer2D("DG", hp.DGSize.Y, hp.DGSize.X, emer.Hidden)
        ca3 = net.AddLayer2D("CA3", hp.CA3Size.Y, hp.CA3Size.X, emer.Hidden)

        ecin.SetClass("EC")
        ecout.SetClass("EC")

        ecin.SetRelPos(relpos.Rel(Rel= relpos.RightOf, Other= "Input", YAlign= relpos.Front, Space= 2))
        ecout.SetRelPos(relpos.Rel(Rel= relpos.RightOf, Other= "ECin", YAlign= relpos.Front, Space= 2))
        dg.SetRelPos(relpos.Rel(Rel= relpos.Above, Other= "Input", YAlign= relpos.Front, XAlign= relpos.Left, Space= 0))
        ca3.SetRelPos(relpos.Rel(Rel= relpos.Above, Other= "DG", YAlign= relpos.Front, XAlign= relpos.Left, Space= 0))
        ca1.SetRelPos(relpos.Rel(Rel= relpos.RightOf, Other= "CA3", YAlign= relpos.Front, Space= 2))

        onetoone = prjn.NewOneToOne()
        pool1to1 = prjn.NewPoolOneToOne()
        full = prjn.NewFull()

        net.ConnectLayers(inl, ecin, onetoone, emer.Forward)
        net.ConnectLayers(ecout, ecin, onetoone, emer.Back)

        # EC <-> CA1 encoder pathways
        pj = net.ConnectLayersPrjn(ecin, ca1, pool1to1, emer.Forward, hip.EcCa1Prjn())
        pj.SetClass("EcCa1Prjn")
        pj = net.ConnectLayersPrjn(ca1, ecout, pool1to1, emer.Forward, hip.EcCa1Prjn())
        pj.SetClass("EcCa1Prjn")
        pj = net.ConnectLayersPrjn(ecout, ca1, pool1to1, emer.Back, hip.EcCa1Prjn())
        pj.SetClass("EcCa1Prjn")

        # Perforant pathway
        ppathDG = prjn.NewUnifRnd()
        ppathDG.PCon = hp.DGPCon
        ppathCA3 = prjn.NewUnifRnd()
        ppathCA3.PCon = hp.CA3PCon

        pj = net.ConnectLayersPrjn(ecin, dg, ppathDG, emer.Forward, hip.CHLPrjn())
        pj.SetClass("HippoCHL")

        if True: # toggle for bcm vs. ppath
            pj = net.ConnectLayersPrjn(ecin, ca3, ppathCA3, emer.Forward, hip.EcCa1Prjn())
            pj.SetClass("PPath")
            pj = net.ConnectLayersPrjn(ca3, ca3, full, emer.Lateral, hip.EcCa1Prjn())
            pj.SetClass("PPath")
        else:
            # so far, this is sig worse, even with error-driven MinusQ1 case (which is better than off)
            pj = net.ConnectLayersPrjn(ecin, ca3, ppathCA3, emer.Forward, hip.CHLPrjn())
            pj.SetClass("PPath")
            pj = net.ConnectLayersPrjn(ca3, ca3, full, emer.Lateral, hip.CHLPrjn())
            pj.SetClass("PPath")

        # always use this for now:
        if True:
            pj = net.ConnectLayersPrjn(ca3, ca1, full, emer.Forward, hip.CHLPrjn())
            pj.SetClass("HippoCHL")
        else:
            # note: this requires lrate = 1.0 or maybe 1.2, doesn't work *nearly* as well
            pj = net.ConnectLayers(ca3, ca1, full, emer.Forward) # default con
            # pj.SetClass("HippoCHL")

        # Mossy fibers
        mossy = prjn.NewUnifRnd()
        mossy.PCon = hp.MossyPCon
        pj = net.ConnectLayersPrjn(dg, ca3, mossy, emer.Forward, hip.CHLPrjn()) # no learning
        pj.SetClass("HippoCHL")

        # using 4 threads total (rest on 0)
        dg.SetThread(1)
        ca3.SetThread(2)
        ca1.SetThread(3) # this has the most

        # note: if you wanted to change a layer type from e.g., Target to Compare, do this:
        # outLay.SetType(emer.Compare)
        # that would mean that the output layer doesn't reflect target values in plus phase
        # and thus removes error-driven learning -- but stats are still computed.

        net.Defaults()
        ss.SetParams("Network", ss.LogSetParams) # only set Network params
        net.Build()
        net.InitWts()

    def ReConfigNet(ss):
        ss.ConfigPats()
        ss.Net = leabra.Network() # start over with new network
        ss.ConfigNet(ss.Net)
        if ss.NetView != 0:
            ss.NetView.SetNet(ss.Net)
            ss.NetView.Update() # issue #41 closed


    def Init(ss):
        """
        Init restarts the run, and initializes everything, including network weights
        and resets the epoch log table
        """
        rand.Seed(ss.RndSeed)
        ss.SetParams("", ss.LogSetParams)
        ss.ReConfigNet()
        ss.ConfigEnv()
        # selected or patterns have been modified etc
        ss.StopNow = False
        ss.NewRun()
        ss.UpdateView(True)

    def NewRndSeed(ss):
        """
        NewRndSeed gets a new random seed based on current time -- otherwise uses
        the same random seed for every run
        """
        ss.RndSeed = int(datetime.now(timezone.utc).timestamp())

    def Counters(ss, train):
        """
        Counters returns a string of the current counter state
        use tabs to achieve a reasonable formatting overall
        and add a few tabs at the end to allow for expansion..
        """
        if train:
            return "Run:\t%d\tEpoch:\t%d\tTrial:\t%d\tCycle:\t%d\tName:\t%s\t\t\t" % (ss.TrainEnv.Run.Cur, ss.TrainEnv.Epoch.Cur, ss.TrainEnv.Trial.Cur, ss.Time.Cycle, ss.TrainEnv.TrialName.Cur)
        else:
            return "Run:\t%d\tEpoch:\t%d\tTrial:\t%d\tCycle:\t%d\tName:\t%s\t\t\t" % (ss.TrainEnv.Run.Cur, ss.TrainEnv.Epoch.Cur, ss.TestEnv.Trial.Cur, ss.Time.Cycle, ss.TestEnv.TrialName.Cur)

    def UpdateView(ss, train):
        if ss.NetView != 0 and ss.NetView.IsVisible():
            ss.NetView.Record(ss.Counters(train))
            ss.NetView.GoUpdate()

    def AlphaCyc(ss, train):
        """
        AlphaCyc runs one alpha-cycle (100 msec, 4 quarters)             of processing.
        External inputs must have already been applied prior to calling,
        using ApplyExt method on relevant layers (see TrainTrial, TestTrial).

        If train is true, then learning DWt or WtFmDWt calls are made.
        Handles netview updating within scope of AlphaCycle
        """

        if ss.Win != 0:
            ss.Win.PollEvents() # this is essential for GUI responsiveness while running
        viewUpdt = ss.TrainUpdt.value
        if not train:
            viewUpdt = ss.TestUpdt.value

        if train:
            ss.Net.WtFmDWt()

        ca1 = leabra.Layer(ss.Net.LayerByName("CA1"))
        ca3 = leabra.Layer(ss.Net.LayerByName("CA3"))
        ecin = leabra.Layer(ss.Net.LayerByName("ECin"))
        ecout = leabra.Layer(ss.Net.LayerByName("ECout"))
        ca1FmECin = hip.EcCa1Prjn(ca1.RcvPrjns.SendName("ECin"))
        ca1FmCa3 = hip.CHLPrjn(ca1.RcvPrjns.SendName("CA3"))
        ca3FmDg = leabra.Prjn(ca3.RcvPrjns.SendName("DG"))

        # First Quarter: CA1 is driven by ECin, not by CA3 recall
        # (which is not really active yet anyway)
        ca1FmECin.WtScale.Abs = 1
        ca1FmCa3.WtScale.Abs = 0

        dgwtscale = ca3FmDg.WtScale.Rel
        ca3FmDg.WtScale.Rel = dgwtscale - ss.Hip.MossyDel

        if train:
            ecout.SetType(emer.Target) # clamp a plus phase during testing
        else:
            ecout.SetType(emer.Compare) # don't clamp

        ecout.UpdateExtFlags() # call this after updating type

        ss.Net.AlphaCycInit()
        ss.Time.AlphaCycStart()
        for qtr in range(4):
            for cyc in range(ss.Time.CycPerQtr):
                ss.Net.Cycle(ss.Time)
                if not train:
                    ss.LogTstCyc(ss.TstCycLog, ss.Time.Cycle)
                ss.Time.CycleInc()
                if ss.ViewOn:
                    if viewUpdt == leabra.Cycle:
                        if cyc != ss.Time.CycPerQtr-1: # will be updated by quarter
                            ss.UpdateView(train)
                    if viewUpdt == leabra.FastSpike:
                        if (cyc+1)%10 == 0:
                            ss.UpdateView(train)
            if qtr == 1: # Second, Third Quarters: CA1 is driven by CA3 recall
                ca1FmECin.WtScale.Abs = 0
                ca1FmCa3.WtScale.Abs = 1
                if train:
                    ca3FmDg.WtScale.Rel = dgwtscale
                else:
                    ca3FmDg.WtScale.Rel = dgwtscale - ss.Hip.MossyDelTest # testing

                ss.Net.GScaleFmAvgAct() # update computed scaling factors
                ss.Net.InitGInc()       # scaling params change, so need to recompute all netins
            if qtr == 3: # Fourth Quarter: CA1 back to ECin drive only
                ca1FmECin.WtScale.Abs = 1
                ca1FmCa3.WtScale.Abs = 0
                ss.Net.GScaleFmAvgAct() # update computed scaling factors
                ss.Net.InitGInc()       # scaling params change, so need to recompute all netins
                if train:               # clamp ECout from ECin
                    ecin.UnitVals(ss.TmpVals, "Act") # note: could use input instead -- not much diff
                    ecout.ApplyExt1D32(ss.TmpVals)
            ss.Net.QuarterFinal(ss.Time)
            if qtr+1 == 3:
                ss.MemStats(train) # must come after QuarterFinal

            ss.Time.QuarterInc()
            if ss.ViewOn:
                if viewUpdt <= leabra.Quarter:
                    ss.UpdateView(train)
                if viewUpdt == leabra.Phase:
                    if qtr >= 2:
                        ss.UpdateView(train)

        ca3FmDg.WtScale.Rel = dgwtscale # restore
        ca1FmCa3.WtScale.Abs = 1

        if train:
            ss.Net.DWt()
        if ss.ViewOn and viewUpdt == leabra.AlphaCycle:
            ss.UpdateView(train)
        if not train:
            if ss.TstCycPlot != 0:
                ss.TstCycPlot.GoUpdate() # make sure up-to-date at end


    def ApplyInputs(ss, en):
        """
        ApplyInputs applies input patterns from given envirbonment.
        It is good practice to have this be a separate method with appropriate
        args so that it can be used for various different contexts
        (training, testing, etc).
        """
        ss.Net.InitExt()

        lays = go.Slice_string(["Input", "ECout"])
        for lnm in lays :
            ly = leabra.Layer(ss.Net.LayerByName(lnm))
            pats = en.State(ly.Nm)
            if pats != 0:
                ly.ApplyExt(pats)

    def TrainTrial(ss):
        """
        TrainTrial runs one trial of training using TrainEnv
        """
        if ss.NeedsNewRun:
            ss.NewRun()

        ss.TrainEnv.Step()

        # Key to query counters FIRST because current state is in NEXT epoch
        # if epoch counter has changed
        epc = env.CounterCur(ss.TrainEnv, env.Epoch)
        chg = env.CounterChg(ss.TrainEnv, env.Epoch)
        if chg:
            ss.LogTrnEpc(ss.TrnEpcLog)
            if ss.ViewOn and ss.TrainUpdt.value > leabra.AlphaCycle:
                ss.UpdateView(True)
            if ss.TestInterval > 0 and epc%ss.TestInterval == 0: # note: epc is *next* so won't trigger first time
                ss.TestAll()
            learned = (ss.NZeroStop > 0 and ss.NZero >= ss.NZeroStop)
            if ss.TrainEnv.Table.Table.MetaData["name"] == "TrainAB" and (learned or epc == ss.MaxEpcs/2):
                ss.TrainEnv.Table = etable.NewIdxView(ss.TrainAC)
                learned = False
            if learned or epc >= ss.MaxEpcs: # done with training..
                ss.RunEnd()
                if ss.TrainEnv.Run.Incr(): # we are done!
                    ss.StopNow = True
                    return
                else:
                    ss.NeedsNewRun = True
                    return

        ss.ApplyInputs(ss.TrainEnv)
        ss.AlphaCyc(True)   # train
        ss.TrialStats(True) # accumulate
        ss.LogTrnTrl(ss.TrnTrlLog)

    def PreTrainTrial(ss):
        """
        PreTrainTrial runs one trial of pretraining using TrainEnv
        """
        if ss.NeedsNewRun:
            ss.NewRun()

        ss.TrainEnv.Step()

        # Key to query counters FIRST because current state is in NEXT epoch
        # if epoch counter has changed
        epc = env.CounterCur(ss.TrainEnv, env.Epoch)
        chg = env.CounterChg(ss.TrainEnv, env.Epoch)
        if chg:
            ss.LogTrnEpc(ss.TrnEpcLog)
            if ss.ViewOn and ss.TrainUpdt.value > leabra.AlphaCycle:
                ss.UpdateView(True)
            if epc >= ss.PreTrainEpcs: # done with training..
                ss.StopNow = True
                return

        ss.ApplyInputs(ss.TrainEnv)
        ss.AlphaCyc(True)   # train
        ss.TrialStats(True) # accumulate
        ss.LogTrnTrl(ss.TrnTrlLog)

    def RunEnd(ss):
        """
        RunEnd is called at the end of a run -- save weights, record final log, etc here
        """
        ss.LogRun(ss.RunLog)
        if ss.SaveWts:
            fnm = ss.WeightsFileName()
            print("Saving Weights to: %s\n" % fnm)
            ss.Net.SaveWtsJSON(gi.FileName(fnm))

    def NewRun(ss):
        """
        NewRun intializes a new run of the model, using the TrainEnv.Run counter
        for the new run value
        """
        run = ss.TrainEnv.Run.Cur
        ss.TrainEnv.Table = etable.NewIdxView(ss.TrainAB)
        ss.TrainEnv.Init(run)
        ss.TestEnv.Init(run)
        ss.Time.Reset()
        ss.Net.InitWts()
        ss.LoadPretrainedWts()
        ss.InitStats()
        ss.TrnTrlLog.SetNumRows(0)
        ss.TrnEpcLog.SetNumRows(0)
        ss.TstEpcLog.SetNumRows(0)
        ss.NeedsNewRun = False

    def LoadPretrainedWts(ss):
        if ss.PreTrainWts == "":
            return False
        ss.Net.OpenWtsJSON(ss.PreTrainWts)
        return True

    def InitStats(ss):
        """
        InitStats initializes all the statistics, especially important for the
        cumulative epoch stats -- called at start of new run
        """

        ss.SumSSE = 0
        ss.SumAvgSSE = 0
        ss.SumCosDiff = 0
        ss.CntErr = 0
        ss.FirstZero = -1
        ss.NZero = 0

        ss.Mem = 0
        ss.TrgOnWasOffAll = 0
        ss.TrgOnWasOffCmp = 0
        ss.TrgOffWasOn = 0
        ss.TrlSSE = 0
        ss.TrlAvgSSE = 0
        ss.EpcSSE = 0
        ss.EpcAvgSSE = 0
        ss.EpcPctErr = 0
        ss.EpcCosDiff = 0

    def MemStats(ss, train):
        """
        MemStats computes ActM vs. Target on ECout with binary counts
        must be called at end of 3rd quarter so that Target values are
        for the entire full pattern as opposed to the plus-phase target
        values clamped from ECin activations
        """
        ecout = leabra.Layer(ss.Net.LayerByName("ECout"))
        ecin = leabra.Layer(ss.Net.LayerByName("ECin"))
        nn = ecout.Shape().Len()
        trgOnWasOffAll = 0.0
        trgOnWasOffCmp = 0.0
        trgOffWasOn = 0.0 # should have been off
        cmpN = 0.0        # completion target
        trgOnN = 0.0
        trgOffN = 0.0
        actMi = ecout.UnitVarIdx("ActM")
        targi = ecout.UnitVarIdx("Target")
        actQ1i = ecout.UnitVarIdx("ActQ1")
        for ni in range(nn):
            actm = ecout.UnitVal1D(actMi, ni)
            trg = ecout.UnitVal1D(targi, ni) # full pattern target
            inact = ecin.UnitVal1D(actQ1i, ni)
            if trg < 0.5: # trgOff
                trgOffN += 1
                if actm > 0.5:
                    trgOffWasOn += 1
            else: # trgOn
                trgOnN += 1
                if inact < 0.5: # missing in ECin -- completion target
                    cmpN += 1
                    if actm < 0.5:
                        trgOnWasOffAll += 1
                        trgOnWasOffCmp += 1
                else:
                    if actm < 0.5:
                        trgOnWasOffAll += 1
        trgOnWasOffAll /= trgOnN
        trgOffWasOn /= trgOffN
        if train: # no cmp
            if trgOnWasOffAll < ss.MemThr and trgOffWasOn < ss.MemThr:
                ss.Mem = 1
            else:
                ss.Mem = 0
        else: # test
            if cmpN > 0: # should be
                trgOnWasOffCmp /= cmpN
                if trgOnWasOffCmp < ss.MemThr and trgOffWasOn < ss.MemThr:
                    ss.Mem = 1
                else:
                    ss.Mem = 0
        ss.TrgOnWasOffAll = trgOnWasOffAll
        ss.TrgOnWasOffCmp = trgOnWasOffCmp
        ss.TrgOffWasOn = trgOffWasOn

    def TrialStats(ss, accum):
        """
        TrialStats computes the trial-level statistics and adds them to the epoch accumulators if
        accum is true.  Note that we're accumulating stats here on the Sim side so the
        core algorithm side remains as simple as possible, and doesn't need to worry about
        different time-scales over which stats could be accumulated etc.
        You can also aggregate directly from log data, as is done for testing stats
        """
        outLay = leabra.Layer(ss.Net.LayerByName("ECout"))
        ss.TrlCosDiff = float(outLay.CosDiff.Cos)
        ss.TrlSSE = outLay.SSE(0.5) # 0.5 = per-unit tolerance -- right side of .5
        ss.TrlAvgSSE = ss.TrlSSE / len(outLay.Neurons)
        if accum:
            ss.SumSSE += ss.TrlSSE
            ss.SumAvgSSE += ss.TrlAvgSSE
            ss.SumCosDiff += ss.TrlCosDiff
            if ss.TrlSSE != 0:
                ss.CntErr += 1
        return

    def TrainEpoch(ss):
        """
        TrainEpoch runs training trials for remainder of this epoch
        """
        ss.StopNow = False
        curEpc = ss.TrainEnv.Epoch.Cur
        while True:
            ss.TrainTrial()
            if ss.StopNow or ss.TrainEnv.Epoch.Cur != curEpc:
                break
        ss.Stopped()

    def TrainRun(ss):
        """
        TrainRun runs training trials for remainder of run
        """
        ss.StopNow = False
        curRun = ss.TrainEnv.Run.Cur
        while True:
            ss.TrainTrial()
            if ss.StopNow or ss.TrainEnv.Run.Cur != curRun:
                break
        ss.Stopped()

    def Train(ss):
        """
        Train runs the full training from this point onward
        """
        ss.StopNow = False
        while True:
            ss.TrainTrial()
            if ss.StopNow:
                break
        ss.Stopped()

    def Stop(ss):
        """
        Stop tells the sim to stop running
        """
        ss.StopNow = True

    def Stopped(ss):
        """
        Stopped is called when a run method stops running -- updates the IsRunning flag and toolbar
        """
        ss.IsRunning = False
        if ss.Win != 0:
            vp = ss.Win.WinViewport2D()
            if ss.ToolBar != 0:
                ss.ToolBar.UpdateActions()
            vp.SetNeedsFullRender()
            ss.UpdateClassView()

    def SaveWeights(ss, filename):
        """
        SaveWeights saves the network weights -- when called with giv.CallMethod
        it will auto-prompt for filename
        """
        ss.Net.SaveWtsJSON(filename)

    def SetDgCa3Off(ss, net, off):
        """
        SetDgCa3Off sets the DG and CA3 layers off (or on)
        """
        ca3 = leabra.Layer(net.LayerByName("CA3"))
        dg = leabra.Layer(net.LayerByName("DG"))
        ca3.Off = off
        dg.Off = off

    def PreTrain(ss):
        """
        PreTrain runs pre-training, saves weights to PreTrainWts
        """
        ss.SetDgCa3Off(ss.Net, True)
        ss.TrainEnv.Table = etable.NewIdxView(ss.TrainAll)

        ss.StopNow = False
        curRun = ss.TrainEnv.Run.Cur
        while True:
            ss.PreTrainTrial()
            if ss.StopNow or ss.TrainEnv.Run.Cur != curRun:
                break
        ss.PreTrainWts = "tmp_pretrained_wts.wts"
        ss.Net.SaveWtsJSON(ss.PreTrainWts)
        ss.TrainEnv.Table = etable.NewIdxView(ss.TrainAB)
        ss.SetDgCa3Off(ss.Net, False)
        ss.Stopped()

    def TestTrial(ss, returnOnChg):
        """
        TestTrial runs one trial of testing -- always sequentially presented inputs
        """
        ss.TestEnv.Step()

        chg = env.CounterChg(ss.TestEnv, env.Epoch)
        if chg:
            if ss.ViewOn and ss.TestUpdt.value > leabra.AlphaCycle:
                ss.UpdateView(False)
            if returnOnChg:
                return

        ss.ApplyInputs(ss.TestEnv)
        ss.AlphaCyc(False)
        ss.TrialStats(False)
        ss.LogTstTrl(ss.TstTrlLog)

    def TestItem(ss, idx):
        """
        TestItem tests given item which is at given index in test item list
        """
        cur = ss.TestEnv.Trial.Cur
        ss.TestEnv.Trial.Cur = idx
        ss.TestEnv.SetTrialName()
        ss.ApplyInputs(ss.TestEnv)
        ss.AlphaCyc(False)
        ss.TrialStats(False)
        ss.TestEnv.Trial.Cur = cur

    def TestAll(ss):
        """
        TestAll runs through the full set of testing items
        """
        ss.TestNm = "AB"
        ss.TestEnv.Table = etable.NewIdxView(ss.TestAB)
        ss.TestEnv.Init(ss.TrainEnv.Run.Cur)
        while True:
            ss.TestTrial(True)
            chg = env.CounterChg(ss.TestEnv, env.Epoch)
            if chg or ss.StopNow:
                break
        if not ss.StopNow:
            ss.TestNm = "AC"
            ss.TestEnv.Table = etable.NewIdxView(ss.TestAC)
            ss.TestEnv.Init(ss.TrainEnv.Run.Cur)
            while True:
                ss.TestTrial(True)
                chg = env.CounterChg(ss.TestEnv, env.Epoch)
                if chg or ss.StopNow:
                    break
            if not ss.StopNow:
                ss.TestNm = "Lure"
                ss.TestEnv.Table = etable.NewIdxView(ss.TestLure)
                ss.TestEnv.Init(ss.TrainEnv.Run.Cur)
                while True:
                    ss.TestTrial(True)
                    chg = env.CounterChg(ss.TestEnv, env.Epoch)
                    if chg or ss.StopNow:
                        break

        ss.LogTstEpc(ss.TstEpcLog)

    def RunTestAll(ss):
        """
        RunTestAll runs through the full set of testing items, has stop running = false at end -- for gui
        """
        ss.StopNow = False
        ss.TestAll()
        ss.Stopped()

    def ParamsName(ss):
        """
        ParamsName returns name of current set of parameters
        """
        if ss.ParamSet == "":
            return "Base"
        return ss.ParamSet

    def SetParams(ss, sheet, setMsg):
        """
        SetParams sets the params for "Base" and then current ParamSet.
        If sheet is empty, then it applies all avail sheets (e.g., Network, Sim)
        otherwise just the named sheet
        if setMsg = true then we output a message for each param that was set.
        """
        if sheet == "":

            ss.Params.ValidateSheets(go.Slice_string(["Network", "Sim", "Hip", "Pat"]))
        ss.SetParamsSet("Base", sheet, setMsg)
        if ss.ParamSet != "" and ss.ParamSet != "Base":
            sps = ss.ParamSet.split()
            for ps in sps:
                ss.SetParamsSet(ps, sheet, setMsg)
                
    def SetParamsSet(ss, setNm, sheet, setMsg):
        """
        SetParamsSet sets the params for given params.Set name.
        If sheet is empty, then it applies all avail sheets (e.g., Network, Sim)
        otherwise just the named sheet
        if setMsg = true then we output a message for each param that was set.
        """
        pset = ss.Params.SetByNameTry(setNm)
        if sheet == "" or sheet == "Network":
            if "Network" in pset.Sheets:
                netp = pset.SheetByNameTry("Network")
                ss.Net.ApplyParams(netp, setMsg)
                
        if sheet == "" or sheet == "Sim":
            if "Sim" in pset.Sheets:
                simp= pset.SheetByNameTry("Sim")
                pyparams.ApplyParams(ss, simp, setMsg)

        if sheet == "" or sheet == "Hip":
            if "Hip" in pset.Sheets:
                simp = pset.SheetByNameTry("Hip")
                pyparams.ApplyParams(ss.Hip, simp, setMsg)

        if sheet == "" or sheet == "Pat":
            if "Pat" in pset.Sheets:
                simp = pset.SheetByNameTry("Pat")
                pyparams.ApplyParams(ss.Pat, simp, setMsg)

    def OpenPat(ss, dt, fname, name, desc):
        err = dt.OpenCSV(gi.FileName(fname), etable.Tab)
        if err != 0:
            log.Println(err)
            return
        dt.SetMetaData("name", name)
        dt.SetMetaData("desc", desc)

    def ConfigPats(ss):
        hp = ss.Hip
        plY = hp.ECPool.Y
        plX = hp.ECPool.X
        npats = ss.Pat.ListSize
        pctAct = hp.ECPctAct
        minDiff = ss.Pat.MinDiffPct
        nOn = patgen.NFmPct(pctAct, plY*plX)
        ctxtflip = patgen.NFmPct(ss.Pat.CtxtFlipPct, nOn)
        patgen.AddVocabEmpty(ss.PoolVocab, "empty", npats, plY, plX)
        patgen.AddVocabPermutedBinary(ss.PoolVocab, "A", npats, plY, plX, pctAct, minDiff)
        patgen.AddVocabPermutedBinary(ss.PoolVocab, "B", npats, plY, plX, pctAct, minDiff)
        patgen.AddVocabPermutedBinary(ss.PoolVocab, "C", npats, plY, plX, pctAct, minDiff)
        patgen.AddVocabPermutedBinary(ss.PoolVocab, "lA", npats, plY, plX, pctAct, minDiff)
        patgen.AddVocabPermutedBinary(ss.PoolVocab, "lB", npats, plY, plX, pctAct, minDiff)
        patgen.AddVocabPermutedBinary(ss.PoolVocab, "ctxt", 3, plY, plX, pctAct, minDiff)

        for i in range(12):
            lst = int(i / 4)
            ctxtNm = "ctxt%d" % (i+1)
            tsr = patgen.AddVocabRepeat(ss.PoolVocab, ctxtNm, npats, "ctxt", lst)
            patgen.FlipBitsRows(tsr, ctxtflip, ctxtflip, 1, 0)
            # todo: also support drifting
            # solution 2: drift based on last trial (will require sequential learning)
            # patgen.VocabDrift(ss.PoolVocab, ss.NFlipBits, "ctxt"+str(i+1))

        ecY = hp.ECSize.Y
        ecX = hp.ECSize.X

        patgen.InitPats(ss.TrainAB, "TrainAB", "TrainAB Pats", "Input", "ECout", npats, ecY, ecX, plY, plX)
        patgen.MixPats(ss.TrainAB, ss.PoolVocab, "Input", go.Slice_string(["A", "B", "ctxt1", "ctxt2", "ctxt3", "ctxt4"]))
        patgen.MixPats(ss.TrainAB, ss.PoolVocab, "ECout", go.Slice_string(["A", "B", "ctxt1", "ctxt2", "ctxt3", "ctxt4"]))

        patgen.InitPats(ss.TestAB, "TestAB", "TestAB Pats", "Input", "ECout", npats, ecY, ecX, plY, plX)
        patgen.MixPats(ss.TestAB, ss.PoolVocab, "Input", go.Slice_string(["A", "empty", "ctxt1", "ctxt2", "ctxt3", "ctxt4"]))
        patgen.MixPats(ss.TestAB, ss.PoolVocab, "ECout", go.Slice_string(["A", "B", "ctxt1", "ctxt2", "ctxt3", "ctxt4"]))

        patgen.InitPats(ss.TrainAC, "TrainAC", "TrainAC Pats", "Input", "ECout", npats, ecY, ecX, plY, plX)
        patgen.MixPats(ss.TrainAC, ss.PoolVocab, "Input", go.Slice_string(["A", "C", "ctxt5", "ctxt6", "ctxt7", "ctxt8"]))
        patgen.MixPats(ss.TrainAC, ss.PoolVocab, "ECout", go.Slice_string(["A", "C", "ctxt5", "ctxt6", "ctxt7", "ctxt8"]))

        patgen.InitPats(ss.TestAC, "TestAC", "TestAC Pats", "Input", "ECout", npats, ecY, ecX, plY, plX)
        patgen.MixPats(ss.TestAC, ss.PoolVocab, "Input", go.Slice_string(["A", "empty", "ctxt5", "ctxt6", "ctxt7", "ctxt8"]))
        patgen.MixPats(ss.TestAC, ss.PoolVocab, "ECout", go.Slice_string(["A", "C", "ctxt5", "ctxt6", "ctxt7", "ctxt8"]))

        patgen.InitPats(ss.TestLure, "TestLure", "TestLure Pats", "Input", "ECout", npats, ecY, ecX, plY, plX)
        patgen.MixPats(ss.TestLure, ss.PoolVocab, "Input", go.Slice_string(["lA", "empty", "ctxt9", "ctxt10", "ctxt11", "ctxt12"])) # arbitrary ctxt here
        patgen.MixPats(ss.TestLure, ss.PoolVocab, "ECout", go.Slice_string(["lA", "lB", "ctxt9", "ctxt10", "ctxt11", "ctxt12"]))    # arbitrary ctxt here

        ss.TrainAll = ss.TrainAB.Clone()
        ss.TrainAll.AppendRows(ss.TrainAC)
        ss.TrainAll.AppendRows(ss.TestLure)

    def ValsTsr(ss, name):
        """
        ValsTsr gets value tensor of given name, creating if not yet made
        """
        if name in ss.ValsTsrs:
            return ss.ValsTsrs[name]
        tsr = etensor.Float32()
        ss.ValsTsrs[name] = tsr
        return tsr

    def RunName(ss):
        """
        RunName returns a name for this run that combines Tag and Params -- add this to
        any file names that are saved.
        """
        if ss.Tag != "":
            pnm = ss.ParamsName()
            if pnm == "Base":
                return ss.Tag
            else:
                return ss.Tag + "_" + pnm
        else:
            return ss.ParamsName()

    def RunEpochName(ss, run, epc):
        """
        RunEpochName returns a string with the run and epoch numbers with leading zeros, suitable
        for using in weights file names.  Uses 3, 5 digits for each.
        """
        return "%03d_%05d" % (run, epc)

    def WeightsFileName(ss):
        """
        WeightsFileName returns default current weights file name
        """
        return ss.Net.Nm + "_" + ss.RunName() + "_" + ss.RunEpochName(ss.TrainEnv.Run.Cur, ss.TrainEnv.Epoch.Cur) + ".wts"

    def LogFileName(ss, lognm):
        """
        LogFileName returns default log file name
        """
        return ss.Net.Nm + "_" + ss.RunName() + "_" + lognm + ".tsv"

    def LogTrnTrl(ss, dt):
        """
        LogTrnTrl adds data from current trial to the TrnTrlLog table.
        log always contains number of testing items
        """
        epc = ss.TrainEnv.Epoch.Cur
        trl = ss.TrainEnv.Trial.Cur

        row = dt.Rows
        if trl == 0:
            row = 0
        dt.SetNumRows(row + 1)

        dt.SetCellFloat("Run", row, float(ss.TrainEnv.Run.Cur))
        dt.SetCellFloat("Epoch", row, float(epc))
        dt.SetCellFloat("Trial", row, float(trl))
        dt.SetCellString("TrialName", row, ss.TestEnv.TrialName.Cur)
        dt.SetCellFloat("SSE", row, ss.TrlSSE)
        dt.SetCellFloat("AvgSSE", row, ss.TrlAvgSSE)
        dt.SetCellFloat("CosDiff", row, ss.TrlCosDiff)

        dt.SetCellFloat("Mem", row, ss.Mem)
        dt.SetCellFloat("TrgOnWasOff", row, ss.TrgOnWasOffAll)
        dt.SetCellFloat("TrgOffWasOn", row, ss.TrgOffWasOn)

        if ss.TrnTrlPlot != 0:
            ss.TrnTrlPlot.GoUpdate()

    def ConfigTrnTrlLog(ss, dt):

        dt.SetMetaData("name", "TrnTrlLog")
        dt.SetMetaData("desc", "Record of training per input pattern")
        dt.SetMetaData("read-only", "true")
        dt.SetMetaData("precision", str(LogPrec))

        nt = ss.TestEnv.Table.Len()
        sch = etable.Schema(
            [etable.Column("Run", etensor.INT64, go.nil, go.nil),
            etable.Column("Epoch", etensor.INT64, go.nil, go.nil),
            etable.Column("Trial", etensor.INT64, go.nil, go.nil),
            etable.Column("TrialName", etensor.STRING, go.nil, go.nil),
            etable.Column("SSE", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("AvgSSE", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("CosDiff", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("Mem", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("TrgOnWasOff", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("TrgOffWasOn", etensor.FLOAT64, go.nil, go.nil)]
        )
        dt.SetFromSchema(sch, nt)

    def ConfigTrnTrlPlot(ss, plt, dt):
        plt.Params.Title = "Hippocampus Train Trial Plot"
        plt.Params.XAxisCol = "Trial"
        plt.SetTable(dt)
        # order of params: on, fixMin, min, fixMax, max
        plt.SetColParams("Run", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("Epoch", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("Trial", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("TrialName", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("SSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("AvgSSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("CosDiff", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)

        plt.SetColParams("Mem", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1)
        plt.SetColParams("TrgOnWasOff", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1)
        plt.SetColParams("TrgOffWasOn", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1)

        return plt

    def LogTrnEpc(ss, dt):
        """
        LogTrnEpc adds data from current epoch to the TrnEpcLog table.
        computes epoch averages prior to logging.

    # this is triggered by increment so use previous value
        """
        row = dt.Rows
        dt.SetNumRows(row + 1)

        epc = ss.TrainEnv.Epoch.Prv
        nt = float(ss.TrainEnv.Table.Len()) # number of trials in view

        ss.EpcSSE = ss.SumSSE / nt
        ss.SumSSE = 0
        ss.EpcAvgSSE = ss.SumAvgSSE / nt
        ss.SumAvgSSE = 0
        ss.EpcPctErr = float(ss.CntErr) / nt
        ss.CntErr = 0
        ss.EpcPctCor = 1 - ss.EpcPctErr
        ss.EpcCosDiff = ss.SumCosDiff / nt
        ss.SumCosDiff = 0

        trlog = ss.TrnTrlLog
        tix = etable.NewIdxView(trlog)

        dt.SetCellFloat("Run", row, float(ss.TrainEnv.Run.Cur))
        dt.SetCellFloat("Epoch", row, float(epc))
        dt.SetCellFloat("SSE", row, ss.EpcSSE)
        dt.SetCellFloat("AvgSSE", row, ss.EpcAvgSSE)
        dt.SetCellFloat("PctErr", row, ss.EpcPctErr)
        dt.SetCellFloat("PctCor", row, ss.EpcPctCor)
        dt.SetCellFloat("CosDiff", row, ss.EpcCosDiff)

        mem = agg.Mean(tix, "Mem")[0]
        dt.SetCellFloat("Mem", row, mem)
        dt.SetCellFloat("TrgOnWasOff", row, agg.Mean(tix, "TrgOnWasOff")[0])
        dt.SetCellFloat("TrgOffWasOn", row, agg.Mean(tix, "TrgOffWasOn")[0])

        for lnm in ss.LayStatNms :
            ly = leabra.Layer(ss.Net.LayerByName(lnm))
            dt.SetCellFloat(ly.Nm+" ActAvg", row, float(ly.Pools[0].ActAvg.ActPAvgEff))

        # note: essential to use Go version of update when called from another goroutine
        if ss.TrnEpcPlot != 0:
            ss.TrnEpcPlot.GoUpdate()
        if ss.TrnEpcFile != 0:
            if not ss.TrnEpcHdrs:
                dt.WriteCSVHeaders(ss.TrnEpcFile, etable.Tab)
                ss.TrnEpcHdrs = True
            dt.WriteCSVRow(ss.TrnEpcFile, row, etable.Tab)

    def ConfigTrnEpcLog(ss, dt):
        dt.SetMetaData("name", "TrnEpcLog")
        dt.SetMetaData("desc", "Record of performance over epochs of training")
        dt.SetMetaData("read-only", "true")
        dt.SetMetaData("precision", str(LogPrec))

        sch = etable.Schema(
            [etable.Column("Run", etensor.INT64, go.nil, go.nil),
            etable.Column("Epoch", etensor.INT64, go.nil, go.nil),
            etable.Column("SSE", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("AvgSSE", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("PctErr", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("PctCor", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("CosDiff", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("Mem", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("TrgOnWasOff", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("TrgOffWasOn", etensor.FLOAT64, go.nil, go.nil)]
        )
        for lnm in ss.LayStatNms :
            sch.append( etable.Column(lnm + " ActAvg", etensor.FLOAT64, go.nil, go.nil))
        dt.SetFromSchema(sch, 0)

    def ConfigTrnEpcPlot(ss, plt, dt):
        plt.Params.Title = "Hippocampus Epoch Plot"
        plt.Params.XAxisCol = "Epoch"
        plt.SetTable(dt)
        # order of params: on, fixMin, min, fixMax, max
        plt.SetColParams("Run", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("Epoch", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("SSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("AvgSSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("PctErr", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
        plt.SetColParams("PctCor", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
        plt.SetColParams("CosDiff", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)

        plt.SetColParams("Mem", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1)         # default plot
        plt.SetColParams("TrgOnWasOff", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1) # default plot
        plt.SetColParams("TrgOffWasOn", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1) # default plot

        for lnm in ss.LayStatNms :
            plt.SetColParams(lnm+" ActAvg", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 0.5)
        return plt

    def LogTstTrl(ss, dt):
        """
        LogTstTrl adds data from current trial to the TstTrlLog table.
    # this is triggered by increment so use previous value
        log always contains number of testing items
        """
        epc = ss.TrainEnv.Epoch.Prv
        trl = ss.TestEnv.Trial.Cur

        row = dt.Rows
        if ss.TestNm == "AB" and trl == 0: # reset at start
            row = 0
        dt.SetNumRows(row + 1)

        dt.SetCellFloat("Run", row, float(ss.TrainEnv.Run.Cur))
        dt.SetCellFloat("Epoch", row, float(epc))
        dt.SetCellString("TestNm", row, ss.TestNm)
        dt.SetCellFloat("Trial", row, float(row))
        dt.SetCellString("TrialName", row, ss.TestEnv.TrialName.Cur)
        dt.SetCellFloat("SSE", row, ss.TrlSSE)
        dt.SetCellFloat("AvgSSE", row, ss.TrlAvgSSE)
        dt.SetCellFloat("CosDiff", row, ss.TrlCosDiff)

        dt.SetCellFloat("Mem", row, ss.Mem)
        dt.SetCellFloat("TrgOnWasOff", row, ss.TrgOnWasOffCmp)
        dt.SetCellFloat("TrgOffWasOn", row, ss.TrgOffWasOn)

        for lnm in ss.LayStatNms :
            ly = leabra.Layer(ss.Net.LayerByName(lnm))
            dt.SetCellFloat(ly.Nm+" ActM.Avg", row, float(ly.Pools[0].ActM.Avg))

        for lnm in ss.LayStatNms :
            ly = leabra.Layer(ss.Net.LayerByName(lnm))
            tsr = ss.ValsTsr(lnm)
            ly.UnitValsTensor(tsr, "Act")
            dt.SetCellTensor(lnm+"Act", row, tsr)

        # note: essential to use Go version of update when called from another goroutine
        if ss.TstTrlPlot != 0:
            ss.TstTrlPlot.GoUpdate()

    def ConfigTstTrlLog(ss, dt):
        # inLay := ss.Net.LayerByName("Input").(leabra.LeabraLayer)
        # outLay := ss.Net.LayerByName("Output").(leabra.LeabraLayer)

        dt.SetMetaData("name", "TstTrlLog")
        dt.SetMetaData("desc", "Record of testing per input pattern")
        dt.SetMetaData("read-only", "true")
        dt.SetMetaData("precision", str(LogPrec))

        nt = ss.TestEnv.Table.Len() # number in view
        sch = etable.Schema(
            [etable.Column("Run", etensor.INT64, go.nil, go.nil),
            etable.Column("Epoch", etensor.INT64, go.nil, go.nil),
            etable.Column("TestNm", etensor.STRING, go.nil, go.nil),
            etable.Column("Trial", etensor.INT64, go.nil, go.nil),
            etable.Column("TrialName", etensor.STRING, go.nil, go.nil),
            etable.Column("SSE", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("AvgSSE", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("CosDiff", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("Mem", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("TrgOnWasOff", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("TrgOffWasOn", etensor.FLOAT64, go.nil, go.nil)]
        )
        for lnm in ss.LayStatNms :
            sch.append( etable.Column(lnm + " ActM.Avg", etensor.FLOAT64, go.nil, go.nil))
        for lnm in ss.LayStatNms :
            ly = leabra.Layer(ss.Net.LayerByName(lnm))
            sch.append( etable.Column(lnm + "Act", etensor.FLOAT64, ly.Shp.Shp, go.nil))

        dt.SetFromSchema(sch, nt)

    def ConfigTstTrlPlot(ss, plt, dt):
        plt.Params.Title = "Hippocampus Test Trial Plot"
        plt.Params.XAxisCol = "TrialName"
        plt.Params.Type = eplot.Bar
        plt.SetTable(dt) # this sets defaults so set params after
        plt.Params.XAxisRot = 45
        # order of params: on, fixMin, min, fixMax, max
        plt.SetColParams("Run", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("Epoch", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("TestNm", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("Trial", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("TrialName", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("SSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("AvgSSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("CosDiff", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)

        plt.SetColParams("Mem", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1)
        plt.SetColParams("TrgOnWasOff", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1)
        plt.SetColParams("TrgOffWasOn", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1)

        for lnm in ss.LayStatNms :
            plt.SetColParams(lnm+" ActM.Avg", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 0.5)
        for lnm in ss.LayStatNms :
            plt.SetColParams(lnm+" Act", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)

        return plt

    def RepsAnalysis(ss):
        """
        RepsAnalysis analyzes representations
        """
        acts = etable.NewIdxView(ss.TstTrlLog)
        for lnm in ss.LayStatNms:
            sm = 0
            if not lnm in ss.SimMats:
                sm = simat.SimMat()
                ss.SimMats[lnm] = sm
            else:
                sm = ss.SimMats[lnm]
            sm.TableColStd(acts, lnm+"Act", "TrialName", True, metric.Correlation)

    def SimMatStat(ss, lnm):
        """
        SimMatStat returns within, between for sim mat statistics
        """
        sm = ss.SimMats[lnm]
        smat = sm.Mat
        nitm = smat.Dim(0)
        ncat = int(nitm / len(ss.TstNms))
        win_sum = float(0)
        win_n = 0
        btn_sum = float(0)
        btn_n = 0
        for y in range(nitm):
            for x in range(y):
                val = smat.FloatVal(go.Slice_int([y, x]))
                same = int((y / ncat)) == int((x / ncat))
                if same:
                    win_sum += val
                    win_n += 1
                else:
                    btn_sum += val
                    btn_n += 1
        if win_n > 0:
            win_sum /= float(win_n)
        if btn_n > 0:
            btn_sum /= float(btn_n)
        return win_sum, btn_sum

    def LogTstEpc(ss, dt):
        row = dt.Rows
        dt.SetNumRows(row + 1)

        ss.RepsAnalysis()

        trl = ss.TstTrlLog
        tix = etable.NewIdxView(trl)
        epc = ss.TrainEnv.Epoch.Prv

        # if ss.LastEpcTime.IsZero():
        #     ss.EpcPerTrlMSec = 0
        # else:
        #     iv = time.Now().Sub(ss.LastEpcTime)
        #     nt = ss.TrainAB.Rows * 4 # 1 train and 3 tests
        #     ss.EpcPerTrlMSec = float(iv) / (float(nt) * float(time.Millisecond))
        # ss.LastEpcTime = time.Now()

        # note: this shows how to use agg methods to compute summary data from another
        # data table, instead of incrementing on the Sim
        dt.SetCellFloat("Run", row, float(ss.TrainEnv.Run.Cur))
        dt.SetCellFloat("Epoch", row, float(epc))
        dt.SetCellFloat("PerTrlMSec", row, ss.EpcPerTrlMSec)
        dt.SetCellFloat("SSE", row, agg.Sum(tix, "SSE")[0])
        dt.SetCellFloat("AvgSSE", row, agg.Mean(tix, "AvgSSE")[0])
        dt.SetCellFloat("PctErr", row, agg.PropIf(tix, "SSE", AggIfGt0)[0])
        dt.SetCellFloat("PctCor", row, agg.PropIf(tix, "SSE", AggIfEq0)[0])
        dt.SetCellFloat("CosDiff", row, agg.Mean(tix, "CosDiff")[0])

        trix = etable.NewIdxView(trl)
        spl = split.GroupBy(trix, go.Slice_string(["TestNm"]))
        for ts in ss.TstStatNms :
            split.Agg(spl, ts, agg.AggMean)
        ss.TstStats = spl.AggsToTable(etable.ColNameOnly)

        for ri in range(ss.TstStats.Rows):
            tst = ss.TstStats.CellString("TestNm", ri)
            for ts in ss.TstStatNms :
                dt.SetCellFloat(tst+" "+ts, row, ss.TstStats.CellFloat(ts, ri))

        for lnm in ss.LayStatNms :
            # win, btn = ss.SimMatStat(lnm)
            win = 0
            btn = 0
            for ts in ss.SimMatStats :
                if ts == "Within":
                    dt.SetCellFloat(lnm+" "+ts, row, win)
                else:
                    dt.SetCellFloat(lnm+" "+ts, row, btn)

        # base zero on testing performance!
        curAB = (ss.TrainEnv.Table.Table.MetaData["name"] == "TrainAB")
        mem = float()
        if curAB:
            mem = dt.CellFloat("AB Mem", row)
        else:
            mem = dt.CellFloat("AC Mem", row)
        if ss.FirstZero < 0 and mem == 1:
            ss.FirstZero = epc
        if mem == 1:
            ss.NZero += 1
        else:
            ss.NZero = 0

        # note: essential to use Go version of update when called from another goroutine
        if ss.TstEpcPlot != 0:
            ss.TstEpcPlot.GoUpdate()
        if ss.TstEpcFile != 0:
            if not ss.TstEpcHdrs:
                dt.WriteCSVHeaders(ss.TstEpcFile, etable.Tab)
                ss.TstEpcHdrs = True
            dt.WriteCSVRow(ss.TstEpcFile, row, etable.Tab)

    def ConfigTstEpcLog(ss, dt):
        dt.SetMetaData("name", "TstEpcLog")
        dt.SetMetaData("desc", "Summary stats for testing trials")
        dt.SetMetaData("read-only", "true")
        dt.SetMetaData("precision", str(LogPrec))

        sch = etable.Schema(
            [etable.Column("Run", etensor.INT64, go.nil, go.nil),
            etable.Column("Epoch", etensor.INT64, go.nil, go.nil),
            etable.Column("PerTrlMSec", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("SSE", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("AvgSSE", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("PctErr", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("PctCor", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("CosDiff", etensor.FLOAT64, go.nil, go.nil)]
        )
        for tn in ss.TstNms :
            for ts in ss.TstStatNms :
                sch.append( etable.Column(tn + " " + ts, etensor.FLOAT64, go.nil, go.nil))
        for lnm in ss.LayStatNms :
            for ts in ss.SimMatStats :
                sch.append( etable.Column(lnm + " " + ts, etensor.FLOAT64, go.nil, go.nil))
        dt.SetFromSchema(sch, 0)

    def ConfigTstEpcPlot(ss, plt, dt):
        plt.Params.Title = "Hippocampus Testing Epoch Plot"
        plt.Params.XAxisCol = "Epoch"
        plt.SetTable(dt) # this sets defaults so set params after
        # order of params: on, fixMin, min, fixMax, max
        plt.SetColParams("Run", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("Epoch", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("PerTrlMSec", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("SSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("AvgSSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("PctErr", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
        plt.SetColParams("PctCor", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
        plt.SetColParams("CosDiff", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)

        for tn in ss.TstNms :
            for ts in ss.TstStatNms :
                if ts == "Mem":
                    plt.SetColParams(tn+" "+ts, eplot.On, eplot.FixMin, 0, eplot.FixMax, 1)
                else:
                    plt.SetColParams(tn+" "+ts, eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
        for lnm in ss.LayStatNms :
            for ts in ss.SimMatStats :
                plt.SetColParams(lnm+" "+ts, eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 1)
        return plt

    def LogTstCyc(ss, dt, cyc):
        """
        LogTstCyc adds data from current trial to the TstCycLog table.
        log just has 100 cycles, is overwritten
        """
        if dt.Rows <= cyc:
            dt.SetNumRows(cyc + 1)

        dt.SetCellFloat("Cycle", cyc, float(cyc))
        for lnm in ss.LayStatNms :
            ly = leabra.Layer(ss.Net.LayerByName(lnm))
            dt.SetCellFloat(ly.Nm+" Ge.Avg", cyc, float(ly.Pools[0].Inhib.Ge.Avg))
            dt.SetCellFloat(ly.Nm+" Act.Avg", cyc, float(ly.Pools[0].Inhib.Act.Avg))

        if cyc%10 == 0: # too slow to do every cyc
            # note: essential to use Go version of update when called from another goroutine
            if ss.TstCycPlot != 0:
                ss.TstCycPlot.GoUpdate()

    def ConfigTstCycLog(ss, dt):
        dt.SetMetaData("name", "TstCycLog")
        dt.SetMetaData("desc", "Record of activity etc over one trial by cycle")
        dt.SetMetaData("read-only", "true")
        dt.SetMetaData("precision", str(LogPrec))

        np = 100 # max cycles
        sch = etable.Schema(
            [etable.Column("Cycle", etensor.INT64, go.nil, go.nil)]
        )
        for lnm in ss.LayStatNms :
            sch.append( etable.Column(lnm + " Ge.Avg", etensor.FLOAT64, go.nil, go.nil))
            sch.append( etable.Column(lnm + " Act.Avg", etensor.FLOAT64, go.nil, go.nil))
        dt.SetFromSchema(sch, np)

    def ConfigTstCycPlot(ss, plt, dt):
        plt.Params.Title = "Hippocampus Test Cycle Plot"
        plt.Params.XAxisCol = "Cycle"
        plt.SetTable(dt)
        # order of params: on, fixMin, min, fixMax, max
        plt.SetColParams("Cycle", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        for lnm in ss.LayStatNms :
            plt.SetColParams(lnm+" Ge.Avg", eplot.On, eplot.FixMin, 0, eplot.FixMax, .5)
            plt.SetColParams(lnm+" Act.Avg", eplot.On, eplot.FixMin, 0, eplot.FixMax, .5)
        return plt

    def LogRun(ss, dt):
        """
        LogRun adds data from current run to the RunLog table.
        """
        epclog = ss.TstEpcLog
        epcix = etable.NewIdxView(epclog)
        if epcix.Len() == 0:
            return

        run = ss.TrainEnv.Run.Cur # this is NOT triggered by increment yet -- use Cur
        row = dt.Rows
        dt.SetNumRows(row + 1)

        # compute mean over last N epochs for run level
        nlast = 1
        if nlast > epcix.Len()-1:
            nlast = epcix.Len() - 1
        epcix.Idxs = epcix.Idxs[epcix.Len()-nlast:]

        params = ss.RunName() # includes tag

        fzero = ss.FirstZero
        if fzero < 0:
            fzero = ss.MaxEpcs

        dt.SetCellFloat("Run", row, float(run))
        dt.SetCellString("Params", row, params)
        dt.SetCellFloat("NEpochs", row, float(ss.TstEpcLog.Rows))
        dt.SetCellFloat("FirstZero", row, float(fzero))
        dt.SetCellFloat("SSE", row, agg.Mean(epcix, "SSE")[0])
        dt.SetCellFloat("AvgSSE", row, agg.Mean(epcix, "AvgSSE")[0])
        dt.SetCellFloat("PctErr", row, agg.Mean(epcix, "PctErr")[0])
        dt.SetCellFloat("PctCor", row, agg.Mean(epcix, "PctCor")[0])
        dt.SetCellFloat("CosDiff", row, agg.Mean(epcix, "CosDiff")[0])

        for tn in ss.TstNms :
            for ts in ss.TstStatNms :
                nm = tn + " " + ts
                dt.SetCellFloat(nm, row, agg.Mean(epcix, nm)[0])
        for lnm in ss.LayStatNms :
            for ts in ss.SimMatStats :
                nm = lnm + " " + ts
                dt.SetCellFloat(nm, row, agg.Mean(epcix, nm)[0])

        ss.LogRunStats()

        # note: essential to use Go version of update when called from another goroutine
        if ss.RunPlot != 0:
            ss.RunPlot.GoUpdate()
        if ss.RunFile != 0:
            if row == 0:
                dt.WriteCSVHeaders(ss.RunFile, etable.Tab)
            dt.WriteCSVRow(ss.RunFile, row, etable.Tab)

    def ConfigRunLog(ss, dt):
        dt.SetMetaData("name", "RunLog")
        dt.SetMetaData("desc", "Record of performance at end of training")
        dt.SetMetaData("read-only", "true")
        dt.SetMetaData("precision", str(LogPrec))

        sch = etable.Schema(
            [etable.Column("Run", etensor.INT64, go.nil, go.nil),
            etable.Column("Params", etensor.STRING, go.nil, go.nil),
            etable.Column("NEpochs", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("FirstZero", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("SSE", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("AvgSSE", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("PctErr", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("PctCor", etensor.FLOAT64, go.nil, go.nil),
            etable.Column("CosDiff", etensor.FLOAT64, go.nil, go.nil)]
        )
        for tn in ss.TstNms :
            for ts in ss.TstStatNms :
                sch.append( etable.Column(tn + " " + ts, etensor.FLOAT64, go.nil, go.nil))
        for lnm in ss.LayStatNms :
            for ts in ss.SimMatStats :
                sch.append( etable.Column(lnm + " " + ts, etensor.FLOAT64, go.nil, go.nil))
        dt.SetFromSchema(sch, 0)

    def ConfigRunPlot(ss, plt, dt):
        plt.Params.Title = "Hippocampus Run Plot"
        plt.Params.XAxisCol = "Run"
        plt.SetTable(dt)
        # order of params: on, fixMin, min, fixMax, max
        plt.SetColParams("Run", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("NEpochs", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("FirstZero", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("SSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("AvgSSE", eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 0)
        plt.SetColParams("PctErr", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
        plt.SetColParams("PctCor", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
        plt.SetColParams("CosDiff", eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)

        for tn in ss.TstNms :
            for ts in ss.TstStatNms :
                if ts == "Mem":
                    plt.SetColParams(tn+" "+ts, eplot.On, eplot.FixMin, 0, eplot.FixMax, 1) # default plot
                else:
                    plt.SetColParams(tn+" "+ts, eplot.Off, eplot.FixMin, 0, eplot.FixMax, 1)
        for lnm in ss.LayStatNms :
            for ts in ss.SimMatStats :
                plt.SetColParams(lnm+" "+ts, eplot.Off, eplot.FixMin, 0, eplot.FloatMax, 1)
        return plt

    def LogRunStats(ss):
        """
        LogRunStats computes RunStats from RunLog data -- can be used for looking at prelim results
        """
        dt = ss.RunLog
        runix = etable.NewIdxView(dt)
        spl = split.GroupBy(runix, go.Slice_string(["Params"]))
        for tn in ss.TstNms :
            nm = tn + " " + "Mem"
            split.Desc(spl, nm)
        split.Desc(spl, "FirstZero")
        split.Desc(spl, "NEpochs")
        for lnm in ss.LayStatNms :
            for ts in ss.SimMatStats :
                split.Desc(spl, lnm+" "+ts)
        ss.RunStats = spl.AggsToTable(etable.AddAggName)
        if ss.RunStatsPlot != 0:
            ss.ConfigRunStatsPlot(ss.RunStatsPlot, ss.RunStats)

    def ConfigRunStatsPlot(ss, plt, dt):
        plt.Params.Title = "Hippocampus Run Stats Plot"
        plt.Params.XAxisCol = "Params"
        plt.SetTable(dt)
        plt.Params.BarWidth = 10
        plt.Params.Type = eplot.Bar
        plt.Params.XAxisRot = 45

        cp = plt.SetColParams("AB Mem:Mean", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1)
        cp.ErrCol = "AB Mem:Sem"
        cp = plt.SetColParams("AC Mem:Mean", eplot.On, eplot.FixMin, 0, eplot.FixMax, 1)
        cp.ErrCol = "AC Mem:Sem"
        cp = plt.SetColParams("FirstZero:Mean", eplot.On, eplot.FixMin, 0, eplot.FixMax, 30)
        cp.ErrCol = "FirstZero:Sem"
        cp = plt.SetColParams("NEpochs:Mean", eplot.On, eplot.FixMin, 0, eplot.FixMax, 30)
        cp.ErrCol = "NEpochs:Sem"
        return plt

    def ConfigGui(ss):
        """
        ConfigGui configures the GoGi gui interface for this simulation,
        """
        width = 1600
        height = 1200

        gi.SetAppName("hip_bench")
        gi.SetAppAbout('This demonstrates a basic Hippocampus model in Leabra. See <a href="https://github.com/emer/emergent">emergent on GitHub</a>.</p>')

        win = gi.NewMainWindow("hip_bench", "Hippocampus AB-AC", width, height)
        ss.Win = win

        vp = win.WinViewport2D()
        ss.vp = vp
        updt = vp.UpdateStart()

        mfr = win.SetMainFrame()

        tbar = gi.AddNewToolBar(mfr, "tbar")
        tbar.SetStretchMaxWidth()
        ss.ToolBar = tbar

        split = gi.AddNewSplitView(mfr, "split")
        split.Dim = mat32.X
        split.SetStretchMax()

        cv = ss.NewClassView("sv")
        cv.AddFrame(split)
        cv.Config()

        tv = gi.AddNewTabView(split, "tv")

        nv = netview.NetView()
        tv.AddTab(nv, "NetView")
        nv.Var = "Act"
        nv.SetNet(ss.Net)
        ss.NetView = nv
        nv.ViewDefaults()

        plt = eplot.Plot2D()
        tv.AddTab(plt, "TrnTrlPlot")
        ss.TrnTrlPlot = ss.ConfigTrnTrlPlot(plt, ss.TrnTrlLog)

        plt = eplot.Plot2D()
        tv.AddTab(plt, "TrnEpcPlot")
        ss.TrnEpcPlot = ss.ConfigTrnEpcPlot(plt, ss.TrnEpcLog)

        plt = eplot.Plot2D()
        tv.AddTab(plt, "TstTrlPlot")
        ss.TstTrlPlot = ss.ConfigTstTrlPlot(plt, ss.TstTrlLog)

        plt = eplot.Plot2D()
        tv.AddTab(plt, "TstEpcPlot")
        ss.TstEpcPlot = ss.ConfigTstEpcPlot(plt, ss.TstEpcLog)

        plt = eplot.Plot2D()
        tv.AddTab(plt, "TstCycPlot")
        ss.TstCycPlot = ss.ConfigTstCycPlot(plt, ss.TstCycLog)

        plt = eplot.Plot2D()
        tv.AddTab(plt, "RunPlot")
        ss.RunPlot = ss.ConfigRunPlot(plt, ss.RunLog)

        split.SetSplitsList(go.Slice_float32([.2, .8]))
        recv = win.This()

        tbar.AddAction(gi.ActOpts(Label="Init", Icon="update", Tooltip="Initialize everything including network weights, and start over.  Also applies current params.", UpdateFunc=UpdtFuncNotRunning), recv, InitCB)

        tbar.AddAction(gi.ActOpts(Label="Train", Icon="run", Tooltip="Starts the network training, picking up from wherever it may have left off.  If not stopped, training will complete the specified number of Runs through the full number of Epochs of training, with testing automatically occuring at the specified interval.", UpdateFunc=UpdtFuncNotRunning), recv, TrainCB)
        
        tbar.AddAction(gi.ActOpts(Label="Stop", Icon="stop", Tooltip="Interrupts running.  Hitting Train again will pick back up where it left off.", UpdateFunc=UpdtFuncRunning), recv, StopCB)
        
        tbar.AddAction(gi.ActOpts(Label="Step Trial", Icon="step-fwd", Tooltip="Advances one training trial at a time.", UpdateFunc=UpdtFuncNotRunning), recv, StepTrialCB)
        
        tbar.AddAction(gi.ActOpts(Label="Step Epoch", Icon="fast-fwd", Tooltip="Advances one epoch (complete set of training patterns) at a time.", UpdateFunc=UpdtFuncNotRunning), recv, StepEpochCB)

        tbar.AddAction(gi.ActOpts(Label="Step Run", Icon="fast-fwd", Tooltip="Advances one full training Run at a time.", UpdateFunc=UpdtFuncNotRunning), recv, StepRunCB)
        
        tbar.AddSeparator("test")
        
        tbar.AddAction(gi.ActOpts(Label="Test Trial", Icon="step-fwd", Tooltip="Runs the next testing trial.", UpdateFunc=UpdtFuncNotRunning), recv, TestTrialCB)
        
        tbar.AddAction(gi.ActOpts(Label="Test Item", Icon="step-fwd", Tooltip="Prompts for a specific input pattern name to run, and runs it in testing mode.", UpdateFunc=UpdtFuncNotRunning), recv, TestItemCB)
        
        tbar.AddAction(gi.ActOpts(Label="Test All", Icon="fast-fwd", Tooltip="Tests all of the testing trials.", UpdateFunc=UpdtFuncNotRunning), recv, TestAllCB)

        tbar.AddSeparator("log")
        
        # tbar.AddAction(gi.ActOpts(Label= "Env", Icon= "gear", Tooltip= "select training input patterns: AB or AC."), win.This(),
        #     funcrecv, send, sig, data:
        #         giv.CallMethod(ss, "SetEnv", vp))

        tbar.AddAction(gi.ActOpts(Label="Reset RunLog", Icon="reset", Tooltip="Resets the accumulated log of all Runs, which are tagged with the ParamSet used"), recv, ResetRunLogCB)

        tbar.AddSeparator("misc")
        
        tbar.AddAction(gi.ActOpts(Label="New Seed", Icon="new", Tooltip="Generate a new initial random seed to get different results.  By default, Init re-establishes the same initial seed every time."), recv, NewRndSeedCB)

        tbar.AddAction(gi.ActOpts(Label="README", Icon="file-markdown", Tooltip="Opens your browser on the README file that contains instructions for how to run this model."), recv, ReadmeCB)

        # main menu
        appnm = gi.AppName()
        mmen = win.MainMenu
        mmen.ConfigMenus(go.Slice_string([appnm, "File", "Edit", "Window"]))

        amen = gi.Action(win.MainMenu.ChildByName(appnm, 0))
        amen.Menu.AddAppMenu(win)

        emen = gi.Action(win.MainMenu.ChildByName("Edit", 1))
        emen.Menu.AddCopyCutPaste(win)

        # note: Command in shortcuts is automatically translated into Control for
        # Linux, Windows or Meta for MacOS
        # fmen := win.MainMenu.ChildByName("File", 0).(*gi.Action)
        # fmen.Menu.AddAction(gi.ActOpts{Label: "Open", Shortcut: "Command+O"},
        #   win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
        #       FileViewOpenSVG(vp)
        #   })
        # fmen.Menu.AddSeparator("csep")
        # fmen.Menu.AddAction(gi.ActOpts{Label: "Close Window", Shortcut: "Command+W"},
        #   win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
        #       win.Close()
        #   })

        win.MainMenuUpdated()
        vp.UpdateEndNoSig(updt)
        win.GoStartEventLoop()

    def TwoFactorRun(ss):
        """
        TwoFactorRun runs outer-loop crossed with inner-loop params
        """
        tag = ss.Tag
        usetag = tag
        if usetag != "":
            usetag += "_"
        for otf in OuterLoopParams:
            for inf in InnerLoopParams:
                ss.Tag = usetag + otf + "_" + inf
                print("running: " + ss.Tag)
                rand.Seed(ss.RndSeed) # each run starts at same seed..
                ss.SetParamsSet(otf, "", ss.LogSetParams)
                ss.SetParamsSet(inf, "", ss.LogSetParams)
                ss.ReConfigNet() # note: this applies Base params to Network
                ss.ConfigEnv()
                ss.StopNow = False
                ss.PreTrain()
                ss.NewRun()
                ss.Train()
        ss.Tag = tag

# TheSim is the overall state for this simulation
TheSim = Sim()
 
def usage():
    print(sys.argv[0] + " --params=<param set> --tag=<extra tag> --setparams --wts --epclog=0 --runlog=0 --nogui")
    print("\t pyleabra -i %s to run in interactive, gui mode" % sys.argv[0])
    print("\t --params=<param set> additional params to apply on top of Base (name must be in loaded Params")
    print("\t --tag=<extra tag>    tag is appended to file names to uniquely identify this run") 
    print("\t --note=<note>        user note -- describe the run params etc") 
    print("\t --runs=<n>           number of runs to do")
    print("\t --epcs=<n>           number of epochs per run")
    print("\t --setparams          show the parameter values that are set")
    print("\t --wts                save final trained weights after every run")
    print("\t --epclog=0/False     turn off save training epoch log data to file named by param set, tag")
    print("\t --runlog=0/False     turn off save run log data to file named by param set, tag")
    print("\t --nogui              if no other args needed, this prevents running under the gui")

def main(argv):
    TheSim.Config()

    # print("n args: %d" % len(argv))
    TheSim.NoGui = len(argv) > 1
    saveEpcLog = True
    saveRunLog = True
        
    try:
        opts, args = getopt.getopt(argv,"h:",["params=","tag=","note=","runs=","epcs=","setparams","wts","epclog=","runlog=","nogui"])
    except getopt.GetoptError:
        usage()
        sys.exit(2)
    for opt, arg in opts:
        # print("opt: %s  arg: %s" % (opt, arg))
        if opt == '-h':
            usage()
            sys.exit()
        elif opt == "--tag":
            TheSim.Tag = arg
        elif opt == "--runs":
            TheSim.MaxRuns = int(arg)
            print("Running %d runs" % TheSim.MaxRuns)
        elif opt == "--epcs":
            TheSim.MaxEpcs = int(arg)
            print("Running %d epochs" % TheSim.MaxEpcs)
        elif opt == "--setparams":
            TheSim.LogSetParams = True
        elif opt == "--wts":
            TheSim.SaveWts = True
            print("Saving final weights per run")
        elif opt == "--epclog":
            if arg.lower() == "false" or arg == "0":
                saveEpcLog = False
        elif opt == "--runlog":
            if arg.lower() == "false" or arg == "0":
                saveRunLog = False
        elif opt == "--nogui":
            TheSim.NoGui = True

    TheSim.Init()
            
    if TheSim.NoGui:
        if saveEpcLog:
            fnm = TheSim.LogFileName("epc") 
            print("Saving test epoch log to: %s" % fnm)
            TheSim.TstEpcFile = efile.Create(fnm)
    
        if saveRunLog:
            fnm = TheSim.LogFileName("run") 
            print("Saving run log to: %s" % fnm)
            TheSim.RunFile = efile.Create(fnm)

        # TheSim.Train()
        TheSim.TwoFactorRun()
        fnm = TheSim.LogFileName("runs")
        TheSim.RunStats.SaveCSV(fnm, etable.Tab, etable.Headers)
        
    else:
        TheSim.ConfigGui()
        print("Note: run pyleabra -i hip_bench.py to run in interactive mode, or just pyleabra, then 'import ra25'")
        print("for non-gui background running, here are the args:")
        usage()
        import code
        code.interact(local=locals())

main(sys.argv[1:])


