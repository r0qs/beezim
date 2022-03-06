package tarball

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
)

func AppendTarData(tarFile string, file *File) error {
	tf, err := os.OpenFile(tarFile, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer tf.Close()

	// https://www.freebsd.org/cgi/man.cgi?query=tar&sektion=5
	// A tar archive consists of a series	of 512-byte records.
	// The end of the archive is indicated by two records consisting entirely of zero bytes.
	// To append to it we start the write 1024 bytes before the end.
	if _, err := tf.Seek(-1<<10, io.SeekEnd); err != nil {
		return err
	}
	tw := tar.NewWriter(tf)

	hdr := &tar.Header{
		Name: file.name,
		Mode: 0600,
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
