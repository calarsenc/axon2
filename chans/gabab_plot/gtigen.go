// Code generated by "goki generate ./..."; DO NOT EDIT.

package main

import (
	"goki.dev/gti"
)

var _ = gti.AddType(&gti.Type{Name: "main.Sim", IDName: "sim", Doc: "Sim holds the params, table, etc", Methods: []gti.Method{{Name: "VmRun", Doc: "VmRun plots the equation as a function of V", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "VGRun", Doc: "VGRun runs the V-G equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "SGRun", Doc: "SGRun runs the spike-g equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}}}, Fields: []gti.Field{{Name: "GABAstd", Doc: "standard chans version of GABAB"}, {Name: "GABAbv", Doc: "multiplier on GABAb as function of voltage"}, {Name: "GABAbo", Doc: "offset of GABAb function"}, {Name: "GABAberev", Doc: "GABAb reversal / driving potential"}, {Name: "Vstart", Doc: "starting voltage"}, {Name: "Vend", Doc: "ending voltage"}, {Name: "Vstep", Doc: "voltage increment"}, {Name: "Smax", Doc: "max number of spikes"}, {Name: "RiseTau", Doc: "rise time constant"}, {Name: "DecayTau", Doc: "decay time constant -- must NOT be same as RiseTau"}, {Name: "GsXInit", Doc: "initial value of GsX driving variable at point of synaptic input onset -- decays expoentially from this start"}, {Name: "MaxTime", Doc: "time when peak conductance occurs, in TimeInc units"}, {Name: "TauFact", Doc: "time constant factor used in integration: (Decay / Rise) ^ (Rise / (Decay - Rise))"}, {Name: "TimeSteps", Doc: "total number of time steps to take"}, {Name: "TimeInc", Doc: "time increment per step"}, {Name: "VGTable", Doc: "table for plot"}, {Name: "SGTable", Doc: "table for plot"}, {Name: "TimeTable", Doc: "table for plot"}, {Name: "VGPlot", Doc: "the plot"}, {Name: "SGPlot", Doc: "the plot"}, {Name: "TimePlot", Doc: "the plot"}}})
