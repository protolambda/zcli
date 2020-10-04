package roots

import (
	"fmt"
	"github.com/protolambda/zcli/cmd/spec_types"
	. "github.com/protolambda/zcli/util"
	"github.com/protolambda/ztyp/tree"
	"github.com/spf13/cobra"
)

var HashTreeRootCmd *cobra.Command

func MakeHashTreeRootCmd(st *spec_types.SpecType) *cobra.Command {
	c := &cobra.Command{
		Use:   fmt.Sprintf("%s [input path]", st.Name),
		Short: fmt.Sprintf("Hash-tree-root a %s, if the input path is not specified, input is read from STDIN", st.TypeName),
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			spec, err := LoadSpec(cmd)
			if Check(err, cmd.ErrOrStderr(), "cannot load spec") {
				return
			}
			var path string
			if len(args) > 0 {
				path = args[0]
			}
			dst := st.Alloc()
			sszObj, err := ItSSZ(dst, spec)
			if Check(err, cmd.ErrOrStderr(), "cannot use type as ssz object") {
				return
			}
			err = LoadSSZInputPath(cmd, path, true, sszObj)
			if Check(err, cmd.ErrOrStderr(), "cannot load input") {
				return
			}

			root := sszObj.HashTreeRoot(tree.GetHashFn())
			Report(cmd.OutOrStdout(), "%s\n", root.String())
		},
	}
	c.Flags().StringP("spec", "s", "mainnet", "The spec configuration to use. Can also be a path to a yaml config file. E.g. 'mainnet', 'minimal', or 'my_yaml_path.yml")
	return c
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
