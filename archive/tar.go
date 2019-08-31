package archive

import (
	"archive/tar"
	"compress/gzip"
	"github.com/fredlahde/kobana/errors"
	"io"
	"os"
	"path/filepath"
)

func Unpack(dst string, r io.Reader) (uint32, error) {
	const op errors.Op = "archive.Unpack"
	gz, err := gzip.NewReader(r)
	if err != nil {
		return 0, errors.E(op, errors.IO, err, errors.C("Unable to create gzip reader"))
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	var numFiles uint32 = 0
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			if err = gz.Close(); err != nil {
				return numFiles, errors.E(op, errors.IO, err)
			}
			return numFiles, nil
		case err != nil:
			return numFiles, errors.E(op, errors.IO, err, errors.C("failed to iterate archive files"))
		// TODO Does this really happen?
		case header == nil:
			continue
		}

		numFiles = numFiles + 1
		target := filepath.Join(dst, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := mkdir(target); err != nil {
				return numFiles, errors.E(op, errors.IO, err)
			}
		case tar.TypeReg:
			if err := writeFile(target, header, tr); err != nil {
				return numFiles, errors.E(op, errors.IO, err)
			}
		}
	}
}

func writeFile(target string, header *tar.Header, tr io.Reader) error {
	const op errors.Op = "archive.writeFile"
	_, err := os.Stat(target)
	var f *os.File

	if os.IsNotExist(err) {
		f, err = os.Create(target)
	} else if err != nil {
		return errors.E(op, errors.IO, err, errors.P("file", target))
	} else {
		f, err = os.OpenFile(target, os.O_WRONLY|os.O_CREATE, os.FileMode(header.Mode))
	}

	// TODO Check for permissions error
	if err != nil {
		return errors.E(
			op, errors.IO, err, errors.C("Unable to create file"), errors.P("file", target))
	}

	if _, err := io.Copy(f, tr); err != nil {
		return errors.E(
			op, errors.IO, err, errors.C("unable to extract file"), errors.P("file", target))
	}
	err = f.Close()
	if err != nil {
		return errors.E(op, errors.IO, err)
	}
	return nil
}

func mkdir(target string) error {
	const op errors.Op = "archive.mkdir"
	_, err := os.Stat(target)
	if os.IsExist(err) {
		return nil
	}

	if os.IsPermission(err) {
		return errors.E(op, errors.Permission, err, errors.C("Not allowed to create directory, check permissions"))
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return errors.E(op, errors.IO, err, errors.C("Error creating directory"))
	}

	return nil
}
