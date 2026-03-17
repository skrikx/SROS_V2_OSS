# Local Bootstrap Flow

W03 bootstrap is in-process and local-only.

1. Parse global flags (`--config`, `--format`, `--workspace`).
2. Load config from explicit path, default path, or defaults.
3. Apply environment overrides.
4. Validate resolved config.
5. Build a local bootstrap bundle with service boundaries.
6. Dispatch command handler.

The bundle exposes compile/runtime/inspection/fabric boundaries without implementing W04-W10 behavior.
