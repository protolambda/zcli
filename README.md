# ZRNT CLI

Debugging command line tool, to work with SSZ files, and process ETH 2.0 state transitions.

Based on the Go-spec: [ZRNT](https://github.com/protolambda/zrnt)

## Installation

### Pre-requisites

- Install Go.
- Add `$HOME/go/bin` to your PATH.

### Install

Options:
 
- `-u` to update dependencies (do not use an old ZRNT or ZSSZ dependency in your debugging CLI)
- `-tags bls_off` to disable BLS for testing purposes (not secure!!!)

```bash
GO111MODULE=on go get github.com/protolambda/zcli
```

## Usage

The `help` commands guide you through the usage

```bash
zcli --help
```

Quick overview of all commands (run `zcli <sub command> --help` to get usage options and info).

```text
zcli
  # these commands all have sub-commands to specify the type of the SSZ data.
  diff             find the differences in SSZ data
  pretty           pretty-print SSZ data (indented JSON)
  check            check SSZ data format
  hash-tree-root   (aliases: hash_tree_root, htr) Compute Hash-Tree-Root, output in hex
  # the type sub-commands:
      attestation
      attestation_data
      attester_slashing
      block
      signed_block
      block_body
      block_header
      signed_block_header
      deposit
      deposit_data
      deposit_message
      eth1_data
      proposer_slashing
      state
      state_dump
      voluntary_exit
      signed_voluntary_exit

  api-util    API utilities for eth2 client users.
      extract-state      Extract the state from an api beacon state (wrapper with root).

  net         Util tools for networking
      enr         Decode ENR record.

  info        Information about eth2 data.
      registry    Print a summary of the validator registry. If the input path is not specified, input is read from STDIN.

  genesis     Generate a genesis state
      mock        Generate a genesis state from a predefined set of keys

  meta        Print meta information of a BeaconState
      committees  Print beacon committees for the given state. For prev, current and next epoch.
      proposers   Print beacon proposer indices for the given state. For current epoch.

  keys        Generate and process keys
      generate    Generate a list of keys
      shard       Shard (split) a YAML list of keys into ranges. Specify sizes as arguments.

  transition  Run a state-transition
      blocks      Process blocks on the pre-state to get a post-state
      slots       Process empty slots on the pre-state to get a post-state
      sub         Run a sub state-transition
          block       Run a block sub state-transition
              attestations       process_attestations sub state-transition
              attester_slashings process_attester_slashings sub state-transition
              block_header       process_block_header sub state-transition
              deposits           process_deposits sub state-transition
              proposer_slashings process_proposer_slashings sub state-transition
              voluntary_exits    process_voluntary_exits sub state-transition

          epoch       Run an epoch sub state-transition
              final_updates                  process_final_updates sub state-transition
              justification_and_finalization process_justification_and_finalization sub state-transition
              registry_updates               process_registry_updates sub state-transition
              slashings                      process_slashings sub state-transition

          op          Process a single operation sub state-transition
              attestation       process_attestation sub state-transition
              attester_slashing process_attester_slashing sub state-transition
              deposit           process_deposit sub state-transition
              proposer_slashing process_proposer_slashing sub state-transition
              voluntary_exit    process_voluntary_exit sub state-transition

  version     Print versions of integrated tools

  help        Help about any command
```

All commands have a `--help` for additional information, flags, etc.

And for many commands, use `--spec` to select a known config (`minimal`, `mainnet`, etc.) or a custom YAML config file!
E.g. `--spec=medalla_config.yaml`


## License

MIT, see [`LICENSE`](./LICENSE) file.
