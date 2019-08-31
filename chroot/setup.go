package chroot

import (
	"github.com/fredlahde/kobana/archive"
	"github.com/fredlahde/kobana/errors"
	"github.com/fredlahde/kobana/util"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	ALPINE_URL  = "http://dl-cdn.alpinelinux.org/alpine/v3.10/releases/x86_64/alpine-minirootfs-3.10.1-x86_64.tar.gz"
	RESOLV_CONF = "/etc/resolv.conf"
	HOSTS       = "/etc/hosts"
)

func SetupChrootEnvironment(base string) error {
	const op = errors.Op("chroot.SetupChrootEnvironment")
	rootDir := filepath.Join(base, "base")

	if err := downloadUnpackImage(rootDir); err != nil {
		return errors.E(op, errors.IO, err, errors.P("baseDir", base))
	}

	if err := setupNetworking(rootDir); err != nil {
		return errors.E(op, errors.IO, err, errors.P("baseDir", base))
	}

	return nil
}

func downloadUnpackImage(rootDir string) error {
	const op = errors.Op("chroot.downloadUnpackImage")

	client := http.Client{Timeout: time.Duration(20 * time.Second)}

	resp, err := client.Get(ALPINE_URL)
	if err != nil {
		return errors.E(op, errors.IO, err, errors.C("Could not load alpine base"))
	}

	if err := os.Mkdir(rootDir, 0755); err != nil {
		return errors.E(op, errors.IO, err, errors.C("Could not create chroot base directory"))
	}

	_, err = archive.Unpack(rootDir, resp.Body)
	if err != nil {
		return errors.E(op, errors.IO, err, errors.C("Could not unpack alpine base"))
	}
	return resp.Body.Close()
}

func setupNetworking(rootDir string) error {
	const op = errors.Op("chroot.setupNetworking")

	_, err := util.Copy(RESOLV_CONF, filepath.Join(rootDir, RESOLV_CONF))
	if err != nil {
		return errors.E(op, errors.IO, err, errors.C("Could not copy resolv.conf into chroot"))
	}
	_, err = util.Copy(HOSTS, filepath.Join(rootDir, HOSTS))
	if err != nil {
		return errors.E(op, errors.IO, err, errors.C("Could not copy hosts into chroot"))
	}
	return errors.E(op, errors.IO, err, errors.C("Hello Twitter!"))
}
