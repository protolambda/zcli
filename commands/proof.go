package commands

import (
	"bytes"
	"context"
	"fmt"
	"github.com/protolambda/zcli/spec_types"
	"github.com/protolambda/zcli/util"
	"github.com/protolambda/ztyp/codec"
	"github.com/protolambda/ztyp/tree"
	"github.com/protolambda/ztyp/view"
	"sort"
	"strconv"
	"strings"
)

type ProofCmd struct{}

func (c *ProofCmd) Help() string {
	return "Create SSZ merkle proofs over any spec object"
}

func (c *ProofCmd) Cmd(route string) (cmd interface{}, err error) {
	phaseTypes, ok := spec_types.TypesByPhase[route]
	if !ok {
		return nil, fmt.Errorf("unrecognized phase: %s", route)
	}
	return &ProofPhaseCmd{PhaseName: route, Types: phaseTypes}, nil
}

func (c *ProofCmd) Routes() []string {
	return spec_types.Phases
}

type ProofPhaseCmd struct {
	PhaseName string
	Types     map[string]spec_types.SpecType
}

func (c *ProofPhaseCmd) Help() string {
	return fmt.Sprintf("Produce and verify arbitrary SSZ merkle proofs over any spec object in %s", c.PhaseName)
}

func (c *ProofPhaseCmd) Cmd(route string) (cmd interface{}, err error) {
	specType, ok := c.Types[route]
	if !ok {
		return nil, fmt.Errorf("unrecognized spec object type: %s", route)
	}
	return &ProofObjCmd{PhaseName: c.PhaseName, TypeName: route, Type: specType}, nil
}

func (c *ProofPhaseCmd) Routes() []string {
	return spec_types.TypeNames(c.Types)
}

type GindicesFlag []tree.Gindex64

func (p *GindicesFlag) String() string {
	if p == nil {
		return ""
	}
	var buf strings.Builder
	for _, v := range *p {
		buf.WriteString(strconv.FormatUint(uint64(v), 10))
	}
	return buf.String()
}

func (p *GindicesFlag) Set(v string) error {
	if p == nil {
		return fmt.Errorf("cannot decode gindices list into nil pointer")
	}
	parts := strings.Split(v, ",")
	*p = make([]tree.Gindex64, 0, len(parts))
	for i, v := range parts {
		s := strings.TrimSpace(v)
		if s == "" {
			continue
		}
		g, err := strconv.ParseUint(s, 0, 64)
		if err != nil {
			return fmt.Errorf("failed to parse gindex list, item %d, got: %q", i, s)
		}
		*p = append(*p, tree.Gindex64(g))
	}
	return nil
}

func (p *GindicesFlag) Type() string {
	return "comma separated list of generalized indices (in any base)"
}

type ProofObjCmd struct {
	PhaseName   string
	TypeName    string
	Type        spec_types.SpecType
	SpecOptions `ask:"."`
	Input       util.ObjInput `ask:"<input>" help:"Input, prefix with format, empty path for STDIN"`
	Gindices    GindicesFlag  `ask:"--gindices" help:"Gindices of leaf values to put in multi-proof"`
	// TODO: path
}

func (c *ProofObjCmd) Help() string {
	return fmt.Sprintf("Produce and verify arbitrary SSZ merkle proofs of type %s (%s)", c.TypeName, c.PhaseName)
}

func (c *ProofObjCmd) Run(ctx context.Context, args ...string) error {
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

	// deduplicate
	leaves := make(map[tree.Gindex64]struct{})
	for _, g := range c.Gindices {
		leaves[g] = struct{}{}
	}
	leavesSorted := make([]tree.Gindex64, 0, len(leaves))
	for g := range leaves {
		leavesSorted = append(leavesSorted, g)
	}
	sort.Slice(leavesSorted, func(i, j int) bool {
		return leavesSorted[i] < leavesSorted[j]
	})

	// mark every gindex that is between the root and the leaves
	interest := make(map[tree.Gindex64]struct{})
	for _, g := range leavesSorted {
		iter, _ := g.BitIter()
		n := tree.Gindex64(1)
		for {
			right, ok := iter.Next()
			if !ok {
				break
			}
			n *= 2
			if right {
				n += 1
			}
			interest[n] = struct{}{}
		}
	}
	witness := make(map[tree.Gindex64]struct{})
	// for every gindex that is covered, check if the sibling is covered, and if not, it's a witness
	for g := range interest {
		if _, ok := interest[g^1]; !ok {
			witness[g^1] = struct{}{}
		}
	}
	witnessSorted := make([]tree.Gindex64, 0, len(witness))
	for g := range witness {
		witnessSorted = append(witnessSorted, g)
	}
	sort.Slice(witnessSorted, func(i, j int) bool {
		return witnessSorted[i] < witnessSorted[j]
	})

	node := v.Backing()
	hFn := tree.GetHashFn()
	root := node.MerkleRoot(hFn)
	fmt.Printf("root %6b: %s\n", 1, root)

	for _, g := range leavesSorted {
		node, err := node.Getter(g)
		if err != nil {
			fmt.Printf("leaf %6b: ? (not available)\n", g)
		} else {
			fmt.Printf("leaf %6b: %s\n", g, node.MerkleRoot(hFn))
		}
	}
	for _, g := range witnessSorted {
		node, err := node.Getter(g)
		if err != nil {
			fmt.Printf("leaf %6b: ? (not available)\n", g)
		} else {
			fmt.Printf("witn %6b: %s\n", g, node.MerkleRoot(hFn))
		}
	}

	return nil
}
