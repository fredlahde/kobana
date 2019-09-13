package cmd

import (
	"github.com/fredlahde/kobana/chroot"
	"github.com/fredlahde/kobana/config"
	"github.com/fredlahde/kobana/mount"
	"github.com/fredlahde/kobana/safety"
	"github.com/spf13/cobra"
	"log"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "creates a new chroot environment",
	Run: run,
}

func run(cmd *cobra.Command, args []string) {
	conf := config.Get()
	base, err := mount.MountRamFs(conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("base is:", base)
	err = chroot.SetupChrootEnvironment(base, conf)
	if err != nil {
		log.Fatal(err)
	}
	err = safety.DropRootPriviliges(conf)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
}
