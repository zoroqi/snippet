package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/zoroqi/snippet/store"
	"log"
	"os"
)

var (
	ERR_FIND_MULTI_SCRIPT = errors.New("find multiple scripts")
	ERR_FIND_NO_SCRIPT    = errors.New("find no script")
)

var rootPath string
var db *store.DB

var RootCmd = &cobra.Command{
	Use:           "snippet",
	Short:         "Simple command-line snippet manager.",
	Long:          `pet - Simple command-line snippet manager.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

var searchFlag store.Search

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&searchFlag.Aliases, "alias", "a", "", "search alias")
	RootCmd.PersistentFlags().StringVarP(&searchFlag.Name, "name", "n", "", "search name")
	RootCmd.PersistentFlags().StringVarP(&searchFlag.ShortName, "short", "s", "", "search short")
	RootCmd.PersistentFlags().StringVar(&searchFlag.Fuzzy, "fuzzy", "", "description fuzzy matching")
	RootCmd.PersistentFlags().StringArrayVar(&searchFlag.Tags, "tag", []string{}, "search tag")
	RootCmd.PersistentFlags().StringVar(&rootPath, "root", "", "root path")
}

func Execute() {
	if e := RootCmd.Execute(); e != nil {
		log.Fatalf("Error: %v\n", e)
		os.Exit(1)
	}

}

func initConfig() {
	if rootPath == "" {
		panic("no root")
	}
	var err error
	db, err = store.Load(rootPath)
	if err != nil {
		log.Fatalf("init error: %v\n", err)
		os.Exit(1)
	}
}
