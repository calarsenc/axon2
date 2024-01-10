// Code generated by "goki generate -add-types"; DO NOT EDIT.

package main

import (
	"goki.dev/gti"
	"goki.dev/ordmap"
)

var _ = gti.AddType(&gti.Type{
	Name:       "main.Sim",
	ShortName:  "main.Sim",
	IDName:     "sim",
	Doc:        "Sim holds the params, table, etc",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"VGCC", &gti.Field{Name: "VGCC", Type: "github.com/emer/axon/v2/chans.VGCCParams", LocalType: "chans.VGCCParams", Doc: "VGCC function", Directives: gti.Directives{}, Tag: ""}},
		{"Vstart", &gti.Field{Name: "Vstart", Type: "float32", LocalType: "float32", Doc: "starting voltage", Directives: gti.Directives{}, Tag: "def:\"-90\""}},
		{"Vend", &gti.Field{Name: "Vend", Type: "float32", LocalType: "float32", Doc: "ending voltage", Directives: gti.Directives{}, Tag: "def:\"0\""}},
		{"Vstep", &gti.Field{Name: "Vstep", Type: "float32", LocalType: "float32", Doc: "voltage increment", Directives: gti.Directives{}, Tag: "def:\"1\""}},
		{"TimeSteps", &gti.Field{Name: "TimeSteps", Type: "int", LocalType: "int", Doc: "number of time steps", Directives: gti.Directives{}, Tag: ""}},
		{"TimeSpike", &gti.Field{Name: "TimeSpike", Type: "bool", LocalType: "bool", Doc: "do spiking instead of voltage ramp", Directives: gti.Directives{}, Tag: ""}},
		{"SpikeFreq", &gti.Field{Name: "SpikeFreq", Type: "float32", LocalType: "float32", Doc: "spiking frequency", Directives: gti.Directives{}, Tag: ""}},
		{"TimeVstart", &gti.Field{Name: "TimeVstart", Type: "float32", LocalType: "float32", Doc: "time-run starting membrane potential", Directives: gti.Directives{}, Tag: ""}},
		{"TimeVend", &gti.Field{Name: "TimeVend", Type: "float32", LocalType: "float32", Doc: "time-run ending membrane potential", Directives: gti.Directives{}, Tag: ""}},
		{"Table", &gti.Field{Name: "Table", Type: "*goki.dev/etable/v2/etable.Table", LocalType: "*etable.Table", Doc: "table for plot", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"Plot", &gti.Field{Name: "Plot", Type: "*goki.dev/etable/v2/eplot.Plot2D", LocalType: "*eplot.Plot2D", Doc: "the plot", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"TimeTable", &gti.Field{Name: "TimeTable", Type: "*goki.dev/etable/v2/etable.Table", LocalType: "*etable.Table", Doc: "table for plot", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"TimePlot", &gti.Field{Name: "TimePlot", Type: "*goki.dev/etable/v2/eplot.Plot2D", LocalType: "*eplot.Plot2D", Doc: "the plot", Directives: gti.Directives{}, Tag: "view:\"-\""}},
	}),
	Embeds: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{
		{"VmRun", &gti.Method{Name: "VmRun", Doc: "VmRun plots the equation as a function of V", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"TimeRun", &gti.Method{Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
	}),
})
