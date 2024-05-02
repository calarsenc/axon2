// Copyright (c) 2020, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// vgcc_plot plots an equation updating over time in a table.Table and PlotView.
package main

//go:generate core generate -add-types

import (
	"strconv"

	"cogentcore.org/core/core"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/plot/plotview"
	"cogentcore.org/core/tensor/table"
	"cogentcore.org/core/views"
	"github.com/emer/axon/v2/chans"
)

func main() {
	sim := &Sim{}
	sim.Config()
	sim.VmRun()
	b := sim.ConfigGUI()
	b.RunMainWindow()
}

// LogPrec is precision for saving float values in logs
const LogPrec = 4

// Sim holds the params, table, etc
type Sim struct {

	// VGCC function
	VGCC chans.VGCCParams

	// starting voltage
	Vstart float32 `default:"-90"`

	// ending voltage
	Vend float32 `default:"0"`

	// voltage increment
	Vstep float32 `default:"1"`

	// number of time steps
	TimeSteps int

	// do spiking instead of voltage ramp
	TimeSpike bool

	// spiking frequency
	SpikeFreq float32

	// time-run starting membrane potential
	TimeVstart float32

	// time-run ending membrane potential
	TimeVend float32

	// table for plot
	Table *table.Table `view:"no-inline"`

	// the plot
	Plot *plotview.PlotView `view:"-"`

	// table for plot
	TimeTable *table.Table `view:"no-inline"`

	// the plot
	TimePlot *plotview.PlotView `view:"-"`
}

// Config configures all the elements using the standard functions
func (ss *Sim) Config() {
	ss.VGCC.Defaults()
	ss.VGCC.Gbar = 1
	ss.Vstart = -90
	ss.Vend = 2
	ss.Vstep = 0.01
	ss.TimeSteps = 200
	ss.TimeSpike = true
	ss.SpikeFreq = 50
	ss.TimeVstart = -70
	ss.TimeVend = -20
	ss.Update()
	ss.Table = &table.Table{}
	ss.ConfigTable(ss.Table)
	ss.TimeTable = &table.Table{}
	ss.ConfigTimeTable(ss.TimeTable)
}

// Update updates computed values
func (ss *Sim) Update() {
}

// VmRun plots the equation as a function of V
func (ss *Sim) VmRun() { //types:add
	ss.Update()
	dt := ss.Table

	nv := int((ss.Vend - ss.Vstart) / ss.Vstep)
	dt.SetNumRows(nv)
	for vi := 0; vi < nv; vi++ {
		v := ss.Vstart + float32(vi)*ss.Vstep
		vnorm := chans.VFromBio(v)
		g := ss.VGCC.GFromV(vnorm)
		m := ss.VGCC.MFromV(v)
		h := ss.VGCC.HFromV(v)
		var dm, dh float32
		ss.VGCC.DMHFromV(vnorm, m, h, &dm, &dh)

		dt.SetFloat("V", vi, float64(v))
		dt.SetFloat("Gvgcc", vi, float64(g))
		dt.SetFloat("M", vi, float64(m))
		dt.SetFloat("H", vi, float64(h))
		dt.SetFloat("dM", vi, float64(dm))
		dt.SetFloat("dH", vi, float64(dh))
	}
	if ss.Plot != nil {
		ss.Plot.UpdatePlot()
	}
}

func (ss *Sim) ConfigTable(dt *table.Table) {
	dt.SetMetaData("name", "VgCcplotTable")
	dt.SetMetaData("read-only", "true")
	dt.SetMetaData("precision", strconv.Itoa(LogPrec))

	dt.AddFloat64Column("V")
	dt.AddFloat64Column("Gvgcc")
	dt.AddFloat64Column("M")
	dt.AddFloat64Column("H")
	dt.AddFloat64Column("dM")
	dt.AddFloat64Column("dH")
	dt.SetNumRows(0)
}

func (ss *Sim) ConfigPlot(plt *plotview.PlotView, dt *table.Table) *plotview.PlotView {
	plt.Params.Title = "VGCC V-G Function Plot"
	plt.Params.XAxisColumn = "V"
	plt.SetTable(dt)
	// order of params: on, fixMin, min, fixMax, max
	plt.SetColParams("V", plotview.Off, plotview.FloatMin, 0, plotview.FloatMax, 0)
	plt.SetColParams("Gvgcc", plotview.On, plotview.FixMin, 0, plotview.FloatMax, 0)
	plt.SetColParams("M", plotview.On, plotview.FixMin, 0, plotview.FloatMax, 0)
	plt.SetColParams("H", plotview.On, plotview.FixMin, 0, plotview.FloatMax, 0)
	return plt
}

/////////////////////////////////////////////////////////////////

// TimeRun runs the equation over time.
func (ss *Sim) TimeRun() { //types:add
	ss.Update()
	dt := ss.TimeTable

	m := float32(0)
	h := float32(1)
	msdt := float32(0.001)
	v := ss.TimeVstart
	vinc := float32(2) * (ss.TimeVend - ss.TimeVstart) / float32(ss.TimeSteps)

	isi := int(1000 / ss.SpikeFreq)

	dt.SetNumRows(ss.TimeSteps)
	var g float32
	for ti := 0; ti < ss.TimeSteps; ti++ {
		vnorm := chans.VFromBio(v)
		t := float32(ti) * msdt
		g = ss.VGCC.Gvgcc(vnorm, m, h)
		var dm, dh float32
		ss.VGCC.DMHFromV(vnorm, m, h, &dm, &dh)
		m += dm
		h += dh

		dt.SetFloat("Time", ti, float64(t))
		dt.SetFloat("V", ti, float64(v))
		dt.SetFloat("Gvgcc", ti, float64(g))
		dt.SetFloat("M", ti, float64(m))
		dt.SetFloat("H", ti, float64(h))
		dt.SetFloat("dM", ti, float64(dm))
		dt.SetFloat("dH", ti, float64(dh))

		if ss.TimeSpike {
			if ti%isi < 3 {
				v = ss.TimeVend
			} else {
				v = ss.TimeVstart
			}
		} else {
			v += vinc
			if v > ss.TimeVend {
				v = ss.TimeVend
			}
		}
	}
	if ss.TimePlot != nil {
		ss.TimePlot.UpdatePlot()
	}
}

func (ss *Sim) ConfigTimeTable(dt *table.Table) {
	dt.SetMetaData("name", "VgCcplotTable")
	dt.SetMetaData("read-only", "true")
	dt.SetMetaData("precision", strconv.Itoa(LogPrec))

	dt.AddFloat64Column("Time")
	dt.AddFloat64Column("V")
	dt.AddFloat64Column("Gvgcc")
	dt.AddFloat64Column("M")
	dt.AddFloat64Column("H")
	dt.AddFloat64Column("dM")
	dt.AddFloat64Column("dH")
	dt.SetNumRows(0)
}

func (ss *Sim) ConfigTimePlot(plt *plotview.PlotView, dt *table.Table) *plotview.PlotView {
	plt.Params.Title = "Time Function Plot"
	plt.Params.XAxisColumn = "Time"
	plt.SetTable(dt)
	// order of params: on, fixMin, min, fixMax, max
	plt.SetColParams("Time", plotview.Off, plotview.FloatMin, 0, plotview.FloatMax, 0)
	plt.SetColParams("Gvgcc", plotview.On, plotview.FixMin, 0, plotview.FloatMax, 0)
	plt.SetColParams("M", plotview.Off, plotview.FixMin, 0, plotview.FloatMax, 0)
	plt.SetColParams("H", plotview.Off, plotview.FixMin, 0, plotview.FloatMax, 0)
	plt.SetColParams("dM", plotview.Off, plotview.FloatMin, 0, plotview.FloatMax, 0)
	plt.SetColParams("dH", plotview.Off, plotview.FloatMin, 0, plotview.FloatMax, 0)
	return plt
}

// ConfigGUI configures the Cogent Core GUI interface for this simulation.
func (ss *Sim) ConfigGUI() *core.Body {
	b := core.NewBody("Vgcc Plot")

	split := core.NewSplits(b, "split")
	sv := views.NewStructView(split, "sv")
	sv.SetStruct(ss)

	tv := core.NewTabs(split, "tv")

	ss.Plot = plotview.NewSubPlot(tv.NewTab("V-G Plot"))
	ss.ConfigPlot(ss.Plot, ss.Table)

	ss.TimePlot = plotview.NewSubPlot(tv.NewTab("TimePlot"))
	ss.ConfigTimePlot(ss.TimePlot, ss.TimeTable)

	split.SetSplits(.3, .7)

	b.AddAppBar(func(tb *core.Toolbar) {
		views.NewFuncButton(tb, ss.VmRun).SetIcon(icons.PlayArrow)
		views.NewFuncButton(tb, ss.TimeRun).SetIcon(icons.PlayArrow)
	})

	return b
}
