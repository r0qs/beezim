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
	"net/url"

	"github.com/r0qs/beezim/internal/httpclient"
)

type DebugAPI struct {
	C *httpclient.Client
}

func NewDebugAPI(beeURL *url.URL, o *httpclient.ClientOptions) (*DebugAPI, error) {
	httpc, err := httpclient.NewClient(beeURL, o)
	if err != nil {
		return nil, err
	}
	return &DebugAPI{C: httpc}, nil
}
