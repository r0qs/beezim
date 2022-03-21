package tarball

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func AppendTarFile(tarFile string, file *File) error {
	tf, err := os.OpenFile(tarFile, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer tf.Close()

	// https://www.freebsd.org/cgi/man.cgi?query=tar&sektion=5
	// A tar archive consists of a series of 512-byte records.
	// The end of the archive is indicated by two records consisting entirely of zero bytes.
	// To append to it we start the write 1024 bytes before the end.
	if _, err := tf.Seek(-1<<10, io.SeekEnd); err != nil {
		return err
	}
	tw := tar.NewWriter(tf)

	hdr := &tar.Header{
		Name: file.name,
		Mode: 0644,
		Size: file.size,
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}

	if _, err := io.Copy(tw, file.DataReader()); err != nil {
		return err
	}

	return tw.Close()
}

// CopyTar copies a tar file from the reader and writes it to the writer
func CopyTar(w io.Writer, tr *tar.Reader) error {
	tw := tar.NewWriter(w)

	for {
		hdr, err := tr.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if _, err := io.Copy(tw, tr); err != nil {
			return err
		}
	}

	if err := tw.Close(); err != nil {
		return err
	}

	return nil
}

func ReadTarBuffer(tarFile string) (*bytes.Buffer, error) {
	f, err := os.Open(tarFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	tr := tar.NewReader(f)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}

		if _, err := io.Copy(tw, tr); err != nil {
			return nil, err
		}
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}

func Untar(tarFile string, targetDir string) error {
	reader, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	tr := tar.NewReader(reader)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		filePath := filepath.Join(targetDir, header.Name)
		baseFilePath := filepath.Dir(filePath)
		if _, err := os.Stat(baseFilePath); os.IsNotExist(err) {
			if err := os.MkdirAll(baseFilePath, 0755); err != nil {
				return err
			}
		}

		switch header.Typeflag {
		case tar.TypeReg:
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			f.Close()
		}
	}
	return nil
}

func Tar(sourceDir string, tarFile string) error {
	if filepath.Ext(tarFile) != ".tar" {
		return errors.New("target tar file must have \".tar\" extention")
	}

	out := filepath.Join(filepath.Dir(sourceDir), tarFile)
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	tw := tar.NewWriter(f)
	defer tw.Close()

	info, err := os.Stat(sourceDir)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(sourceDir)
	}

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		hdr, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		hdr.Name = filepath.Join(baseDir, strings.TrimPrefix(path, sourceDir))
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		// skip write if is not a regular file
		if !info.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tw, file)
		if err != nil {
			return err
		}

		return nil
	})
}
