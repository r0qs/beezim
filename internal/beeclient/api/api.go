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
	"net/url"

	"github.com/r0qs/beezim/internal/httpclient"
)

const (
	GasPriceHeader            = "Gas-Price"
	SwarmPinHeader            = "Swarm-Pin"
	SwarmTagHeader            = "Swarm-Tag"
	SwarmEncryptHeader        = "Swarm-Encrypt"
	SwarmIndexDocumentHeader  = "Swarm-Index-Document"
	SwarmErrorDocumentHeader  = "Swarm-Error-Document"
	SwarmFeedIndexHeader      = "Swarm-Feed-Index"
	SwarmFeedIndexNextHeader  = "Swarm-Feed-Index-Next"
	SwarmCollectionHeader     = "Swarm-Collection"
	SwarmPostageBatchIdHeader = "Swarm-Postage-Batch-Id"
	SwarmDeferredUploadHeader = "Swarm-Deferred-Upload"
)

type Api struct {
	C     *httpclient.Client
	Bytes *BytesService
	Chunk *ChunkService
	Dirs  *DirsService
}

func NewAPI(beeURL *url.URL, o *httpclient.ClientOptions) (*Api, error) {
	httpc, err := httpclient.NewClient(beeURL, o)
	if err != nil {
		return nil, err
	}
	a := &Api{C: httpc}
	a.Bytes = newBytesService(a)
	a.Chunk = newChunkService(a)
	a.Dirs = newDirsService(a)
	return a, nil
}

type UploadOptions struct {
	Pin     bool
	Tag     uint32
	BatchID string
}

type UploadCollectionOptions struct {
	Pin                 bool
	Tag                 uint32
	BatchID             string
	IndexDocumentHeader string
	ErrorDocumentHeader string
}
