// Copyright 2020 Ethersphere. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This code is based on the beekeeper beeclient api

package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ethersphere/bee/pkg/swarm"
)

type DirsService struct {
	api *Api
}

func newDirsService(a *Api) *DirsService {
	return &DirsService{api: a}
}

// Download downloads data from the node
func (ds *DirsService) Download(ctx context.Context, addr swarm.Address, path string) (resp io.ReadCloser, err error) {
	return ds.api.C.RequestData(ctx, http.MethodGet, fmt.Sprintf("/bzz/%s/%s", addr.String(), path), nil)
}

// DirsUploadResponse represents Upload's response
type DirsUploadResponse struct {
	Reference swarm.Address `json:"reference"`
}

// Upload uploads TAR collection to the node
func (ds *DirsService) Upload(ctx context.Context, data io.Reader, size int64, o UploadCollectionOptions) (DirsUploadResponse, error) {
	var resp DirsUploadResponse

	var indexDocumentHeader, errorDocumentHeader string
	if o.IndexDocumentHeader != "" {
		indexDocumentHeader = o.IndexDocumentHeader
	} else {
		indexDocumentHeader = "index.html"
	}

	if o.ErrorDocumentHeader != "" {
		errorDocumentHeader = o.ErrorDocumentHeader
	} else {
		errorDocumentHeader = "error.html"
	}

	header := make(http.Header)
	header.Set("Content-Type", "application/x-tar")
	header.Set("Content-Length", strconv.FormatInt(size, 10))
	header.Set(SwarmCollectionHeader, "true")
	header.Set(SwarmDeferredUploadHeader, "true")
	header.Set(SwarmIndexDocumentHeader, indexDocumentHeader)
	header.Set(SwarmErrorDocumentHeader, errorDocumentHeader)
	header.Set(SwarmPostageBatchIdHeader, o.BatchID)

	if o.Pin {
		header.Set(SwarmPinHeader, "true")
	}

	if o.Tag != 0 {
		header.Set(SwarmTagHeader, strconv.FormatUint(uint64(o.Tag), 10))
	}

	err := ds.api.C.RequestWithHeader(ctx, http.MethodPost, "/bzz", header, data, &resp)
	return resp, err
}
