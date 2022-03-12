package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/r0qs/beezim/indexer"

	"github.com/spf13/cobra"
)

func newParserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "parse zim file",
		Long:  "\nThe default behavior is to parse the ZIM and convert it to a tar file ready for upload.\nIf you only want to extract its content, use this command with the option --extract-only.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if optionZimFile != "" {
				if filepath.Ext(optionZimFile) != ".zim" {
					return fmt.Errorf("file must has .zim extention")
				}
				return parse(optionDataDir, optionZimFile)
			}
			return fmt.Errorf("zim file not provided")
		},
	}
	cmd.Flags().StringVar(&optionZimFile, optionNameZimFile, "", "path for the zim file")
	cmd.Flags().BoolVar(&optionExtractOnly, optionNameExtractOnly, false, "parse and extract the zim file to the datadir")

	return cmd
}

func parse(dataDir string, zimFile string) error {
	zimPath := filepath.Join(dataDir, zimFile)
	dirName := strings.TrimSuffix(filepath.Base(zimPath), ".zim")

	sidx, err := indexer.New(zimPath)
	if err != nil {
		return err
	}

	// Parse zim file
	zimArticles := sidx.ParseZIM()
	// TODO: kill the parser goroutine is there is nothing to do: zim was already parsed.

	if optionExtractOnly {
		outputDir := filepath.Join(optionDataDir, dirName)
		return sidx.UnZim(outputDir, zimArticles)
	} else {
		// TODO: what should be the default policy? check if file already exists and
		// do not build the tar, or overwrite it everytime?
		tarFile := filepath.Join(dataDir, fmt.Sprintf("%s.tar", dirName))
		// Build tar
		if err := sidx.TarZim(tarFile, zimArticles); err != nil {
			return err
		}

		if optionEnableSearch {
			// Append index page with search tool
			if err := sidx.MakeIndexSearchPage(tarFile); err != nil {
				return fmt.Errorf("Failed to copy index.html page to tar file: %v", err)
			}
		} else {
			// Append redirected index page
			if err := sidx.MakeRedirectIndexPage(tarFile); err != nil {
				return fmt.Errorf("Failed to copy index.html page to tar file: %v", err)
			}
		}

		// Append 404 page
		if err := sidx.MakeErrorPage(tarFile); err != nil {
			return fmt.Errorf("Failed to copy error.html page to tar file: %v", err)
		}
	}

	return nil
}
