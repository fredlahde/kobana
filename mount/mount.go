package mount

import (
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	MOUNT_BASE   string = "/mnt/kobana"
	RAND_LEN     int    = 32
	LETTER_BYTES string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func MountRamFs() (string, error) {
	base, err := makeMountBaseDir()
	if err != nil {
		return "", errors.Wrap(err, "unable to create mount base")
	}
	var (
		flags   uintptr
		options string = ""
	)
	flags = syscall.MS_NOATIME | syscall.MS_SILENT
	err = syscall.Mount("ramfs", base, "ramfs", flags, options)
	if err != nil {
		return "", errors.Wrap(err, "unable to mount ramfs")
	}
	return base, nil
}

func UmountRamFs(base string) error {
	if err := syscall.Unmount(base, 0); err != nil {
		return err
	}
	return os.RemoveAll(base)
}

func makeMountBaseDir() (string, error) {
	path := filepath.Join(MOUNT_BASE, randStringBytes(RAND_LEN))
	return path, os.MkdirAll(path, 0755)
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = LETTER_BYTES[rand.Intn(len(LETTER_BYTES))]
	}
	return string(b)
}
