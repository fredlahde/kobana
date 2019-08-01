package chroot

import (
	"github.com/fredlahde/kobana/archive"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const ALPINE_URL = "http://dl-cdn.alpinelinux.org/alpine/v3.10/releases/x86_64/alpine-minirootfs-3.10.1-x86_64.tar.gz"

func SetupChrootEnvironment(base string) error {
	client := http.Client{Timeout: time.Duration(20 * time.Second)}

	resp, err := client.Get(ALPINE_URL)
	if err != nil {
		return errors.Wrap(err, "Could not load alpine base")
	}
	chrootBase := filepath.Join(base, "base")
	if err := os.Mkdir(chrootBase, 0755); err != nil {
		return errors.Wrap(err, "Could not create chroot base directory")
	}
	_, err = archive.Unpack(chrootBase, resp.Body)
	if err != nil {
		return errors.Wrap(err, "Could not unpack alpine base")
	}
	return resp.Body.Close()
}
