package mount

import (
	names "github.com/docker/docker/pkg/namesgenerator"
	"github.com/fredlahde/kobana/config"
	"github.com/fredlahde/kobana/errors"
	syscall "golang.org/x/sys/unix"
	"os"
	"path/filepath"
)

func MountRamFs(conf *config.Config) (string, error) {
	var op = errors.Op("mount.MountRamFs")
	base, err := makeMountBaseDir(conf.BaseDir)
	if err != nil {
		return "", errors.E(op, errors.IO, err, errors.C("unable to create mount base"))
	}
	var (
		flags   uintptr
		options string = ""
	)
	flags = syscall.MS_NOATIME | syscall.MS_SILENT
	err = syscall.Mount("ramfs", base, "ramfs", flags, options)
	if err != nil {
		return "", errors.E(op, errors.KindFromSyscallErrno(err), err, errors.C("unable to mount ramfs"))
	}
	return base, nil
}

func UmountRamFs(base string) error {
	op := errors.Op("mount.UmountRamFs")
	if err := syscall.Unmount(base, 0); err != nil {
		return errors.E(op, errors.IO, err, errors.C("unable to unmount ram fs"))
	}
	if err := os.RemoveAll(base); err != nil {
		return errors.E(op, errors.IO, err, errors.C("could not delete ram fs base dir"))
	}
	return nil
}

func makeMountBaseDir(root string) (string, error) {
	op := errors.Op("mount.makeMountBaseDir")

	path := filepath.Join(root, names.GetRandomName(0))
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", errors.E(op, errors.IO, err, errors.P("mountBaseDir", root))
	}
	return path, nil
}
