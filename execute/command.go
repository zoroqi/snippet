package execute

import (
	"github.com/zoroqi/snippet/store"
	"golang.design/x/clipboard"
	"io"
	"os"
	"os/exec"
	"runtime"
)

func CopyScriptToClipboard(s store.Snippet) error {
	if err := clipboard.Init(); err != nil {
		return err
	}
	bs, err := os.ReadFile(s.Path)
	if err != nil {
		return err
	}
	clipboard.Write(clipboard.FmtText, bs)
	return err
}

func EditFileWithVim(file string) error {
	command := "vim " + file
	return runCommand(command, os.Stdin, os.Stdout)
}

func runCommand(command string, r io.Reader, w io.Writer) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = w
	cmd.Stdin = r
	return cmd.Run()
}
