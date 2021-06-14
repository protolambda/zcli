package commands

import (
	"context"
	"github.com/protolambda/ask"
	"github.com/protolambda/zcli/util"
	"time"
)

type TransitionCmd struct{}

func (c *TransitionCmd) Help() string {
	return "Run state transitions and sub-processes"
}

func (c *TransitionCmd) Cmd(route string) (cmd interface{}, err error) {
	switch route {
	case "phase0", "altair", "merge", "sharding":
		return &TransitionSubCmd{PreFork: route}, nil
	default:
		return nil, ask.UnrecognizedErr
	}
}

type TransitionSubCmd struct {
	PreFork string
}

func (c *TransitionSubCmd) Cmd(route string) (cmd interface{}, err error) {
	switch route {
	case "slots":
		return &TransitionSlotsCmd{PreFork: c.PreFork}, nil
	case "blocks":
		return &TransitionBlocksCmd{PreFork: c.PreFork}, nil
	case "sub":
		return &TransitionSubRouterCmd{PreFork: c.PreFork}, nil
	}
	return nil, ask.UnrecognizedErr
}

type TransitionSlotsCmd struct {
	PreFork          string
	Slots            uint64        `ask:"<slots>" help:"Number of slots to process"`
	Timeout          time.Duration `ask:"--timeout" help:"Timeout, e.g. 100ms"`
	util.SpecOptions `ask:"."`
	Pre              util.ObjInput `ask:"--pre" help:"Pre-state"`
	Post             util.ObjInput `ask:"--post" help:"Post-state"`
	// TODO: maybe fork-override, to transition between forks?
}

type TransitionBlocksCmd struct {
	PreFork          string
	VerifyStateRoot  bool          `ask:"--verify-state-root" help:"Verify the state root of each block"`
	Timeout          time.Duration `ask:"--timeout" help:"Timeout, e.g. 100ms"`
	util.SpecOptions `ask:"."`
	Pre              util.ObjInput `ask:"--pre" help:"Pre-state"`
	Post             util.ObjInput `ask:"--post" help:"Post-state"`
	// TODO: maybe fork-override, to transition between forks?
}

func (c *TransitionBlocksCmd) Run(ctx context.Context, args ...string) error {
	// 1 block per arg
	// TODO
	return nil
}

type TransitionSubRouterCmd struct {
	PreFork string
}

func (c *TransitionSubRouterCmd) Cmd(route string) (cmd interface{}, err error) {
	// TODO
	switch c.PreFork {
	case "phase0":
		switch route {
		case "justification_and_finalization":
		case "rewards_and_penalties":
		case "registry_updates":
		case "slashings":
		case "final_updates": // legacy combination of below processes
		case "effective_balance_updates":
		case "slashings_reset":
		case "randao_mixes_reset":
		case "historical_roots_update":
		case "participation_record_updates":
		case "block_header":
		case "randao":
		case "eth1_data":
		case "proposer_slashings":
		case "proposer_slashing":
		case "attester_slashings":
		case "attester_slashing":
		case "attestations":
		case "attestation":
		case "deposits":
		case "deposit":
		case "voluntary_exits":
		case "voluntary_exit":
		}
		// TODO more forks
	}

	return nil, ask.UnrecognizedErr
}

type TransitionEpochSubCmd struct {
	PreFork          string
	Transition       string
	Timeout          time.Duration `ask:"--timeout" help:"Timeout, e.g. 100ms"`
	util.SpecOptions `ask:"."`
	Pre              util.ObjInput `ask:"--pre" help:"Pre-state"`
	Post             util.ObjInput `ask:"--post" help:"Post-state"`
}

func (c *TransitionEpochSubCmd) Run(ctx context.Context, args ...string) error {
	// TODO
	return nil
}

type TransitionBlockSubCmd struct {
	PreFork          string
	Transition       string
	Timeout          time.Duration `ask:"--timeout" help:"Timeout, e.g. 100ms"`
	util.SpecOptions `ask:"."`
	Pre              util.ObjInput `ask:"--pre" help:"Pre-state"`
	Op               util.ObjInput `ask:"<op>" help:"Block operation input"`
	Post             util.ObjInput `ask:"--post" help:"Post-state"`
}

func (c *TransitionBlockSubCmd) Run(ctx context.Context, args ...string) error {
	// TODO
	return nil
}

type TransitionBlockOpsSubCmd struct {
	PreFork          string
	Transition       string
	Timeout          time.Duration `ask:"--timeout" help:"Timeout, e.g. 100ms"`
	util.SpecOptions `ask:"."`
	Pre              util.ObjInput `ask:"--pre" help:"Pre-state"`
	Post             util.ObjInput `ask:"--post" help:"Post-state"`
}

func (c *TransitionBlockOpsSubCmd) Run(ctx context.Context, args ...string) error {
	// 1 op per arg
	// TODO
	return nil
}
