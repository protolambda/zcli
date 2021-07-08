package commands

import (
	"context"
	"fmt"
	"github.com/protolambda/zcli/spec_types"
	"github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2/configs"
)

type ConvertCmd struct{}

func (c *ConvertCmd) Help() string {
	return "Convert spec object from one format to another"
}

func (c *ConvertCmd) Cmd(route string) (cmd interface{}, err error) {
	phaseTypes, ok := spec_types.TypesByPhase[route]
	if !ok {
		return nil, fmt.Errorf("unrecognized phase: %s", route)
	}
	return &ConvertPhaseCmd{PhaseName: route, Types: phaseTypes}, nil
}

func (c *ConvertCmd) Routes() []string {
	return spec_types.Phases
}

type ConvertPhaseCmd struct {
	PhaseName string
	Types     map[string]spec_types.SpecType
}

func (c *ConvertPhaseCmd) Cmd(route string) (cmd interface{}, err error) {
	specType, ok := c.Types[route]
	if !ok {
		return nil, fmt.Errorf("unrecognized spec object type: %s", route)
	}
	return &ConvertObjCmd{PhaseName: c.PhaseName, TypeName: route, Type: specType}, nil
}

func (c *ConvertPhaseCmd) Routes() []string {
	return spec_types.TypeNames(c.Types)
}

type ConvertObjCmd struct {
	PhaseName           string
	TypeName            string
	Type                spec_types.SpecType
	configs.SpecOptions `ask:"."`
	Input               util.ObjInput  `ask:"<input>" help:"Input path, prefix with format, empty path for STDIN"`
	Output              util.ObjOutput `ask:"<output>" help:"Output path, prefix with format, empty path for STDOUT"`
}

func (c *ConvertObjCmd) Run(ctx context.Context, args ...string) error {
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	obj := c.Type.Alloc(spec)
	if err := c.Input.Read(obj); err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}
	if err := c.Output.Write(obj); err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}
	return nil
}
