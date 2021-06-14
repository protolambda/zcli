package commands

import (
	"context"
	"fmt"
	"github.com/protolambda/ask"
	"github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2/beacon/common"
)

type MetaCmd struct{}

func (c *MetaCmd) Help() string {
	return "List metadata of beacon state"
}

func (c *MetaCmd) Cmd(route string) (cmd interface{}, err error) {
	switch route {
	case "phase0", "altair", "merge", "sharding":
		return &MetaPhaseCmd{Phase: route}, nil
	}
	return nil, ask.UnrecognizedErr
}

type MetaPhaseCmd struct {
	Phase string
}

func (c *MetaPhaseCmd) Help() string {
	return fmt.Sprintf("List metadata of beacon state (phase %s)", c.Phase)
}

func (c *MetaPhaseCmd) Cmd(route string) (cmd interface{}, err error) {
	switch route {
	case "committees":
		return &CommitteesCmd{Phase: route}, nil
	case "proposers":
		return &ProposersCmd{Phase: route}, nil
	case "sync_committees", "sync-committees", "synccommittees":
		return &SyncCommitteesCmd{Phase: route}, nil
	}
	return nil, ask.UnrecognizedErr
}

type CommitteesCmd struct {
	Phase       string
	SpecOptions `ask:"."`
	State       util.ObjInput `ask:"<state>" help:"BeaconState, prefix with format, empty path for STDIN"`
}

func (c *CommitteesCmd) Run(ctx context.Context, args ...string) error {
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	// TODO state
	var state common.BeaconState
	switch c.Phase {
	case "phase0":
		// TODO
		//c.State.Read()
	}
	epc, err := common.NewEpochsContext(spec, state)
	if err != nil {
		return fmt.Errorf("cannot compute state epochs context: %v", err)
	}
	currentEpoch := epc.CurrentEpoch.Epoch

	for epoch := currentEpoch.Previous(); epoch <= currentEpoch+1; epoch++ {
		committeesPerSlot, err := epc.GetCommitteeCountPerSlot(epoch)
		if err != nil {
			return err
		}
		start, err := spec.EpochStartSlot(epoch)
		if err != nil {
			return err
		}
		end, err := spec.EpochStartSlot(epoch)
		if err != nil {
			return err
		}
		for slot := start; slot < end; slot++ {
			for i := common.CommitteeIndex(0); i < common.CommitteeIndex(committeesPerSlot); i++ {
				committee, err := epc.GetBeaconCommittee(slot, i)
				if err != nil {
					return fmt.Errorf("cannot get committee for slot %d committee index %d", slot, i)
				}
				fmt.Printf(`epoch: %7d    slot: %9d    committee index: %4d (out of %4d)   size: %5d    indices: %v\n`,
					spec.SlotToEpoch(slot), slot, i, committeesPerSlot, len(committee), committee)
			}
		}
	}
	return nil
}

type ProposersCmd struct {
	Phase       string
	SpecOptions `ask:"."`
	State       util.ObjInput `ask:"<state>" help:"BeaconState, prefix with format, empty path for STDIN"`
}

func (c *ProposersCmd) Run(ctx context.Context, args ...string) error {
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	// TODO state
	var state common.BeaconState
	switch c.Phase {
	case "phase0":
		// TODO
		//c.State.Read()
	}
	epc, err := common.NewEpochsContext(spec, state)
	if err != nil {
		return fmt.Errorf("cannot compute state epochs context: %v", err)
	}
	currentEpoch := epc.CurrentEpoch.Epoch
	start, err := spec.EpochStartSlot(currentEpoch)
	if err != nil {
		return err
	}
	end, err := spec.EpochStartSlot(currentEpoch + 1)
	if err != nil {
		return err
	}
	for slot := start; slot < end; slot++ {
		proposerIndex, err := epc.GetBeaconProposer(slot)
		if err != nil {
			return fmt.Errorf("cannot compute proposer index for slot %d: %v", slot, err)
		}
		fmt.Printf(`epoch: %7d    slot: %9d    proposer index: %4d\n`, spec.SlotToEpoch(slot), slot, proposerIndex)
	}
	return nil
}

type SyncCommitteesCmd struct {
	Phase       string
	SpecOptions `ask:"."`
	State       util.ObjInput `ask:"<state>" help:"BeaconState, prefix with format, empty path for STDIN"`
}

func (c *SyncCommitteesCmd) Run(ctx context.Context, args ...string) error {
	// TODO
	return nil
}
