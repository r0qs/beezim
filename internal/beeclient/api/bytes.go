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

package api

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ethersphere/bee/pkg/swarm"
)

// Download downloads data from the node
func (a *Api) DownloadBytes(ctx context.Context, addr swarm.Address) (resp io.ReadCloser, err error) {
	return a.C.RequestData(ctx, http.MethodGet, fmt.Sprintf("/bytes/%s", addr.String()), nil)
}

// BytesUploadResponse represents Upload's response
type BytesUploadResponse struct {
	Reference swarm.Address `json:"reference"`
}

// Upload uploads bytes to the node
func (a *Api) UploadBytes(ctx context.Context, data io.Reader, o UploadOptions) (BytesUploadResponse, error) {
	var resp BytesUploadResponse

	header := make(http.Header)
	header.Set("Content-Type", "application/octet-stream")
	if o.Pin {
		header.Add(SwarmPinHeader, "true")
	}
	err := a.C.RequestWithHeader(ctx, http.MethodPost, "/bytes", header, data, &resp)
	return resp, err
}
