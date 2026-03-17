# Compile Edge

Operator command:

- `sros compile --intent "..."`
- `sros compile --input examples/intents/minimal.txt`

Optional flags:

- `--emit-srxml`
- `--format text|json`
- `--receipt-path <path>`
- `--operator`, `--tenant`, `--workspace`

Compile output includes normalized intent, classification, topology, canonical run contract, and compile receipt.
