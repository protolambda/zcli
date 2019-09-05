package util

import (
	"fmt"
	"github.com/protolambda/zrnt/eth2/phase0"
	"github.com/protolambda/zssz"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func Report(out io.Writer, msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(out, msg, args...)
}

func Check(err error, out io.Writer, msg string, args ...interface{}) bool {
	if err != nil {
		Report(out, msg, args...)
		Report(out, "\n%v", err)
		return true
	} else {
		return false
	}
}

func WriteState(cmd *cobra.Command, outKey string, state *phase0.BeaconState) error {
	outPath, err := cmd.Flags().GetString(outKey)
	if err != nil {
		return fmt.Errorf("post path could not be parsed: %v", err)
	}

	var w io.Writer
	if outPath == "" {
		w = os.Stdout
	} else {
		w, err = os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	}

	_, err = zssz.Encode(w, state, phase0.BeaconStateSSZ)
	if err != nil {
		return fmt.Errorf("cannot encode post-state: %v", err)
	}
	return nil
}
