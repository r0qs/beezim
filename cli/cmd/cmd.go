package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"swiki/indexer"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const (
	kiwixZimURL string = "https://download.kiwix.org/zim"
)

var (
	baseDir  string
	uploader *indexer.BeeUploader
)

var (
	optionKiwix          string
	optionGasPrice       string
	optionBeeApiUrl      string
	optionBeeDebugApiUrl string
	optionBatchID        string
	optionBatchDepth     uint64
	optionBatchAmount    int64
	optionGatewayMode    bool
	optionDataDir        string
	optionClean          bool
	optionZimFile        string
	optionZimURL         string
	optionTarFile        string
)

const (
	optionNameKiwix          = "kiwix"
	optionNameGasPrice       = "gas-price"
	optionNameBeeApiUrl      = "bee-api-url"
	optionNameBeeDebugApiUrl = "bee-debug-api-url"
	optionNameBatchID        = "batch-id"
	optionNameBatchDepth     = "batch-depth"
	optionNameBatchAmount    = "batch-amount"
	optionNameGatewayMode    = "gateway"
	optionNameDataDir        = "datadir"
	optionNameClean          = "clean"
	optionNameZimFile        = "zim"
	optionNameZimURL         = "url"
	optionNameTarFile        = "tar"
)

func init() {
	_, pwd, _, _ := runtime.Caller(0)
	baseDir = filepath.Join(path.Dir(pwd), "../..")

	// TODO: load from config (use viper)
	if err := godotenv.Load(filepath.Join(baseDir, ".env")); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	rootCmd.PersistentFlags().StringVar(&optionKiwix, optionNameKiwix, "wikipedia", "name of the compressed website hosted by Kiwix. Run \"list\" to see all available options")
	rootCmd.PersistentFlags().StringVar(&optionGasPrice, optionNameGasPrice, "", "gas price for postage stamps purchase")
	rootCmd.PersistentFlags().StringVar(&optionBeeApiUrl, optionNameBeeApiUrl, os.Getenv("BEE_API_URL"), "bee api url")
	rootCmd.PersistentFlags().StringVar(&optionBeeDebugApiUrl, optionNameBeeDebugApiUrl, os.Getenv("BEE_DEBUG_API_URL"), "bee debug api url")
	rootCmd.PersistentFlags().StringVar(&optionBatchID, optionNameBatchID, "", "bee postage batch ID")
	rootCmd.PersistentFlags().Uint64Var(&optionBatchDepth, optionNameBatchDepth, 30, "bee postage batch depth")
	rootCmd.PersistentFlags().Int64Var(&optionBatchAmount, optionNameBatchAmount, 100000000, "bee postage batch amount")
	rootCmd.PersistentFlags().BoolVar(&optionGatewayMode, optionNameGatewayMode, false, fmt.Sprintf("connect to the swarm public gateway (default \"%s\")", os.Getenv("BEE_GATEWAY")))
	rootCmd.PersistentFlags().StringVar(&optionDataDir, optionNameDataDir, "", "path to datadir directory (default \"./datadir\")")
	rootCmd.PersistentFlags().BoolVar(&optionClean, optionNameClean, false, "delete all downloaded zim and generated tar files")
}

var rootCmd = &cobra.Command{
	Use:           "szim",
	Short:         "Swarm zim mirror command-line tool",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(_ *cobra.Command, _ []string) (err error) {
		if optionGatewayMode {
			optionBeeApiUrl = os.Getenv("BEE_GATEWAY")
			optionBeeDebugApiUrl = ""
		}

		uploader, err = indexer.NewUploader(optionBeeApiUrl, optionBeeDebugApiUrl, optionBatchDepth, optionBatchAmount)
		if err != nil {
			return err
		}

		return setDataDir()
	},
}

func Execute() (err error) {
	rootCmd.AddCommand(
		listWebCmd,
		newDownloadCmd(),
		newUploadCmd(),
		newParserCmd(),
		newMirrorCmd(),
	)

	return rootCmd.Execute()
}

// setDataDir sets the data directory to the root directory
// of the project by default.
func setDataDir() error {

	if optionDataDir == "" {
		optionDataDir = filepath.Join(baseDir, "datadir")
	}

	// ensure datadir dir exists
	if _, err := os.Stat(optionDataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(optionDataDir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func makeURL(filePath string) string {
	return path.Join(optionBeeApiUrl, "bzz", filePath)
}
