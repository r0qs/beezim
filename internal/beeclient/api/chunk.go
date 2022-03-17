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
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ethersphere/bee/pkg/swarm"
)

type ChunkService struct {
	api *Api
}

func newChunkService(a *Api) *ChunkService {
	return &ChunkService{api: a}
}

func (cs *ChunkService) Download(ctx context.Context, addr swarm.Address, targets ...string) (resp io.ReadCloser, err error) {
	url := fmt.Sprintf("/chunks/%s", addr.String())
	if len(targets) != 0 {
		url = fmt.Sprintf("%s?targets=%s", url, strings.Join(targets, ","))
	}
	return cs.api.C.RequestData(ctx, http.MethodGet, url, nil)
}

type ChunksUploadResponse struct {
	Reference swarm.Address `json:"reference"`
}

func (cs *ChunkService) Upload(ctx context.Context, data []byte, o UploadOptions) (ChunksUploadResponse, error) {
	var resp ChunksUploadResponse

	header := make(http.Header)
	header.Set("Content-Type", "application/octet-stream")
	if o.Pin {
		header.Add(SwarmPinHeader, "true")
	}
	header.Add(SwarmPostageBatchIdHeader, o.BatchID)

	err := cs.api.C.RequestWithHeader(ctx, http.MethodPost, "/chunks", header, bytes.NewReader(data), &resp)
	return resp, err
}
