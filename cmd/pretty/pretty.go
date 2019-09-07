package pretty

import (
	"fmt"
	"github.com/protolambda/zcli/cmd/spec_types"
	. "github.com/protolambda/zcli/util"
	"github.com/protolambda/zssz"
	"github.com/spf13/cobra"
)

var PrettyCmd *cobra.Command

func MakeCmd(st *spec_types.SpecType) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s [input path]", st.Name),
		Short: fmt.Sprintf("Pretty print a %s, if the input path is not specified, input is read from STDIN", st.TypeName),
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var path string
			if len(args) > 0 {
				path = args[0]
			}
			dst := st.Alloc()

			err := LoadSSZInputPath(cmd, path, true, dst, st.SSZTyp)
			if Check(err, cmd.ErrOrStderr(), "cannot load input") {
				return
			}

			zssz.Pretty(cmd.OutOrStdout(), "  ", dst, st.SSZTyp)
		},
	}
}

func init() {
	PrettyCmd = &cobra.Command{
		Use:   "pretty",
		Short: "pretty-print SSZ data",
	}

	for _, st := range spec_types.SpecTypes {
		PrettyCmd.AddCommand(MakeCmd(st))
	}
}
