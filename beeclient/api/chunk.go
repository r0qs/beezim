// Copyright 2020 Ethersphere. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This code is based on the beekeeper beeclient api

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

func (c *ChunkService) Download(ctx context.Context, addr swarm.Address, targets ...string) (resp io.ReadCloser, err error) {
	url := fmt.Sprintf("/chunks/%s", addr.String())
	if len(targets) != 0 {
		url = fmt.Sprintf("%s?targets=%s", url, strings.Join(targets, ","))
	}
	return c.api.RequestData(ctx, http.MethodGet, url, nil)
}

type ChunksUploadResponse struct {
	Reference swarm.Address `json:"reference"`
}

func (c *ChunkService) Upload(ctx context.Context, data []byte, o UploadOptions) (ChunksUploadResponse, error) {
	var resp ChunksUploadResponse

	header := make(http.Header)
	header.Set("Content-Type", "application/octet-stream")
	if o.Pin {
		header.Add(SwarmPinHeader, "true")
	}
	header.Add(SwarmPostageBatchIdHeader, o.BatchID)

	err := c.api.RequestWithHeader(ctx, http.MethodPost, "/chunks", header, bytes.NewReader(data), &resp)
	return resp, err
}
