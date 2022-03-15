package indexer

import (
	"archive/tar"
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
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
)

//go:embed assets/*
var assetsFS embed.FS

//go:embed templates/*
var templateFS embed.FS

var assets []fs.DirEntry

func init() {
	var err error

	assets, err = assetsFS.ReadDir("assets")
	if err != nil {
		log.Fatal("error parsing assets:", err)
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

type IndexMetadata struct {
	Title    string
	MimeType string
	Redirect bool
}

type SwarmWikiIndexer struct {
	mu      sync.Mutex
	ZimPath string
	Z       *zim.ZimReader
	entries map[string]IndexEntry // RELATIVE_PATH or ArticleID -> METADATA ?
}

// TODO: store root in a local kv db pointing to the metadata in swarm
// or maybe in a feed and parse the feed on load to collect all root pages and their metadata.

type IndexEntry struct {
	Path     string
	Metadata IndexMetadata
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

func (idx *SwarmWikiIndexer) AddEntry(entryPath string, metadata IndexMetadata) {
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
				// TODO: when is enable-search true append xapian scripts; exclude it from tar (only keep the xapianasm.data), and modify index template to load the js.

				if a.EntryType == zim.RedirectEntry {
					ridx, err := a.RedirectIndex()
					if err != nil {
						return
					}

					ra, err := idx.Z.ArticleAtURLIdx(ridx)
					if err != nil {
						return
					}

					buf, err := buildRedirectPage(path.Base(ra.FullURL()))
					if err != nil {
						log.Fatalf("error building redirect page: %v", err)
					}
					data = buf.Bytes()

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
				idx.AddEntry(a.FullURL(), IndexMetadata{
					Title:    a.Title,
					MimeType: a.MimeType(),
					Redirect: a.EntryType == zim.RedirectEntry,
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

func (idx *SwarmWikiIndexer) UnZim(outputDir string, files <-chan Article) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}
	}

	for file := range files {
		filePath := filepath.Join(outputDir, file.path)
		fileDirPath := filepath.Dir(filePath)

		if _, err := os.Stat(fileDirPath); os.IsNotExist(err) {
			if err := os.MkdirAll(fileDirPath, 0755); err != nil {
				return err
			}
		}

		f, err := os.Create(filePath)
		if err != nil {
			return err
		}

		if _, err := f.Write(file.data); err != nil {
			return err
		}

		f.Close()
	}

	return nil
}

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

func buildRedirectPage(pagePath string) (*bytes.Buffer, error) {
	tmplData := map[string]interface{}{
		"Path": pagePath,
	}

	redirectTmpl, err := template.ParseFS(templateFS, "templates/index-redirect.html")
	if err != nil {
		return nil, fmt.Errorf("error parsing index redirect template: %v", err)
	}

	var buf bytes.Buffer
	if err := redirectTmpl.ExecuteTemplate(&buf, "index-redirect.html", tmplData); err != nil {
		return nil, err
	}
	return &buf, nil
}

// MakeRedirectIndexPage creates an redirect index to the main page
// when it exists in the zim archive.
func (idx *SwarmWikiIndexer) MakeRedirectIndexPage(tarFile string) error {
	log.Printf("Appending redirect index.html to %s", filepath.Base(tarFile))

	mainPage, err := idx.Z.MainPage()
	if err != nil {
		return err
	}

	// TODO: handle the case where there is no main page in the article.
	// Should we add an index and browse all articles?
	if mainPage == nil {
		return errors.New("no index found in the ZIM")
	}

	buf, err := buildRedirectPage(mainPage.FullURL())
	if err != nil {
		return err
	}

	return tarball.AppendTarFile(tarFile, tarball.NewBufferFile("index.html", buf))
}

// parseTemplate parses a given template and replace content when requested
func parseTemplate(contentTmpl string, data interface{}) (*bytes.Buffer, error) {
	baseTmpl, err := template.ParseFS(templateFS, "templates/page/*.html")
	if err != nil {
		return nil, fmt.Errorf("error parsing base templates: %v", err)
	}

	// add dynamic content to pages
	// FIXME: current we only support replace the content. Maybe we can improve that in the future do to something like Hugo does, or use Hugo instead.
	if contentTmpl != "" {
		tmpl, err := template.New("content").ParseFS(templateFS, fmt.Sprintf("templates/%s", contentTmpl))
		if err != nil {
			return nil, err
		}

		// don't attempt to add in the tree if their is nothing to be added
		if tmpl != nil && tmpl.Tree != nil {
			_, err = baseTmpl.AddParseTree("content", tmpl.Tree)
			if err != nil {
				return nil, err
			}
		}
	}

	var buf bytes.Buffer
	if err := baseTmpl.ExecuteTemplate(&buf, "page", data); err != nil {
		return nil, err
	}

	return &buf, nil
}

// makePage creates a page with a given template data
func makePage(name, template string, tmplData map[string]interface{}, tarFile string) error {
	log.Printf("Appending %s page to %s", name, filepath.Base(tarFile))

	buf, err := parseTemplate(template, tmplData)
	if err != nil {
		return err
	}

	return tarball.AppendTarFile(tarFile, tarball.NewBufferFile(name, buf))
}

type Node struct {
	Path     string  `json:"path"`
	Icon     string  `json:"icon"`
	MimeType string  `json:"mimeType"`
	Title    string  `json:"title"`
	Redirect bool    `json:"redirect"`
	Nodes    []*Node `json:"nodes"`
}

func groupDataByPrefix(idxEntries map[string]IndexEntry) map[string]*Node {
	m := make(map[string]*Node)
	for p, entry := range idxEntries {
		n := &Node{
			Path:     entry.Path,
			MimeType: entry.Metadata.MimeType,
			Title:    entry.Metadata.Title,
			Redirect: entry.Metadata.Redirect,
			Icon:     "",
		}
		var id string
		switch path.Dir(p)[0] {
		case '-':
			id = "Assets"
		case 'A':
			id = "Articles"
		case 'I':
			id = "Media"
		case 'M':
			id = "Metadata"
		case 'X': // FIXME: remove when add wasm
			id = "Indexes"
		}

		if _, ok := m[id]; !ok {
			m[id] = &Node{
				Path:  id,
				Icon:  "", //TODO add icons style
				Nodes: make([]*Node, 0),
			}
		}
		m[id].Nodes = append(m[id].Nodes, n)
	}
	return m
}

// MakeIndexSearchPage creates a custom index with the text search tool and
// embed the current main page in the new index.
func (idx *SwarmWikiIndexer) MakeIndexSearchPage(tarFile string) error {
	mainPage, err := idx.Z.MainPage()
	if err != nil {
		return err
	}

	mainURL := ""
	if mainPage != nil {
		mainURL = mainPage.FullURL()
	}

	tmplData := map[string]interface{}{
		"File":        filepath.Base(idx.ZimPath),
		"Count":       strconv.Itoa(int(idx.Z.ArticleCount)),
		"Articles":    groupDataByPrefix(idx.entries),
		"HasMainPage": (mainURL != ""),
		"MainURL":     mainURL,
	}

	// make about's page using about template
	if err = makePage("about.html", "about.html", tmplData, tarFile); err != nil {
		return err
	}

	// make browse files page using files template
	if err = makePage("files.html", "files.html", tmplData, tarFile); err != nil {
		return err
	}

	// make files page in JSON format
	if file, err := json.Marshal(idx.entries); err == nil {
		if err = tarball.AppendTarFile(tarFile, tarball.NewBufferFile("files.json", bytes.NewBuffer(file))); err != nil {
			return err
		}
	}

	// make index page using index-search template
	return makePage("index.html", "index-search.html", tmplData, tarFile)
}

// MakeErrorPage creates an error page
func (idx *SwarmWikiIndexer) MakeErrorPage(tarFile string) error {
	data, err := fs.ReadFile(templateFS, "templates/error.html")
	if err != nil {
		return err
	}

	return tarball.AppendTarFile(tarFile, tarball.NewBytesFile("error.html", data))
}

func AddAssets(tarFile string) error {
	log.Printf("Appending assets to %s", filepath.Base(tarFile))

	return fs.WalkDir(assetsFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		data, err := fs.ReadFile(assetsFS, path)
		if err != nil {
			return err
		}

		if err = tarball.AppendTarFile(tarFile, tarball.NewBytesFile(path, data)); err != nil {
			return err
		}

		return nil
	})
}
