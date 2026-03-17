# Error Model

CLI errors should tell the operator:

- what failed
- likely cause
- next action

Examples:

- missing input flags point back to the matching `--help`
- config and environment failures point to `sros verify`
- deferred or unwired paths point back to boundary inspection
