package archive

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestUnpack(t *testing.T) {
	unpackDir := filepath.Join("../test_fixtures", "unpack_test")
	tarFile, err := os.Open(filepath.Join("../test_fixtures", "test.tar.gz"))
	if err != nil {
		t.Fatal("unable to open test archive", err)
	}
	numFiles, err := Unpack(unpackDir, tarFile)
	if err != nil {
		t.Fatal(err)
	}
	if numFiles != 1 {
		t.Errorf("Got %d want %d files", numFiles, 1)
	}

	want := []byte{'A', 'A', 'A', 'A', 'A', 'A', '\n'}
	got, err := ioutil.ReadFile(filepath.Join(unpackDir, "foo.txt"))
	if err != nil {
		t.Fatal("Unable to read extracted file", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Unpacked file has wrong contents, got %v want %v", got, want)
	}

}
