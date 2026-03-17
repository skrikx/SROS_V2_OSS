package commands

import (
	"os"
	"os/exec"
	"path/filepath"
)

func runLocalScript(ctx *Context, rel string) error {
	script := filepath.Join(ctx.CWD, filepath.FromSlash(rel))
	if _, err := os.Stat(script); err != nil {
		return EnvironmentError(err.Error())
	}
	candidates := [][]string{
		{"C:\\Program Files\\Git\\bin\\bash.exe", script},
		{"C:\\Program Files\\Git\\usr\\bin\\bash.exe", script},
		{"bash", script},
	}
	var lastErr error
	for _, candidate := range candidates {
		if candidate[0] != "bash" {
			if _, err := os.Stat(candidate[0]); err != nil {
				continue
			}
		}
		cmd := exec.Command(candidate[0], candidate[1:]...)
		cmd.Dir = ctx.CWD
		cmd.Stdout = ctx.Stdout
		cmd.Stderr = ctx.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}
	if lastErr != nil {
		return EnvironmentError(lastErr.Error())
	}
	return EnvironmentError("no usable bash runtime found")
}
