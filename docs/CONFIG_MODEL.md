# Config Model

## Search Order

1. `--config <path>`
2. `SROS_CONFIG`
3. `<workspace>/sros.yaml`
4. `<workspace>/.sros/config.yaml`
5. defaults

## Supported Keys

- `mode`
- `workspace_root`
- `artifact_root`
- `policy_bundle_path`
- `memory_store_path`
- `trace_store_path`
- `output_format`

## Environment Overrides

- `SROS_MODE`
- `SROS_WORKSPACE_ROOT`
- `SROS_ARTIFACT_ROOT`
- `SROS_POLICY_BUNDLE_PATH`
- `SROS_MEMORY_STORE_PATH`
- `SROS_TRACE_STORE_PATH`
- `SROS_OUTPUT_FORMAT`

Only local CLI mode (`local_cli`) is valid in v2.
