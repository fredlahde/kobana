package mount

import (
	names "github.com/docker/docker/pkg/namesgenerator"
	"github.com/pkg/errors"
	syscall "golang.org/x/sys/unix"
	"os"
	"path/filepath"
)

const (
	MOUNT_BASE   string = "/mnt/kobana"
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
	path := filepath.Join(MOUNT_BASE, names.GetRandomName(0))
	return path, os.MkdirAll(path, 0755)
}
