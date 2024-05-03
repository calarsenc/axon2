// Copyright (c) 2023, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// effort_plot plots the Rubicon effort cost equations.
package main

import (
	"math/rand"
	"strconv"

	"cogentcore.org/core/base/num"
	"cogentcore.org/core/base/randx"
	"cogentcore.org/core/core"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/math32/minmax"
	"cogentcore.org/core/plot/plotview"
	"cogentcore.org/core/tensor/table"
	"cogentcore.org/core/views"
	"github.com/emer/axon/v2/axon"
)

func DriveEffortGUI() {
	ep := &DrEffPlot{}
	ep.Config()
	b := ep.ConfigGUI()
	b.RunMainWindow()
}

// LogPrec is precision for saving float values in logs
const LogPrec = 4

// DrEffPlot holds the params, table, etc
type DrEffPlot struct {

	// context just for plotting
	Context axon.Context

	// Rubicon params
	Rubicon axon.Rubicon

	// total number of time steps to simulate
	TimeSteps int

	// range for number of time steps between US receipt
	USTime minmax.Int

	// range for random effort per step
	Effort minmax.F32

	// table for plot
	Table *table.Table `view:"no-inline"`

	// the plot
	Plot *plotview.PlotView `view:"-"`

	// table for plot
	TimeTable *table.Table `view:"no-inline"`

	// the plot
	TimePlot *plotview.PlotView `view:"-"`

	// random number generator
	Rand randx.SysRand `view:"-"`
}

// Config configures all the elements using the standard functions
func (ss *DrEffPlot) Config() {
	ss.Context.Defaults()
	pp := &ss.Rubicon
	pp.SetNUSs(&ss.Context, 1, 1)
	pp.Defaults()
	pp.Drive.DriveMin = 0
	pp.Drive.Base[0] = 1
	pp.Drive.Tau[0] = 100
	pp.Drive.Satisfaction[0] = 1
	pp.Drive.Update()
	ss.TimeSteps = 100
	ss.USTime.Set(2, 20)
	ss.Effort.Set(0.5, 1.5)
	ss.Update()
	ss.Table = &table.Table{}
	ss.ConfigTable(ss.Table)
	ss.TimeTable = &table.Table{}
	ss.ConfigTimeTable(ss.TimeTable)
}

// Update updates computed values
func (ss *DrEffPlot) Update() {
}

// EffortPlot plots the equation as a function of effort / time
func (ss *DrEffPlot) EffortPlot() { //types:add
	ss.Update()
	ctx := &ss.Context
	pp := &ss.Rubicon
	dt := ss.Table
	nv := 100
	dt.SetNumRows(nv)
	pp.TimeEffortReset(ctx, 0)
	for vi := 0; vi < nv; vi++ {
		ev := 1 - axon.RubiconNormFun(0.02)
		dt.SetFloat("X", vi, float64(vi))
		dt.SetFloat("Y", vi, float64(ev))

		pp.AddTimeEffort(ctx, 0, 1) // unit
	}
	ss.Plot.Update()
}

// UrgencyPlot plots the equation as a function of effort / time
func (ss *DrEffPlot) UrgencyPlot() { //types:add
	ctx := &ss.Context
	pp := &ss.Rubicon
	ss.Update()
	dt := ss.Table
	nv := 100
	dt.SetNumRows(nv)
	pp.Urgency.Reset(ctx, 0)
	for vi := 0; vi < nv; vi++ {
		ev := pp.Urgency.Urge(ctx, 0)
		dt.SetFloat("X", vi, float64(vi))
		dt.SetFloat("Y", vi, float64(ev))

		pp.Urgency.AddEffort(ctx, 0, 1) // unit
	}
	ss.Plot.Update()
}

func (ss *DrEffPlot) ConfigTable(dt *table.Table) {
	dt.SetMetaData("name", "PlotTable")
	dt.SetMetaData("read-only", "true")
	dt.SetMetaData("precision", strconv.Itoa(LogPrec))

	dt.AddFloat64Column("X")
	dt.AddFloat64Column("Y")
	dt.SetNumRows(0)
}

func (ss *DrEffPlot) ConfigPlot(plt *plotview.PlotView, dt *table.Table) *plotview.PlotView {
	plt.Params.Title = "Effort Discount or Urgency Function Plot"
	plt.Params.XAxisColumn = "X"
	plt.SetTable(dt)
	// order of params: on, fixMin, min, fixMax, max
	plt.SetColParams("X", plotview.Off, plotview.FloatMin, 0, plotview.FloatMax, 0)
	plt.SetColParams("Y", plotview.On, plotview.FixMin, 0, plotview.FixMax, 1)
	return plt
}

/////////////////////////////////////////////////////////////////

// TimeRun runs the equation over time.
func (ss *DrEffPlot) TimeRun() { //types:add
	ss.Update()
	dt := ss.TimeTable
	pp := &ss.Rubicon
	ctx := &ss.Context
	pp.TimeEffortReset(ctx, 0)
	pp.Urgency.Reset(ctx, 0)
	ut := ss.USTime.Min + rand.Intn(ss.USTime.Range())
	dt.SetNumRows(ss.TimeSteps)
	axon.SetGlbUSposV(ctx, 0, axon.GvUSpos, 1, 0)
	pp.Drive.ToBaseline(ctx, 0)
	// pv.Update()
	lastUS := 0
	for ti := 0; ti < ss.TimeSteps; ti++ {
		ev := 1 - axon.RubiconNormFun(0.02)
		urg := pp.Urgency.Urge(ctx, 0)
		ei := ss.Effort.Min + rand.Float32()*ss.Effort.Range()
		dr := axon.GlbUSposV(ctx, 0, axon.GvDrives, 0)
		usv := float32(0)
		if ti == lastUS+ut {
			ei = 0 // don't update on us trial
			lastUS = ti
			ut = ss.USTime.Min + rand.Intn(ss.USTime.Range())
			usv = 1
		}
		dt.SetFloat("T", ti, float64(ti))
		dt.SetFloat("Eff", ti, float64(ev))
		dt.SetFloat("EffInc", ti, float64(ei))
		dt.SetFloat("Urge", ti, float64(urg))
		dt.SetFloat("US", ti, float64(usv))
		dt.SetFloat("Drive", ti, float64(dr))

		axon.SetGlbUSposV(ctx, 0, axon.GvUSpos, 1, usv)
		axon.SetGlbV(ctx, 0, axon.GvHadRew, num.FromBool[float32](usv > 0))
		pp.EffortUrgencyUpdate(ctx, 0, 0)
		pp.DriveUpdate(ctx, 0)
	}
	ss.TimePlot.Update()
}

func (ss *DrEffPlot) ConfigTimeTable(dt *table.Table) {
	dt.SetMetaData("name", "TimeTable")
	dt.SetMetaData("read-only", "true")
	dt.SetMetaData("precision", strconv.Itoa(LogPrec))

	dt.AddFloat64Column("T")
	dt.AddFloat64Column("Eff")
	dt.AddFloat64Column("EffInc")
	dt.AddFloat64Column("Urge")
	dt.AddFloat64Column("US")
	dt.AddFloat64Column("Drive")
	dt.SetNumRows(0)
}

func (ss *DrEffPlot) ConfigTimePlot(plt *plotview.PlotView, dt *table.Table) *plotview.PlotView {
	plt.Params.Title = "Effort / Drive over Time Plot"
	plt.Params.XAxisColumn = "T"
	plt.SetTable(dt)
	// order of params: on, fixMin, min, fixMax, max
	plt.SetColParams("T", plotview.Off, plotview.FloatMin, 0, plotview.FloatMax, 0)
	plt.SetColParams("Eff", plotview.On, plotview.FixMin, 0, plotview.FixMax, 1)
	plt.SetColParams("EffInc", plotview.Off, plotview.FixMin, 0, plotview.FixMax, ss.Effort.Max)
	plt.SetColParams("Urge", plotview.On, plotview.FixMin, 0, plotview.FixMax, 1)
	plt.SetColParams("US", plotview.On, plotview.FixMin, 0, plotview.FixMax, 1)
	plt.SetColParams("Drive", plotview.On, plotview.FixMin, 0, plotview.FixMax, 1)
	return plt
}

// ConfigGUI configures the Cogent Core GUI interface for this simulation.
func (ss *DrEffPlot) ConfigGUI() *core.Body {
	b := core.NewBody("Drive / Effort / Urgency Plotting")

	split := core.NewSplits(b, "split")
	sv := views.NewStructView(split, "sv")
	sv.SetStruct(ss)

	tv := core.NewTabs(split, "tv")

	ss.Plot = plotview.NewSubPlot(tv.NewTab("Effort Plot"))
	ss.ConfigPlot(ss.Plot, ss.Table)

	ss.TimePlot = plotview.NewSubPlot(tv.NewTab("TimePlot"))
	ss.ConfigTimePlot(ss.TimePlot, ss.TimeTable)

	split.SetSplits(.3, .7)

	b.AddAppBar(func(tb *core.Toolbar) {
		views.NewFuncButton(tb, ss.EffortPlot).SetIcon(icons.PlayArrow)
		views.NewFuncButton(tb, ss.UrgencyPlot).SetIcon(icons.PlayArrow)
		views.NewFuncButton(tb, ss.TimeRun).SetIcon(icons.PlayArrow)
	})

	return b
}
