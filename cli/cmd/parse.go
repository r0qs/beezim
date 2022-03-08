package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"swiki/indexer"

	"github.com/spf13/cobra"
)

func newParserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "parse zim file and convert it to a tar file ready for upload",
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

	return cmd
}

func parse(dataDir string, zimFile string) error {
	zimPath := filepath.Join(dataDir, zimFile)
	tarDirName := strings.TrimSuffix(filepath.Base(zimPath), ".zim")

	sidx, err := indexer.New(zimPath)
	if err != nil {
		return err
	}

	// Parse zim file
	zimArticles := sidx.ParseZIM()

	// Build tar
	tarFile := filepath.Join(dataDir, fmt.Sprintf("%s.tar", tarDirName))
	if err := sidx.TarZim(tarFile, zimArticles); err != nil {
		return err
	}

	// Append index page
	if err := sidx.MakeIndexPage(tarFile); err != nil {
		return fmt.Errorf("Failed to copy index.html page to tar file: %v", err)
	}

	// Append 404 page
	if err := sidx.MakeErrorPage(tarFile); err != nil {
		return fmt.Errorf("Failed to copy error.html page to tar file: %v", err)
	}

	return nil
}
