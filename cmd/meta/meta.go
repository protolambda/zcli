package meta

import (
	. "github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2/core"
	"github.com/protolambda/zrnt/eth2/phase0"
	"github.com/spf13/cobra"
)

var MetaCmd, CommitteesCmd, ProposersCommand *cobra.Command

func init() {
	MetaCmd = &cobra.Command{
		Use:   "meta",
		Short: "Print meta information of a BeaconState",
	}

	CommitteesCmd = &cobra.Command{
		Use:   "committees [BeaconState ssz input path]",
		Short: "Print beacon committees for the given state. For prev, current and next epoch. If the input path is not specified, input is read from STDIN",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			var path string
			if len(args) > 0 {
				path = args[0]
			}
			var state phase0.BeaconState
			err := LoadSSZInputPath(cmd, path, true, &state, phase0.BeaconStateSSZ)
			if Check(err, cmd.ErrOrStderr(), "cannot load input") {
				return
			}

			full := phase0.NewFullFeaturedState(&state)
			full.LoadPrecomputedData()

			start := state.PreviousEpoch().GetStartSlot()
			end := (state.CurrentEpoch() + 1).GetStartSlot()
			for slot := start; slot < end; slot++ {
				committeesPerSlot := core.CommitteeIndex(full.GetCommitteeCountAtSlot(slot))
				for i := core.CommitteeIndex(0); i < committeesPerSlot; i++ {
					committee := full.GetBeaconCommittee(slot, i)
					Report(cmd.OutOrStdout(), `epoch: %7d    slot: %9d    committee index: %4d (out of %4d)   size: %5d    indices: %v`,
						slot.ToEpoch(), slot, i, committeesPerSlot, len(committee), committee)
				}
			}
		},
	}

	ProposersCommand = &cobra.Command{
		Use:   "proposers [BeaconState ssz input path]",
		Short: "Print beacon proposer indices for the given state. For prev, current and next epoch. If the input path is not specified, input is read from STDIN",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			var path string
			if len(args) > 0 {
				path = args[0]
			}
			var state phase0.BeaconState
			err := LoadSSZInputPath(cmd, path, true, &state, phase0.BeaconStateSSZ)
			if Check(err, cmd.ErrOrStderr(), "cannot load input") {
				return
			}

			full := phase0.NewFullFeaturedState(&state)
			full.LoadPrecomputedData()

			start := state.PreviousEpoch().GetStartSlot()
			end := (state.CurrentEpoch() + 1).GetStartSlot()
			for slot := start; slot < end; slot++ {
				proposerIndex := full.GetBeaconProposerIndex(slot)
				Report(cmd.OutOrStdout(), `epoch: %7d    slot: %9d    proposer index: %4d`, slot.ToEpoch(), slot, proposerIndex)
			}
		},
	}

	MetaCmd.AddCommand(CommitteesCmd, ProposersCommand)
}
