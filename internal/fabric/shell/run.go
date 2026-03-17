package shell

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func Run(ctx context.Context, name string, args ...string) (Capture, error) {
	if !SafeCommand(name) {
		return Capture{}, fmt.Errorf("unsafe shell command blocked")
	}
	cmd := exec.CommandContext(ctx, name, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	capture := Capture{Stdout: stdout.String(), Stderr: stderr.String()}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			capture.Code = exitErr.ExitCode()
		}
		return capture, err
	}
	return capture, nil
}
