package cmd

import (
	"errors"
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
	snippets := db.Find(searchFlag)
	if len(snippets) != 1 {
		return errors.New("find multiple scripts")
	}
	snippet := snippets[0]
	return execute.CopyScriptToClipboard(snippet)
}
