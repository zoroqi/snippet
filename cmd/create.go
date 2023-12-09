package cmd

import (
	"errors"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"github.com/zoroqi/snippet/store"
	"golang.design/x/clipboard"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create a snippet",
	Long:    `Create a snippet`,
	RunE:    create,
}

var createFlag struct {
	clip       bool
	filePath   string
	targetPath string
}

func init() {
	createCmd.Flags().BoolVarP(&createFlag.clip, "clip", "c",
		false, "create a snippet from clipboard")
	createCmd.Flags().StringVarP(&createFlag.filePath, "file", "f",
		"", "create a snippet from file")
	createCmd.Flags().StringVarP(&createFlag.targetPath, "target", "t",
		"", "create a snippet, and save it to directory")
	RootCmd.AddCommand(createCmd)
}

func id[T any](t T) T {
	return t
}

func create(cmd *cobra.Command, args []string) error {

	if !createFlag.clip && createFlag.filePath == "" {
		return errors.New("no snippet, use clip or file flag")
	}

	snippet := store.Snippet{}

	var err error
	if snippet.Name, err = readLine("Name> ", false, id[string]); err != nil {
		return err
	}
	snippet.ShortName = strings.ReplaceAll(snippet.Name, " ", "_")
	if snippet.Description, err = readLine("Description> ", false, id[string]); err != nil {
		return err
	}
	if alias, err := readLine("Aliases> ", true, id[string]); err != nil {
		return err
	} else {
		snippet.Aliases = strings.Split(strings.ReplaceAll(alias, " ", ""), ",")
	}
	if tags, err := readLine("Tags> ", true, id[string]); err != nil {
		return err
	} else {
		snippet.Tags = strings.Split(strings.ReplaceAll(tags, " ", ""), ",")
	}

	if snippet.CanExec, err = readLine[bool]("CanExec> ", false, func(s string) bool {
		if s == "t" || strings.ToLower(s) == "true" {
			return true
		}
		return false
	}); err != nil {
		return err
	}

	targetPath := filepath.Join(rootPath, createFlag.targetPath)
	var script []byte
	if createFlag.clip {
		path := filepath.Join(targetPath, snippet.ShortName+".sh")
		err = clipboard.Init()
		if err != nil {
			return err
		}
		snippet.Path = path
		snippet.Language = store.SH
		script = clipboard.Read(clipboard.FmtText)
	} else {
		path := filepath.Join(targetPath, fmt.Sprintf("%s%s", snippet.ShortName, filepath.Ext(createFlag.filePath)))
		script, err = os.ReadFile(createFlag.filePath)
		if err != nil {
			return err
		}
		snippet.Path = path
		snippet.Language = filepath.Ext(path)
		if snippet.Language == "" {
			snippet.Language = store.SH
		}
		snippet.Language = strings.TrimPrefix(snippet.Language, ".")
	}
	return db.Add(snippet, script)
}

func readLine[T any](message string, canBeNull bool, mapper func(string) T) (T, error) {
	l, err := readline.NewEx(&readline.Config{
		Prompt:            message,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	})
	var t T
	if err != nil {
		return t, err
	}
	defer l.Close()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" && canBeNull {
			return mapper(""), nil
		}
		return mapper(line), nil
	}
	return t, errors.New("canceled")
}
