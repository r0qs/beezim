package cmd

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newMirrorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mirror",
		Short: "mirror kiwix zim repositories to swarm",
		RunE: func(cmd *cobra.Command, args []string) error {
			zimPath, err := download(optionDataDir, optionZimFile, optionZimURL)
			if err != nil {
				return err
			}

			zimFile := filepath.Base(zimPath)
			err = parse(optionDataDir, zimFile)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			ext := filepath.Ext(zimFile)
			tarFile := fmt.Sprintf("%s.tar", zimFile[:len(zimFile)-len(ext)])
			addr, err := upload(ctx, optionDataDir, tarFile, optionBatchID)
			if err != nil {
				return err
			}
			log.Printf("collection %v uploaded with reference: %v", tarFile, addr)
			fmt.Printf("\nTry the link: %s\n", makeURL(addr.String()))
			return nil
		},
	}
	cmd.Flags().StringVar(&optionZimFile, optionNameZimFile, "", "path to the zim file")
	cmd.Flags().StringVar(&optionZimURL, optionNameZimURL, "", "download URL for the zim files")

	return cmd
}
