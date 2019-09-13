package cmd

import (
	"fmt"
	"github.com/fredlahde/kobana/config"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kobana",
	Short: "Instant chroots!",
	Long:  `kobana is a tool to create docker-like envrionments, but with chroots`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func initConfig() {
	if cfgFile == "" {
		log.Fatal("You need to pass a config via --config")
	}

	err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatal(err)
	}
}
