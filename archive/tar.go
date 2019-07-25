package archive

import (
	"archive/tar"
	"compress/gzip"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

func Unpack(dst string, r io.Reader) (uint32, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return 0, err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	var numFiles uint32 = 0
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return numFiles, nil
		case err != nil:
			return numFiles, err
		// TODO Does this really happen?
		case header == nil:
			continue
		}

		numFiles = numFiles + 1
		target := filepath.Join(dst, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := mkdir(target); err != nil {
				return numFiles, err
			}
		case tar.TypeReg:
			if err := writeFile(target, header, tr); err != nil {
				return numFiles, err
			}
		}
	}
}

func writeFile(target string, header *tar.Header, tr io.Reader) error {
	f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, os.FileMode(header.Mode))
	// TODO Check for permissions error
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, tr); err != nil {
		return err
	}
	return f.Close()
}

func mkdir(target string) error {
	_, err := os.Stat(target)

	if os.IsNotExist(err) {
		return nil
	}

	if os.IsPermission(err) {
		return errors.Wrapf(err, "Not allowed to create directory %s, check permissions", target)
	}

	if err := os.Mkdir(target, 0755); err != nil {
		return errors.Wrapf(err, "Error creating directory %s", target)
	}
	return nil
}
