package checker

import (
	"bytes"
	"fmt"
	"github.com/protolambda/zcli/cmd/spec_types"
	. "github.com/protolambda/zcli/util"
	"github.com/protolambda/zssz"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var CheckCmd *cobra.Command

func MakeCmd(st *spec_types.SpecType) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s [input path]", st.Name),
		Short: fmt.Sprintf("Check if the input is a valid serialized %s, if the input path is not specified, input is read from STDIN", st.TypeName),
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var path string
			if len(args) > 0 {
				path = args[0]
			}

			var r io.Reader
			if path == "" {
				r = cmd.InOrStdin()
			} else {
				var err error
				r, err = os.Open(path)
				if Check(err, cmd.ErrOrStderr(), "cannot read ssz from input path") {
					return
				}
			}
			var buf bytes.Buffer
			_, err := buf.ReadFrom(r)
			if Check(err, cmd.ErrOrStderr(), "cannot read ssz into buffer") {
				return
			}

			if err := zssz.DryCheck(bytes.NewReader(buf.Bytes()), uint64(buf.Len()), st.SSZTyp); Check(err, cmd.ErrOrStderr(), "cannot verify input") {
				return
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "valid %s!\n", st.TypeName)
		},
	}
}

func init() {
	CheckCmd = &cobra.Command{
		Use:   "check",
		Short: "check SSZ data format",
	}

	for _, st := range spec_types.SpecTypes {
		CheckCmd.AddCommand(MakeCmd(st))
	}
}
