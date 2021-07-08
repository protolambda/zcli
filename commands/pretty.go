package commands

import (
	"context"
	"fmt"
	"github.com/protolambda/zcli/spec_types"
	"github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2/configs"
)

type PrettyCmd struct{}

func (c *PrettyCmd) Help() string {
	return "Pretty-print spec object (output indented JSON)"
}

func (c *PrettyCmd) Cmd(route string) (cmd interface{}, err error) {
	phaseTypes, ok := spec_types.TypesByPhase[route]
	if !ok {
		return nil, fmt.Errorf("unrecognized phase: %s", route)
	}
	return &PrettyPhaseCmd{PhaseName: route, Types: phaseTypes}, nil
}

func (c *PrettyCmd) Routes() []string {
	return spec_types.Phases
}

type PrettyPhaseCmd struct {
	PhaseName string
	Types     map[string]spec_types.SpecType
}

func (c *PrettyPhaseCmd) Cmd(route string) (cmd interface{}, err error) {
	specType, ok := c.Types[route]
	if !ok {
		return nil, fmt.Errorf("unrecognized spec object type: %s", route)
	}
	return &PrettyObjCmd{PhaseName: c.PhaseName, TypeName: route, Type: specType}, nil
}

func (c *PrettyPhaseCmd) Routes() []string {
	return spec_types.TypeNames(c.Types)
}

type PrettyObjCmd struct {
	PhaseName           string
	TypeName            string
	Type                spec_types.SpecType
	configs.SpecOptions `ask:"."`
	Input               util.ObjInput `ask:"<input>" help:"Input path, prefix with format, empty path for STDIN"`
	Output              string        `ask:"[output]" help:"Output path, empty path for STDOUT"`
}

func (c *PrettyObjCmd) Run(ctx context.Context, args ...string) error {
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	obj := c.Type.Alloc(spec)
	if err := c.Input.Read(obj); err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}
	out := util.ObjOutput("pretty:" + c.Output)
	if err := out.Write(obj); err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}
	return nil
}
