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

	opts.DebugAPIURL, err = url.Parse(beeDebugApiUrl)
	if err != nil {
		return nil, fmt.Errorf("error parsing debug api url: %v", err)
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

// UploadTar uploads a zim tarball to swarm
// TODO: upload by filename
func (b BeeUploader) UploadTar(outputDir string, opts api.UploadCollectionOptions) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".tar" {
			buf, err := tarball.ReadTarBuffer(path)
			if err != nil {
				return err
			}
			tarFile := tarball.NewBufferFile(info.Name(), buf)
			if err := b.client.UploadCollection(ctx, tarFile, opts); err != nil {
				return err
			}
			log.Printf("collection %v uploaded with reference: %v", info.Name(), tarFile.Address())
		}
		return nil
	})
}

//TODO: Buy stamps
//TODO: Make manifest metadata
