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
	diffCmd := &cobra.Command{
		Use:   fmt.Sprintf("%s <path A> <path B>", st.Name),
		Short: fmt.Sprintf("Diff two %s objects", st.TypeName),
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			spec, err := LoadSpec(cmd)
			if Check(err, cmd.ErrOrStderr(), "cannot load spec") {
				return
			}
			rA, err := os.Open(args[0])
			if Check(err, cmd.ErrOrStderr(), "cannot open SSZ input A") {
				return
			}
			objA := st.Alloc()
			objASSZ, err := ItSSZ(objA, spec)
			if Check(err, cmd.ErrOrStderr(), "cannot open decode input A") {
				return
			}
			if Check(LoadSSZInput(rA, objASSZ), cmd.ErrOrStderr(), "cannot load SSZ input A") {
				return
			}
			rB, err := os.Open(args[1])
			if Check(err, cmd.ErrOrStderr(), "cannot open SSZ input A") {
				return
			}
			objB := st.Alloc()
			objBSSZ, err := ItSSZ(objB, spec)
			if Check(err, cmd.ErrOrStderr(), "cannot open decode input B") {
				return
			}
			if Check(LoadSSZInput(rB, objBSSZ), cmd.ErrOrStderr(), "cannot load SSZ input B") {
				return
			}
			if diff, equal := messagediff.PrettyDiff(objA, objB, messagediff.SliceWeakEmptyOption{}); equal {
				Report(cmd.OutOrStdout(), "%s objects A and B are equal\n", st.TypeName)
			} else {
				Report(cmd.OutOrStdout(), "%s objects A and B are different:\n%s\n", st.TypeName, diff)
			}
		},
	}
	diffCmd.Flags().StringP("spec", "s", "mainnet", "The spec configuration to use. Can also be a path to a yaml config file. E.g. 'mainnet', 'minimal', or 'my_yaml_path.yml")
	return diffCmd
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
