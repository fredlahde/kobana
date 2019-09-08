package chroot

import (
	"github.com/fredlahde/kobana/archive"
	"github.com/fredlahde/kobana/config"
	"github.com/fredlahde/kobana/errors"
	"github.com/fredlahde/kobana/util"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	RESOLV_CONF = "/etc/resolv.conf"
	HOSTS       = "/etc/hosts"
)

func SetupChrootEnvironment(base string, config *config.Config) error {
	const op = errors.Op("chroot.SetupChrootEnvironment")
	rootDir := filepath.Join(base, "base")

	if err := downloadUnpackImage(rootDir, config); err != nil {
		return errors.E(op, errors.IO, err, errors.P("baseDir", base))
	}

	if err := setupNetworking(rootDir); err != nil {
		return errors.E(op, errors.IO, err, errors.P("baseDir", base))
	}

	return nil
}

func downloadUnpackImage(rootDir string, config *config.Config) error {
	const op = errors.Op("chroot.downloadUnpackImage")

	client := http.Client{Timeout: time.Duration(20 * time.Second)}

	log.Println("Downloading", config.BaseImage)
	resp, err := client.Get(config.BaseImage)
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
	return nil
}
