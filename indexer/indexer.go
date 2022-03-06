package indexer

import (
	"archive/tar"
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"swiki/internal/tarball"
	"sync"
	"time"

	zim "github.com/akhenakh/gozim"
	"github.com/cheggaaa/pb/v3"
	"github.com/ethersphere/bee/pkg/swarm"
)

//go:embed assets/*
var assetsFS embed.FS

//go:embed templates/*
var templateFS embed.FS

type Article struct {
	path     string
	title    string
	data     []byte
	mimeType string
}

func (a Article) Path() string {
	return a.path
}

func (a Article) Data() []byte {
	return a.Data()
}

type SwarmWikiIndexer struct {
	mu        sync.Mutex
	outputDir string
	ZimPath   string
	Z         *zim.ZimReader
	m         map[string]ManifestEntry // TODO: replace by a real swarm manifest?
	uploaded  string                   // TODO: hash of the root manifest metadata (if empty, not uploaded)
	templates *template.Template
}

// ManifestEntry abstract a swarm manifest entry for a local
type ManifestEntry struct {
	Reference string
	Path      string
	Metadata  map[string]string // TODO: define all metadata entries
}

func New(zimPath string, outputDir string) (*SwarmWikiIndexer, error) {
	// load base templates
	templates, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, err
	}

	z, err := zim.NewReader(zimPath, false)
	if err != nil {
		return nil, err
	}

	return &SwarmWikiIndexer{
		outputDir: outputDir,
		ZimPath:   zimPath,
		Z:         z,
		m:         make(map[string]ManifestEntry),
		templates: templates,
	}, nil
}

func (idx *SwarmWikiIndexer) AddEntry(entry ManifestEntry) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.m[entry.Path] = ManifestEntry{
		Reference: entry.Reference,
		Path:      entry.Path,
		Metadata:  entry.Metadata,
	}
}

func (idx *SwarmWikiIndexer) Manifest() map[string]ManifestEntry {
	return idx.m
}

func (idx *SwarmWikiIndexer) ParseZIM() chan Article {
	zimArticles := make(chan Article)
	go func() {
		defer close(zimArticles)
		progressBar := pb.New(int(idx.Z.ArticleCount))
		progressBar.Set(pb.Bytes, true)
		progressBar.Start()

		log.Printf("Parsing zim file: %s", filepath.Base(idx.ZimPath))
		start := time.Now()
		idx.Z.ListTitlesPtrIterator(func(i uint32) {
			a, err := idx.Z.ArticleAtURLIdx(i)
			if err != nil || a.EntryType == zim.DeletedEntry {
				log.Fatalf("Error or deleted entry: %v - %v", err, a.EntryType)
			}

			// FIXME: for now, all namespaces are considered equal when parsing
			// https://openzim.org/wiki/ZIM_file_format
			// TODO: add search indexes
			switch a.Namespace {
			case '-', // Assets (CSS, JS, Favicon)
				'A', // Text files (Article Format) ?
				'I', // Media files?
				'M', // ZIM Metadata
				'X': // Search indexes?
				data, err := a.Data()
				if err != nil {
					log.Fatal(err)
				}
				article := Article{
					path:     a.FullURL(),
					title:    a.Title,
					data:     data,
					mimeType: a.MimeType(),
				}
				zimArticles <- article
				// TODO: For now we are ignoring some cases, but we should create "_exceptions/" directory in case of errors extracting the files like is done by the zim-tools.
				// https://github.com/openzim/zim-tools/blob/a26a450110e9ca2ec1b20de8237a3bd382af71f5/src/zimdump.cpp#L214
			default:
				fmt.Println("ignoring entry:", a.FullURL())
			}
			progressBar.Increment()
		})
		progressBar.Finish()
		elapsed := time.Since(start)
		log.Printf("File processed in %v", elapsed)
	}()
	return zimArticles
}

// makeManifest consumes zimArticles and create manifest entries
func (idx *SwarmWikiIndexer) MakeManifest(ctx context.Context, files <-chan Article, quit chan struct{}) {
	for a := range files {
		buf := bytes.NewBuffer(a.data)
		file := tarball.NewBufferFile(a.path, buf)
		file.CalculateHash()
		reference := swarm.NewAddress(file.Hash())

		entry := ManifestEntry{
			Reference: reference.String(),
			Path:      a.path,
			Metadata: map[string]string{
				"Title":    a.title,
				"MimeType": a.mimeType,
				// TODO: add addresses and searchable data
			},
		}

		idx.AddEntry(entry)

		select {
		case <-quit:
			return
		default:
		}
	}
}

// TODO: move file operations to its own package
func (idx *SwarmWikiIndexer) TarZim(tarDir string, files <-chan Article) error {
	_, err := os.Stat(idx.outputDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(idx.outputDir, 0755); err != nil {
			return err
		}
	}

	tarFilename := filepath.Join(idx.outputDir, fmt.Sprintf("%s.tar", tarDir))
	tarFile, err := os.Create(tarFilename)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	tw := tar.NewWriter(tarFile)
	for file := range files {
		hdr := &tar.Header{
			Name: file.path,
			Mode: 0600,
			Size: int64(len(file.data)),
		}
		// TODO: add file.mimeType ?
		// TODO: set header.Typeflag (tar.TypeDir, tar.TypeReg)

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if _, err := tw.Write(file.data); err != nil {
			return err
		}
	}

	if err := tw.Close(); err != nil {
		return err
	}
	return nil
}

func (idx *SwarmWikiIndexer) MakeIndexPage(tarDir string) error {

	tmplData := map[string]interface{}{
		"Path":     path.Base(idx.ZimPath),
		"Count":    strconv.Itoa(int(idx.Z.ArticleCount)),
		"Manifest": idx.m,
	}

	mainPage, err := idx.Z.MainPage()
	if err != nil {
		return err
	}

	if mainPage != nil {
		tmplData["HasMainPage"] = true
		tmplData["MainURL"] = mainPage.FullURL()
	}

	var buf bytes.Buffer
	if err := idx.templates.ExecuteTemplate(&buf, "index.html", tmplData); err != nil {
		return err
	}
	tarFile := filepath.Join(idx.outputDir, fmt.Sprintf("%s.tar", tarDir))
	return tarball.AppendTarData(tarFile, tarball.NewBufferFile("index.html", &buf))
}

func (idx *SwarmWikiIndexer) MakeErrorPage(tarDir string) error {
	data := map[string]interface{}{
		"Path": path.Base(idx.ZimPath),
	}
	var buf bytes.Buffer
	if err := idx.templates.ExecuteTemplate(&buf, "error.html", data); err != nil {
		return err
	}
	tarFile := filepath.Join(idx.outputDir, fmt.Sprintf("%s.tar", tarDir))
	return tarball.AppendTarData(tarFile, tarball.NewBufferFile("error.html", &buf))
}
