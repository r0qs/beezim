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

	"github.com/cheggaaa/pb/v3"
	"github.com/joho/godotenv"
)

const (
	baseWikiPath string = "https://download.kiwix.org/zim/wikipedia"
	swarmGateway string = "https://gateway.ethswarm.org"
	depth        uint64 = 30
	amount       int64  = 100000000
)

var (
	gasPrice       string
	zimFile        string
	manifestHash   string
	batchID        string
	beeApiUrl      string
	beeDebugApiUrl string
	gatewayMode    bool
	downloadEnable bool
	_, pwd, _, _   = runtime.Caller(0)
	basepath       = path.Join(path.Dir(pwd), "..")
	outputDir      = filepath.Join(basepath, "output")
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	flag.BoolVar(&gatewayMode, "gateway", false, fmt.Sprintf("connect to the swarm public gateway (default \"%s\")", os.Getenv("BEE_GATEWAY")))
	flag.StringVar(&beeApiUrl, "bee-api-url", os.Getenv("BEE_API_URL"), "bee api url")
	flag.StringVar(&beeDebugApiUrl, "bee-debug-api-url", os.Getenv("BEE_DEBUG_API_URL"), "bee debug api url")
	flag.StringVar(&batchID, "batch", "", "batch ID")
	flag.StringVar(&zimFile, "zim", "", "path for the zim file")
	flag.BoolVar(&downloadEnable, "download", true, fmt.Sprintf("download zim file from: \"%s\"", baseWikiPath))
	flag.StringVar(&manifestHash, "manifest", "", "the manifest containing the zim search indexes")
	flag.Parse()

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
				targetZimURL := fmt.Sprintf("%s/%s", baseWikiPath, zimFile)
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

	// Parse zim file
	zimArticles := sidx.ParseZIM()

	// Build tar
	tarDirName := strings.TrimSuffix(filepath.Base(sidx.ZimPath), ".zim")
	if err := sidx.TarZim(tarDirName, zimArticles); err != nil {
		log.Fatal(err)
	}

	if err := sidx.MakeIndexPage(tarDirName); err != nil {
		log.Fatal(err)
	}

	if err := sidx.MakeErrorPage(tarDirName); err != nil {
		log.Fatal(err)
	}

	// TODO: load node config from flags/config file
	c, err := indexer.NewUploader(beeApiUrl, beeDebugApiUrl, depth, amount)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: get tag and buy batchID
	err = c.UploadTar(outputDir, api.UploadCollectionOptions{
		Pin:     true,
		BatchID: batchID,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func downloadZim(targetURL string, dstFile string) (string, error) {
	resp, err := http.Get(targetURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error status: %v\n", resp.Status)
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
