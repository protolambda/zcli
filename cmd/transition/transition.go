package transition

import (
	"bytes"
	"fmt"
	"github.com/protolambda/zrnt/eth2/beacon/attestations"
	"github.com/protolambda/zrnt/eth2/beacon/deposits"
	"github.com/protolambda/zrnt/eth2/beacon/exits"
	"github.com/protolambda/zrnt/eth2/beacon/header"
	"github.com/protolambda/zrnt/eth2/beacon/slashings/attslash"
	"github.com/protolambda/zrnt/eth2/beacon/slashings/propslash"
	"github.com/protolambda/zrnt/eth2/beacon/transfers"
	"github.com/protolambda/zrnt/eth2/core"
	"github.com/protolambda/zrnt/eth2/phase0"
	"github.com/protolambda/zssz"
	"github.com/protolambda/zssz/types"
	"github.com/spf13/cobra"
	"io"
	"os"
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
			if check(err, cmd.ErrOrStderr(), "delta flag could not be parsed") {
				return
			}

			slots, _ := strconv.ParseUint(args[0], 10, 64)

			state, err := loadPreFull(cmd)
			if check(err, cmd.ErrOrStderr(), "pre state could not be loaded") {
				return
			}

			to := core.Slot(slots)
			if isDelta {
				to += state.Slot
			} else if to < state.Slot {
				report(cmd.ErrOrStderr(), "to slot is lower than pre-state slot")
				return
			}

			state.ProcessSlots(to)
			err = writePost(cmd, state.BeaconState)
			if check(err, cmd.ErrOrStderr(), "could not write post-state") {
				return
			}
		},
	}
	SlotsCmd.Flags().Bool("delta", false, "to interpret the slot number as a delta from the pre-state")
	TransitionCmd.AddCommand(SlotsCmd)

	BlocksCmd = &cobra.Command{
		Use:   "blocks",
		Short: "Process blocks on the pre-state to get a post-state",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			verifyStateRoot, err := cmd.Flags().GetBool("verify-state-root")
			if check(err, cmd.ErrOrStderr(), "verify-state-root could not be parsed") {
				return
			}

			state, err := loadPreFull(cmd)
			if check(err, cmd.ErrOrStderr(), "could not load pre-state") {
				return
			}

			for i := 0; i < len(args); i++ {
				var b phase0.BeaconBlock
				err := loadSSZ(args[i], &b, phase0.BeaconBlockSSZ)
				if check(err, cmd.ErrOrStderr(), "could not load block: %s", args[i]) {
					return
				}

				blockProc := &phase0.BlockProcessFeature{Block: &b, Meta: state}

				err = state.StateTransition(blockProc, verifyStateRoot)
				if check(err, cmd.ErrOrStderr(), "failed block transition to block %s", args[i]) {
					return
				}
			}

			err = writePost(cmd, state.BeaconState)
			if check(err, cmd.ErrOrStderr(), "could not write post-state") {
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
		if check(err, cmd.ErrOrStderr(), "pre state could not be loaded") {
			return
		}
		change(state)
		err = writePost(cmd, state.BeaconState)
		if check(err, cmd.ErrOrStderr(), "could not write post-state") {
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
		Use:   "attestation",
		Short: "process_attestation sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				var op attestations.Attestation
				err := loadSSZ(args[0], &op, attestations.AttestationSSZ)
				if check(err, cmd.ErrOrStderr(), "could not load attestation") {
					return
				}
				err = state.ProcessAttestation(&op)
				if check(err, cmd.ErrOrStderr(), "failed to process attestation") {
					return
				}
			})
		},
	}
	AttesterSlashingCmd = &cobra.Command{
		Use:   "attester_slashing",
		Short: "process_attester_slashing sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				var op attslash.AttesterSlashing
				err := loadSSZ(args[0], &op, attslash.AttesterSlashingSSZ)
				if check(err, cmd.ErrOrStderr(), "could not load attester slashing") {
					return
				}
				err = state.ProcessAttesterSlashing(&op)
				if check(err, cmd.ErrOrStderr(), "failed to process attester slashing") {
					return
				}
			})
		},
	}
	ProposerSlashingCmd = &cobra.Command{
		Use:   "proposer_slashing",
		Short: "process_proposer_slashing sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				var op propslash.ProposerSlashing
				err := loadSSZ(args[0], &op, propslash.ProposerSlashingSSZ)
				if check(err, cmd.ErrOrStderr(), "could not load proposer slashing") {
					return
				}
				err = state.ProcessProposerSlashing(&op)
				if check(err, cmd.ErrOrStderr(), "failed to process proposer slashing") {
					return
				}
			})
		},
	}
	DepositCmd = &cobra.Command{
		Use:   "deposit",
		Short: "process_deposit sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				var op deposits.Deposit
				err := loadSSZ(args[0], &op, deposits.DepositSSZ)
				if check(err, cmd.ErrOrStderr(), "could not load deposit") {
					return
				}
				err = state.ProcessDeposit(&op)
				if check(err, cmd.ErrOrStderr(), "failed to process deposit") {
					return
				}
			})
		},
	}
	TransferCmd = &cobra.Command{
		Use:   "transfer",
		Short: "process_transfer sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				var op transfers.Transfer
				err := loadSSZ(args[0], &op, transfers.TransferSSZ)
				if check(err, cmd.ErrOrStderr(), "could not load transfer") {
					return
				}
				err = state.ProcessTransfer(&op)
				if check(err, cmd.ErrOrStderr(), "failed to process transfer") {
					return
				}
			})
		},
	}
	VoluntaryExitCmd = &cobra.Command{
		Use:   "voluntary_exit",
		Short: "process_voluntary_exit sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				var op exits.VoluntaryExit
				err := loadSSZ(args[0], &op, exits.VoluntaryExitSSZ)
				if check(err, cmd.ErrOrStderr(), "could not load voluntary exit") {
					return
				}
				err = state.ProcessVoluntaryExit(&op)
				if check(err, cmd.ErrOrStderr(), "failed to process voluntary exit") {
					return
				}
			})
		},
	}
	OpCmd.AddCommand(AttestationCmd, AttesterSlashingCmd, ProposerSlashingCmd, DepositCmd, TransferCmd, VoluntaryExitCmd)

	BlockHeaderCmd = &cobra.Command{
		Use:   "block_header",
		Short: "process_block_header sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				var bh header.BeaconBlockHeader
				err := loadSSZ(args[0], &bh, header.BeaconBlockHeaderSSZ)
				if check(err, cmd.ErrOrStderr(), "could not load block header") {
					return
				}
				err = state.ProcessHeader(&bh)
				if check(err, cmd.ErrOrStderr(), "failed to process block header") {
					return
				}
			})
		},
	}
	AttestationsCmd = &cobra.Command{
		Use:   "attestations",
		Short: "process_attestations sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				if uint64(len(args)) > ((*phase0.Attestations)(nil)).Limit() {
					report(cmd.ErrOrStderr(), "too many attestations")
					return
				}
				ops := make(phase0.Attestations, len(args), len(args))
				for i, arg := range args {
					err := loadSSZ(args[0], &ops[i], attestations.AttestationSSZ)
					if check(err, cmd.ErrOrStderr(), "could not load attestation %d %s", i, arg) {
						return
					}
				}
				err := state.ProcessAttestations(ops)
				if check(err, cmd.ErrOrStderr(), "failed to process attestations") {
					return
				}
			})
		},
	}
	AttesterSlashingsCmd = &cobra.Command{
		Use:   "attester_slashings",
		Short: "process_attester_slashings sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				if uint64(len(args)) > ((*phase0.AttesterSlashings)(nil)).Limit() {
					report(cmd.ErrOrStderr(), "too many attester slashings")
					return
				}
				ops := make(phase0.AttesterSlashings, len(args), len(args))
				for i, arg := range args {
					err := loadSSZ(args[0], &ops[i], attslash.AttesterSlashingSSZ)
					if check(err, cmd.ErrOrStderr(), "could not load attester slashing %d %s", i, arg) {
						return
					}
				}
				err := state.ProcessAttesterSlashings(ops)
				if check(err, cmd.ErrOrStderr(), "failed to process attester slashings") {
					return
				}
			})
		},
	}
	ProposerSlashingsCmd = &cobra.Command{
		Use:   "proposer_slashings",
		Short: "process_proposer_slashings sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				if uint64(len(args)) > ((*phase0.ProposerSlashings)(nil)).Limit() {
					report(cmd.ErrOrStderr(), "too many proposer slashings")
					return
				}
				ops := make(phase0.ProposerSlashings, len(args), len(args))
				for i, arg := range args {
					err := loadSSZ(args[0], &ops[i], propslash.ProposerSlashingSSZ)
					if check(err, cmd.ErrOrStderr(), "could not load proposer slashing %d %s", i, arg) {
						return
					}
				}
				err := state.ProcessProposerSlashings(ops)
				if check(err, cmd.ErrOrStderr(), "failed to process proposer slashings") {
					return
				}
			})
		},
	}
	DepositsCmd = &cobra.Command{
		Use:   "deposits",
		Short: "process_deposits sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				if uint64(len(args)) > ((*phase0.Deposits)(nil)).Limit() {
					report(cmd.ErrOrStderr(), "too many deposits")
					return
				}
				ops := make(phase0.Deposits, len(args), len(args))
				for i, arg := range args {
					err := loadSSZ(args[0], &ops[i], deposits.DepositSSZ)
					if check(err, cmd.ErrOrStderr(), "could not load deposit %d %s", i, arg) {
						return
					}
				}
				err := state.ProcessDeposits(ops)
				if check(err, cmd.ErrOrStderr(), "failed to process deposits") {
					return
				}
			})
		},
	}
	TransfersCmd = &cobra.Command{
		Use:   "transfers",
		Short: "process_transfers sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				if uint64(len(args)) > ((*phase0.Transfers)(nil)).Limit() {
					report(cmd.ErrOrStderr(), "too many transfers")
					return
				}
				ops := make(phase0.Transfers, len(args), len(args))
				for i, arg := range args {
					err := loadSSZ(args[0], &ops[i], transfers.TransferSSZ)
					if check(err, cmd.ErrOrStderr(), "could not load transfer %d %s", i, arg) {
						return
					}
				}
				err := state.ProcessTransfers(ops)
				if check(err, cmd.ErrOrStderr(), "failed to process transfers") {
					return
				}
			})
		},
	}
	VoluntaryExitsCmd = &cobra.Command{
		Use:   "voluntary_exits",
		Short: "process_voluntary_exits sub state-transition",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			transition(cmd, func(state *phase0.FullFeaturedState) {
				if uint64(len(args)) > ((*phase0.VoluntaryExits)(nil)).Limit() {
					report(cmd.ErrOrStderr(), "too many voluntary exits")
					return
				}
				ops := make(phase0.VoluntaryExits, len(args), len(args))
				for i, arg := range args {
					err := loadSSZ(args[0], &ops[i], exits.VoluntaryExitSSZ)
					if check(err, cmd.ErrOrStderr(), "could not load voluntary exit %d %s", i, arg) {
						return
					}
				}
				err := state.ProcessVoluntaryExits(ops)
				if check(err, cmd.ErrOrStderr(), "failed to process voluntary exits") {
					return
				}
			})
		},
	}
	BlockCmd.AddCommand(BlockHeaderCmd, AttestationsCmd,
		AttesterSlashingsCmd, ProposerSlashingsCmd,
		DepositsCmd, TransfersCmd, VoluntaryExitsCmd)
}

func report(out io.Writer, msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(out, msg, args...)
}

func check(err error, out io.Writer, msg string, args ...interface{}) bool {
	if err != nil {
		report(out, msg, args...)
		report(out, "%v", err)
		return true
	} else {
		return false
	}
}

func loadSSZ(path string, dst interface{}, ssz types.SSZ) error {
	r, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot read SSZ from input path: %s\n%v", path, err)
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		return fmt.Errorf("cannot read SSZ into buffer: %s\n%v", path, err)
	}
	err = zssz.Decode(&buf, uint64(buf.Len()), dst, ssz)
	if err != nil {
		return fmt.Errorf("cannot decode SSZ: %s\n%v", path, err)
	}
	return nil
}

func loadPreFull(cmd *cobra.Command) (*phase0.FullFeaturedState, error) {
	inPath, err := cmd.Flags().GetString("pre")
	if err != nil {
		return nil, fmt.Errorf("pre path could not be parsed")
	}

	var r io.Reader
	if inPath == "" {
		r = os.Stdin
	} else {
		r, err = os.Open(inPath)
		if err != nil {
			return nil, fmt.Errorf("cannot read pre from input path: %v", err)
		}
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read pre-state into buffer: %v", err)
	}

	var pre phase0.BeaconState
	err = zssz.Decode(&buf, uint64(buf.Len()), &pre, phase0.BeaconStateSSZ)
	if err != nil {
		return nil, fmt.Errorf("cannot decode pre-state: %v", err)
	}

	preFull := phase0.NewFullFeaturedState(&pre)
	preFull.LoadPrecomputedData()

	return preFull, nil
}

func writePost(cmd *cobra.Command, state *phase0.BeaconState) error {
	outPath, err := cmd.Flags().GetString("post")
	if err != nil {
		return fmt.Errorf("post path could not be parsed: %v", err)
	}

	var w io.Writer
	if outPath == "" {
		w = os.Stdout
	} else {
		w, err = os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	}

	_, err = zssz.Encode(w, state, phase0.BeaconStateSSZ)
	if err != nil {
		return fmt.Errorf("cannot encode post-state: %v", err)
	}
	return nil
}
