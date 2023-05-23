// Copyright (c) 2022, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// performs the SynCa synaptic Ca integration function on all sending projections

// note: all must be visible always because accessor methods refer to them
[[vk::binding(0, 1)]] StructuredBuffer<uint> NeuronIxs; // [Neurons][Idxs]
[[vk::binding(1, 1)]] StructuredBuffer<uint> SynapseIxs;  // [Layer][SendPrjns][SendNeurons][Syns]
[[vk::binding(1, 2)]] RWStructuredBuffer<float> Neurons; // [Neurons][Vars][Data]
[[vk::binding(2, 2)]] RWStructuredBuffer<float> NeuronAvgs; // [Neurons][Vars]
[[vk::binding(0, 3)]] RWStructuredBuffer<float> Synapses;  // [Layer][SendPrjns][SendNeurons][Syns]
[[vk::binding(1, 3)]] RWStructuredBuffer<float> SynapseCas;  // [Layer][SendPrjns][SendNeurons][Syns][Data]

#include "context.hlsl"
#include "layerparams.hlsl"
#include "prjnparams.hlsl"

// note: binding is var, set

// Set 0: uniform layer params -- could not have prjns also be uniform..
[[vk::binding(0, 0)]] uniform LayerParams Layers[]; // [Layer]

// Set 1: effectively uniform indexes and prjn params as structured buffers in storage
[[vk::binding(2, 1)]] StructuredBuffer<PrjnParams> Prjns; // [Layer][SendPrjns]
[[vk::binding(3, 1)]] StructuredBuffer<StartN> SendCon; // [Layer][SendPrjns][SendNeurons]
[[vk::binding(4, 1)]] StructuredBuffer<uint> RecvPrjnIdxs; // [Layer][RecvPrjns][RecvNeurons]
[[vk::binding(5, 1)]] StructuredBuffer<StartN> RecvCon; // [Layer][RecvPrjns][RecvNeurons]
[[vk::binding(6, 1)]] StructuredBuffer<uint> RecvSynIdxs; // [Layer][RecvPrjns][RecvNeurons][Syns]

// Set 2: main network structs and vals -- all are writable
[[vk::binding(0, 2)]] StructuredBuffer<Context> Ctx; // [0]


void SynCaSendSyn(in Context ctx, in PrjnParams pj, uint syni, uint di, float snCaSyn, float updtThr) {
	uint ri = SynI(ctx, syni, SynRecvIdx);
	pj.SynCaSyn(ctx, syni, ri, di, snCaSyn, updtThr);
}

void SynCaRecvSyn(in Context ctx, in PrjnParams pj, uint syni, uint di, float rnCaSyn, float updtThr) {
	uint si = SynI(ctx, syni, SynSendIdx);
	if (NrnV(ctx, si, di, Spike) > 0) { // already processed in Send
		return;
	}
	pj.SynCaSyn(ctx, syni, si, di, rnCaSyn, updtThr);
}

void SynCaSendPrjn(in Context ctx, in PrjnParams pj, in LayerParams ly, uint ni, uint lni, uint di, float updtThr) {
	if (pj.Learn.Learn == 0) {
		return;
	}
	
	float snCaSyn = pj.Learn.KinaseCa.SpikeG * NrnV(ctx, ni, di, CaSyn);
	uint cni = pj.Idxs.SendConSt + lni;
	uint synst = pj.Idxs.SynapseSt + SendCon[cni].Start;
	uint synn = SendCon[cni].N;
	
	for (uint ci = 0; ci < synn; ci++) {
		SynCaSendSyn(ctx, pj, synst + ci, di, snCaSyn, updtThr);
	}
}

void SynCaRecvPrjn(in Context ctx, in PrjnParams pj, in LayerParams ly, uint ni, uint lni, uint di, float updtThr) {
	if (pj.Learn.Learn == 0) {
		return;
	}
	
	float rnCaSyn = pj.Learn.KinaseCa.SpikeG * NrnV(ctx, ni, di, CaSyn);
	uint cni = pj.Idxs.RecvConSt + lni;
	uint synst = pj.Idxs.RecvSynSt + RecvCon[cni].Start;
	uint synn = RecvCon[cni].N;
	
	for (uint ci = 0; ci < synn; ci++) {
		SynCaRecvSyn(ctx, pj, RecvSynIdxs[synst + ci], di, rnCaSyn, updtThr);
	}
}

void SynCa2(in Context ctx, in LayerParams ly, uint ni, uint di) {
	float updtThr = ly.Learn.CaLearn.UpdtThr;

	if (NrnV(ctx, ni, di, CaSpkP) < updtThr && NrnV(ctx, ni, di, CaSpkD) < updtThr) {
		return;
	}
	uint lni = ni - ly.Idxs.NeurSt; // layer-based as in Go
	
	for (uint pi = 0; pi < ly.Idxs.SendN; pi++) {
		SynCaSendPrjn(ctx, Prjns[ly.Idxs.SendSt + pi], ly, ni, lni, di, updtThr);
	}
	for (uint pi = 0; pi < ly.Idxs.RecvN; pi++) {
		SynCaRecvPrjn(ctx, Prjns[RecvPrjnIdxs[ly.Idxs.RecvSt + pi]], ly, ni, lni, di, updtThr);
	}
}

void SynCa(in Context ctx, uint ni, uint di) {
	if (NrnV(ctx, ni, di, Spike) == 0) {
		return;
	}
	uint li = NrnI(ctx, ni, NrnLayIdx);
	SynCa2(ctx, Layers[li], ni, di);
}

[numthreads(64, 1, 1)]
void main(uint3 idx : SV_DispatchThreadID) { // over Neurons * Data
	uint ni = Ctx[0].NetIdxs.ItemIdx(idx.x);
	if (!Ctx[0].NetIdxs.NeurIdxIsValid(ni)) {
		return;
	}
	uint di = Ctx[0].NetIdxs.DataIdx(idx.x);
	SynCa(Ctx[0], ni, di);
}


