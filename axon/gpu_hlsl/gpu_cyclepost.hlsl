// Copyright (c) 2022, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// does CyclePost: iterates over data parallel -- handles all special context updates

#include "synmem.hlsl"

// note: all must be visible always because accessor methods refer to them
[[vk::binding(0, 1)]] StructuredBuffer<uint> NeuronIxs; // [Neurons][Idxs]
[[vk::binding(1, 1)]] StructuredBuffer<uint> SynapseIxs;  // [Layer][SendPrjns][SendNeurons][Syns]
[[vk::binding(1, 2)]] RWStructuredBuffer<float> Neurons; // [Neurons][Vars][Data]
[[vk::binding(2, 2)]] RWStructuredBuffer<float> NeuronAvgs; // [Neurons][Vars]
[[vk::binding(0, 3)]] RWStructuredBuffer<SynMemBlock> Synapses;  // [Layer][SendPrjns][SendNeurons][Syns]
[[vk::binding(1, 3)]] RWStructuredBuffer<SynMemBlock> SynapseCas;  // [Layer][SendPrjns][SendNeurons][Syns][Data]

#include "context.hlsl"
#include "layerparams.hlsl"

// note: binding is var, set

// Set 0: uniform layer params -- could not have prjns also be uniform..
[[vk::binding(0, 0)]] StructuredBuffer<LayerParams> Layers; // [Layer]

// Set 1: effectively uniform indexes and prjn params as structured buffers in storage

// Set 2: main network structs and vals -- all are writable
[[vk::binding(0, 2)]] RWStructuredBuffer<Context> Ctx; // [0]
[[vk::binding(3, 2)]] RWStructuredBuffer<Pool> Pools; // [Layer][Pools][Data]
[[vk::binding(4, 2)]] RWStructuredBuffer<LayerVals> LayVals; // [Layer][Data]


void CyclePostVSPatch(inout Context ctx, in LayerParams ly, uint li, uint di, int pi, in Pool pl) {
	ly.CyclePostVSPatchLayer(ctx, di, pi, pl);
}

float LDTSrcLayAct(int layIdx, uint di) {
	if (layIdx < 0) {
		return 0.0;
	}
	return Pools[Layers[layIdx].Idxs.PoolIdx(0, di)].AvgMax.CaSpkP.Cycle.Avg;
}

void CyclePostLDT(inout Context ctx, uint di, in LayerParams ly, inout LayerVals vals) {
	float srcLay1Act = LDTSrcLayAct(ly.LDT.SrcLay1Idx, di);
	float srcLay2Act = LDTSrcLayAct(ly.LDT.SrcLay2Idx, di);
	float srcLay3Act = LDTSrcLayAct(ly.LDT.SrcLay3Idx, di);
	float srcLay4Act = LDTSrcLayAct(ly.LDT.SrcLay4Idx, di);
	ly.CyclePostLDTLayer(ctx, di, vals, srcLay1Act, srcLay2Act, srcLay3Act, srcLay4Act);
}

void CyclePost2(inout Context ctx, in LayerParams ly, uint li, uint di, inout LayerVals vals, in Pool lpl) {
	switch (ly.LayType) {
	case PTNotMaintLayer: {
		ly.CyclePostPTNotMaintLayer(ctx, di, lpl);
		break;
	}
	case CeMLayer: {
		ly.CyclePostCeMLayer(ctx, di, lpl);
		break;
	}
	case VSPatchLayer: {
		int npl = ly.Idxs.ShpPlY * ly.Idxs.ShpPlX;
		for (int pi = 0; pi < npl; pi++) {
			CyclePostVSPatch(ctx, ly, li, di, pi+1, Pools[ly.Idxs.PoolIdx(1+pi, di)]);
		}
		break;
	}
	case LDTLayer: {
		CyclePostLDT(ctx, di, ly, vals);
		break;
	}
	case VTALayer: {
		ly.CyclePostVTALayer(ctx, di);
		break;
	}
	case RWDaLayer: {
		ly.CyclePostRWDaLayer(ctx, di, vals, LayVals[ctx.NetIdxs.ValsIdx(ly.RWDa.RWPredLayIdx, di)]);
		break;
	}
	case TDPredLayer: {
		ly.CyclePostTDPredLayer(ctx, di, vals);
		break;
	}
	case TDIntegLayer: {
		ly.CyclePostTDIntegLayer(ctx, di, vals, LayVals[ctx.NetIdxs.ValsIdx(ly.TDInteg.TDPredLayIdx, di)]);
		break;
	}
	case TDDaLayer: {
		ly.CyclePostTDDaLayer(ctx, di, vals, LayVals[ctx.NetIdxs.ValsIdx(ly.TDDa.TDIntegLayIdx, di)]);
		break;
	}
	}
}

void CyclePost(inout Context ctx, in LayerParams ly, int li, uint di) {
	CyclePost2(ctx, ly, uint(li), di, LayVals[ly.Idxs.ValsIdx(di)], Pools[ly.Idxs.PoolIdx(0, di)]);
}

[numthreads(1, 1, 1)]
void main(uint3 idx : SV_DispatchThreadID) { // todo: iterate over global Data parallel
	if (idx.x >= Ctx[0].NetIdxs.NData) {
		return;
	}

	uint di = idx.x;
	
	// note: this bizarre logic is only way to get multiple writes to Context
	// to actually stick -- needs to be done sequentially within one thread
	// and not even in a for loop for some reason.
	int pnmi = -1;
	int cmpi = -1;
	int cmni = -1;
	int ldti = -1;
	int vspi = -1;
	int vtai = -1;
	int rwdi = -1;
	int tdpi = -1;
	int tdii = -1;
	int tddi = -1;
	for (int li = 0; li < Ctx[0].NetIdxs.NLayers; li++) {
		switch (Layers[li].LayType) {
		case PTNotMaintLayer:
			pnmi = li;
			break;
		case CeMLayer:
			if (Layers[li].Learn.NeuroMod.Valence == Positive) {
				cmpi = li;
			} else {
				cmni = li;
			}
			break;
		case VSPatchLayer:
			vspi = li;
			break;
		case LDTLayer:
			ldti = li;
			break;
		case VTALayer:
			vtai = li;
			break;
		case RWDaLayer:
			rwdi = li;
			break;
		case TDPredLayer:
			tdpi = li;
			break;
		case TDIntegLayer:
			tdii = li;
			break;
		case TDDaLayer:
			tddi = li;
			break;
		}
	}
	if (pnmi >= 0) {
		CyclePost(Ctx[0], Layers[pnmi], pnmi, di);
	}
	if (cmpi >= 0) {
		CyclePost(Ctx[0], Layers[cmpi], cmpi, di);
	}                                      
	if (cmni >= 0) {                       
		CyclePost(Ctx[0], Layers[cmni], cmni, di);
	}                                      
	if (ldti >= 0) { // depends on pn note:mi
		CyclePost(Ctx[0], Layers[ldti], ldti, di);
	}                                      
	if (vspi >= 0) {                       
		CyclePost(Ctx[0], Layers[vspi], vspi, di);
	}                                      
	if (rwdi >= 0) {                       
		CyclePost(Ctx[0], Layers[rwdi], rwdi, di);
	}                                      
	if (tdpi >= 0) {                       
		CyclePost(Ctx[0], Layers[tdpi], tdpi, di);
	}                                      
	if (tdii >= 0) {                       
		CyclePost(Ctx[0], Layers[tdii], tdii, di);
	}                                      
	if (tddi >= 0) {                       
		CyclePost(Ctx[0], Layers[tddi], tddi, di);
	}                                      
	// note: this depends vspi, cm*i, ldds on ti
	if (vtai >= 0) {                       
		CyclePost(Ctx[0], Layers[vtai], vtai, di);
	}
}

