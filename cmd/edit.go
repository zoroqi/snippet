package cmd

import (
	"errors"
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

func edit(cmd *cobra.Command, args []string) (err error) {
	snippets := db.Find(searchFlag)
	if len(snippets) != 1 {
		return errors.New("find multiple scripts")
	}
	return execute.EditFileWithVim(snippets[0].Path)
}
