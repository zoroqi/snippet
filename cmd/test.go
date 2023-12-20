package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/zoroqi/snippet/store"
	"os"
	"path/filepath"
	"strings"
)

var testCmd = &cobra.Command{
	Use:  "test",
	Long: `Test script`,
	RunE: testF,
}

var testFlag struct {
	file     string
	language string
	text     string
}

func init() {
	testCmd.Flags().StringVarP(&testFlag.file, "file", "f", "", "file path")
	testCmd.Flags().StringVarP(&testFlag.language, "language", "l", "", "language")
	testCmd.Flags().StringVarP(&testFlag.text, "text", "", "", "text")
	RootCmd.AddCommand(testCmd)
}

func testF(cmd *cobra.Command, args []string) error {

	if testFlag.file == "" && testFlag.text == "" {
		return errors.New("no file and text")
	}
	script := testFlag.text
	if testFlag.file != "" {
		if testFlag.language == "" {
			ext := filepath.Ext(testFlag.file)
			testFlag.language = strings.TrimPrefix(ext, ".")
		}
	} else {
		tmp := os.TempDir()
		name := "snippet_test_script"
		if testFlag.language != "" {
			name = name + "." + testFlag.language
		} else {
			name = name + ".sh"
		}
		f, err := os.CreateTemp(tmp, name)
		defer f.Close()
		if err != nil {
			return err
		}
		_, err = f.WriteString(script)
		if err != nil {
			return err
		}
		testFlag.file = f.Name()
	}
	snippet := store.Snippet{
		Name:        "test",
		ShortName:   "test",
		Description: "test",
		Language:    testFlag.language,
		CanExec:     true,
		Path:        testFlag.file,
	}
	return executeSnippet(snippet, args)
}
