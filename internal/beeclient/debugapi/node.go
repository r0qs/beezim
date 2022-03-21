// Copyright 2020 Ethersphere.
// Copyright 2022 Beezim Authors.
// All rights reserved.
// Use of this source code is originally governed by
// BSD 3-Clause and our modifications by GPLv3
// license that can be found in the LICENSE file.
//
// This code is based on the beekeeper beeclient api.
// The http client was split to its own package.
// The bee api and debug api were modified and
// simplified to fit the purposes of Beezim.

package debugapi

import (
	"context"
	"net/http"

	"github.com/ethersphere/bee/pkg/swarm"
)

type Addresses struct {
	Ethereum     string        `json:"ethereum"`
	Overlay      swarm.Address `json:"overlay"`
	PublicKey    string        `json:"public_key"`
	Underlay     []string      `json:"underlay"`
	PSSPublicKey string        `json:"pss_public_key"`
}

func (d *DebugAPI) Addresses(ctx context.Context) (Addresses, error) {
	var resp Addresses
	err := d.C.RequestJSON(ctx, http.MethodGet, "/addresses", nil, &resp)
	return resp, err
}

type Peer struct {
	Address swarm.Address `json:"address"`
}

type Peers struct {
	Peers []Peer `json:"peers"`
}

func (d *DebugAPI) Peers(ctx context.Context) (resp Peers, err error) {
	err = d.C.RequestJSON(ctx, http.MethodGet, "/peers", nil, &resp)
	return
}
