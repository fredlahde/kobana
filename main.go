package main

import (
	"flag"
	"github.com/fredlahde/kobana/chroot"
	"github.com/fredlahde/kobana/config"
	"github.com/fredlahde/kobana/mount"
	"github.com/fredlahde/kobana/safety"
	"log"
	"os"
)

var (
	configFile string
)

func userHome() string {
	return os.Getenv("HOME")
}

func init() {
	flag.StringVar(&configFile, "config", "", "path to config")
}

func main() {
	flag.Parse()

	if configFile == "" {
		log.Fatal("You need to specify a path to a config")
	}

	conf, err := config.ParseConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	base, err := mount.MountRamFs(conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("base is:", base)
	err = chroot.SetupChrootEnvironment(base, conf)
	if err != nil {
		log.Fatal(err)
	}
	err = safety.DropRootPriviliges("lothar", "lothar")
	if err != nil {
		log.Fatal(err)
	}
}
