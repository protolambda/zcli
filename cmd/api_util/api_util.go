package api_util

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

var ExtractState *cobra.Command

func MakeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("extract_state [input (APIBeaconState) path] [output (BeaconState) path]"),
		Short: fmt.Sprintf("Extract the state from an api beacon state (wrapper with root). If the input path is not specified, input is read from STDIN. If the output is not specified, output is written to STDOUT. The root of the wrapper is output to STDERR."),
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var inPath, outPath string
			if len(args) > 0 {
				if len(args) == 1 {
					Report(cmd.ErrOrStderr(), "received two arguments but need either 0 or 2.")
					return
				}
				inPath = args[0]
				outPath = args[1]
			}

			var r io.Reader
			if inPath == "" {
				r = cmd.InOrStdin()
			} else {
				var err error
				r, err = os.Open(inPath)
				if Check(err, cmd.ErrOrStderr(), "cannot read ssz from input path") {
					return
				}
			}
			var buf bytes.Buffer
			_, err := buf.ReadFrom(r)
			if Check(err, cmd.ErrOrStderr(), "cannot read ssz into buffer") {
				return
			}

			var apiState spec_types.APIBeaconState
			if err := zssz.Decode(bytes.NewReader(buf.Bytes()), uint64(buf.Len()), &apiState, spec_types.APIBeaconStateSSZ);
				Check(err, cmd.ErrOrStderr(), "cannot verify input") {
				return
			}

			if Check(WriteStateOutputFile(cmd, outPath, &apiState.BeaconState), cmd.ErrOrStderr(), "cannot output state") {
				return
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "API provided root: 0x%x\n", apiState.Root)
		},
	}
}

func init() {
	ExtractState = MakeCmd()
}
