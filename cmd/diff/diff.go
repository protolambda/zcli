package diff

import (
	"fmt"
	"github.com/protolambda/messagediff"
	"github.com/protolambda/zcli/cmd/spec_types"
	. "github.com/protolambda/zcli/util"
	"github.com/spf13/cobra"
	"os"
)

var DiffCmd *cobra.Command


func MakeCmd(st *spec_types.SpecType) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s <path A> <path B>", st.Name),
		Short: fmt.Sprintf("Diff two %s objects", st.TypeName),
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			rA, err := os.Open(args[0])
			if Check(err, cmd.ErrOrStderr(), "cannot open SSZ input A") {
				return
			}
			objA := st.Alloc()
			if Check(LoadSSZInput(rA, objA, st.SSZTyp), cmd.ErrOrStderr(), "cannot load SSZ input A") {
				return
			}
			rB, err := os.Open(args[1])
			if Check(err, cmd.ErrOrStderr(), "cannot open SSZ input A") {
				return
			}
			objB := st.Alloc()
			if Check(LoadSSZInput(rB, objB, st.SSZTyp), cmd.ErrOrStderr(), "cannot load SSZ input B") {
				return
			}
			if diff, equal := messagediff.PrettyDiff(objA, objB, messagediff.SliceWeakEmptyOption{}); equal {
				Report(cmd.OutOrStdout(), "%s objects A and B are equal\n", st.TypeName)
			} else {
				Report(cmd.OutOrStdout(), "%s objects A and B are different:\n%s\n", st.TypeName, diff)
			}
		},
	}
}


func init() {
	DiffCmd = &cobra.Command{
		Use:   "diff",
		Short: "find the differences in SSZ data",
	}

	for _, st := range spec_types.SpecTypes {
		DiffCmd.AddCommand(MakeCmd(st))
	}
}
