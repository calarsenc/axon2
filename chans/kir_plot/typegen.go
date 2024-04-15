// Code generated by "core generate -add-types"; DO NOT EDIT.

package main

import (
	"cogentcore.org/core/types"
)

var _ = types.AddType(&types.Type{Name: "main.Sim", IDName: "sim", Doc: "Sim holds the params, table, etc", Methods: []types.Method{{Name: "VmRun", Doc: "VmRun plots the equation as a function of V", Directives: []types.Directive{{Tool: "types", Directive: "add"}}}, {Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: []types.Directive{{Tool: "types", Directive: "add"}}}}, Fields: []types.Field{{Name: "Kir", Doc: "kIR function"}, {Name: "Vstart", Doc: "starting voltage"}, {Name: "Vend", Doc: "ending voltage"}, {Name: "Vstep", Doc: "voltage increment"}, {Name: "TimeSteps", Doc: "number of time steps"}, {Name: "TimeSpike", Doc: "do spiking instead of voltage ramp"}, {Name: "SpikeFreq", Doc: "spiking frequency"}, {Name: "TimeVstart", Doc: "time-run starting membrane potential"}, {Name: "TimeVend", Doc: "time-run ending membrane potential"}, {Name: "Table", Doc: "table for plot"}, {Name: "Plot", Doc: "the plot"}, {Name: "TimeTable", Doc: "table for plot"}, {Name: "TimePlot", Doc: "the plot"}}})
