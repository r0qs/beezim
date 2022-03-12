package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/r0qs/beezim/internal/beeclient/api"

	"github.com/ethersphere/bee/pkg/swarm"
	"github.com/spf13/cobra"
)

func newUploadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "upload zim file to swarm",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkTarFileName(optionTarFile); err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			addr, err := upload(ctx, optionDataDir, optionTarFile, optionBatchID)
			if err != nil {
				return err
			}
			log.Printf("collection %v uploaded with reference: %v", optionTarFile, addr)
			return nil
		},
	}
	cmd.Flags().StringVar(&optionTarFile, optionNameTarFile, "", "tar file name")
	// TODO: add upload all option
	cmd.AddCommand(
		newUploadAllCmd(),
	)

	return cmd
}

func checkTarFileName(tarFile string) error {
	if tarFile == "" {
		return fmt.Errorf("please provide a tar file")
	}
	if filepath.Ext(tarFile) != ".tar" {
		return fmt.Errorf("file must has .tar extention")
	}
	return nil
}

func upload(ctx context.Context, dataDir string, tarFile string, batchID string) (swarm.Address, error) {
	tarPath := filepath.Join(dataDir, tarFile)
	if _, err := os.Stat(tarPath); os.IsNotExist(err) {
		return swarm.Address{}, fmt.Errorf("tar file %s not found", tarFile)
	}
	// TODO: get tag and pin option.
	// TODO: buy batchID as needed:
	// estimate the cost of upload the zim and the size of the stamp to do so.
	// allow users to agree/deny with the cost before proceed with each upload
	// the stamps will be automatically bought
	// TODO: keep address for local metadata
	// TODO: command to buy stamps and check if stamp they are usable
	// --wait-usable-stamp (keep waiting until bought stamp is ready)
	return uploader.UploadTarFile(ctx, tarPath, tarFile, api.UploadCollectionOptions{
		Pin:                 true,
		BatchID:             batchID,
		IndexDocumentHeader: "index.html",
		ErrorDocumentHeader: "error.html",
	})
}

// Upload Subcommands
func newUploadAllCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "upload all zim files to swarm of a specifc kiwix mirror",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			addrs, err := uploadAllFrom(ctx, optionDataDir, optionKiwix, optionBatchID)
			if err != nil {
				return err
			}
			for name, addr := range addrs {
				log.Printf("collection %v uploaded with reference: %v", name, addr)
			}
			return nil
		},
	}
}

func uploadAllFrom(ctx context.Context, dataDir string, kiwixMirror string, batchID string) (map[string]swarm.Address, error) {
	filter := func(filename string) bool {
		return strings.Contains(filename, kiwixMirror)
	}
	return uploader.UploadMatchTar(ctx, dataDir, filter, api.UploadCollectionOptions{
		Pin:                 true,
		BatchID:             batchID,
		IndexDocumentHeader: "index.html",
		ErrorDocumentHeader: "error.html",
	})
}
