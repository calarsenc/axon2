// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deep

import (
	"fmt"
	"log"
	"math"

	"github.com/emer/axon/axon"
	"github.com/goki/ki/kit"
	"github.com/goki/mat32"
)

// Driver describes the source of driver inputs from cortex into TRC (pulvinar)
type Driver struct {
	Driver string `desc:"driver layer"`
	Off    int    `inactive:"-" desc:"offset into TRC pool"`
}

// Drivers are a list of drivers
type Drivers []*Driver

// Add adds new driver(s)
func (dr *Drivers) Add(laynms ...string) {
	for _, laynm := range laynms {
		d := &Driver{}
		d.Driver = laynm
		*dr = append(*dr, d)
	}
}

// AddOne adds one new driver -- python does not work with varargs
func (dr *Drivers) AddOne(laynm string) {
	d := &Driver{}
	d.Driver = laynm
	*dr = append(*dr, d)
}

// TRCParams provides parameters for how the plus-phase (outcome) state of thalamic relay cell
// (e.g., Pulvinar) neurons is computed from the corresponding driver neuron Burst activation.
// Drivers are hard clamped using Clamp.Rate.
type TRCParams struct {
	DriversOff bool          `def:"false" desc:"Turn off the driver inputs, in which case this layer behaves like a standard layer"`
	BurstQtr   axon.Quarters `desc:"Quarter(s) when bursting occurs -- typically Q4 but can also be Q2 and Q4 for beta-frequency updating.  Note: this is a bitflag and must be accessed using its Set / Has etc routines"`
	DriveScale float32       `def:"0.05" min:"0.0" desc:"multiplier on driver input strength, multiplies activation of driver layer"`
	MaxInhib   float32       `def:"0.6" min:"0.01" desc:"Level of Max driver layer activation at which the predictive non-burst inputs are fully inhibited.  Computationally, it is essential that driver inputs inhibit effect of predictive non-driver (CTLayer) inputs, so that the plus phase is not always just the minus phase plus something extra (the error will never go to zero then).  When max driver act input exceeds this value, predictive non-driver inputs are fully suppressed.  If there is only weak burst input however, then the predictive inputs remain and this critically prevents the network from learning to turn activation off, which is difficult and severely degrades learning."`
	HardClamp  bool          `desc:"driver plus phase is hard clamped instead of soft clamped"`
	DriveMin   float32       `def:"0.2" min:"0.01" desc:"Minimum layer Max activation for a layer to produce driving activation -- if a layer has less activation than this, it is not capable of driving."`
	NoTopo     bool          `desc:"Do not treat the pools in this layer as topographically organized relative to driver inputs -- all drivers compress down to give same input to all pools"`
	AvgMix     float32       `min:"0" max:"1" desc:"proportion of average across driver pools that is combined with Max to provide some graded tie-breaker signal -- especially important for large pool downsampling, e.g., when doing NoTopo"`
	Binarize   bool          `desc:"Apply threshold to driver burst input for computing plus-phase activations -- above BinThr, then Act = BinOn, below = BinOff.  This is beneficial for layers with weaker graded activations, such as V1 or other perceptual inputs."`
	BinThr     float32       `viewif:"Binarize" desc:"Threshold for binarizing in terms of sending Burst activation"`
	BinOn      float32       `def:"0.3" viewif:"Binarize" desc:"Resulting driver Ge value for units above threshold -- lower value around 0.3 or so seems best (DriveScale is NOT applied -- generally same range as that)."`
	BinOff     float32       `def:"0" viewif:"Binarize" desc:"Resulting driver Ge value for units below threshold -- typically 0."`
}

func (tp *TRCParams) Update() {
}

func (tp *TRCParams) Defaults() {
	tp.BurstQtr.Set(int(axon.Q4))
	tp.DriveScale = 0.05
	tp.MaxInhib = 0.6
	tp.DriveMin = 0.2
	tp.Binarize = false
	tp.BinThr = 0.4
	tp.BinOn = 0.3
	tp.BinOff = 0
}

// DriveGe returns effective excitatory conductance to use for given driver input Burst activation
func (tp *TRCParams) DriveGe(act float32) float32 {
	if tp.Binarize {
		if act >= tp.BinThr {
			return tp.BinOn
		} else {
			return tp.BinOff
		}
	} else {
		if tp.HardClamp {
			return act
		} else {
			return tp.DriveScale * act
		}
	}
}

// GeFmMaxAvg returns the drive Ge value as function of max and average
func (tp *TRCParams) GeFmMaxAvg(max, avg float32) float32 {
	deff := (1-tp.AvgMix)*max + tp.AvgMix*avg
	return tp.DriveGe(deff)
}

// TRCLayer is the thalamic relay cell layer for DeepAxon.
// It has normal activity during the minus phase, as activated by CT etc inputs,
// and is then driven by strong 5IB driver inputs in the plus phase.
// For attentional modulation, TRC maintains pool-level correspondence with CT inputs
// which creates challenges for aligning with driver inputs.
// * Max operation used to integrate across multiple drivers, where necessary,
//   e.g., multiple driver pools map onto single TRC pool (common feedforward theme),
//   *even when there is no logical connection for the i'th unit in each pool* --
//   to make this dimensionality reduction more effective, using lateral connectivity
//   between pools that favors this correspondence is beneficial.  Overall, this is
//   consistent with typical DCNN max pooling organization.
// * Typically, pooled 4D TRC layers should have fewer pools than driver layers,
//   in which case the respective pool geometry is interpolated.  Ideally, integer size
//   differences are best (e.g., driver layer has 2x pools vs TRC).
// * Pooled 4D TRC layer should in general not predict flat 2D drivers, but if so
//   the drivers are replicated for each pool.
// * Similarly, there shouldn't generally be more TRC pools than driver pools, but
//   if so, drivers replicate across pools.
type TRCLayer struct {
	TopoInhibLayer           // access as .TopoInhibLayer
	TRC            TRCParams `view:"inline" desc:"parameters for computing TRC plus-phase (outcome) activations based on Burst activation from corresponding driver neuron"`
	Drivers        Drivers   `desc:"name of SuperLayer that sends 5IB Burst driver inputs to this layer"`
}

var KiT_TRCLayer = kit.Types.AddType(&TRCLayer{}, LayerProps)

func (ly *TRCLayer) Defaults() {
	ly.TopoInhibLayer.Defaults()
	ly.Act.Init.Decay = 0.5
	ly.Act.GABAB.Gbar = 0.005 // output layer settings
	ly.Act.NMDA.Gbar = 0.01
	ly.TRC.Defaults()
	ly.TopoInhib.Defaults()
	ly.Typ = TRC
}

// UpdateParams updates all params given any changes that might have been made to individual values
// including those in the receiving projections of this layer
func (ly *TRCLayer) UpdateParams() {
	ly.TopoInhibLayer.UpdateParams()
	ly.TRC.Update()
	ly.TopoInhib.Update()
}

func (ly *TRCLayer) Class() string {
	return "TRC " + ly.Cls
}

func (ly *TRCLayer) IsTarget() bool {
	return true // We are a Target-like layer: don't do various adaptive steps
}

///////////////////////////////////////////////////////////////////////////////////////
// Drivers

func (ly *TRCLayer) InitWts() {
	ly.TopoInhibLayer.InitWts()
	ly.SetDriverOffs()
}

// UnitsSize returns the dimension of the units, either within a pool for 4D, or layer for 2D
func UnitsSize(ly *axon.Layer) (x, y int) {
	if ly.Is4D() {
		y = ly.Shp.Dim(2)
		x = ly.Shp.Dim(3)
	} else {
		y = ly.Shp.Dim(0)
		x = ly.Shp.Dim(1)
	}
	return
}

// DriverLayer returns the driver layer for given Driver
func (ly *TRCLayer) DriverLayer(drv *Driver) (*axon.Layer, error) {
	tly, err := ly.Network.LayerByNameTry(drv.Driver)
	if err != nil {
		err = fmt.Errorf("TRCLayer %s: Driver Layer: %v", ly.Name(), err)
		log.Println(err)
		return nil, err
	}
	return tly.(axon.AxonLayer).AsAxon(), nil
}

// SetDriverOffs sets the driver offsets
func (ly *TRCLayer) SetDriverOffs() error {
	mx, my := UnitsSize(&ly.Layer)
	mn := my * mx
	off := 0
	var err error
	for _, drv := range ly.Drivers {
		var dl *axon.Layer
		dl, err = ly.DriverLayer(drv)
		if err != nil {
			continue
		}
		drv.Off = off
		x, y := UnitsSize(dl)
		off += y * x
	}
	if off > mn {
		err = fmt.Errorf("TRCLayer %s: size of drivers: %d is greater than units: %d", ly.Name(), off, mn)
		log.Println(err)
	}
	return err
}

func DriveAct(dni int, dly *axon.Layer, sly *SuperLayer, issuper bool) float32 {
	act := float32(0)
	if issuper {
		act = sly.SuperNeurs[dni].Burst
	} else {
		act = dly.Neurons[dni].Act
	}
	lmax := dly.Pools[0].Inhib.Act.Max // normalize by drive layer max act
	if lmax > 0.1 {                    // this puts all layers on equal footing for driving..
		return act / lmax
	}
	return act
}

// SetDriverNeuron sets the driver activation for given Neuron,
// based on given Ge driving value (use DriveFmMaxAvg) from driver layer (Burst or Act)
func (ly *TRCLayer) SetDriverNeuron(tni int, drvGe, drvInhib float32, cyc int) {
	if tni >= len(ly.Neurons) {
		return
	}
	nrn := &ly.Neurons[tni]
	if nrn.IsOff() {
		return
	}
	if ly.TRC.HardClamp {
		nrn.Ext = drvGe
		clrmsk, setmsk, _ := ly.ApplyExtFlags()
		nrn.ClearMask(clrmsk)
		nrn.SetMask(setmsk)
	} else {
		geRaw := (1-drvInhib)*nrn.GeRaw + drvGe

		nrn.NMDA = ly.Act.NMDA.NMDA(nrn.NMDA, geRaw, nrn.NMDASyn)
		nrn.Gnmda = ly.Act.NMDA.Gnmda(nrn.NMDA, nrn.VmDend)
		// note: GABAB integrated in ActFmG one timestep behind, b/c depends on integrated Gi inhib

		// note: each step broken out here so other variants can add extra terms to Raw
		ly.Act.GeFmRaw(nrn, geRaw+nrn.Gnmda, cyc, nrn.ActM)
		nrn.GeRaw = 0
		ly.Act.GiFmRaw(nrn, nrn.GiRaw)
		nrn.GiRaw = 0
	}
}

// SetDriverActs sets the driver activations, integrating across all the driver layers
func (ly *TRCLayer) SetDriverActs(ltime *axon.Time) {
	cyc := ltime.Cycle // for bursting
	if ly.IsTarget() {
		cyc = ltime.QuarterCycle()
	}
	nux, nuy := UnitsSize(&ly.Layer)
	nun := nux * nuy
	pyn := ly.Shp.Dim(0)
	pxn := ly.Shp.Dim(1)
	for _, drv := range ly.Drivers {
		dly, err := ly.DriverLayer(drv)
		if err != nil {
			continue
		}
		sly, issuper := dly.AxonLay.(*SuperLayer)
		drvMax := dly.Pools[0].Inhib.Act.Max
		if ly.TRC.HardClamp && drvMax < ly.TRC.DriveMin {
			continue
		}
		drvInhib := mat32.Min(1, drvMax/ly.TRC.MaxInhib)

		if dly.Is2D() {
			if ly.Is2D() {
				for dni := range dly.Neurons {
					tni := drv.Off + dni
					drvAct := DriveAct(dni, dly, sly, issuper)
					ly.SetDriverNeuron(tni, ly.TRC.GeFmMaxAvg(drvAct, drvAct), drvInhib, cyc)
				}
			} else { // copy flat to all pools -- not typical
				for dni := range dly.Neurons {
					drvAct := DriveAct(dni, dly, sly, issuper)
					tni := drv.Off + dni
					for py := 0; py < pyn; py++ {
						for px := 0; px < pxn; px++ {
							pni := (py*pxn+px)*nun + tni
							ly.SetDriverNeuron(pni, ly.TRC.GeFmMaxAvg(drvAct, drvAct), drvInhib, cyc)
						}
					}
				}
			}
		} else { // dly is 4D
			dpyn := dly.Shp.Dim(0)
			dpxn := dly.Shp.Dim(1)
			duxn, duyn := UnitsSize(dly)
			dnun := duxn * duyn
			if ly.Is2D() {
				for dni := 0; dni < dnun; dni++ {
					max := float32(0)
					avg := float32(0)
					avgn := 0
					for py := 0; py < dpyn; py++ {
						for px := 0; px < dpxn; px++ {
							pi := (py*dpxn + px)
							pni := pi*dnun + dni
							act := DriveAct(pni, dly, sly, issuper)
							max = mat32.Max(max, act)
							pmax := dly.Pools[1+pi].Inhib.Act.Max
							if pmax > 0.5 {
								avg += act
								avgn++
							}
						}
					}
					if avgn > 0 {
						avg /= float32(avgn)
					}
					tni := drv.Off + dni
					ly.SetDriverNeuron(tni, ly.TRC.GeFmMaxAvg(max, avg), drvInhib, cyc)
				}
			} else if ly.TRC.NoTopo { // ly is 4D
				for dni := 0; dni < dnun; dni++ {
					max := float32(0)
					avg := float32(0)
					avgn := 0
					for py := 0; py < dpyn; py++ {
						for px := 0; px < dpxn; px++ {
							pi := (py*dpxn + px)
							pni := pi*dnun + dni
							act := DriveAct(pni, dly, sly, issuper)
							max = mat32.Max(max, act)
							pmax := dly.Pools[1+pi].Inhib.Act.Max
							if pmax > 0.5 {
								avg += act
								avgn++
							}
						}
					}
					if avgn > 0 {
						avg /= float32(avgn)
					}
					drvGe := ly.TRC.GeFmMaxAvg(max, avg)
					tni := drv.Off + dni
					for py := 0; py < pyn; py++ {
						for px := 0; px < pxn; px++ {
							pni := (py*pxn+px)*nun + tni
							ly.SetDriverNeuron(pni, drvGe, drvInhib, cyc)
						}
					}
				}
			} else { // ly is 4D
				pyr := float64(dpyn) / float64(pyn)
				pxr := float64(dpxn) / float64(pxn)
				for py := 0; py < pyn; py++ {
					sdpy := int(math.Round(float64(py) * pyr))
					edpy := int(math.Round(float64(py+1) * pyr))
					for px := 0; px < pxn; px++ {
						sdpx := int(math.Round(float64(px) * pxr))
						edpx := int(math.Round(float64(px+1) * pxr))
						pni := (py*pxn + px) * nun
						for dni := 0; dni < dnun; dni++ {
							max := float32(0)
							avg := float32(0)
							avgn := 0
							for dpy := sdpy; dpy < edpy; dpy++ {
								for dpx := sdpx; dpx < edpx; dpx++ {
									pi := (dpy*dpxn + dpx)
									dpni := pi*dnun + dni
									act := DriveAct(dpni, dly, sly, issuper)
									max = mat32.Max(max, act)
									pmax := dly.Pools[1+pi].Inhib.Act.Max
									if pmax > 0.5 {
										avg += act
										avgn++
									}
								}
							}
							if avgn > 0 {
								avg /= float32(avgn)
							}
							tni := pni + drv.Off + dni
							ly.SetDriverNeuron(tni, ly.TRC.GeFmMaxAvg(max, avg), drvInhib, cyc)
						}
					}
				}
			}
		}
	}
}

// GFmInc integrates new synaptic conductances from increments sent during last SendGDelta.
func (ly *TRCLayer) GFmInc(ltime *axon.Time) {
	ly.RecvGInc(ltime)
	if ly.TRC.HardClamp || ly.TRC.DriversOff || !ly.TRC.BurstQtr.Has(ltime.Quarter) {
		ly.GFmIncNeur(ltime) // regular
		return
	}
	ly.SetDriverActs(ltime)
}

// ActFmG computes rate-code activation from Ge, Gi, Gl conductances
// and updates learning running-average activations from that Act
func (ly *TRCLayer) ActFmG(ltime *axon.Time) {
	if !ly.TRC.DriversOff && ly.TRC.BurstQtr.Has(ltime.Quarter) {
		ly.SetDriverActs(ltime)
	}
	ly.TopoInhibLayer.ActFmG(ltime)
}

// QuarterFinal does updating after end of a quarter
func (ly *TRCLayer) QuarterFinal(ltime *axon.Time) {
	if ly.TRC.HardClamp && ltime.Quarter == 2 {
		for ni := range ly.Neurons {
			nrn := &ly.Neurons[ni]
			if nrn.IsOff() {
				continue
			}
			// reset in prep for clamping next quarter
			nrn.ISI = -1
			nrn.ISIAvg = -1
		}
	}
	ly.TopoInhibLayer.QuarterFinal(ltime)
}
