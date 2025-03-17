package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zoroqi/snippet/execute"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit snippet file",
	Long:  `Edit snippet file (opened by vim)`,
	RunE:  edit,
}

func init() {
	RootCmd.AddCommand(editCmd)
}

func edit(cmd *cobra.Command, args []string) error {
	snippets := db.Find(searchFlag, args)
	snippet, err := findOnlySnippet(snippets)
	if err != nil {
		return err
	}
	return execute.EditFileWithVim(snippet.Path)
}
