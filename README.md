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
- `-tags preset_minimal` to compile the minimal spec preset into the CLI

```bash
go get -u -tags preset_minimal github.com/protolambda/zcli
```

## Usage

The `help` commands guide you through the usage

```bash
zcli --help
```

Quick overview of all commands (run `zcli <sub command> --help` to get usage options and info).

```text
zcli
  diff        find the differences in SSZ data
      attestation
      attester_slashing
      block
      block_body
      block_header
      deposit
      deposit_data
      eth1_data
      proposer_slashing
      state
      transfer
      voluntary_exit

  genesis     Generate a genesis state
      mock        Generate a genesis state from a predefined set of keys

  pretty      pretty-print SSZ data
      attestation
      attester_slashing
      block
      block_body
      block_header
      deposit
      deposit_data
      eth1_data
      proposer_slashing
      state
      transfer
      voluntary_exit

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
              transfers          process_transfers sub state-transition
              voluntary_exits    process_voluntary_exits sub state-transition

          epoch       Run an epoch sub state-transition
              crosslinks                     process_crosslinks sub state-transition
              final_updates                  process_final_updates sub state-transition
              justification_and_finalization process_justification_and_finalization sub state-transition
              registry_updates               process_registry_updates sub state-transition
              slashings                      process_slashings sub state-transition

          op          Process a single operation sub state-transition
              attestation       process_attestation sub state-transition
              attester_slashing process_attester_slashing sub state-transition
              deposit           process_deposit sub state-transition
              proposer_slashing process_proposer_slashing sub state-transition
              transfer          process_transfer sub state-transition
              voluntary_exit    process_voluntary_exit sub state-transition

  help        Help about any command
```
