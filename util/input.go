package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/snappy"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/codec"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

type ObjInput string

func (p *ObjInput) String() string {
	if p == nil {
		return "ssz:"
	}
	return string(*p)
}

func (p *ObjInput) Set(v string) error {
	*p = ObjInput(v)
	return nil
}

func (p *ObjInput) Type() string {
	return "object input (prefix with 'ssz:', 'ssz_snappy:', 'json:' or 'yaml:')"
}

func (p *ObjInput) Read(dest common.SSZObj) error {
	if p == nil {
		return fmt.Errorf("no input specified")
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
			return fmt.Errorf("failed to read std-in as input data: %v", err)
		}
		data = buf.Bytes()
	} else {
		var err error
		data, err = ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to open input file: '%s': %v", path, err)
		}
	}

	switch typ {
	case "ssz_snappy", "ssz-snappy":
		uncompressed, err := snappy.Decode(nil, data)
		if err != nil {
			return fmt.Errorf("failed to uncompress ssz_snappy input: %v", err)
		}
		data = uncompressed
		dec := codec.NewDecodingReader(bytes.NewReader(data), uint64(len(data)))
		sszDest, ok := dest.(codec.Deserializable)
		if !ok {
			return fmt.Errorf("cannot decode SSZ-snappy input into destination type %T", dest)
		}
		return sszDest.Deserialize(dec)
	case "ssz":
		dec := codec.NewDecodingReader(bytes.NewReader(data), uint64(len(data)))
		sszDest, ok := dest.(codec.Deserializable)
		if !ok {
			return fmt.Errorf("cannot decode SSZ input into destination type %T", dest)
		}
		return sszDest.Deserialize(dec)
	case "json":
		return json.Unmarshal(data, dest)
	case "yaml":
		return yaml.Unmarshal(data, dest)
	default:
		return fmt.Errorf("unrecognized data type, prefix input value with 'ssz:', 'json:' or 'yaml:'. Got: %q", typ+":")
	}
}
