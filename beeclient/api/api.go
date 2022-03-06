// Copyright 2020 Ethersphere. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This code is based on the beekeeper beeclient api

package api

import (
	"net/url"
	"swiki/httpclient"
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
	*httpclient.Client
	Bytes *BytesService
	Chunk *ChunkService
	Dirs  *DirsService
}

func NewAPI(beeURL *url.URL, o *httpclient.ClientOptions) (*Api, error) {
	httpc, err := httpclient.NewClient(beeURL, o)
	if err != nil {
		return nil, err
	}
	a := &Api{httpc, nil, nil, nil}
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
