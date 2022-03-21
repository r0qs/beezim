// Copyright 2021 Ethersphere.
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
	"fmt"
	"net/http"

	"github.com/r0qs/beezim/internal/beeclient/api"

	"github.com/ethersphere/bee/pkg/bigint"
)

type postageResponse struct {
	BatchID string `json:"batchID"`
}

type PostageOptions struct {
	GasPrice string
}

// CreatePostageBatch sends a create postage request to a node that returns the bactchID
func (d *DebugAPI) CreatePostageBatch(ctx context.Context, amount int64, depth uint64, label string, o PostageOptions) (string, error) {
	h := http.Header{}

	if o.GasPrice != "" {
		h.Add(api.GasPriceHeader, o.GasPrice)
	}

	url := fmt.Sprintf("/stamps/%d/%d?label=%s", amount, depth, label)
	var resp postageResponse
	err := d.C.RequestWithHeader(ctx, http.MethodPost, url, h, nil, &resp)
	if err != nil {
		return "", err
	}
	return resp.BatchID, err
}

type PostageStampResponse struct {
	BatchID       string         `json:"batchID"`
	Utilization   uint32         `json:"utilization"`
	Usable        bool           `json:"usable"`
	Label         string         `json:"label"`
	Depth         uint8          `json:"depth"`
	Amount        *bigint.BigInt `json:"amount"`
	BucketDepth   uint8          `json:"bucketDepth"`
	BlockNumber   uint64         `json:"blockNumber"`
	ImmutableFlag bool           `json:"immutableFlag"`
	Exists        bool           `json:"exists"`
	BatchTTL      int64          `json:"batchTTL"`
}

type postageStampsResponse struct {
	Stamps []PostageStampResponse `json:"stamps"`
}

// Fetches the list postage stamp batches
func (d *DebugAPI) PostageBatches(ctx context.Context) ([]PostageStampResponse, error) {
	var resp postageStampsResponse
	err := d.C.Request(ctx, http.MethodGet, "/stamps", nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Stamps, nil
}
