package transition

import (
	"fmt"
	. "github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2/beacon/attestations"
	"github.com/protolambda/zrnt/eth2/beacon/deposits"
	"github.com/protolambda/zrnt/eth2/beacon/exits"
	"github.com/protolambda/zrnt/eth2/beacon/header"
	"github.com/protolambda/zrnt/eth2/beacon/slashings/attslash"
	"github.com/protolambda/zrnt/eth2/beacon/slashings/propslash"
	"github.com/protolambda/zrnt/eth2/beacon/transfers"
	"github.com/protolambda/zrnt/eth2/core"
	"github.com/protolambda/zrnt/eth2/phase0"
	"github.com/spf13/cobra"
	"strconv"
)

var (
	TransitionCmd *cobra.Command
	BlocksCmd     *cobra.Command
	SlotsCmd      *cobra.Command
	SubCmd        *cobra.Command
)

var (
	EpochCmd                        *cobra.Command
	CrosslinksCmd                   *cobra.Command
	FinalUpdatesCmd                 *cobra.Command
	JustificationAndFinalizationCmd *cobra.Command
	RegistryUpdatesCmd              *cobra.Command
	SlashingsCmd                    *cobra.Command
)

var (
	OpCmd               *cobra.Command
	AttestationCmd      *cobra.Command
	AttesterSlashingCmd *cobra.Command
	ProposerSlashingCmd *cobra.Command
	DepositCmd          *cobra.Command
	TransferCmd         *cobra.Command
	VoluntaryExitCmd    *cobra.Command
)

var (
	BlockCmd             *cobra.Command
	BlockHeaderCmd       *cobra.Command
	AttestationsCmd      *cobra.Command
	AttesterSlashingsCmd *cobra.Command
	ProposerSlashingsCmd *cobra.Command
	DepositsCmd          *cobra.Command
	TransfersCmd         *cobra.Command
	VoluntaryExitsCmd    *cobra.Command
)

func init() {
	TransitionCmd = &cobra.Command{
		Use:   "transition",
		Short: "Run a state-transition",
	}
	TransitionCmd.PersistentFlags().StringP("pre", "i", "", "Pre (Input) path. If none is specified, input is read from STDIN")
	TransitionCmd.PersistentFlags().StringP("post", "o", "", "Post (Output) path. If none is specified, output is written to STDOUT")

	SlotsCmd = &cobra.Command{
		Use:   "slots <number>",
		Short: "Process empty slots on the pre-state to get a post-state",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("expected one argument: <number>")
			}
			_, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("argument %v is a not a valid number", args[0])
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			isDelta, err := cmd.Flags().GetBool("delta")
			if Check(err, cmd.ErrOrStderr(), "delta flag could not be parsed") {
				return
			}

			slots, _ := strconv.ParseUint(args[0], 10, 64)

			state, err := loadPreFull(cmd)
			if Check(err, cmd.ErrOrStderr(), "pre state could not be loaded") {
				return
			}

			to := core.Slot(slots)
			if isDelta {
				to += state.Slot
			} else if to < state.Slot {
				Report(cmd.ErrOrStderr(), "to slot is lower than pre-state slot")
				return
			}

			state.ProcessSlots(to)
			err = writePost(cmd, state.BeaconState)
			if Check(err, cmd.ErrOrStderr(), "could not write post-state") {
				return
			}
		},
	}
	SlotsCmd.Flags().Bool("delta", false, "to interpret the slot number as a delta from the pre-state")
	TransitionCmd.AddCommand(SlotsCmd)

	BlocksCmd = &cobra.Command{
		Use:   "blocks [<block 0.ssz> [<block 1.ssz> [<block 2.ssz> [ ... ]]]",
		Short: "Process blocks on the pre-state to get a post-state",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			verifyStateRoot, err := cmd.Flags().GetBool("verify-state-root")
			if Check(err, cmd.ErrOrStderr(), "verify-state-root could not be parsed") {
				return
			}

			state, err := loadPreFull(cmd)
			if Check(err, cmd.ErrOrStderr(), "could not load pre-state") {
				return
			}

			for i := 0; i < len(args); i++ {
				var b phase0.BeaconBlock
				err := LoadSSZ(args[i], &b, phase0.BeaconBlockSSZ)
				if Check(err, cmd.ErrOrStderr(), "could not load block: %s", args[i]) {
					return
				}

				blockProc := &phase0.BlockProcessFeature{Block: &b, Meta: state}

				err = state.StateTransition(blockProc, verifyStateRoot)
				if Check(err, cmd.ErrOrStderr(), "failed block transition to block %s", args[i]) {
					// still output the state, just stop processing blocks
					break
				}
			}

			err = writePost(cmd, state.BeaconState)
			if Check(err, cmd.ErrOrStderr(), "could not write post-state") {
				return
			}
		},
	}

	BlocksCmd.Flags().Bool("verify-state-root", true, "change the state-root verification step")

	TransitionCmd.AddCommand(BlocksCmd)

	SubCmd = &cobra.Command{
		Use:   "sub",
		Short: "Run a sub state-transition",
	}
	TransitionCmd.AddCommand(SubCmd)

	EpochCmd = &cobra.Command{
		Use:   "epoch",
		Short: "Run an epoch sub state-transition",
	}
	OpCmd = &cobra.Command{
		Use:   "op",
		Short: "Process a single operation sub state-transition",
	}
	BlockCmd = &cobra.Command{
		Use:   "block",
		Short: "Run a block sub state-transition",
	}
	SubCmd.AddCommand(EpochCmd, OpCmd, BlockCmd)

	transition := func(cmd *cobra.Command, change func(state *phase0.FullFeaturedState)) {
		state, err := loadPreFull(cmd)
		if Check(err, cmd.ErrOrStderr(), "pre state could not be loaded") {
			return
		}
		change(state)
		err = writePost(cmd, state.BeaconState)
		if Check(err, cmd.ErrOrStderr(), "could not write post-state") {
			return
		}
	}
	CrosslinksCmd = &cobra.Command{
		Use:   "crosslinks",
		Short: "process_crosslinks sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				state.ProcessEpochCrosslinks()
			})
		},
	}
	FinalUpdatesCmd = &cobra.Command{
		Use:   "final_updates",
		Short: "process_final_updates sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				state.ProcessEpochFinalUpdates()
			})
		},
	}
	JustificationAndFinalizationCmd = &cobra.Command{
		Use:   "justification_and_finalization",
		Short: "process_justification_and_finalization sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				state.ProcessEpochJustification()
			})
		},
	}
	RegistryUpdatesCmd = &cobra.Command{
		Use:   "registry_updates",
		Short: "process_registry_updates sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				state.ProcessEpochRegistryUpdates()
			})
		},
	}
	SlashingsCmd = &cobra.Command{
		Use:   "slashings",
		Short: "process_slashings sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				state.ProcessEpochSlashings()
			})
		},
	}
	EpochCmd.AddCommand(CrosslinksCmd, FinalUpdatesCmd, JustificationAndFinalizationCmd, RegistryUpdatesCmd, SlashingsCmd)

	AttestationCmd = &cobra.Command{
		Use:   "attestation <data.ssz>",
		Short: "process_attestation sub state-transition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				var op attestations.Attestation
				err := LoadSSZ(args[0], &op, attestations.AttestationSSZ)
				if Check(err, cmd.ErrOrStderr(), "could not load attestation") {
					return
				}
				err = state.ProcessAttestation(&op)
				if Check(err, cmd.ErrOrStderr(), "failed to process attestation") {
					return
				}
			})
		},
	}
	AttesterSlashingCmd = &cobra.Command{
		Use:   "attester_slashing <data.ssz>",
		Short: "process_attester_slashing sub state-transition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				var op attslash.AttesterSlashing
				err := LoadSSZ(args[0], &op, attslash.AttesterSlashingSSZ)
				if Check(err, cmd.ErrOrStderr(), "could not load attester slashing") {
					return
				}
				err = state.ProcessAttesterSlashing(&op)
				if Check(err, cmd.ErrOrStderr(), "failed to process attester slashing") {
					return
				}
			})
		},
	}
	ProposerSlashingCmd = &cobra.Command{
		Use:   "proposer_slashing <data.ssz>",
		Short: "process_proposer_slashing sub state-transition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				var op propslash.ProposerSlashing
				err := LoadSSZ(args[0], &op, propslash.ProposerSlashingSSZ)
				if Check(err, cmd.ErrOrStderr(), "could not load proposer slashing") {
					return
				}
				err = state.ProcessProposerSlashing(&op)
				if Check(err, cmd.ErrOrStderr(), "failed to process proposer slashing") {
					return
				}
			})
		},
	}
	DepositCmd = &cobra.Command{
		Use:   "deposit <data.ssz>",
		Short: "process_deposit sub state-transition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				var op deposits.Deposit
				err := LoadSSZ(args[0], &op, deposits.DepositSSZ)
				if Check(err, cmd.ErrOrStderr(), "could not load deposit") {
					return
				}
				err = state.ProcessDeposit(&op)
				if Check(err, cmd.ErrOrStderr(), "failed to process deposit") {
					return
				}
			})
		},
	}
	TransferCmd = &cobra.Command{
		Use:   "transfer <data.ssz>",
		Short: "process_transfer sub state-transition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				var op transfers.Transfer
				err := LoadSSZ(args[0], &op, transfers.TransferSSZ)
				if Check(err, cmd.ErrOrStderr(), "could not load transfer") {
					return
				}
				err = state.ProcessTransfer(&op)
				if Check(err, cmd.ErrOrStderr(), "failed to process transfer") {
					return
				}
			})
		},
	}
	VoluntaryExitCmd = &cobra.Command{
		Use:   "voluntary_exit <data.ssz>",
		Short: "process_voluntary_exit sub state-transition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				var op exits.VoluntaryExit
				err := LoadSSZ(args[0], &op, exits.VoluntaryExitSSZ)
				if Check(err, cmd.ErrOrStderr(), "could not load voluntary exit") {
					return
				}
				err = state.ProcessVoluntaryExit(&op)
				if Check(err, cmd.ErrOrStderr(), "failed to process voluntary exit") {
					return
				}
			})
		},
	}
	OpCmd.AddCommand(AttestationCmd, AttesterSlashingCmd, ProposerSlashingCmd, DepositCmd, TransferCmd, VoluntaryExitCmd)

	BlockHeaderCmd = &cobra.Command{
		Use:   "block_header <data.ssz>",
		Short: "process_block_header sub state-transition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				var bh header.BeaconBlockHeader
				err := LoadSSZ(args[0], &bh, header.BeaconBlockHeaderSSZ)
				if Check(err, cmd.ErrOrStderr(), "could not load block header") {
					return
				}
				err = state.ProcessHeader(&bh)
				if Check(err, cmd.ErrOrStderr(), "failed to process block header") {
					return
				}
			})
		},
	}
	AttestationsCmd = &cobra.Command{
		Use:   "attestations [<data 0.ssz> [<data 1.ssz> [<data 2.ssz> [ ... ]]]]",
		Short: "process_attestations sub state-transition",
		Args:  cobra.RangeArgs(0, core.MAX_ATTESTATIONS),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				if uint64(len(args)) > ((*phase0.Attestations)(nil)).Limit() {
					Report(cmd.ErrOrStderr(), "too many attestations")
					return
				}
				ops := make(phase0.Attestations, len(args), len(args))
				for i, arg := range args {
					err := LoadSSZ(arg, &ops[i], attestations.AttestationSSZ)
					if Check(err, cmd.ErrOrStderr(), "could not load attestation %d %s", i, arg) {
						return
					}
				}
				err := state.ProcessAttestations(ops)
				if Check(err, cmd.ErrOrStderr(), "failed to process attestations") {
					return
				}
			})
		},
	}
	AttesterSlashingsCmd = &cobra.Command{
		Use:   "attester_slashings [<data 0.ssz> [<data 1.ssz> [<data 2.ssz> [ ... ]]]]",
		Short: "process_attester_slashings sub state-transition",
		Args:  cobra.RangeArgs(0, core.MAX_ATTESTER_SLASHINGS),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				if uint64(len(args)) > ((*phase0.AttesterSlashings)(nil)).Limit() {
					Report(cmd.ErrOrStderr(), "too many attester slashings")
					return
				}
				ops := make(phase0.AttesterSlashings, len(args), len(args))
				for i, arg := range args {
					err := LoadSSZ(arg, &ops[i], attslash.AttesterSlashingSSZ)
					if Check(err, cmd.ErrOrStderr(), "could not load attester slashing %d %s", i, arg) {
						return
					}
				}
				err := state.ProcessAttesterSlashings(ops)
				if Check(err, cmd.ErrOrStderr(), "failed to process attester slashings") {
					return
				}
			})
		},
	}
	ProposerSlashingsCmd = &cobra.Command{
		Use:   "proposer_slashings [<data 0.ssz> [<data 1.ssz> [<data 2.ssz> [ ... ]]]]",
		Short: "process_proposer_slashings sub state-transition",
		Args:  cobra.RangeArgs(0, core.MAX_PROPOSER_SLASHINGS),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				if uint64(len(args)) > ((*phase0.ProposerSlashings)(nil)).Limit() {
					Report(cmd.ErrOrStderr(), "too many proposer slashings")
					return
				}
				ops := make(phase0.ProposerSlashings, len(args), len(args))
				for i, arg := range args {
					err := LoadSSZ(arg, &ops[i], propslash.ProposerSlashingSSZ)
					if Check(err, cmd.ErrOrStderr(), "could not load proposer slashing %d %s", i, arg) {
						return
					}
				}
				err := state.ProcessProposerSlashings(ops)
				if Check(err, cmd.ErrOrStderr(), "failed to process proposer slashings") {
					return
				}
			})
		},
	}
	DepositsCmd = &cobra.Command{
		Use:   "deposits [<data 0.ssz> [<data 1.ssz> [<data 2.ssz> [ ... ]]]]",
		Short: "process_deposits sub state-transition",
		Args:  cobra.RangeArgs(0, core.MAX_DEPOSITS),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				if uint64(len(args)) > ((*phase0.Deposits)(nil)).Limit() {
					Report(cmd.ErrOrStderr(), "too many deposits")
					return
				}
				ops := make(phase0.Deposits, len(args), len(args))
				for i, arg := range args {
					err := LoadSSZ(arg, &ops[i], deposits.DepositSSZ)
					if Check(err, cmd.ErrOrStderr(), "could not load deposit %d %s", i, arg) {
						return
					}
				}
				err := state.ProcessDeposits(ops)
				if Check(err, cmd.ErrOrStderr(), "failed to process deposits") {
					return
				}
			})
		},
	}
	TransfersCmd = &cobra.Command{
		Use:   "transfers [<data 0.ssz> [<data 1.ssz> [<data 2.ssz> [ ... ]]]]",
		Short: "process_transfers sub state-transition",
		Args:  cobra.RangeArgs(0, core.MAX_TRANSFERS),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				if uint64(len(args)) > ((*phase0.Transfers)(nil)).Limit() {
					Report(cmd.ErrOrStderr(), "too many transfers")
					return
				}
				ops := make(phase0.Transfers, len(args), len(args))
				for i, arg := range args {
					err := LoadSSZ(arg, &ops[i], transfers.TransferSSZ)
					if Check(err, cmd.ErrOrStderr(), "could not load transfer %d %s", i, arg) {
						return
					}
				}
				err := state.ProcessTransfers(ops)
				if Check(err, cmd.ErrOrStderr(), "failed to process transfers") {
					return
				}
			})
		},
	}
	VoluntaryExitsCmd = &cobra.Command{
		Use:   "voluntary_exits [<data 0.ssz> [<data 1.ssz> [<data 2.ssz> [ ... ]]]]",
		Short: "process_voluntary_exits sub state-transition",
		Args:  cobra.RangeArgs(0, core.MAX_VOLUNTARY_EXITS),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				if uint64(len(args)) > ((*phase0.VoluntaryExits)(nil)).Limit() {
					Report(cmd.ErrOrStderr(), "too many voluntary exits")
					return
				}
				ops := make(phase0.VoluntaryExits, len(args), len(args))
				for i, arg := range args {
					err := LoadSSZ(arg, &ops[i], exits.VoluntaryExitSSZ)
					if Check(err, cmd.ErrOrStderr(), "could not load voluntary exit %d %s", i, arg) {
						return
					}
				}
				err := state.ProcessVoluntaryExits(ops)
				if Check(err, cmd.ErrOrStderr(), "failed to process voluntary exits") {
					return
				}
			})
		},
	}
	BlockCmd.AddCommand(BlockHeaderCmd, AttestationsCmd,
		AttesterSlashingsCmd, ProposerSlashingsCmd,
		DepositsCmd, TransfersCmd, VoluntaryExitsCmd)
}

func loadPreFull(cmd *cobra.Command) (*phase0.FullFeaturedState, error) {
	pre, err := LoadStateInputFlag(cmd, "pre", true)
	if err != nil {
		return nil, err
	}
	preFull := phase0.NewFullFeaturedState(pre)
	preFull.LoadPrecomputedData()

	return preFull, nil
}

func writePost(cmd *cobra.Command, state *phase0.BeaconState) error {
	return WriteStateOutput(cmd, "post", state)
}
