package commands

import (
	"context"
	"fmt"
	"github.com/protolambda/ask"
	"github.com/protolambda/zcli/spec_types"
	"github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/configs"
)

type MetaCmd struct{}

func (c *MetaCmd) Help() string {
	return "List metadata of beacon state"
}

func (c *MetaCmd) Cmd(route string) (cmd interface{}, err error) {
	switch route {
	case "phase0", "altair", "bellatrix", "capella":
		return &MetaPhaseCmd{Phase: route}, nil
	}
	return nil, ask.UnrecognizedErr
}

func (c *MetaCmd) Routes() []string {
	return spec_types.Phases
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
		return &CommitteesCmd{Phase: c.Phase}, nil
	case "proposers":
		return &ProposersCmd{Phase: c.Phase}, nil
	case "sync_committees", "sync-committees", "synccommittees":
		if c.Phase != "altair" {
			return nil, ask.UnrecognizedErr
		}
		return &SyncCommitteesCmd{Phase: c.Phase}, nil
	}
	return nil, ask.UnrecognizedErr
}

func (c *MetaPhaseCmd) Routes() []string {
	out := []string{"committees", "proposers"}
	if c.Phase == "altair" {
		out = append(out, "sync-committees")
	}
	return out
}

type CommitteesCmd struct {
	Phase               string
	configs.SpecOptions `ask:"."`
	State               util.StateInput `ask:"<state>" help:"BeaconState, prefix with format, empty path for STDIN"`
}

func (c *CommitteesCmd) Help() string {
	return fmt.Sprintf("List previous/current/next beacon committees (phase %s)", c.Phase)
}

func (c *CommitteesCmd) Run(ctx context.Context, args ...string) error {
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	state, err := c.State.Read(spec, c.Phase)
	if err != nil {
		return err
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
		end, err := spec.EpochStartSlot(epoch + 1)
		if err != nil {
			return err
		}
		for slot := start; slot < end; slot++ {
			for i := common.CommitteeIndex(0); i < common.CommitteeIndex(committeesPerSlot); i++ {
				committee, err := epc.GetBeaconCommittee(slot, i)
				if err != nil {
					return fmt.Errorf("cannot get committee for slot %d committee index %d", slot, i)
				}
				fmt.Printf("epoch: %7d    slot: %9d    committee index: %4d (out of %2d)   size: %3d    indices: %v\n",
					spec.SlotToEpoch(slot), slot, i, committeesPerSlot, len(committee), committee)
			}
		}
	}
	return nil
}

type ProposersCmd struct {
	Phase               string
	configs.SpecOptions `ask:"."`
	State               util.StateInput `ask:"<state>" help:"BeaconState, prefix with format, empty path for STDIN"`
}

func (c *ProposersCmd) Help() string {
	return fmt.Sprintf("List current beacon propoosers (phase %s)", c.Phase)
}

func (c *ProposersCmd) Run(ctx context.Context, args ...string) error {
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	state, err := c.State.Read(spec, c.Phase)
	if err != nil {
		return err
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
		fmt.Printf("epoch: %7d    slot: %9d    proposer index: %4d\n", spec.SlotToEpoch(slot), slot, proposerIndex)
	}
	return nil
}

type SyncCommitteesCmd struct {
	Phase               string
	configs.SpecOptions `ask:"."`
	State               util.StateInput `ask:"<state>" help:"BeaconState, prefix with format, empty path for STDIN"`
}

func (c *SyncCommitteesCmd) Help() string {
	return fmt.Sprintf("List current/next sync-committee members (phase %s)", c.Phase)
}

func (c *SyncCommitteesCmd) Run(ctx context.Context, args ...string) error {
	if c.Phase != "altair" {
		return fmt.Errorf("only Altair is supported for looking up the sync committee indices, not %q", c.Phase)
	}
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	state, err := c.State.Read(spec, c.Phase)
	if err != nil {
		return err
	}
	epc, err := common.NewEpochsContext(spec, state)
	if err != nil {
		return fmt.Errorf("cannot compute state epochs context: %v", err)
	}
	fmt.Println("--- current sync committee ---")
	for i, vi := range epc.CurrentSyncCommittee.Indices {
		fmt.Printf("current[%4d] => %4d: %s\n", i, vi, epc.CurrentSyncCommittee.CachedPubkeys[i].Compressed.String())
	}
	fmt.Println("--- next sync committee ---")
	for i, vi := range epc.NextSyncCommittee.Indices {
		fmt.Printf("next[%4d] => %4d: %s\n", i, vi, epc.NextSyncCommittee.CachedPubkeys[i].Compressed.String())
	}
	return nil
}
