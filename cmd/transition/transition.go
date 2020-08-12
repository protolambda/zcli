package transition

import (
	"context"
	"fmt"
	. "github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2/beacon"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

var (
	TransitionCmd *cobra.Command
	BlocksCmd     *cobra.Command
	SlotsCmd      *cobra.Command
	SubCmd        *cobra.Command
)

var (
	EpochCmd                        *cobra.Command
	FinalUpdatesCmd                 *cobra.Command
	JustificationAndFinalizationCmd *cobra.Command
	RewardsAndPenalties             *cobra.Command
	RegistryUpdatesCmd              *cobra.Command
	SlashingsCmd                    *cobra.Command
)

var (
	OpCmd               *cobra.Command
	AttestationCmd      *cobra.Command
	AttesterSlashingCmd *cobra.Command
	ProposerSlashingCmd *cobra.Command
	DepositCmd          *cobra.Command
	VoluntaryExitCmd    *cobra.Command
)

var (
	BlockCmd             *cobra.Command
	BlockHeaderCmd       *cobra.Command
	AttestationsCmd      *cobra.Command
	AttesterSlashingsCmd *cobra.Command
	ProposerSlashingsCmd *cobra.Command
	DepositsCmd          *cobra.Command
	VoluntaryExitsCmd    *cobra.Command
)

func init() {
	TransitionCmd = &cobra.Command{
		Use:   "transition",
		Short: "Run a state-transition",
	}
	TransitionCmd.PersistentFlags().StringP("pre", "i", "", "Pre (Input) path. If none is specified, input is read from STDIN")
	TransitionCmd.PersistentFlags().StringP("post", "o", "", "Post (Output) path. If none is specified, output is written to STDOUT")

	withTimeout := func(cmd *cobra.Command) *cobra.Command {
		cmd.Flags().Uint64("timeout", 0, "timeout in milliseconds, 0 to disable timeout")
		return cmd
	}
	SlotsCmd = withTimeout(&cobra.Command{
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
			timeout, err := cmd.Flags().GetUint64("timeout")
			if Check(err, cmd.ErrOrStderr(), "timeout flag could not be parsed") {
				return
			}

			slots, _ := strconv.ParseUint(args[0], 10, 64)

			state, epc, err := loadPreFull(cmd)
			if Check(err, cmd.ErrOrStderr(), "pre state could not be loaded") {
				return
			}

			slot, err := state.Slot()
			if Check(err, cmd.ErrOrStderr(), "could not read slot") {
				return
			}

			to := beacon.Slot(slots)
			if isDelta {
				to += slot
			} else if to <= slot {
				Report(cmd.ErrOrStderr(), "to slot is lower or equal to pre-state slot")
				return
			}
			ctx := context.Background()
			if timeout != 0 {
				ctx, _ = context.WithTimeout(ctx, time.Duration(timeout)*time.Microsecond)
			}
			err = state.ProcessSlots(ctx, epc, to)
			if Check(err, cmd.ErrOrStderr(), "Failed transition, could not compute post-state") {
				return
			}

			err = writePost(cmd, state)
			if Check(err, cmd.ErrOrStderr(), "could not write post-state") {
				return
			}
		},
	})
	SlotsCmd.Flags().Bool("delta", false, "to interpret the slot number as a delta from the pre-state")
	TransitionCmd.AddCommand(SlotsCmd)

	BlocksCmd = withTimeout(&cobra.Command{
		Use:   "blocks [<block 0.ssz> [<block 1.ssz> [<block 2.ssz> [ ... ]]]",
		Short: "Process signed blocks on the pre-state to get a post-state",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			verifyStateRoot, err := cmd.Flags().GetBool("verify-state-root")
			if Check(err, cmd.ErrOrStderr(), "verify-state-root could not be parsed") {
				return
			}

			timeout, err := cmd.Flags().GetUint64("timeout")
			if Check(err, cmd.ErrOrStderr(), "timeout flag could not be parsed") {
				return
			}

			state, epc, err := loadPreFull(cmd)
			if Check(err, cmd.ErrOrStderr(), "could not load pre-state") {
				return
			}

			ctx := context.Background()
			if timeout != 0 {
				ctx, _ = context.WithTimeout(ctx, time.Duration(timeout)*time.Microsecond)
			}

			for i := 0; i < len(args); i++ {
				var b beacon.SignedBeaconBlock
				err := LoadSSZ(args[i], &b, beacon.SignedBeaconBlockSSZ)
				if Check(err, cmd.ErrOrStderr(), "could not load block: %s", args[i]) {
					return
				}

				err = state.StateTransition(ctx, epc, &b, verifyStateRoot)
				if Check(err, cmd.ErrOrStderr(), "failed block transition to block %s", args[i]) {
					// still output the state, just stop processing blocks
					break
				}
			}

			err = writePost(cmd, state)
			if Check(err, cmd.ErrOrStderr(), "could not write post-state") {
				return
			}
		},
	})
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

	transition := func(cmd *cobra.Command, change func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error) {
		ctx := context.Background()
		timeout, err := cmd.Flags().GetUint64("timeout")
		if Check(err, cmd.ErrOrStderr(), "timeout flag could not be parsed") {
			return
		}
		if timeout != 0 {
			ctx, _ = context.WithTimeout(ctx, time.Duration(timeout)*time.Microsecond)
		}
		state, epc, err := loadPreFull(cmd)
		if Check(err, cmd.ErrOrStderr(), "pre state could not be loaded") {
			return
		}
		err = change(ctx, state, epc)
		// Check and report the error, but still write the post state (even if an incomplete transition), it may be useful for debugging.
		_ = Check(err, cmd.ErrOrStderr(), "failed transition, could not compute post-state")
		err = writePost(cmd, state)
		if Check(err, cmd.ErrOrStderr(), "could not write post-state") {
			return
		}
	}
	FinalUpdatesCmd = withTimeout(&cobra.Command{
		Use:   "final_updates",
		Short: "process_final_updates sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				process, err := state.PrepareEpochProcess(ctx, epc)
				if err != nil {
					return err
				}
				return state.ProcessEpochFinalUpdates(ctx, epc, process)
			})
		},
	})
	JustificationAndFinalizationCmd = withTimeout(&cobra.Command{
		Use:   "justification_and_finalization",
		Short: "process_justification_and_finalization sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				process, err := state.PrepareEpochProcess(ctx, epc)
				if err != nil {
					return err
				}
				return state.ProcessEpochJustification(ctx, epc, process)
			})
		},
	})
	RewardsAndPenalties = withTimeout(&cobra.Command{
		Use:   "rewards_and_penalties",
		Short: "process_rewards_and_penalties sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				process, err := state.PrepareEpochProcess(ctx, epc)
				if err != nil {
					return err
				}
				return state.ProcessEpochRewardsAndPenalties(ctx, epc, process)
			})
		},
	})
	RegistryUpdatesCmd = withTimeout(&cobra.Command{
		Use:   "registry_updates",
		Short: "process_registry_updates sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				process, err := state.PrepareEpochProcess(ctx, epc)
				if err != nil {
					return err
				}
				return state.ProcessEpochRegistryUpdates(ctx, epc, process)
			})
		},
	})
	SlashingsCmd = withTimeout(&cobra.Command{
		Use:   "slashings",
		Short: "process_slashings sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				process, err := state.PrepareEpochProcess(ctx, epc)
				if err != nil {
					return err
				}
				return state.ProcessEpochSlashings(ctx, epc, process)
			})
		},
	})
	EpochCmd.AddCommand(FinalUpdatesCmd, JustificationAndFinalizationCmd,
		RewardsAndPenalties, RegistryUpdatesCmd, SlashingsCmd)

	AttestationCmd = withTimeout(&cobra.Command{
		Use:   "attestation <data.ssz>",
		Short: "process_attestation sub state-transition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				var op beacon.Attestation
				err := LoadSSZ(args[0], &op, beacon.AttestationSSZ)
				if err != nil {
					return err
				}
				return state.ProcessAttestation(epc, &op)
			})
		},
	})
	AttesterSlashingCmd = withTimeout(&cobra.Command{
		Use:   "attester_slashing <data.ssz>",
		Short: "process_attester_slashing sub state-transition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				var op beacon.AttesterSlashing
				err := LoadSSZ(args[0], &op, beacon.AttesterSlashingSSZ)
				if err != nil {
					return err
				}
				return state.ProcessAttesterSlashing(epc, &op)
			})
		},
	})
	ProposerSlashingCmd = withTimeout(&cobra.Command{
		Use:   "proposer_slashing <data.ssz>",
		Short: "process_proposer_slashing sub state-transition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				var op beacon.ProposerSlashing
				err := LoadSSZ(args[0], &op, beacon.ProposerSlashingSSZ)
				if err != nil {
					return err
				}
				return state.ProcessProposerSlashing(epc, &op)
			})
		},
	})
	DepositCmd = withTimeout(&cobra.Command{
		Use:   "deposit <data.ssz>",
		Short: "process_deposit sub state-transition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				var op beacon.Deposit
				err := LoadSSZ(args[0], &op, beacon.DepositSSZ)
				if err != nil {
					return err
				}
				return state.ProcessDeposit(epc, &op, false)
			})
		},
	})
	VoluntaryExitCmd = withTimeout(&cobra.Command{
		Use:   "voluntary_exit <data.ssz>",
		Short: "process_voluntary_exit sub state-transition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				var op beacon.SignedVoluntaryExit
				err := LoadSSZ(args[0], &op, beacon.VoluntaryExitSSZ)
				if err != nil {
					return err
				}
				return state.ProcessVoluntaryExit(epc, &op)
			})
		},
	})
	OpCmd.AddCommand(AttestationCmd, AttesterSlashingCmd, ProposerSlashingCmd, DepositCmd, VoluntaryExitCmd)

	BlockHeaderCmd = withTimeout(&cobra.Command{
		Use:   "block_header <data.ssz>",
		Short: "process_block_header sub state-transition (block input, not header)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				var bh beacon.BeaconBlock
				err := LoadSSZ(args[0], &bh, beacon.BeaconBlockSSZ)
				if err != nil {
					return err
				}
				return state.ProcessHeader(ctx, epc, &bh)
			})
		},
	})
	AttestationsCmd = withTimeout(&cobra.Command{
		Use:   "attestations [<data 0.ssz> [<data 1.ssz> [<data 2.ssz> [ ... ]]]]",
		Short: "process_attestations sub state-transition",
		Args:  cobra.RangeArgs(0, beacon.MAX_ATTESTATIONS),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				ops := make(beacon.Attestations, len(args), len(args))
				for i, arg := range args {
					err := LoadSSZ(arg, &ops[i], beacon.AttestationSSZ)
					if err != nil {
						return fmt.Errorf("could not load operation %d %s: %v", i, arg, err)
					}
				}
				return state.ProcessAttestations(ctx, epc, ops)
			})
		},
	})
	AttesterSlashingsCmd = withTimeout(&cobra.Command{
		Use:   "attester_slashings [<data 0.ssz> [<data 1.ssz> [<data 2.ssz> [ ... ]]]]",
		Short: "process_attester_slashings sub state-transition",
		Args:  cobra.RangeArgs(0, beacon.MAX_ATTESTER_SLASHINGS),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				ops := make(beacon.AttesterSlashings, len(args), len(args))
				for i, arg := range args {
					err := LoadSSZ(arg, &ops[i], beacon.AttesterSlashingSSZ)
					if err != nil {
						return fmt.Errorf("could not load operation %d %s: %v", i, arg, err)
					}
				}
				return state.ProcessAttesterSlashings(ctx, epc, ops)
			})
		},
	})
	ProposerSlashingsCmd = withTimeout(&cobra.Command{
		Use:   "proposer_slashings [<data 0.ssz> [<data 1.ssz> [<data 2.ssz> [ ... ]]]]",
		Short: "process_proposer_slashings sub state-transition",
		Args:  cobra.RangeArgs(0, beacon.MAX_PROPOSER_SLASHINGS),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				ops := make(beacon.ProposerSlashings, len(args), len(args))
				for i, arg := range args {
					err := LoadSSZ(arg, &ops[i], beacon.ProposerSlashingSSZ)
					if err != nil {
						return fmt.Errorf("could not load operation %d %s: %v", i, arg, err)
					}
				}
				return state.ProcessProposerSlashings(ctx, epc, ops)
			})
		},
	})
	DepositsCmd = withTimeout(&cobra.Command{
		Use:   "deposits [<data 0.ssz> [<data 1.ssz> [<data 2.ssz> [ ... ]]]]",
		Short: "process_deposits sub state-transition",
		Args:  cobra.RangeArgs(0, beacon.MAX_DEPOSITS),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				ops := make(beacon.Deposits, len(args), len(args))
				for i, arg := range args {
					err := LoadSSZ(arg, &ops[i], beacon.DepositSSZ)
					if err != nil {
						return fmt.Errorf("could not load operation %d %s: %v", i, arg, err)
					}
				}
				return state.ProcessDeposits(ctx, epc, ops)
			})
		},
	})
	VoluntaryExitsCmd = withTimeout(&cobra.Command{
		Use:   "voluntary_exits [<data 0.ssz> [<data 1.ssz> [<data 2.ssz> [ ... ]]]]",
		Short: "process_voluntary_exits sub state-transition",
		Args:  cobra.RangeArgs(0, beacon.MAX_VOLUNTARY_EXITS),
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(ctx context.Context, state *beacon.BeaconStateView, epc *beacon.EpochsContext) error {
				ops := make(beacon.VoluntaryExits, len(args), len(args))
				for i, arg := range args {
					err := LoadSSZ(arg, &ops[i], beacon.VoluntaryExitSSZ)
					if err != nil {
						return fmt.Errorf("could not load operation %d %s: %v", i, arg, err)
					}
				}
				return state.ProcessVoluntaryExits(ctx, epc, ops)
			})
		},
	})
	BlockCmd.AddCommand(BlockHeaderCmd, AttestationsCmd,
		AttesterSlashingsCmd, ProposerSlashingsCmd,
		DepositsCmd, VoluntaryExitsCmd)
}

func loadPreFull(cmd *cobra.Command) (*beacon.BeaconStateView, *beacon.EpochsContext, error) {
	pre, err := LoadStateViewInputFlag(cmd, "pre", true)
	if err != nil {
		return nil, nil, err
	}
	epc, err := pre.NewEpochsContext()
	if err != nil {
		return nil, nil, err
	}
	return pre, epc, nil
}

func writePost(cmd *cobra.Command, state *beacon.BeaconStateView) error {
	return WriteStateViewOutput(cmd, "post", state)
}
