package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/protolambda/ztyp/codec"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
)

type ObjOutput string

func (p *ObjOutput) String() string {
	if p == nil {
		return "ssz:"
	}
	return string(*p)
}

func (p *ObjOutput) Set(v string) error {
	*p = ObjOutput(v)
	return nil
}

func (p *ObjOutput) Type() string {
	return "object output (prefix with 'ssz:', 'json:', 'pretty:' or 'yaml:')"
}

func (p *ObjOutput) Write(obj interface{}) error {
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
		sszObj, ok := obj.(codec.Serializable)
		if !ok {
			return fmt.Errorf("cannot encode SSZ object type %T to output", obj)
		}
		return sszObj.Serialize(enc)
	case "json":
		return json.NewEncoder(out).Encode(obj)
	case "pretty":
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(obj)
	case "yaml":
		return yaml.NewEncoder(out).Encode(obj)
	default:
		return fmt.Errorf("unrecognized data type, prefix output value with 'ssz:', 'json:', 'pretty:' or 'yaml:'. Got: %q", typ+":")
	}
}
