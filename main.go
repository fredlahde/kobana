package main

import (
	"github.com/fredlahde/kobana/chroot"
	"github.com/fredlahde/kobana/mount"
	"github.com/fredlahde/kobana/safety"
	"log"
)

func main() {
	base, err := mount.MountRamFs()
	if err != nil {
		log.Fatal("unable to mount ramfs: ", err)
	}
	log.Println("base is:", base)
	err = chroot.SetupChrootEnvironment(base)
	if err != nil {
		log.Fatal("unable to create chroot environment: ", err)
	}
	err = safety.DropRootPriviliges("lothar", "lothar")
	if err != nil {
		log.Fatal("Unable to drop root priviliges", err)
	}
}
