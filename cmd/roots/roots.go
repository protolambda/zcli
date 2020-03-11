package roots

import (
	"fmt"
	"github.com/protolambda/zcli/cmd/spec_types"
	. "github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2/util/ssz"
	"github.com/spf13/cobra"
)

var HashTreeRootCmd *cobra.Command

func MakeHashTreeRootCmd(st *spec_types.SpecType) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s [input path]", st.Name),
		Short: fmt.Sprintf("Hash-tree-root a %s, if the input path is not specified, input is read from STDIN", st.TypeName),
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

			root := ssz.HashTreeRoot(dst, st.SSZTyp)
			Report(cmd.OutOrStdout(), "%x\n", root)
		},
	}
}

func init() {
	HashTreeRootCmd = &cobra.Command{
		Use:     "hash-tree-root",
		Aliases: []string{"hash_tree_root", "htr"},
		Short:   "hash-tree-root SSZ data",
	}

	for _, st := range spec_types.SpecTypes {
		HashTreeRootCmd.AddCommand(MakeHashTreeRootCmd(st))
	}
}
