package cmd

import (
	"errors"
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
	alias string
}

func init() {
	RootCmd.AddCommand(execCmd)
}

func exec(cmd *cobra.Command, args []string) error {
	snippets := db.Find(searchFlag)
	if len(snippets) != 1 {
		return errors.New("find multiple scripts")
	}
	snippet := snippets[0]
	if !snippet.CanExec {
		return errors.New("execute can't execute")
	}
	switch snippet.Language {
	case store.ANKO:
		return anko.Execute(snippet.Path, args)
	default:
		return nil
	}
}
