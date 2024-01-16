// Code generated by "goki generate ./..."; DO NOT EDIT.

package main

import (
	"goki.dev/gti"
)

var _ = gti.AddType(&gti.Type{Name: "main.Sim", IDName: "sim", Doc: "Sim holds the params, table, etc", Methods: []gti.Method{{Name: "VmRun", Doc: "VmRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "VGRun", Doc: "VGRun runs the V-G equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "SGRun", Doc: "SGRun runs the spike-g equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "VmRun", Doc: "VmRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "Run", Doc: "Run runs the equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "CaRun", Doc: "CaRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "CamRun", Doc: "CamRun plots the equation as a function of Ca", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "VmRun", Doc: "VmRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "Run", Doc: "Run runs the equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}}, Fields: []gti.Field{{Name: "CaDt", Doc: "Ca time constants"}, {Name: "Minit"}, {Name: "Pinit"}, {Name: "Dinit"}, {Name: "MdtAdj", Doc: "adjustment to dt to account for discrete time updating"}, {Name: "PdtAdj", Doc: "adjustment to dt to account for discrete time updating"}, {Name: "DdtAdj", Doc: "adjustment to dt to account for discrete time updating"}, {Name: "TimeSteps", Doc: "number of time steps"}, {Name: "Table", Doc: "table for plot"}, {Name: "Plot", Doc: "the plot"}, {Name: "TimeTable", Doc: "table for plot"}, {Name: "TimePlot", Doc: "the plot"}}})
