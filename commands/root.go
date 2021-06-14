package commands

import (
	"context"
	"fmt"
	"github.com/protolambda/zcli/spec_types"
	"github.com/protolambda/zcli/util"
	"github.com/protolambda/ztyp/tree"
)

type RootCmd struct{}

func (c *RootCmd) Help() string {
	return "Compute the SSZ hash-tree-root of a spec object"
}

func (c *RootCmd) Cmd(route string) (cmd interface{}, err error) {
	phaseTypes, ok := spec_types.TypesByPhase[route]
	if !ok {
		return nil, fmt.Errorf("unrecognized phase: %s", route)
	}
	return &RootPhaseCmd{PhaseName: route, Types: phaseTypes}, nil
}

func (c *RootCmd) Routes() []string {
	return spec_types.Phases
}

type RootPhaseCmd struct {
	PhaseName string
	Types     map[string]spec_types.SpecType
}

func (c *RootPhaseCmd) Help() string {
	return fmt.Sprintf("Compute the SSZ hash-tree-root of any %s spec object", c.PhaseName)
}

func (c *RootPhaseCmd) Cmd(route string) (cmd interface{}, err error) {
	specType, ok := c.Types[route]
	if !ok {
		return nil, fmt.Errorf("unrecognized spec object type: %s", route)
	}
	return &RootObjCmd{PhaseName: c.PhaseName, TypeName: route, Type: specType}, nil
}

func (c *RootPhaseCmd) Routes() []string {
	return spec_types.TypeNames(c.Types)
}

type RootObjCmd struct {
	PhaseName        string
	TypeName         string
	Type             spec_types.SpecType
	util.SpecOptions `ask:"."`
	Input            util.ObjInput `ask:"<input>" help:"Input, prefix with format, empty path for STDIN"`
	// TODO: path
}

func (c *RootObjCmd) Help() string {
	return fmt.Sprintf("Compute SSZ hash-tree-root of type %s (%s)", c.TypeName, c.PhaseName)
}

func (c *RootObjCmd) Run(ctx context.Context, args ...string) error {
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	objA := c.Type.Alloc(spec)
	if err := c.Input.Read(objA); err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}
	root := objA.HashTreeRoot(tree.GetHashFn())
	fmt.Println(root.String())
	return nil
}
