package commands

func newTestCommand() *Command {
	cmd := &Command{
		Name:    "test",
		Summary: "Local smoke and showcase test surfaces",
		Usage:   "sros test <smoke|first-run>",
		Examples: []string{
			"sros test smoke",
			"sros test first-run",
		},
	}
	cmd.Subcommands = []*Command{
		{
			Name:    "smoke",
			Summary: "Run the local smoke script",
			Usage:   "sros test smoke",
			Run: func(ctx *Context, args []string) error {
				if err := requireNoArgs(args); err != nil {
					return err
				}
				return runLocalScript(ctx, "scripts/test_smoke.sh")
			},
		},
		{
			Name:    "first-run",
			Summary: "Run the first-run smoke path used by the front door",
			Usage:   "sros test first-run",
			Run: func(ctx *Context, args []string) error {
				if err := requireNoArgs(args); err != nil {
					return err
				}
				return runLocalScript(ctx, "scripts/first_run_smoke.sh")
			},
		},
	}
	return cmd
}
