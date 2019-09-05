package pretty

import (
	"fmt"
	. "github.com/protolambda/zcli/util"
	"github.com/spf13/cobra"
)

var PrettyCmd, StateCmd *cobra.Command

func init() {
	PrettyCmd = &cobra.Command{
		Use:   "pretty",
		Short: "pretty-print SSZ data",
	}

	StateCmd = &cobra.Command{
		Use:   "state [input path]",
		Short: "Pretty print a BeaconState",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			state, err := LoadStateInput(cmd, args[0])
			if Check(err, cmd.ErrOrStderr(), "cannot load state input") {
				return
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%v", state)
		},
	}

	PrettyCmd.AddCommand(StateCmd)
}
