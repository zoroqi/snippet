package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zoroqi/snippet/execute"
)

var copyCmd = &cobra.Command{
	Use:     "copy",
	Aliases: []string{"cp"},
	Short:   "Copy a snippet",
	Long:    `Copy a snippet`,
	RunE:    copyFunc,
}

func init() {
	RootCmd.AddCommand(copyCmd)
}

func copyFunc(cmd *cobra.Command, args []string) error {
	snippets := db.Find(searchFlag, args)
	snippet, err := findOnlySnippet(snippets)
	if err != nil {
		return err
	}
	return execute.CopyScriptToClipboard(snippet)
}
