package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/protolambda/zrnt/eth2/beacon"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/merge"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/zrnt/eth2/beacon/sharding"
	"github.com/protolambda/ztyp/codec"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
)

type StateOutput string

func (p *StateOutput) String() string {
	if p == nil {
		return "ssz:"
	}
	return string(*p)
}

func (p *StateOutput) Set(v string) error {
	*p = StateOutput(v)
	return nil
}

func (p *StateOutput) Type() string {
	return "BeaconState output (prefix with 'ssz:', 'json:', 'pretty:' or 'yaml:')"
}

func (p *StateOutput) Write(spec *common.Spec, obj common.BeaconState) error {
	if p == nil {
		return fmt.Errorf("no output specified")
	}
	full := string(*p)
	partIndex := strings.Index(full, ":")
	var typ, path string
	if partIndex >= 0 {
		typ = full[:partIndex]
		path = full[partIndex+1:]
	} else {
		// default to ssz input
		typ = "ssz"
		path = full
	}

	var out io.Writer
	if path == "" {
		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()
		out = w
	} else {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			return err
		}
		defer f.Close()
		w := bufio.NewWriter(f)
		defer w.Flush()
		out = w
	}

	switch typ {
	case "ssz":
		enc := codec.NewEncodingWriter(out)
		return obj.Serialize(enc)
	case "json", "pretty", "yaml":
		var tmp bytes.Buffer
		enc := codec.NewEncodingWriter(&tmp)
		if err := obj.Serialize(enc); err != nil {
			return err
		}
		var flat common.SpecObj
		if up, ok := obj.(*beacon.StandardUpgradeableBeaconState); ok {
			obj = up.BeaconState
		}
		switch obj.(type) {
		case *phase0.BeaconStateView:
			flat = new(phase0.BeaconState)
		case *altair.BeaconStateView:
			flat = new(altair.BeaconState)
		case *merge.BeaconStateView:
			flat = new(merge.BeaconState)
		case *sharding.BeaconStateView:
			flat = new(sharding.BeaconState)
		default:
			return fmt.Errorf("failed to detect state type for output: %T", obj)
		}
		data := tmp.Bytes()
		if err := flat.Deserialize(spec, codec.NewDecodingReader(bytes.NewReader(data), uint64(len(data)))); err != nil {
			return err
		}
		switch typ {
		case "json":
			return json.NewEncoder(out).Encode(flat)
		case "pretty":
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(flat)
		case "yaml":
			return yaml.NewEncoder(out).Encode(flat)
		default:
			panic("unreachable")
		}
	default:
		return fmt.Errorf("unrecognized data type, prefix output value with 'ssz:', 'json:', 'pretty:' or 'yaml:'. Got: %q", typ+":")
	}
}
