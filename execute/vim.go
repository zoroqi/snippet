package execute

import (
	"io"
	"os"
	"os/exec"
	"runtime"
)

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
