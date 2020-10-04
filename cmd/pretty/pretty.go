package pretty

import (
	"encoding/json"
	"fmt"
	"github.com/protolambda/zcli/cmd/spec_types"
	. "github.com/protolambda/zcli/util"
	"github.com/spf13/cobra"
)

var PrettyCmd *cobra.Command

func MakeCmd(st *spec_types.SpecType) *cobra.Command {
	c := &cobra.Command{
		Use:   fmt.Sprintf("%s [input path]", st.Name),
		Short: fmt.Sprintf("Pretty print a %s as formatted JSON, if the input path is not specified, input is read from STDIN", st.TypeName),
		Args:  cobra.MaximumNArgs(1),
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

			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			enc.SetEscapeHTML(false)
			err = enc.Encode(dst)
			if Check(err, cmd.ErrOrStderr(), "cannot encode to formatted json") {
				return
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "")
		},
	}
	c.Flags().StringP("spec", "s", "mainnet", "The spec configuration to use. Can also be a path to a yaml config file. E.g. 'mainnet', 'minimal', or 'my_yaml_path.yml")
	return c
}

func init() {
	PrettyCmd = &cobra.Command{
		Use:   "pretty",
		Short: "pretty-print SSZ data, to indented JSON",
	}

	for _, st := range spec_types.SpecTypes {
		PrettyCmd.AddCommand(MakeCmd(st))
	}
}
