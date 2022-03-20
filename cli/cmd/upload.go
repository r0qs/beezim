package cmd

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/r0qs/beezim/internal/beeclient/api"
	"github.com/r0qs/beezim/internal/tarball"

	"github.com/ethersphere/bee/pkg/swarm"
	"github.com/spf13/cobra"
)

func newUploadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload tar file to swarm",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkTarFileName(optionTarFile); err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			addr, err := upload(ctx, optionDataDir, optionTarFile, optionBeeBatchID)
			if err != nil {
				return err
			}
			log.Printf("collection %v uploaded with reference: %v", optionTarFile, addr)
			fmt.Printf("\nTry the link: %s\n", makeURL(addr.String()))
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
	addr, err := uploadTarFile(ctx, tarPath, tarFile, api.UploadCollectionOptions{
		MimeType:            api.ContentTypeTar,
		Tag:                 optionBeeTag,
		Pin:                 optionBeePin,
		BatchID:             batchID,
		IndexDocumentHeader: "index.html",
		ErrorDocumentHeader: "error.html",
	})
	if err != nil {
		return swarm.Address{}, err
	}

	if optionClean {
		cleanDatadir()
	}
	return addr, nil
}

// Upload Subcommands
func newUploadAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Upload all tar files of a specifc kiwix mirror to swarm",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			addrs, err := uploadAllFrom(ctx, optionDataDir, optionKiwix, optionBeeBatchID)
			if err != nil {
				return err
			}
			for name, addr := range addrs {
				log.Printf("collection %v uploaded with reference: %v", name, addr)
			}
			return nil
		},
	}
	cmd.MarkFlagRequired(optionKiwix)

	return cmd
}

func uploadAllFrom(ctx context.Context, dataDir string, kiwixMirror string, batchID string) (map[string]swarm.Address, error) {
	filter := func(filename string) bool {
		return strings.Contains(filename, kiwixMirror)
	}

	addrs, err := uploadMatchTar(ctx, dataDir, filter, api.UploadCollectionOptions{
		MimeType:            api.ContentTypeTar,
		Tag:                 optionBeeTag,
		Pin:                 optionBeePin,
		BatchID:             batchID,
		IndexDocumentHeader: "index.html",
		ErrorDocumentHeader: "error.html",
	})
	if err != nil {
		return nil, err
	}

	if optionClean {
		cleanDatadir()
	}
	return addrs, nil
}

func uploadMatchTar(ctx context.Context, targetDir string, filter func(x string) bool, opts api.UploadCollectionOptions) (map[string]swarm.Address, error) {
	files := make(map[string]swarm.Address)
	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".tar" && filter(info.Name()) {
			addr, err := uploadTarFile(ctx, path, info.Name(), opts)
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

// TODO: make generic for any type of collection
func uploadTarFile(ctx context.Context, path string, name string, opts api.UploadCollectionOptions) (addr swarm.Address, err error) {
	f, err := os.Open(path)
	if err != nil {
		return swarm.Address{}, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return swarm.Address{}, err
	}

	header := fmt.Sprintf("Uploading tar file: %s", name)
	progressBar := newNetProgressBar(header, int(info.Size()), false)
	progressBar.Start()
	defer progressBar.Finish()

	r, w := io.Pipe()
	done := make(chan error)
	defer close(done)

	go func() {
		if addr, err = bee.UploadCollection(ctx, progressBar.NewProxyReader(r), info.Size(), opts); err != nil {
			done <- err
			return
		}
		done <- nil
	}()

	tr := tar.NewReader(f)
	if err = tarball.CopyTar(w, tr); err != nil {
		return swarm.Address{}, err
	}

	if err = w.Close(); err != nil {
		return swarm.Address{}, err
	}
	err = <-done

	return addr, err
}
