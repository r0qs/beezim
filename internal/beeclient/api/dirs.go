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

// UploadTar uploads collection to the node
func (ds *DirsService) UploadCollection(ctx context.Context, data io.Reader, size int64, o UploadCollectionOptions) (DirsUploadResponse, error) {
	var resp DirsUploadResponse

	header := make(http.Header)
	// Default to tar collection
	if o.MimeType == "" {
		header.Set("Content-Type", ContentTypeTar)
	} else {
		header.Set("Content-Type", o.MimeType)
	}
	header.Set("Content-Length", strconv.FormatInt(size, 10))
	header.Set(SwarmCollectionHeader, "true")
	header.Set(SwarmDeferredUploadHeader, "true")
	header.Set(SwarmPostageBatchIdHeader, o.BatchID)

	if o.IndexDocumentHeader != "" {
		header.Set(SwarmIndexDocumentHeader, o.IndexDocumentHeader)
	}

	if o.ErrorDocumentHeader != "" {
		header.Set(SwarmErrorDocumentHeader, o.ErrorDocumentHeader)
	}

	if o.Pin {
		header.Set(SwarmPinHeader, "true")
	}

	if o.Tag != 0 {
		header.Set(SwarmTagHeader, strconv.FormatUint(uint64(o.Tag), 10))
	}

	err := ds.api.C.RequestWithHeader(ctx, http.MethodPost, "/bzz", header, data, &resp)
	return resp, err
}
