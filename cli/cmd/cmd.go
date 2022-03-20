package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"

	pb "github.com/cheggaaa/pb/v3"
	"github.com/r0qs/beezim/internal/beeclient"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const (
	kiwixZimURL string = "https://download.kiwix.org/zim"
)

var (
	baseDir      string
	bee          *beeclient.BeeClient
	cpuprofile   *os.File
	memprofile   *os.File
	blockprofile *os.File
)

var (
	optionKiwix          string
	optionGasPrice       string
	optionBeeApiUrl      string
	optionBeeDebugApiUrl string
	optionBeeBatchID     string
	optionBeeBatchDepth  uint64
	optionBeeBatchAmount int64
	optionBeeTag         uint32
	optionBeePin         bool
	optionGatewayMode    bool
	optionDataDir        string
	optionClean          bool
	optionZimFile        string
	optionZimURL         string
	optionTarFile        string
	optionExtractOnly    bool
	optionEnableSearch   bool
	optionCPUProfile     string
	optionMEMProfile     string
	optionBlockProfile   string
)

const (
	optionNameKiwix          = "kiwix"
	optionNameGasPrice       = "gas-price"
	optionNameBeeApiUrl      = "bee-api-url"
	optionNameBeeDebugApiUrl = "bee-debug-api-url"
	optionNameBeeBatchID     = "batch-id"
	optionNameBeeBatchDepth  = "batch-depth"
	optionNameBeeBatchAmount = "batch-amount"
	optionNameBeeTag         = "tag"
	optionNameBeePin         = "pin"
	optionNameGatewayMode    = "gateway"
	optionNameDataDir        = "datadir"
	optionNameClean          = "clean"
	optionNameZimFile        = "zim"
	optionNameZimURL         = "url"
	optionNameTarFile        = "tar"
	optionNameExtractOnly    = "extract-only"
	optionNameEnableSearch   = "enable-search"
	optionNameCPUProfile     = "cpuprofile"
	optionNameMEMProfile     = "memprofile"
	optionNameBlockProfile   = "blockprofile"
)

func init() {
	_, pwd, _, _ := runtime.Caller(0)
	baseDir = filepath.Join(filepath.Dir(pwd), "../..")

	// TODO: load from config (use viper)
	// FIXME: this approach currently does not work with make install.
	// TODO: move config files to home over ~/.beezim
	if err := godotenv.Load(filepath.Join(baseDir, ".env")); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	rootCmd.PersistentFlags().StringVar(&optionKiwix, optionNameKiwix, "wikipedia", "name of the compressed website hosted by Kiwix. Run \"list\" to see all available options")
	rootCmd.PersistentFlags().StringVar(&optionGasPrice, optionNameGasPrice, "", "gas price for postage stamps purchase")
	rootCmd.PersistentFlags().StringVar(&optionBeeApiUrl, optionNameBeeApiUrl, os.Getenv("BEE_API_URL"), "bee api url")
	rootCmd.PersistentFlags().StringVar(&optionBeeDebugApiUrl, optionNameBeeDebugApiUrl, os.Getenv("BEE_DEBUG_API_URL"), "bee debug api url")
	rootCmd.PersistentFlags().StringVar(&optionBeeBatchID, optionNameBeeBatchID, "", "bee postage batch ID")
	rootCmd.PersistentFlags().Uint64Var(&optionBeeBatchDepth, optionNameBeeBatchDepth, 30, "bee postage batch depth")
	rootCmd.PersistentFlags().Int64Var(&optionBeeBatchAmount, optionNameBeeBatchAmount, 100000000, "bee postage batch amount")
	rootCmd.PersistentFlags().Uint32Var(&optionBeeTag, optionNameBeeTag, 0, "bee tag UID to the attached to the uploaded data")
	rootCmd.PersistentFlags().BoolVar(&optionBeePin, optionNameBeePin, false, "whether the uploaded data should be locally pinned on a node")
	rootCmd.PersistentFlags().BoolVar(&optionGatewayMode, optionNameGatewayMode, false, fmt.Sprintf("connect to the swarm public gateway (default \"%s\")", os.Getenv("BEE_GATEWAY")))
	rootCmd.PersistentFlags().StringVar(&optionDataDir, optionNameDataDir, "", "path to datadir directory (default \"./datadir\")")
	rootCmd.PersistentFlags().BoolVar(&optionClean, optionNameClean, false, "delete all downloaded zim and generated tar files")
	rootCmd.PersistentFlags().BoolVar(&optionEnableSearch, optionNameEnableSearch, false, "enable search index")
	rootCmd.PersistentFlags().StringVar(&optionCPUProfile, optionNameCPUProfile, "", "write cpu profile to file")
	rootCmd.PersistentFlags().StringVar(&optionMEMProfile, optionNameMEMProfile, "", "write memory profile to file")
	rootCmd.PersistentFlags().StringVar(&optionBlockProfile, optionNameBlockProfile, "", "write goroutines blocking profile to file")
}

var rootCmd = &cobra.Command{
	Use:           "beezim",
	Short:         "Swarm zim mirror command-line tool",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(_ *cobra.Command, _ []string) (err error) {
		if optionGatewayMode {
			optionBeeApiUrl = os.Getenv("BEE_GATEWAY")
			optionBeeDebugApiUrl = ""
		}

		bee, err = NewBeeClient(optionBeeApiUrl, optionBeeDebugApiUrl)
		if err != nil {
			return err
		}

		if optionCPUProfile != "" || optionMEMProfile != "" || optionBlockProfile != "" {
			if err = startProfiling(); err != nil {
				return err
			}
		}

		return setDataDir()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if optionCPUProfile != "" || optionMEMProfile != "" || optionBlockProfile != "" {
			stopProfiling()
		}
	},
}

func Execute() (err error) {
	rootCmd.AddCommand(
		listWebCmd,
		newDownloadCmd(),
		newUploadCmd(),
		newParserCmd(),
		newMirrorCmd(),
		newCleanCmd(),
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

func NewBeeClient(beeApiUrl string, beeDebugApiUrl string) (*beeclient.BeeClient, error) {
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

	return beeclient.NewBee(opts)
}

func newNetProgressBar(headerText string, size int, eta bool) *pb.ProgressBar {
	var tmpl strings.Builder
	tmpl.WriteString(`{{ string . "header" }} | {{ counters . }} {{ bar . "[" "=" ">" " " "]" }} {{ percent . }} {{ speed . }} `)

	if eta {
		tmpl.WriteString(`{{ rtime . "eta %s" }}`)
	} else {
		tmpl.WriteString(`{{ etime . }}`)
	}

	bar := pb.ProgressBarTemplate(tmpl.String()).New(size)
	bar.Set("header", headerText)
	return bar
}

func startProfiling() (err error) {
	if optionCPUProfile != "" {
		cpuprofile, err = os.Create(optionCPUProfile)
		if err != nil {
			return fmt.Errorf("could not create cpu profile: %v ", err)
		}
		if err := pprof.StartCPUProfile(cpuprofile); err != nil {
			return fmt.Errorf("could not start cpu profile: %v ", err)
		}
	}

	if optionMEMProfile != "" {
		memprofile, err = os.Create(optionMEMProfile)
		if err != nil {
			return fmt.Errorf("could not create memory profile: %v ", err)
		}
	}

	if optionBlockProfile != "" {
		blockprofile, err = os.Create(optionBlockProfile)
		if err != nil {
			return fmt.Errorf("could not create block profile: %v ", err)
		}
		runtime.SetBlockProfileRate(1)
	}
	return nil
}

func stopProfiling() {
	if cpuprofile != nil {
		pprof.StopCPUProfile()
		cpuprofile.Close()
		cpuprofile = nil
	}

	if memprofile != nil {
		pprof.Lookup("heap").WriteTo(memprofile, 0)
		memprofile.Close()
		memprofile = nil
	}

	if blockprofile != nil {
		pprof.Lookup("block").WriteTo(blockprofile, 0)
		blockprofile.Close()
		blockprofile = nil
		runtime.SetBlockProfileRate(0)
	}
}

//TODO: Buy stamps
//TODO: Make manifest metadata
