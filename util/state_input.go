package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/merge"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/zrnt/eth2/beacon/sharding"
	"github.com/protolambda/ztyp/codec"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

type StateInput string

func (p *StateInput) String() string {
	if p == nil {
		return "ssz:"
	}
	return string(*p)
}

func (p *StateInput) Set(v string) error {
	*p = StateInput(v)
	return nil
}

func (p *StateInput) Type() string {
	return "BeaconState input (prefix with 'ssz:', 'json:' or 'yaml:')"
}

func (p *StateInput) Read(spec *common.Spec, phase string) (common.BeaconState, error) {
	if p == nil {
		return nil, fmt.Errorf("no input specified")
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

	var data []byte
	if path == "" {
		var buf bytes.Buffer
		_, err := buf.ReadFrom(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("failed to read std-in as input data: %v", err)
		}
		data = buf.Bytes()
	} else {
		var err error
		data, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open input file: '%s': %v", path, err)
		}
	}

	switch typ {
	case "ssz":
		// nothing to do, already ssz bytes
		break
	case "json", "yaml":
		// convert to SSZ (we can't decode json or yaml directly into a tree-backed structure)
		var flat common.SpecObj
		switch phase {
		case "phase0":
			flat = new(phase0.BeaconState)
		case "altair":
			flat = new(altair.BeaconState)
		case "merge":
			flat = new(merge.BeaconState)
		case "sharding":
			flat = new(sharding.BeaconState)
		default:
			return nil, fmt.Errorf("unrecognized phase: %s", phase)
		}
		if typ == "json" {
			if err := json.Unmarshal(data, &flat); err != nil {
				return nil, fmt.Errorf("failed to decode JSON beacon state into flat structure: %v", err)
			}
		} else {
			if err := yaml.Unmarshal(data, &flat); err != nil {
				return nil, fmt.Errorf("failed to decode YAML beacon state into flat structure: %v", err)
			}
		}
		var buf bytes.Buffer
		if err := flat.Serialize(spec, codec.NewEncodingWriter(&buf)); err != nil {
			return nil, err
		}
		data = buf.Bytes()
	default:
		return nil, fmt.Errorf("unrecognized data type, prefix input value with 'ssz:', 'json:' or 'yaml:'. Got: %q", typ+":")
	}

	dec := codec.NewDecodingReader(bytes.NewReader(data), uint64(len(data)))
	switch phase {
	case "phase0":
		return phase0.AsBeaconStateView(phase0.BeaconStateType(spec).Deserialize(dec))
	case "altair":
		return altair.AsBeaconStateView(altair.BeaconStateType(spec).Deserialize(dec))
	case "merge":
		return merge.AsBeaconStateView(merge.BeaconStateType(spec).Deserialize(dec))
	case "sharding":
		return sharding.AsBeaconStateView(sharding.BeaconStateType(spec).Deserialize(dec))
	default:
		return nil, fmt.Errorf("unrecognized phase: %s", phase)
	}
}
