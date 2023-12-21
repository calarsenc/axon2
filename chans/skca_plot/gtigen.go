// Code generated by "goki generate ./..."; DO NOT EDIT.

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
		{"SKCa", &gti.Field{Name: "SKCa", Type: "invalid type", LocalType: "chans.SKCaParams", Doc: "SKCa params", Directives: gti.Directives{}, Tag: ""}},
		{"CaParams", &gti.Field{Name: "CaParams", Type: "invalid type", LocalType: "kinase.CaParams", Doc: "time constants for integrating Ca from spiking across M, P and D cascading levels", Directives: gti.Directives{}, Tag: ""}},
		{"NoSpikeThr", &gti.Field{Name: "NoSpikeThr", Type: "float32", LocalType: "float32", Doc: "threshold of SK M gating factor above which the neuron cannot spike", Directives: gti.Directives{}, Tag: "def:\"0.5\""}},
		{"CaStep", &gti.Field{Name: "CaStep", Type: "float32", LocalType: "float32", Doc: "Ca conc increment for M gating func plot", Directives: gti.Directives{}, Tag: "def:\"0.05\""}},
		{"TimeSteps", &gti.Field{Name: "TimeSteps", Type: "int", LocalType: "int", Doc: "number of time steps", Directives: gti.Directives{}, Tag: ""}},
		{"TimeSpike", &gti.Field{Name: "TimeSpike", Type: "bool", LocalType: "bool", Doc: "do spiking instead of Ca conc ramp", Directives: gti.Directives{}, Tag: ""}},
		{"SpikeFreq", &gti.Field{Name: "SpikeFreq", Type: "float32", LocalType: "float32", Doc: "spiking frequency", Directives: gti.Directives{}, Tag: ""}},
		{"Table", &gti.Field{Name: "Table", Type: "*goki.dev/etable/v2/etable.Table", LocalType: "*etable.Table", Doc: "table for plot", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"Plot", &gti.Field{Name: "Plot", Type: "*goki.dev/etable/v2/eplot.Plot2D", LocalType: "*eplot.Plot2D", Doc: "the plot", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"TimeTable", &gti.Field{Name: "TimeTable", Type: "*goki.dev/etable/v2/etable.Table", LocalType: "*etable.Table", Doc: "table for plot", Directives: gti.Directives{}, Tag: "view:\"no-inline\""}},
		{"TimePlot", &gti.Field{Name: "TimePlot", Type: "*goki.dev/etable/v2/eplot.Plot2D", LocalType: "*eplot.Plot2D", Doc: "the plot", Directives: gti.Directives{}, Tag: "view:\"-\""}},
	}),
	Embeds: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{
		{"CamRun", &gti.Method{Name: "CamRun", Doc: "CamRun plots the equation as a function of Ca", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"TimeRun", &gti.Method{Name: "TimeRun", Doc: "TimeRun runs the equation over time.", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
	}),
})
