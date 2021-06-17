# ZRNT CLI

Debugging command line tool, to work with SSZ files, process ETH 2.0 state transitions, and compute proofs and meta-data.

Based on the Go-spec: [ZRNT](https://github.com/protolambda/zrnt)

## Installation

### Pre-requisites

- Install Go 1.16+
- Add `$HOME/go/bin` to your PATH.

### Install

Options:
 
- `-u` to force-update dependencies
- `-tags bls_off` to disable BLS for testing purposes (not secure!!!)

```bash
# outside of an existing go module directory
go install github.com/protolambda/zcli@latest
```

## Usage

The `help` commands guide you through the usage

```bash
zcli --help
```

Quick overview of all commands (run `zcli <sub command> --help` to get usage options and info).

```text
zcli
  pretty <phase> <type> <input>                  Pretty-print spec object (output indented JSON)
  convert <phase> <type> <input> <output>        Convert spec object from one format to another
  diff <phase> <type> <a> <b>                    Diff spec data
  meta <phase> <subcmd>                          List metadata of beacon state
  proof <phase> <type> <input> --gindices        Create SSZ merkle proofs over any spec object
  root <phase> <type> <input>                    Compute the SSZ hash-tree-root of a spec object
  transition <pre-phase> <slots/blocks/sub>      Run state transitions and sub-processes
  tree <phase> <type>                            Dump SSZ merkle tree of any spec object
  version                                        Print ZCLI and ZRNT version
```

All commands have a `--help` for additional information, flags, etc.

And for many commands, use `--config` and `--preset-{forkname}` to select a known (`minimal`, `mainnet`, etc.) or custom YAML config/preset file!
E.g. `--config=local_testnet.yaml`

Inputs/outputs can:
- be specified as empty `""`, to read from STDIN/STDOUT
- be specified with a prefix `json:`, `yaml:` or `ssz:` to read/write that format. Writing can also use `pretty:` (indented JSON).

## License

MIT, see [`LICENSE`](./LICENSE) file.
