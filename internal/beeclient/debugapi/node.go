// Copyright 2020 Ethersphere. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This code is based on the beekeeper beeclient api

package debugapi

import (
	"context"
	"net/http"

	"github.com/ethersphere/bee/pkg/swarm"
)

type NodeService struct {
	debugAPI *DebugAPI
}

func newNodeService(d *DebugAPI) *NodeService {
	return &NodeService{debugAPI: d}
}

type Addresses struct {
	Ethereum     string        `json:"ethereum"`
	Overlay      swarm.Address `json:"overlay"`
	PublicKey    string        `json:"public_key"`
	Underlay     []string      `json:"underlay"`
	PSSPublicKey string        `json:"pss_public_key"`
}

func (ns *NodeService) Addresses(ctx context.Context) (Addresses, error) {
	var resp Addresses
	err := ns.debugAPI.C.RequestJSON(ctx, http.MethodGet, "/addresses", nil, &resp)
	return resp, err
}

type Peer struct {
	Address swarm.Address `json:"address"`
}

type Peers struct {
	Peers []Peer `json:"peers"`
}

func (ns *NodeService) Peers(ctx context.Context) (resp Peers, err error) {
	err = ns.debugAPI.C.RequestJSON(ctx, http.MethodGet, "/peers", nil, &resp)
	return
}
