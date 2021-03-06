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

package beeclient

import (
	"math"

	"github.com/ethersphere/bee/pkg/swarm"
)

const MinimumBatchDepth = 11

func EstimatePostageBatchDepth(contentLength int64, isEncrypted bool) (uint64, int64) {
	totalChunks := CalculateNumberOfChunks(contentLength, isEncrypted)
	depth := uint64(math.Log2(float64(totalChunks)))
	if depth < MinimumBatchDepth {
		depth = MinimumBatchDepth
	}
	return depth, totalChunks
}

// CalculateNumberOfChunks calculates the number of chunks in an arbitrary
// content length.
func CalculateNumberOfChunks(contentLength int64, isEncrypted bool) int64 {
	if contentLength <= swarm.ChunkSize {
		return 1
	}
	branchingFactor := swarm.Branches
	if isEncrypted {
		branchingFactor = swarm.EncryptedBranches
	}

	dataChunks := math.Ceil(float64(contentLength) / float64(swarm.ChunkSize))
	totalChunks := dataChunks
	intermediate := dataChunks / float64(branchingFactor)

	for intermediate > 1 {
		totalChunks += math.Ceil(intermediate)
		intermediate = intermediate / float64(branchingFactor)
	}

	return int64(totalChunks) + 1
}
