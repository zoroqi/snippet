package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zoroqi/snippet/execute/anko"
	"github.com/zoroqi/snippet/store"
)

var execCmd = &cobra.Command{
	Use:     "exec",
	Aliases: []string{"e"},
	Short:   "Execute a snippet.",
	Long:    `Execute a snippet.`,
	RunE:    exec,
}

var execFlag struct {
}

func init() {
	RootCmd.AddCommand(execCmd)
}

func exec(cmd *cobra.Command, args []string) error {
	snippets := db.Find(searchFlag)
	snippet, err := findOnlySnippet(snippets)
	if err != nil {
		return err
	}
	if !snippet.CanExec {
		script, err := store.ReadScript(snippet)
		if err != nil {
			return err
		}
		fmt.Print(script)
		return nil
	}
	switch snippet.Language {
	case store.ANKO:
		return anko.Execute(snippet.Path, args)
	default:
		return nil
	}
}
