package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/protolambda/zcli/spec_types"
	"github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2/configs"
	"github.com/protolambda/ztyp/codec"
	"github.com/protolambda/ztyp/tree"
	"github.com/protolambda/ztyp/view"
	"io"
	"os"
)

type TreeCmd struct{}

func (c *TreeCmd) Help() string {
	return "Dump SSZ merkle tree of any spec object"
}

func (c *TreeCmd) Cmd(route string) (cmd interface{}, err error) {
	phaseTypes, ok := spec_types.TypesByPhase[route]
	if !ok {
		return nil, fmt.Errorf("unrecognized phase: %s", route)
	}
	return &TreePhaseCmd{PhaseName: route, Types: phaseTypes}, nil
}

func (c *TreeCmd) Routes() []string {
	return spec_types.Phases
}

type TreePhaseCmd struct {
	PhaseName string
	Types     map[string]spec_types.SpecType
}

func (c *TreePhaseCmd) Help() string {
	return fmt.Sprintf("Dump SSZ merkle tree of any spec object in %s", c.PhaseName)
}

func (c *TreePhaseCmd) Cmd(route string) (cmd interface{}, err error) {
	specType, ok := c.Types[route]
	if !ok {
		return nil, fmt.Errorf("unrecognized spec object type: %s", route)
	}
	return &TreeObjCmd{PhaseName: c.PhaseName, TypeName: route, Type: specType}, nil
}

func (c *TreePhaseCmd) Routes() []string {
	return spec_types.TypeNames(c.Types)
}

type TreeObjCmd struct {
	PhaseName           string
	TypeName            string
	Type                spec_types.SpecType
	configs.SpecOptions `ask:"."`
	Input               util.ObjInput `ask:"<input>" help:"Input, prefix with format, empty path for STDIN"`
}

func (c *TreeObjCmd) Help() string {
	return fmt.Sprintf("Dump merkle tree of type %s (%s)", c.TypeName, c.PhaseName)
}

func (c *TreeObjCmd) Run(ctx context.Context, args ...string) error {
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	objA := c.Type.Alloc(spec)
	if err := c.Input.Read(objA); err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}

	var v view.View
	// Hack to convert flat-structure input to tree-structure
	{
		var buf bytes.Buffer
		if err := objA.Serialize(codec.NewEncodingWriter(&buf)); err != nil {
			return err
		}
		data := buf.Bytes()
		dec := codec.NewDecodingReader(bytes.NewReader(data), uint64(len(data)))
		treeType := c.Type.TypeDef(spec)
		v, err = treeType.Deserialize(dec)
		if err != nil {
			return err
		}
	}

	node := v.Backing()
	hFn := tree.GetHashFn()
	// run merkleization upfront (caching will ease work later)
	_ = node.MerkleRoot(hFn)
	return dumpTree(os.Stdout, node)
}

type nodeStruct struct {
	Left  *nodeStruct `json:"left,omitempty"`
	Right *nodeStruct `json:"right,omitempty"`
	Root  tree.Root   `json:"root"`
}

func dumpTree(w io.Writer, node tree.Node) error {
	treeDumpStruct := toJsonStruct(node, tree.GetHashFn())
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(treeDumpStruct)
}

func toJsonStruct(node tree.Node, hFn tree.HashFn) *nodeStruct {
	if node.IsLeaf() {
		return &nodeStruct{Root: node.MerkleRoot(hFn)}
	} else {
		left, _ := node.Left()
		right, _ := node.Right()
		// avoid non-tree structure memory blow up danger (some ztyp internals reference sub-trees twice)
		if left == right {
			panic("don't")
		}
		return &nodeStruct{Left: toJsonStruct(left, hFn), Right: toJsonStruct(right, hFn), Root: node.MerkleRoot(hFn)}
	}
}
