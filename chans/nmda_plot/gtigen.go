// Code generated by "core generate ./..."; DO NOT EDIT.

package main

import (
	"cogentcore.org/core/gti"
)

var _ = gti.AddType(&gti.Type{Name: "main.Sim", IDName: "sim", Doc: "Sim holds the params, table, etc", Methods: []gti.Method{{Name: "VmRun", Doc: "VmRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "VGRun", Doc: "VGRun runs the V-G equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "SGRun", Doc: "SGRun runs the spike-g equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "VmRun", Doc: "VmRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "Run", Doc: "Run runs the equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}}, Fields: []gti.Field{{Name: "NMDAStd", Doc: "standard NMDA implementation in chans"}, {Name: "NMDAv", Doc: "multiplier on NMDA as function of voltage"}, {Name: "MgC", Doc: "magnesium ion concentration -- somewhere between 1 and 1.5"}, {Name: "NMDAd", Doc: "denominator of NMDA function"}, {Name: "NMDAerev", Doc: "NMDA reversal / driving potential"}, {Name: "BugVoff", Doc: "for old buggy NMDA: voff value to use"}, {Name: "Vstart", Doc: "starting voltage"}, {Name: "Vend", Doc: "ending voltage"}, {Name: "Vstep", Doc: "voltage increment"}, {Name: "Tau", Doc: "decay time constant for NMDA current -- rise time is 2 msec and not worth extra effort for biexponential"}, {Name: "TimeSteps", Doc: "number of time steps"}, {Name: "TimeV", Doc: "voltage for TimeRun"}, {Name: "TimeGin", Doc: "NMDA Gsyn current input at every time step"}, {Name: "Table", Doc: "table for plot"}, {Name: "Plot", Doc: "the plot"}, {Name: "TimeTable", Doc: "table for plot"}, {Name: "TimePlot", Doc: "the plot"}}})
