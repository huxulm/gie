package internal

import (
	"os"
	"os/exec"
)

func GitExec(args ...string) error {
	c := exec.Command("git", args...)
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	return c.Run()
}
