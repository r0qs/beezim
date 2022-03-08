package indexer

import (
	"archive/tar"
	"bytes"
	"embed"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/r0qs/beezim/internal/tarball"

	zim "github.com/akhenakh/gozim"
	"github.com/cheggaaa/pb/v3"
	"github.com/ethersphere/bee/pkg/swarm"
)

//go:embed assets/*
var assetsFS embed.FS

//go:embed templates/*
var templateFS embed.FS

var templates *template.Template

func init() {
	var err error

	templates, err = template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		log.Fatal("error parsing templates:", err)
	}
}

type Article struct {
	path string
	data []byte
}

func (a Article) Path() string {
	return a.path
}

func (a Article) Data() []byte {
	return a.Data()
}

type SwarmWikiIndexer struct {
	mu      sync.Mutex
	ZimPath string
	Z       *zim.ZimReader
	entries map[string]IndexEntry // RELATIVE_PATH or ArticleID -> METADATA ?
	root    swarm.Address         // TODO: hash of the root manifest metadata (if empty, not uploaded)
}

// TODO: store root in a local kv db pointing to the metadata in swarm
// or maybe in a feed and parse the feed on load to collect all root pages and their metadata.

type IndexEntry struct {
	Path     string
	Metadata map[string]string
}

func New(zimPath string) (*SwarmWikiIndexer, error) {
	z, err := zim.NewReader(zimPath, false)
	if err != nil {
		return nil, err
	}

	return &SwarmWikiIndexer{
		ZimPath: zimPath,
		Z:       z,
		entries: make(map[string]IndexEntry),
	}, nil
}

func (idx *SwarmWikiIndexer) AddEntry(entryPath string, metadata map[string]string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.entries[entryPath] = IndexEntry{
		Path:     entryPath,
		Metadata: metadata,
	}
}

func (idx *SwarmWikiIndexer) Entries() map[string]IndexEntry {
	return idx.entries
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
				return
			}

			// FIXME: for now, all namespaces are considered equal when parsing
			// https://openzim.org/wiki/ZIM_file_format
			var data []byte
			switch a.Namespace {
			case '-', // Assets (CSS, JS, Favicon)
				'A', // Text files (Article Format)
				'I', // Media files
				'M', // ZIM Metadata
				'X': // Search indexes (Xapian db)

				if a.EntryType == zim.RedirectEntry {
					ridx, err := a.RedirectIndex()
					if err != nil {
						return
					}
					ra, err := idx.Z.ArticleAtURLIdx(ridx)
					if err != nil {
						return
					}
					data, err = buildRedirectPage(filepath.Base(ra.FullURL()))
					if err != nil {
						log.Fatalf("error building redirect page: %v", err)
					}
				} else {
					data, err = a.Data()
					if err != nil {
						return
					}
				}

				zimArticles <- Article{
					path: a.FullURL(),
					data: data,
				}

				// TODO: add addresses and searchable data
				idx.AddEntry(a.FullURL(), map[string]string{
					"Title":    a.Title,
					"MimeType": a.MimeType(),
				})

				// TODO: For now we are ignoring some cases, but we should create "_exceptions/" directory in case of errors extracting the files like is done by the zim-tools.
				// https://github.com/openzim/zim-tools/blob/a26a450110e9ca2ec1b20de8237a3bd382af71f5/src/zimdump.cpp#L214
			default:
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
func (idx *SwarmWikiIndexer) TarZim(tarFile string, files <-chan Article) error {
	f, err := os.Create(tarFile)
	if err != nil {
		return err
	}
	defer f.Close()

	tw := tar.NewWriter(f)
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

func buildRedirectPage(pagePath string) ([]byte, error) {
	tmplData := map[string]interface{}{
		"URL": pagePath,
	}

	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, "redirect.html", tmplData); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (idx *SwarmWikiIndexer) MakeIndexPage(tarFile string) error {
	tmplData := map[string]interface{}{
		"File":     path.Base(idx.ZimPath),
		"Count":    strconv.Itoa(int(idx.Z.ArticleCount)),
		"Articles": idx.entries,
	}

	mainPage, err := idx.Z.MainPage()
	if err != nil {
		return err
	}

	if mainPage != nil {
		tmplData["HasMainPage"] = true
		tmplData["MainURL"] = mainPage.FullURL()
	}
	// TODO: handle the case where there is no main page.
	// Use the default index and browse all articles

	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, "index.html", tmplData); err != nil {
		return err
	}
	return tarball.AppendTarData(tarFile, tarball.NewBufferFile("index.html", &buf))
}

func (idx *SwarmWikiIndexer) MakeErrorPage(tarFile string) error {
	tmplData := map[string]interface{}{
		"File": path.Base(idx.ZimPath),
	}
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, "error.html", tmplData); err != nil {
		return err
	}
	return tarball.AppendTarData(tarFile, tarball.NewBufferFile("error.html", &buf))
}
