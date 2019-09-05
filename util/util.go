package util

import (
	"bytes"
	"fmt"
	"github.com/protolambda/zrnt/eth2/phase0"
	"github.com/protolambda/zssz"
	"github.com/protolambda/zssz/types"
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

func LoadSSZ(path string, dst interface{}, ssz types.SSZ) error {
	r, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot read SSZ from input path: %s\n%v", path, err)
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		return fmt.Errorf("cannot read SSZ into buffer: %s\n%v", path, err)
	}
	err = zssz.Decode(bytes.NewReader(buf.Bytes()), uint64(buf.Len()), dst, ssz)
	if err != nil {
		return fmt.Errorf("cannot decode SSZ: %s\n%v", path, err)
	}
	return nil
}

func LoadSSZInputPath(cmd *cobra.Command, inPath string, stdInFallback bool, dst interface{}, ssz types.SSZ) error {
	var r io.Reader
	if stdInFallback && inPath == "" {
		r = cmd.InOrStdin()
	} else {
		var err error
		r, err = os.Open(inPath)
		if err != nil {
			return fmt.Errorf("cannot read ssz from input path: %v", err)
		}
	}
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		return fmt.Errorf("cannot read ssz into buffer: %v", err)
	}

	err = zssz.Decode(bytes.NewReader(buf.Bytes()), uint64(buf.Len()), dst, ssz)
	if err != nil {
		return fmt.Errorf("cannot decode ssz: %v", err)
	}
	return nil
}

func LoadStateInputFlag(cmd *cobra.Command, inputKey string, stdInFallback bool) (*phase0.BeaconState, error) {
	inPath, err := cmd.Flags().GetString(inputKey)
	if err != nil {
		return nil, fmt.Errorf("state path could not be parsed")
	}
	return LoadStateInputPath(cmd, inPath, stdInFallback)
}

func LoadStateInputPath(cmd *cobra.Command, inPath string, stdInFallback bool) (*phase0.BeaconState, error) {
	var r io.Reader
	if stdInFallback && inPath == "" {
		r = cmd.InOrStdin()
	} else {
		var err error
		r, err = os.Open(inPath)
		if err != nil {
			return nil, fmt.Errorf("cannot read state from input path: %v", err)
		}
	}
	return LoadStateInput(r)
}

func LoadStateInput(r io.Reader) (*phase0.BeaconState, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read state into buffer: %v", err)
	}

	var pre phase0.BeaconState
	err = zssz.Decode(bytes.NewReader(buf.Bytes()), uint64(buf.Len()), &pre, phase0.BeaconStateSSZ)
	if err != nil {
		return nil, fmt.Errorf("cannot decode state: %v", err)
	}

	return &pre, nil
}

func WriteStateOutput(cmd *cobra.Command, outKey string, state *phase0.BeaconState) error {
	outPath, err := cmd.Flags().GetString(outKey)
	if err != nil {
		return fmt.Errorf("post path could not be parsed: %v", err)
	}

	var w io.Writer
	if outPath == "" {
		w = cmd.OutOrStdout()
	} else {
		w, err = os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	}

	_, err = zssz.Encode(w, state, phase0.BeaconStateSSZ)
	if err != nil {
		return fmt.Errorf("cannot encode post-state: %v", err)
	}
	return nil
}
