package util

import (
	"bytes"
	"fmt"
	"github.com/protolambda/zrnt/eth2/beacon"
	"github.com/protolambda/zrnt/eth2/configs"
	"github.com/protolambda/ztyp/codec"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Report(out io.Writer, msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(out, msg+"\n", args...)
}

func Check(err error, out io.Writer, msg string, args ...interface{}) bool {
	if err != nil {
		Report(out, msg, args...)
		Report(out, "%v", err)
		return true
	} else {
		return false
	}
}

func ItSSZ(obj interface{}, spec *beacon.Spec) (beacon.SSZObj, error) {
	if v, ok := obj.(beacon.SpecObj); ok {
		return spec.Wrap(v), nil
	}
	if v, ok := obj.(beacon.SSZObj); ok {
		return v, nil
	}
	return nil, fmt.Errorf("object has no ssz methods")
}

func LoadSpec(cmd *cobra.Command) (*beacon.Spec, error) {
	specPath, err := cmd.Flags().GetString("spec")
	if err != nil {
		return nil, fmt.Errorf("spec name/path could not be parsed")
	}
	ext := filepath.Ext(specPath)
	if ext == ".yaml" || ext == ".yml" {
		data, err := ioutil.ReadFile(specPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load spec yaml: %v", err)
		}
		var conf beacon.Phase0Config
		if err := yaml.Unmarshal(data, &conf); err != nil {
			return nil, fmt.Errorf("failed to decode spec yaml: %v", err)
		}
		name := filepath.Base(specPath)
		return &beacon.Spec{
			PRESET_NAME:  name[:strings.LastIndex(name, ext)],
			Phase0Config: conf,
			Phase1Config: configs.Mainnet.Phase1Config, // TODO: ZCLI may support phase 1 changes later.
		}, nil
	} else {
		switch specPath {
		case "mainnet":
			return configs.Mainnet, nil
		case "minimal":
			return configs.Minimal, nil
		default:
			return nil, fmt.Errorf("invalid/unknown spec: '%s'", specPath)
		}
	}
}

func LoadSSZ(path string, dst codec.Deserializable) error {
	r, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot read SSZ from input path: %s\n%v", path, err)
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		return fmt.Errorf("cannot read SSZ into buffer: %s\n%v", path, err)
	}
	err = dst.Deserialize(codec.NewDecodingReader(bytes.NewReader(buf.Bytes()), uint64(buf.Len())))
	if err != nil {
		return fmt.Errorf("cannot decode SSZ: %s\n%v", path, err)
	}
	return nil
}

func LoadSSZInput(r io.Reader, dst codec.Deserializable) error {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		return fmt.Errorf("cannot read SSZ into buffer: \n%v", err)
	}
	err = dst.Deserialize(codec.NewDecodingReader(bytes.NewReader(buf.Bytes()), uint64(buf.Len())))
	if err != nil {
		return fmt.Errorf("cannot decode SSZ: \n%v", err)
	}
	return nil
}

func LoadSSZInputPath(cmd *cobra.Command, inPath string, stdInFallback bool, dst codec.Deserializable) error {
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

	err = dst.Deserialize(codec.NewDecodingReader(bytes.NewReader(buf.Bytes()), uint64(buf.Len())))
	if err != nil {
		return fmt.Errorf("cannot decode ssz: %v", err)
	}
	return nil
}

func LoadStateViewInputPath(cmd *cobra.Command, inPath string, stdInFallback bool, spec *beacon.Spec) (*beacon.BeaconStateView, error) {
	var r io.Reader
	if stdInFallback && inPath == "" {
		r = cmd.InOrStdin()
	} else {
		var err error
		r, err = os.Open(inPath)
		if err != nil {
			return nil, fmt.Errorf("cannot read state view from input path: %v", err)
		}
	}
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read ssz into buffer: %v", err)
	}

	v, err := spec.BeaconState().Deserialize(codec.NewDecodingReader(bytes.NewReader(buf.Bytes()), uint64(buf.Len())))
	return beacon.AsBeaconStateView(v, nil)
}

func LoadStateViewInputFlag(cmd *cobra.Command, inputKey string, stdInFallback bool, spec *beacon.Spec) (*beacon.BeaconStateView, error) {
	inPath, err := cmd.Flags().GetString(inputKey)
	if err != nil {
		return nil, fmt.Errorf("state path could not be parsed")
	}
	return LoadStateViewInputPath(cmd, inPath, stdInFallback, spec)
}

func LoadStateInputPath(cmd *cobra.Command, inPath string, stdInFallback bool, spec *beacon.Spec) (*beacon.BeaconState, error) {
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
	return LoadStateInput(r, spec)
}

func LoadStateInput(r io.Reader, spec *beacon.Spec) (*beacon.BeaconState, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read state into buffer: %v", err)
	}

	var pre beacon.BeaconState
	err = pre.Deserialize(spec, codec.NewDecodingReader(bytes.NewReader(buf.Bytes()), uint64(buf.Len())))
	if err != nil {
		return nil, fmt.Errorf("cannot decode state: %v", err)
	}

	return &pre, nil
}

func WriteStateViewOutput(cmd *cobra.Command, outKey string, state *beacon.BeaconStateView) error {
	outPath, err := cmd.Flags().GetString(outKey)
	if err != nil {
		return err
	}
	return WriteStateViewOutputFile(cmd, outPath, state)
}

func WriteStateViewOutputFile(cmd *cobra.Command, outPath string, state *beacon.BeaconStateView) (err error) {
	var w io.Writer
	if outPath == "" {
		w = cmd.OutOrStdout()
	} else {
		w, err = os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if err := state.Serialize(codec.NewEncodingWriter(w)); err != nil {
		return fmt.Errorf("cannot encode post-state: %v", err)
	}
	return nil
}

func WriteStateOutputFile(cmd *cobra.Command, outPath string, state *beacon.BeaconState, spec *beacon.Spec) (err error) {
	var w io.Writer
	if outPath == "" {
		w = cmd.OutOrStdout()
	} else {
		w, err = os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = state.Serialize(spec, codec.NewEncodingWriter(w))
	if err != nil {
		return fmt.Errorf("cannot encode post-state: %v", err)
	}
	return nil
}
