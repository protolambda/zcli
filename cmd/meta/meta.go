package meta

import (
	. "github.com/protolambda/zcli/util"
	. "github.com/protolambda/zrnt/eth2/beacon"
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
			spec, err := LoadSpec(cmd)
			if Check(err, cmd.ErrOrStderr(), "cannot load spec") {
				return
			}
			var path string
			if len(args) > 0 {
				path = args[0]
			}

			state, err := LoadStateViewInputPath(cmd, path, true, spec)
			if Check(err, cmd.ErrOrStderr(), "cannot load input") {
				return
			}
			epc, err := spec.NewEpochsContext(state)
			if Check(err, cmd.ErrOrStderr(), "cannot compute state epochs context") {
				return
			}
			currentEpoch := epc.CurrentEpoch.Epoch

			start := spec.EpochStartSlot(currentEpoch.Previous())
			end := spec.EpochStartSlot(currentEpoch + 1)
			for slot := start; slot < end; slot++ {
				committeesPerSlot, err := epc.GetCommitteeCountAtSlot(slot)
				if Check(err, cmd.ErrOrStderr(), "cannot compute committee count for slot %d", slot) {
					return
				}
				for i := CommitteeIndex(0); i < CommitteeIndex(committeesPerSlot); i++ {
					committee, err := epc.GetBeaconCommittee(slot, i)
					if Check(err, cmd.ErrOrStderr(), "cannot get committee for slot %d committee index %d", slot, i) {
						return
					}
					Report(cmd.OutOrStdout(), `epoch: %7d    slot: %9d    committee index: %4d (out of %4d)   size: %5d    indices: %v`,
						spec.SlotToEpoch(slot), slot, i, committeesPerSlot, len(committee), committee)
				}
			}
		},
	}
	CommitteesCmd.Flags().StringP("spec", "s", "mainnet", "The spec configuration to use. Can also be a path to a yaml config file. E.g. 'mainnet', 'minimal', or 'my_yaml_path.yml")

	ProposersCommand = &cobra.Command{
		Use:   "proposers [BeaconState ssz input path]",
		Short: "Print beacon proposer indices for the given state. For current epoch. If the input path is not specified, input is read from STDIN",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			spec, err := LoadSpec(cmd)
			if Check(err, cmd.ErrOrStderr(), "cannot load spec") {
				return
			}
			var path string
			if len(args) > 0 {
				path = args[0]
			}
			state, err := LoadStateViewInputPath(cmd, path, true, spec)
			if Check(err, cmd.ErrOrStderr(), "cannot load input") {
				return
			}
			epc, err := spec.NewEpochsContext(state)
			if Check(err, cmd.ErrOrStderr(), "cannot compute state epochs context") {
				return
			}
			currentEpoch := epc.CurrentEpoch.Epoch
			start := spec.EpochStartSlot(currentEpoch)
			end := spec.EpochStartSlot(currentEpoch + 1)
			for slot := start; slot < end; slot++ {
				proposerIndex, err := epc.GetBeaconProposer(slot)
				if Check(err, cmd.ErrOrStderr(), "cannot compute proposer index for slot %d", slot) {
					return
				}
				Report(cmd.OutOrStdout(), `epoch: %7d    slot: %9d    proposer index: %4d`, spec.SlotToEpoch(slot), slot, proposerIndex)
			}
		},
	}
	ProposersCommand.Flags().StringP("spec", "s", "mainnet", "The spec configuration to use. Can also be a path to a yaml config file. E.g. 'mainnet', 'minimal', or 'my_yaml_path.yml")

	MetaCmd.AddCommand(CommitteesCmd, ProposersCommand)
}
