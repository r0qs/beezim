package indexer

import (
	"archive/tar"
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"swiki/internal/tarball"
	"text/template"
	"time"

	zim "github.com/akhenakh/gozim"
	"github.com/cheggaaa/pb/v3"
	"github.com/ethersphere/bee/pkg/manifest/simple"
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
	metadata map[string]string
}

func (a Article) Path() string {
	return a.path
}

func (a Article) Data() []byte {
	return a.Data()
}

func (a Article) Metadata() map[string]string {
	return a.metadata
}

type SwarmWikiIndexer struct {
	outputDir string
	ZimPath   string
	Z         *zim.ZimReader
	m         simple.Manifest
	templates *template.Template
}

type ManifestEntry struct {
	Reference swarm.Address
	Path      string
	Metadata  map[string]string
}

func New(zimPath string, outputDir string) (*SwarmWikiIndexer, error) {
	// TODO: load assets for the searcher
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
		m:         simple.NewManifest(),
		templates: templates,
	}, nil
}

func (idx *SwarmWikiIndexer) AddEntry(entry ManifestEntry) {
	idx.m.Add(entry.Path, entry.Reference.String(), entry.Metadata)
}

func (idx *SwarmWikiIndexer) Manifest() simple.Manifest {
	return idx.m
}

func (idx *SwarmWikiIndexer) ParseZIM() chan *Article {
	zimArticles := make(chan *Article)
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
				article := &Article{
					path:     a.FullURL(),
					title:    a.Title,
					data:     data,
					mimeType: a.MimeType(),
					metadata: make(map[string]string), // TODO: add search info
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

// TODO: move file operations to its own package
func (idx *SwarmWikiIndexer) TarZim(tarDir string, files <-chan *Article) error {
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
		"Path":  path.Base(idx.ZimPath),
		"Count": strconv.Itoa(int(idx.Z.ArticleCount)),
		// "Manifest":    idx.m, // TODO: for searching
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
