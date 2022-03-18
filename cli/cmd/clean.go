package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean files in datadir",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cleanDatadir()
		},
	}

	return cmd
}

// TODO: add option for no confirmation?
func cleanDatadir() error {
	if optionDataDir == "" || optionDataDir == "/" {
		return nil
	}

	baseDir := optionDataDir
	files, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		action := fmt.Sprintf("Are you sure you want delete %s?", file.Name())
		confirmationReader := NewConfirmationInputReader(action, func() error {
			filePath := filepath.Join(baseDir, file.Name())
			if file.IsDir() {
				fmt.Println("deleting directory...", filePath)
				return os.RemoveAll(filePath)
			}
			fmt.Println("deleting file...", filePath)
			return os.Remove(filePath)
		})

		_, err := confirmationReader.ReadInput()
		if err != nil && err != AbortCmd {
			return err
		}
		if err == AbortCmd {
			return nil
		}
	}

	return nil
}
