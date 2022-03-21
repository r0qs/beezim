package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

func newDownloadCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download zim file",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := download(optionDataDir, optionZimFile, optionZimURL)
			return err
		},
	}
	cmd.Flags().StringVar(&optionZimFile, optionNameZimFile, "", "path to the zim file")
	cmd.Flags().StringVar(&optionZimURL, optionNameZimURL, "", "download URL for the zim files")
	// TODO: add download all option

	return cmd
}

func download(dataDir, zimFile, zimURL string) (string, error) {
	if zimFile != "" && zimURL == "" {
		if filepath.Ext(zimFile) != ".zim" {
			return "", fmt.Errorf("file must has .zim extention")
		}
		zimURL = fmt.Sprintf("%s/%s/%s", kiwixZimURL, optionKiwix, zimFile)
	} else if zimFile == "" && zimURL != "" {
		zimFile = path.Base(zimURL)
	} else if zimFile != "" && zimURL != "" {
		return "", fmt.Errorf("--zim and --url are mutually exclusive. Please use --zim with --kiwix or just --url")
	} else {
		return "", fmt.Errorf("--zim or --url should be provided")
	}

	zimDownloadPath := fmt.Sprintf("%s/%s", dataDir, zimFile)
	if _, err := os.Stat(zimDownloadPath); os.IsNotExist(err) {
		if err := downloadZim(zimURL, zimDownloadPath); err != nil {
			return "", err
		}
	}
	return zimDownloadPath, nil
}

// TODO: continue the download from where it stopped in case of crash
// TODO: keep track of already uploaded files (in the metadata kv)
func downloadZim(targetURL string, dstFile string) error {
	resp, err := http.Get(targetURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %v [status: %v]\n", targetURL, resp.Status)
	}

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	dest, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dest.Close()

	header := fmt.Sprintf("Downloading zim file: %s", filepath.Base(dstFile))
	progressBar := newNetProgressBar(header, size, true)
	progressBar.Start()

	io.Copy(dest, progressBar.NewProxyReader(resp.Body))

	progressBar.Finish()

	// TODO: use a proper logger and make log messages optional by level (info, debug, etc)
	log.Printf("Zim file saved to: %s \n", dstFile)
	return nil
}
