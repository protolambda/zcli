package commands

import (
	"context"
	"fmt"
	"github.com/protolambda/messagediff"
	"github.com/protolambda/zcli/spec_types"
	"github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2/configs"
)

type DiffCmd struct{}

func (c *DiffCmd) Help() string {
	return "Diff spec data"
}

func (c *DiffCmd) Cmd(route string) (cmd interface{}, err error) {
	phaseTypes, ok := spec_types.TypesByPhase[route]
	if !ok {
		return nil, fmt.Errorf("unrecognized phase: %s", route)
	}
	return &DiffPhaseCmd{PhaseName: route, Types: phaseTypes}, nil
}

func (c *DiffCmd) Routes() []string {
	return spec_types.Phases
}

type DiffPhaseCmd struct {
	PhaseName string
	Types     map[string]spec_types.SpecType
}

func (c *DiffPhaseCmd) Cmd(route string) (cmd interface{}, err error) {
	specType, ok := c.Types[route]
	if !ok {
		return nil, fmt.Errorf("unrecognized spec object type: %s", route)
	}
	return &DiffObjCmd{PhaseName: c.PhaseName, TypeName: route, Type: specType}, nil
}

func (c *DiffPhaseCmd) Routes() []string {
	return spec_types.TypeNames(c.Types)
}

type DiffObjCmd struct {
	PhaseName           string
	TypeName            string
	Type                spec_types.SpecType
	configs.SpecOptions `ask:"."`
	InputA              util.ObjInput `ask:"<a>" help:"Input A, prefix with format, empty path for STDIN"`
	InputB              util.ObjInput `ask:"<b>" help:"Input B, prefix with format, empty path for STDIN"`
}

func (c *DiffObjCmd) Run(ctx context.Context, args ...string) error {
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	objA := c.Type.Alloc(spec)
	if err := c.InputA.Read(objA); err != nil {
		return fmt.Errorf("failed to read input A: %v", err)
	}
	objB := c.Type.Alloc(spec)
	if err := c.InputB.Read(objB); err != nil {
		return fmt.Errorf("failed to read input B: %v", err)
	}
	if diff, equal := messagediff.PrettyDiff(objA, objB, messagediff.SliceWeakEmptyOption{}); equal {
		fmt.Errorf("%s (%s) objects A and B are equal\n", c.TypeName, c.PhaseName)
	} else {
		fmt.Printf("%s (%s) objects A and B are different:\n%s\n", c.TypeName, c.PhaseName, diff)
	}
	return nil
}
