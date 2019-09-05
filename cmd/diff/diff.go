package diff

import (
	"github.com/protolambda/messagediff"
	. "github.com/protolambda/zcli/util"
	"github.com/spf13/cobra"
	"os"
)

var DiffCmd, StateCmd *cobra.Command

func init() {
	DiffCmd = &cobra.Command{
		Use:   "diff",
		Short: "find the differences in SSZ data",
	}

	StateCmd = &cobra.Command{
		Use:   "state <path A> <path B>",
		Short: "Diff two BeaconState objects",
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			rA, err := os.Open(args[0])
			if Check(err, cmd.ErrOrStderr(), "cannot open state input A") {
				return
			}
			stateA, err := LoadStateInput(rA)
			if Check(err, cmd.ErrOrStderr(), "cannot load state input A") {
				return
			}
			rB, err := os.Open(args[1])
			if Check(err, cmd.ErrOrStderr(), "cannot open state input A") {
				return
			}
			stateB, err := LoadStateInput(rB)
			if Check(err, cmd.ErrOrStderr(), "cannot load state input A") {
				return
			}
			if diff, equal := messagediff.PrettyDiff(stateA, stateB, messagediff.SliceWeakEmptyOption{}); equal {
				Report(cmd.OutOrStdout(), "states A and B are equal")
			} else {
				Report(cmd.OutOrStdout(), "states A and B are different:\n%s", diff)
			}
		},
	}

	DiffCmd.AddCommand(StateCmd)
}
