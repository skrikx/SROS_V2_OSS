package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"srosv2/internal/shared/config"
)

type doctorCheck struct {
	Name   string `json:"name"`
	Pass   bool   `json:"pass"`
	Detail string `json:"detail"`
}

func newDoctorCommand() *Command {
	return &Command{
		Name:    "doctor",
		Summary: "Run local environment and repo integrity checks",
		Usage:   "sros doctor",
		Run: func(ctx *Context, args []string) error {
			if err := requireNoArgs(args); err != nil {
				return err
			}

			checks := []doctorCheck{
				checkConfig(ctx.Config),
				checkDirectory("workspace_root", ctx.Config.WorkspaceRoot),
				ensureDirectory("artifact_root", ctx.Config.ArtifactRoot),
				checkDirectory("cmd/sros", filepath.Join(ctx.Config.WorkspaceRoot, "cmd", "sros")),
				checkDirectory("contracts", filepath.Join(ctx.Config.WorkspaceRoot, "contracts")),
				checkDirectory("docs", filepath.Join(ctx.Config.WorkspaceRoot, "docs")),
			}

			failed := 0
			lines := []string{"doctor report:"}
			for _, item := range checks {
				state := "OK"
				if !item.Pass {
					state = "FAIL"
					failed++
				}
				lines = append(lines, fmt.Sprintf("- [%s] %s: %s", state, item.Name, item.Detail))
			}

			payload := map[string]any{"checks": checks, "failed": failed}
			_ = writeOutput(ctx, stringsJoin(lines), payload)
			if failed > 0 {
				return EnvironmentError("doctor detected local environment issues")
			}
			return nil
		},
	}
}

func checkConfig(cfg config.Config) doctorCheck {
	if err := config.Validate(cfg); err != nil {
		return doctorCheck{Name: "config_validation", Pass: false, Detail: err.Error()}
	}
	return doctorCheck{Name: "config_validation", Pass: true, Detail: "configuration is valid"}
}

func checkDirectory(name, path string) doctorCheck {
	info, err := os.Stat(path)
	if err != nil {
		return doctorCheck{Name: name, Pass: false, Detail: err.Error()}
	}
	if !info.IsDir() {
		return doctorCheck{Name: name, Pass: false, Detail: "path is not a directory"}
	}
	return doctorCheck{Name: name, Pass: true, Detail: path}
}

func ensureDirectory(name, path string) doctorCheck {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return doctorCheck{Name: name, Pass: false, Detail: err.Error()}
	}
	return doctorCheck{Name: name, Pass: true, Detail: path}
}
