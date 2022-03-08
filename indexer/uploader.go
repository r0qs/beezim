package indexer

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"swiki/beeclient"
	"swiki/beeclient/api"
	"swiki/internal/tarball"

	"github.com/ethersphere/bee/pkg/swarm"
)

type BeeUploader struct {
	client      *beeclient.BeeClient
	stampDepth  uint64
	stampAmount *big.Int
}

func NewUploader(beeApiUrl string, beeDebugApiUrl string, depth uint64, amount int64) (*BeeUploader, error) {
	var err error
	opts := beeclient.ClientOptions{}

	opts.APIURL, err = url.Parse(beeApiUrl)
	if err != nil {
		return nil, fmt.Errorf("error parsing api url: %v", err)
	}

	if beeDebugApiUrl != "" {
		opts.DebugAPIURL, err = url.Parse(beeDebugApiUrl)
		if err != nil {
			return nil, fmt.Errorf("error parsing debug api url: %v", err)
		}
	}

	c, err := beeclient.NewBee(opts)
	if err != nil {
		return nil, err
	}

	return &BeeUploader{
		client:      c,
		stampDepth:  depth,
		stampAmount: big.NewInt(amount),
	}, nil
}

func (b BeeUploader) UploadMatchTar(ctx context.Context, outputDir string, filter func(x string) bool, opts api.UploadCollectionOptions) (map[string]swarm.Address, error) {
	files := make(map[string]swarm.Address)
	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".tar" && filter(info.Name()) {
			addr, err := b.UploadTarFile(ctx, path, info.Name(), opts)
			if err != nil {
				return err
			}
			files[info.Name()] = addr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		log.Println("no tar files found for the given filter")
	}
	return files, nil
}

func (b BeeUploader) UploadTarFile(ctx context.Context, path string, name string, opts api.UploadCollectionOptions) (swarm.Address, error) {
	buf, err := tarball.ReadTarBuffer(path)
	if err != nil {
		return swarm.Address{}, err
	}
	tarFile := tarball.NewBufferFile(name, buf)
	if err := b.client.UploadCollection(ctx, tarFile, opts); err != nil {
		return swarm.Address{}, err
	}
	return tarFile.Address(), nil
}

//TODO: Buy stamps
//TODO: Make manifest metadata
