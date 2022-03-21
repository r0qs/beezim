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
package beeclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/r0qs/beezim/internal/beeclient/api"
	"github.com/r0qs/beezim/internal/beeclient/debugapi"
	"github.com/r0qs/beezim/internal/httpclient"
	"github.com/r0qs/beezim/internal/tarball"

	"github.com/ethersphere/bee/pkg/swarm"
)

type ClientOptions struct {
	APIURL              *url.URL
	APIInsecureTLS      bool
	DebugAPIURL         *url.URL
	DebugAPIInsecureTLS bool
}

type BeeClient struct {
	api   *api.Api
	debug *debugapi.DebugAPI
}

func NewBee(opts ClientOptions) (c *BeeClient, err error) {
	c = &BeeClient{}

	if opts.APIURL != nil {
		c.api, err = api.NewAPI(opts.APIURL, &httpclient.ClientOptions{
			HTTPClient: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: opts.APIInsecureTLS,
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}
	}
	if opts.DebugAPIURL != nil {
		c.debug, err = debugapi.NewDebugAPI(opts.DebugAPIURL, &httpclient.ClientOptions{
			HTTPClient: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: opts.DebugAPIInsecureTLS,
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *BeeClient) DownloadChunk(ctx context.Context, addr swarm.Address, targets ...string) (io.ReadCloser, error) {
	return c.api.Chunk.Download(ctx, addr, targets...)
}

func (c *BeeClient) UploadChunk(ctx context.Context, data []byte, o api.UploadOptions) (swarm.Address, error) {
	resp, err := c.api.Chunk.Upload(ctx, data, o)
	return resp.Reference, err
}

func (c *BeeClient) DownloadBytes(ctx context.Context, addr swarm.Address) (io.ReadCloser, error) {
	return c.api.Bytes.Download(ctx, addr)
}

func (c *BeeClient) UploadBytes(ctx context.Context, data io.Reader, o api.UploadOptions) (swarm.Address, error) {
	resp, err := c.api.Bytes.Upload(ctx, data, o)
	return resp.Reference, err
}

// UploadCollection uploads collection from a given reader
func (c *BeeClient) UploadCollection(ctx context.Context, r io.Reader, size int64, o api.UploadCollectionOptions) (swarm.Address, error) {
	resp, err := c.api.Dirs.UploadCollection(ctx, r, size, o)
	if err != nil {
		return swarm.Address{}, fmt.Errorf("upload collection: %v", err)
	}
	return resp.Reference, nil
}

// DownloadManifestFile downloads manifest file from the node and returns it's size and hash
func (c *BeeClient) DownloadManifestFile(ctx context.Context, addr swarm.Address, path string) (size int64, hash []byte, err error) {
	r, err := c.api.Dirs.Download(ctx, addr, path)
	if err != nil {
		return 0, nil, fmt.Errorf("download manifest file %s: %w", path, err)
	}

	h := tarball.FileHasher()
	size, err = io.Copy(h, r)
	if err != nil {
		return 0, nil, fmt.Errorf("download manifest file %s: %w", path, err)
	}

	return size, h.Sum(nil), nil
}

func (c *BeeClient) DownloadManifestBytes(ctx context.Context, addr swarm.Address, path string) ([]byte, error) {
	r, err := c.api.Dirs.Download(ctx, addr, path)
	if err != nil {
		return nil, fmt.Errorf("download manifest %s: %w", path, err)
	}
	defer r.Close()

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("download manifest %s: %w", path, err)
	}

	return buf, nil
}

func (c *BeeClient) Addresses(ctx context.Context) (debugapi.Addresses, error) {
	return c.debug.Node.Addresses(ctx)
}

// CreatePostageBatch returns the batchID of a batch of postage stamps
func (c *BeeClient) CreatePostageBatch(ctx context.Context, amount int64, depth uint64, label string, o debugapi.PostageOptions) (string, error) {
	if depth < MinimumBatchDepth {
		depth = MinimumBatchDepth
	}
	return c.debug.Postage.CreatePostageBatch(ctx, amount, depth, label, o)
}

// PostageBatches returns the list of batches of node
func (c *BeeClient) PostageBatches(ctx context.Context) ([]debugapi.PostageStampResponse, error) {
	return c.debug.Postage.PostageBatches(ctx)
}
