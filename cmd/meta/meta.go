package meta

import (
	. "github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2/core"
	"github.com/protolambda/zrnt/eth2/phase0"
	"github.com/spf13/cobra"
)

var MetaCmd, CommitteesCmd *cobra.Command

func init() {
	MetaCmd = &cobra.Command{
		Use:   "meta",
		Short: "Print meta information of a BeaconState",
	}

	CommitteesCmd = &cobra.Command{
		Use:   "committees [BeaconState ssz input path]",
		Short: "Print shard committees for the given state. For prev, current and next epoch. If the input path is not specified, input is read from STDIN",
		Args:  cobra.ExactArgs(1),
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

			for i := core.Shard(0); i < core.SHARD_COUNT; i++ {
				prevEpoch := state.PreviousEpoch()
				currEpoch := state.CurrentEpoch()
				nextEpoch := state.CurrentEpoch() + 1
				prev := full.GetCrosslinkCommittee(prevEpoch, i)
				curr := full.GetCrosslinkCommittee(currEpoch, i)
				next := full.GetCrosslinkCommittee(nextEpoch, i)
				Report(cmd.OutOrStdout(), `shard %d:
  previous (%d): %v
  current  (%d): %v
  next     (%d): %v
`, i, prevEpoch, prev, currEpoch, curr, nextEpoch, next)
			}
		},
	}

	MetaCmd.AddCommand(CommitteesCmd)
}
