package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zoroqi/snippet/store"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Show all snippets",
	Long:    `Show all snippets`,
	RunE:    list,
}

func init() {
	RootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	snippets := db.Snippets
	store.SnippetPrintTable(snippets)
	return nil
}
