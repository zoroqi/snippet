package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zoroqi/snippet/store"
	"log"
	"strings"
)

var showCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"s"},
	Short:   "Show snippet detail",
	Long:    `Show snippet detail`,
	RunE:    show,
}

func init() {
	RootCmd.AddCommand(showCmd)
}

func show(cmd *cobra.Command, args []string) error {
	snippets := db.Find(searchFlag, args)
	for _, snippet := range snippets {
		s, err := store.ShowSnippet(snippet)
		fmt.Println(s)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Println(strings.Repeat("-", 80))
	}
	return nil
}
