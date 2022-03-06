package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"swiki/beeclient/api"
	"swiki/indexer"
	"text/tabwriter"

	"github.com/cheggaaa/pb/v3"
	"github.com/joho/godotenv"
)

const (
	baseZimPath  string = "https://download.kiwix.org/zim"
	swarmGateway string = "https://gateway.ethswarm.org"
	depth        uint64 = 30
	amount       int64  = 100000000
)

// List of compressed websites currently maintained by Kiwix
var zims = []string{
	"gutenberg",
	"mooc",
	"other",
	"phet",
	"psiram",
	"stack_exchange",
	"ted",
	"videos",
	"vikidia",
	"wikibooks",
	"wikihow",
	"wikinews",
	"wikipedia",
	"wikiquote",
	"wikisource",
	"wikiversity",
	"wikivoyage",
	"wiktionary",
	"zimit",
}

var (
	gasPrice       string
	zimFile        string
	manifestHash   string // TODO
	batchID        string
	beeApiUrl      string
	beeDebugApiUrl string
	website        string
	gatewayMode    bool
	downloadEnable bool
	mirror         bool // TODO
	listWebsites   bool
	_, pwd, _, _   = runtime.Caller(0)
	basepath       = path.Join(path.Dir(pwd), "..")
	outputDir      = filepath.Join(basepath, "output")
)

func init() {
	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	flag.BoolVar(&gatewayMode, "gateway", false, fmt.Sprintf("connect to the swarm public gateway (default \"%s\")", os.Getenv("BEE_GATEWAY")))
	flag.StringVar(&beeApiUrl, "bee-api-url", os.Getenv("BEE_API_URL"), "bee api url")
	flag.StringVar(&beeDebugApiUrl, "bee-debug-api-url", os.Getenv("BEE_DEBUG_API_URL"), "bee debug api url")
	flag.StringVar(&batchID, "batch", "", "batch ID")
	flag.StringVar(&zimFile, "zim", "", "path for the zim file")
	flag.BoolVar(&downloadEnable, "download", false, fmt.Sprintf("download zim file from: \"%s\"", baseZimPath))
	flag.StringVar(&manifestHash, "manifest", "", "the manifest containing the zim search indexes")
	flag.StringVar(&website, "website", "wikipedia", "compressed website hosted by Kiwix to be mirrored")
	flag.BoolVar(&mirror, "mirror", false, "mirror all compressed websites hosted by Kiwix")
	flag.BoolVar(&listWebsites, "show-sites", false, "shows a list of compressed websites currently maintained by Kiwix")
	flag.Parse()

	// TODO: use cobra command line lib
	if listWebsites {
		printWebsiteList()
		return
	}

	// TODO: download and process multiple zims in parallel
	// create a indexer for each and keep the root manifests in a local kv
	// key: main page title/ or some user defined tag, value: manifest root hash
	// TODO: make a search engine in JS, embed it in the indexes creation per zim

	// wikipedia_es_climate_change_mini_2022-02.zim
	// wikipedia_bm_all_maxi_2022-02.zim
	var zimPath string
	if zimFile != "" {
		if zimFile == "ALL" {
			// TODO: loop downloading and processing all in separate go routines
			// store the manifests metadata
			return
		}
		log.Printf("Looking for zim file at: %s", zimFile)
		if _, err := os.Stat(zimFile); os.IsNotExist(err) {
			zimDownloadPath := fmt.Sprintf("%s/%s", outputDir, filepath.Base(zimFile))
			if _, errOut := os.Stat(zimDownloadPath); os.IsNotExist(errOut) && downloadEnable {
				// download zim
				targetZimURL := fmt.Sprintf("%s/%s", websitePath(website), zimFile)
				zimPath, err = downloadZim(targetZimURL, zimDownloadPath)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				// use existent downloaded zim
				zimPath = zimDownloadPath
			}
		} else {
			// ignores download flag and use given zim
			zimPath = zimFile
		}
	} else {
		log.Fatal("Please inform the path to the zim file to open/download")
	}

	sidx, err := indexer.New(zimPath, outputDir)
	if err != nil {
		log.Fatal(err)
	}

	// var wg sync.WaitGroup
	// quit := make(chan struct{})
	// defer close(quit)

	// Parse zim file
	zimArticles := sidx.ParseZIM()

	// go func() {
	// 	// Build manifest metadata
	// 	wg.Add(1)
	// 	ctx, cancel := context.WithCancel(context.Background())
	// 	defer cancel()

	// 	sidx.MakeManifest(ctx, zimArticles, quit)
	// 	wg.Done()
	// }()

	// Build tar
	tarDirName := strings.TrimSuffix(filepath.Base(sidx.ZimPath), ".zim")
	if err := sidx.TarZim(tarDirName, zimArticles); err != nil {
		log.Fatal(err)
	}

	// Index template uses metadata, so wait for it.
	// wg.Wait()

	if err := sidx.MakeIndexPage(tarDirName); err != nil {
		log.Fatalf("Failed to copy index.html page to tar file: %v", err)
	}

	if err := sidx.MakeErrorPage(tarDirName); err != nil {
		log.Fatalf("Failed to copy error.html page to tar file: %v", err)
	}

	// TODO: load node config from flags/config file
	c, err := indexer.NewUploader(beeApiUrl, beeDebugApiUrl, depth, amount)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: get tag and buy batchID
	// TODO: estimate the cost of upload the zim and the size of the stamp to do so.
	// allow users to agree/deny with the cost before proceed with each upload
	// the stamps will be automatically bought
	err = c.UploadTar(outputDir, api.UploadCollectionOptions{
		Pin:     true,
		BatchID: batchID,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: continue the download from where it stopped in case of crash
// TODO: keep track of already uploaded files (in the metadata kv)
func downloadZim(targetURL string, dstFile string) (string, error) {
	resp, err := http.Get(targetURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error trying to download %v [status: %v]\n", targetURL, resp.Status)
	}

	log.Printf("Downloading zim file to: %v\n", filepath.Base(dstFile))
	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return "", err
	}

	dest, err := os.Create(dstFile)
	if err != nil {
		return "", err
	}
	defer dest.Close()

	progressBar := pb.Full.New(int(size))
	progressBar.Start()

	io.Copy(dest, progressBar.NewProxyReader(resp.Body))

	progressBar.Finish()

	log.Printf("Zim file saved to: %s \n", dstFile)
	return dstFile, nil
}

func websitePath(website string) string {
	return fmt.Sprintf("%s/%s", baseZimPath, website)
}

func printWebsiteList() {
	const sep = "======="

	w := tabwriter.NewWriter(os.Stdout, 2, 8, 2, ' ', 0)
	fmt.Fprintf(w, "%s Kiwix Zims: %d available compressed websites %s\n", sep, len(zims), sep)
	fmt.Fprintf(w, "#\tWebsite\tURL\t\n")
	for i, site := range zims {
		fmt.Fprintf(w, "%00d\t%s\t%s\t", i+1, site, websitePath(site))
		fmt.Fprintln(w, "")
	}
	w.Flush()
}
