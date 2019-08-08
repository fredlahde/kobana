package util

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
)

// Copy copies the file at src into a file at dst
// It returns the number of bytes copied or an error
func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, errors.Wrap(err, "Could not load info from file")
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, errors.Wrap(err, "Could not open source file")
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, errors.Wrap(err, "Could not create destination file")
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)

	// We can safely call Close() twice, since it's a no-op when the file is already closed
	// It'll return an error, but we ignore it in the defer
	if err := source.Close(); err != nil {
		return 0, errors.Wrap(err, "Could not close source file")
	}
	if err := destination.Close(); err != nil {
		return 0, errors.Wrap(err, "Could not close destination file")
	}

	return nBytes, err
}
