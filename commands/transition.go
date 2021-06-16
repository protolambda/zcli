package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/protolambda/ask"
	"github.com/protolambda/zcli/spec_types"
	"github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2/beacon"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/merge"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/zrnt/eth2/beacon/sharding"
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

func (c *TransitionCmd) Routes() []string {
	return spec_types.Phases
}

type TransitionSubCmd struct {
	PreFork string
}

func (c *TransitionSubCmd) Help() string {
	return fmt.Sprintf("Run state sub-processing (%s pre-state)", c.PreFork)
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

func (c *TransitionSubCmd) Routes() []string {
	return []string{"slots", "blocks", "sub"}
}

type TransitionSlotsCmd struct {
	PreFork          string
	Slots            uint64        `ask:"<slots>" help:"Number of slots to process"`
	Timeout          time.Duration `ask:"--timeout" help:"Timeout, e.g. 100ms"`
	util.SpecOptions `ask:"."`
	Pre              util.StateInput  `ask:"--pre" help:"Pre-state"`
	Post             util.StateOutput `ask:"--post" help:"Post-state"`
	// TODO: maybe fork-override, to transition between forks?
}

func (c *TransitionSlotsCmd) Help() string {
	return fmt.Sprintf("Process empty slots (%s pre-state)", c.PreFork)
}

func (c *TransitionSlotsCmd) Run(ctx context.Context, args ...string) error {
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	pre, err := c.Pre.Read(spec, c.PreFork)
	if err != nil {
		return err
	}
	state := &beacon.StandardUpgradeableBeaconState{BeaconState: pre}
	epc, err := common.NewEpochsContext(spec, state)
	if err != nil {
		return err
	}
	slot, err := state.Slot()
	if err != nil {
		return err
	}
	if err := common.ProcessSlots(ctx, spec, epc, state, slot + common.Slot(c.Slots)); err != nil {
		return err
	}
	return c.Post.Write(spec, state)
}

type TransitionBlocksCmd struct {
	PreFork          string
	VerifyStateRoot  bool          `ask:"--verify-state-root" help:"Verify the state root of each block"`
	Timeout          time.Duration `ask:"--timeout" help:"Timeout, e.g. 100ms"`
	util.SpecOptions `ask:"."`
	Pre              util.StateInput  `ask:"--pre" help:"Pre-state"`
	Post             util.StateOutput `ask:"--post" help:"Post-state"`
	// TODO: maybe fork-override, to transition between forks?
}

func (c *TransitionBlocksCmd) Help() string {
	return fmt.Sprintf("Process blocks (%s pre-state)", c.PreFork)
}

func (c *TransitionBlocksCmd) Run(ctx context.Context, args ...string) error {
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	pre, err := c.Pre.Read(spec, c.PreFork)
	if err != nil {
		return err
	}
	state := &beacon.StandardUpgradeableBeaconState{BeaconState: pre}
	epc, err := common.NewEpochsContext(spec, state)
	if err != nil {
		return err
	}
	genesisValRoot, err := state.GenesisValidatorsRoot()
	if err != nil {
		return err
	}
	phase := c.PreFork
	for i, arg := range args {
		var obj common.EnvelopeBuilder
		var digest common.ForkDigest
		switch phase {
		case "phase0":
			obj = new(phase0.SignedBeaconBlock)
			digest = common.ComputeForkDigest(spec.GENESIS_FORK_VERSION, genesisValRoot)
		case "altair":
			obj = new(altair.SignedBeaconBlock)
			digest = common.ComputeForkDigest(spec.ALTAIR_FORK_VERSION, genesisValRoot)
		case "merge":
			obj = new(merge.SignedBeaconBlock)
			digest = common.ComputeForkDigest(spec.MERGE_FORK_VERSION, genesisValRoot)
		case "sharding":
			obj = new(sharding.SignedBeaconBlock)
			digest = common.ComputeForkDigest(spec.SHARDING_FORK_VERSION, genesisValRoot)
		}
		input := util.ObjInput(arg)
		if err := input.Read(obj); err != nil {
			return fmt.Errorf("failed to read block %d: %v", i, err)
		}
		benv := obj.Envelope(spec, digest)
		if err := common.StateTransition(ctx, spec, epc, state, benv, c.VerifyStateRoot); err != nil {
			return fmt.Errorf("failed to process block %d: %v", i, err)
		}
	}
	return c.Post.Write(spec, state)
}

type TransitionSubRouterCmd struct {
	PreFork string
}

func (c *TransitionSubRouterCmd) Help() string {
	return fmt.Sprintf("Run state-transition sub-process (%s pre-state)", c.PreFork)
}

func (c *TransitionEpochSubCmd) Routes() []string {
	return append(append([]string{}, epochSubProcessingByPhase[c.PreFork]...), blockOpSubProcessingByPhase[c.PreFork]...)
}

func (c *TransitionSubRouterCmd) Cmd(route string) (cmd interface{}, err error) {
	if checkAny(epochSubProcessingByPhase[c.PreFork], route) {
		return &TransitionEpochSubCmd{PreFork: c.PreFork, Transition: route}, nil
	}
	if checkAny(blockOpSubProcessingByPhase[c.PreFork], route) {
		return &TransitionBlockSubCmd{PreFork: c.PreFork, Transition: route}, nil
	}
	return nil, ask.UnrecognizedErr
}

func checkAny(hay []string, needle string) bool {
	for _, h := range hay {
		if h == needle {
			return true
		}
	}
	return false
}

type TransitionEpochSubCmd struct {
	PreFork          string
	Transition       string
	Timeout          time.Duration `ask:"--timeout" help:"Timeout, e.g. 100ms"`
	util.SpecOptions `ask:"."`
	Pre              util.StateInput  `ask:"--pre" help:"Pre-state"`
	Post             util.StateOutput `ask:"--post" help:"Post-state"`
}

func (c *TransitionEpochSubCmd) Help() string {
	return fmt.Sprintf("Run epoch-sub-process %s (%s pre-state)", c.Transition, c.PreFork)
}

func (c *TransitionEpochSubCmd) Run(ctx context.Context, args ...string) error {
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	state, err := c.Pre.Read(spec, c.PreFork)
	if err != nil {
		return err
	}
	epc, err := common.NewEpochsContext(spec, state)
	if err != nil {
		return err
	}
	vals, err := state.Validators()
	if err != nil {
		return err
	}
	flats, err := common.FlattenValidators(vals)
	if err != nil {
		return err
	}
	maybeOutput := func(err error) error {
		if err != nil {
			return err
		}
		return c.Post.Write(spec, state)
	}
	switch c.Transition {
	case "pending_shard_confirmations":
		if c.PreFork != "sharding" {
			return errors.New("process_pending_shard_confirmations is only in Sharding")
		}
		return maybeOutput(sharding.ProcessPendingShardConfirmations(ctx, spec, state.(*sharding.BeaconStateView)))
	case "charge_confirmed_shard_fees":
		if c.PreFork != "sharding" {
			return errors.New("charge_confirmed_shard_fees is only in Sharding")
		}
		return maybeOutput(sharding.ChargeConfirmedShardFees(ctx, spec, epc, state.(*sharding.BeaconStateView)))
	case "reset_pending_shard_work":
		if c.PreFork != "sharding" {
			return errors.New("reset_pending_shard_work is only in Sharding")
		}
		return maybeOutput(sharding.ResetPendingShardWork(ctx, spec, epc, state.(*sharding.BeaconStateView)))
	case "justification_and_finalization", "inactivity_updates", "rewards_and_penalties":
		switch c.PreFork {
		case "phase0", "merge", "sharding":
			attesterData, err := phase0.ComputeEpochAttesterData(ctx, spec, epc, flats, state.(phase0.Phase0PendingAttestationsBeaconState))
			if err != nil {
				return err
			}
			switch c.Transition {
			case "justification_and_finalization":
				just := phase0.JustificationStakeData{
					CurrentEpoch:                  epc.CurrentEpoch.Epoch,
					TotalActiveStake:              epc.TotalActiveStake,
					PrevEpochUnslashedTargetStake: attesterData.PrevEpochUnslashedStake.TargetStake,
					CurrEpochUnslashedTargetStake: attesterData.CurrEpochUnslashedTargetStake,
				}
				return maybeOutput(phase0.ProcessEpochJustification(ctx, spec, &just, state))
			case "inactivity_updates":
				return errors.New("inactivity_updates only runs in Altair")
			case "rewards_and_penalties":
				return maybeOutput(phase0.ProcessEpochRewardsAndPenalties(ctx, spec, epc, attesterData, state.(phase0.BalancesBeaconState)))
			}
		case "altair":
			attesterData, err := altair.ComputeEpochAttesterData(ctx, spec, epc, flats, state.(*altair.BeaconStateView))
			if err != nil {
				return err
			}
			switch c.Transition {
			case "justification_and_finalization":
				just := phase0.JustificationStakeData{
					CurrentEpoch:                  epc.CurrentEpoch.Epoch,
					TotalActiveStake:              epc.TotalActiveStake,
					PrevEpochUnslashedTargetStake: attesterData.PrevEpochUnslashedStake.TargetStake,
					CurrEpochUnslashedTargetStake: attesterData.CurrEpochUnslashedTargetStake,
				}
				return maybeOutput(phase0.ProcessEpochJustification(ctx, spec, &just, state))
			case "inactivity_updates":
				return maybeOutput(altair.ProcessInactivityUpdates(ctx, spec, attesterData, state.(*altair.BeaconStateView)))
			case "rewards_and_penalties":
				return maybeOutput(altair.ProcessEpochRewardsAndPenalties(ctx, spec, epc, attesterData, state.(*altair.BeaconStateView)))
			}
		}
	case "registry_updates":
		return phase0.ProcessEpochRegistryUpdates(ctx, spec, epc, flats, state)
	case "slashings":
		return phase0.ProcessEpochSlashings(ctx, spec, epc, flats, state)
	case "final_updates": // legacy combination of below processes
		if c.PreFork != "phase0" {
			return errors.New("final_updates is a legacy combination of multiple processing functions, only available in phase0")
		}
		if err := phase0.ProcessEth1DataReset(ctx, spec, epc, state); err != nil {
			return err
		}
		if err := phase0.ProcessEffectiveBalanceUpdates(ctx, spec, epc, flats, state); err != nil {
			return err
		}
		if err := phase0.ProcessSlashingsReset(ctx, spec, epc, state); err != nil {
			return err
		}
		if err := phase0.ProcessRandaoMixesReset(ctx, spec, epc, state); err != nil {
			return err
		}
		if err := phase0.ProcessHistoricalRootsUpdate(ctx, spec, epc, state); err != nil {
			return err
		}
		return maybeOutput(phase0.ProcessParticipationRecordUpdates(ctx, spec, epc, state.(phase0.Phase0PendingAttestationsBeaconState)))
	case "eth1_data_reset":
		return maybeOutput(phase0.ProcessEth1DataReset(ctx, spec, epc, state))
	case "effective_balance_updates":
		return maybeOutput(phase0.ProcessEffectiveBalanceUpdates(ctx, spec, epc, flats, state))
	case "slashings_reset":
		return maybeOutput(phase0.ProcessSlashingsReset(ctx, spec, epc, state))
	case "randao_mixes_reset":
		return maybeOutput(phase0.ProcessRandaoMixesReset(ctx, spec, epc, state))
	case "historical_roots_update":
		return maybeOutput(phase0.ProcessHistoricalRootsUpdate(ctx, spec, epc, state))
	case "participation_record_updates":
		if c.PreFork == "altair" {
			return errors.New("participation_record_updates was removed in Altair")
		}
		return maybeOutput(phase0.ProcessParticipationRecordUpdates(ctx, spec, epc, state.(phase0.Phase0PendingAttestationsBeaconState)))
	case "participation_flag_updates":
		if c.PreFork != "altair" {
			return errors.New("participation_flag_updates is only in Altair")
		}
		return maybeOutput(altair.ProcessParticipationFlagUpdates(ctx, spec, state.(*altair.BeaconStateView)))
	case "sync_committee_updates":
		if c.PreFork != "altair" {
			return errors.New("sync_committee_updates is only in Altair")
		}
		return maybeOutput(altair.ProcessSyncCommitteeUpdates(ctx, spec, epc, state.(*altair.BeaconStateView)))
	case "shard_epoch_increment":
		if c.PreFork != "sharding" {
			return errors.New("process_shard_epoch_increment is only in Sharding")
		}
		return maybeOutput(sharding.ProcessShardEpochIncrement(ctx, spec, epc, state.(*sharding.BeaconStateView)))
	}
	return ask.UnrecognizedErr
}

type TransitionBlockSubCmd struct {
	PreFork          string
	Transition       string
	Timeout          time.Duration `ask:"--timeout" help:"Timeout, e.g. 100ms"`
	util.SpecOptions `ask:"."`
	Pre              util.StateInput  `ask:"--pre" help:"Pre-state"`
	Op               util.ObjInput    `ask:"<op>" help:"Block operation input"`
	Post             util.StateOutput `ask:"--post" help:"Post-state"`
}

func (c *TransitionBlockSubCmd) Help() string {
	return fmt.Sprintf("Run block-sub-process %s (%s pre-state)", c.Transition, c.PreFork)
}

func (c *TransitionBlockSubCmd) Run(ctx context.Context, args ...string) error {
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	state, err := c.Pre.Read(spec, c.PreFork)
	if err != nil {
		return err
	}
	epc, err := common.NewEpochsContext(spec, state)
	if err != nil {
		return err
	}
	maybeOutput := func(err error) error {
		if err != nil {
			return err
		}
		return c.Post.Write(spec, state)
	}
	switch c.Transition {
	case "block_header":
		slot, err := state.Slot()
		if err != nil {
			return err
		}
		proposerIndex, err := epc.GetBeaconProposer(slot)
		if err != nil {
			return err
		}
		var header common.BeaconBlockHeader
		if err := c.Op.Read(&header); err != nil {
			return err
		}
		return maybeOutput(common.ProcessHeader(ctx, spec, state, &header, proposerIndex))
	case "randao":
		var reveal common.BLSSignature
		if err := c.Op.Read(&reveal); err != nil {
			return err
		}
		return maybeOutput(phase0.ProcessRandaoReveal(ctx, spec, epc, state, reveal))
	case "eth1_data":
		var eth1Data common.Eth1Data
		if err := c.Op.Read(&eth1Data); err != nil {
			return err
		}
		return maybeOutput(phase0.ProcessEth1Vote(ctx, spec, epc, state, eth1Data))
	case "proposer_slashing":
		var propSl phase0.ProposerSlashing
		if err := c.Op.Read(&propSl); err != nil {
			return err
		}
		return maybeOutput(phase0.ProcessProposerSlashing(spec, epc, state, &propSl))
	case "attester_slashing":
		switch c.PreFork {
		case "phase0", "altair", "merge":
			var attSl phase0.AttesterSlashing
			if err := c.Op.Read(&attSl); err != nil {
				return err
			}
			return maybeOutput(phase0.ProcessAttesterSlashing(spec, epc, state, &attSl))
		case "sharding":
			var attSl sharding.AttesterSlashing
			if err := c.Op.Read(&attSl); err != nil {
				return err
			}
			return maybeOutput(sharding.ProcessAttesterSlashing(spec, epc, state, &attSl))
		}
	case "attestation":
		switch c.PreFork {
		case "phase0", "merge":
			var att phase0.Attestation
			if err := c.Op.Read(&att); err != nil {
				return err
			}
			return maybeOutput(phase0.ProcessAttestation(spec, epc, state.(phase0.Phase0PendingAttestationsBeaconState), &att))
		case "altair":
			var att phase0.Attestation
			if err := c.Op.Read(&att); err != nil {
				return err
			}
			return maybeOutput(altair.ProcessAttestation(spec, epc, state.(*altair.BeaconStateView), &att))
		case "sharding":
			var att sharding.Attestation
			if err := c.Op.Read(&att); err != nil {
				return err
			}
			return maybeOutput(sharding.ProcessAttestation(spec, epc, state.(*sharding.BeaconStateView), &att))
		}
	case "deposit":
		var dep common.Deposit
		if err := c.Op.Read(&dep); err != nil {
			return err
		}
		return maybeOutput(phase0.ProcessDeposit(spec, epc, state, &dep, false))
	case "voluntary_exit":
		var exit phase0.SignedVoluntaryExit
		if err := c.Op.Read(&exit); err != nil {
			return err
		}
		return maybeOutput(phase0.ProcessVoluntaryExit(spec, epc, state, &exit))
	case "sync_aggregate":
		if c.PreFork != "altair" {
			return fmt.Errorf("fork %s does not have sync_aggregate processing", c.PreFork)
		}
		var agg altair.SyncAggregate
		if err := c.Op.Read(&agg); err != nil {
			return err
		}
		return maybeOutput(altair.ProcessSyncAggregate(ctx, spec, epc, state.(*altair.BeaconStateView), &agg))
	case "execution_payload":
		switch c.PreFork {
		case "phase0", "altair":
			return fmt.Errorf("fork %s does not have execution_payload processing", c.PreFork)
		case "merge", "sharding":
			var payload common.ExecutionPayload
			if err := c.Op.Read(&payload); err != nil {
				return err
			}
			return maybeOutput(merge.ProcessExecutionPayload(ctx, spec, state.(merge.ExecutionTrackingBeaconState),
				&payload, new(NoOpExecutionEngine)))
		}
	case "shard_proposer_slashing":
		if c.PreFork != "sharding" {
			return fmt.Errorf("fork %s does not have shard_proposer_slashing processing", c.PreFork)
		}
		var sl sharding.ShardProposerSlashing
		if err := c.Op.Read(&sl); err != nil {
			return err
		}
		return maybeOutput(sharding.ProcessShardProposerSlashing(spec, epc, state, &sl))
	case "shard_header":
		if c.PreFork != "sharding" {
			return fmt.Errorf("fork %s does not have shard_header processing", c.PreFork)
		}
		var h sharding.SignedShardBlobHeader
		if err := c.Op.Read(&h); err != nil {
			return err
		}
		return maybeOutput(sharding.ProcessShardHeader(spec, epc, state.(*sharding.BeaconStateView), &h))
	}
	return ask.UnrecognizedErr
}

type NoOpExecutionEngine struct{}

func (m *NoOpExecutionEngine) NewBlock(ctx context.Context, executionPayload *common.ExecutionPayload) (success bool, err error) {
	return true, nil
}

var _ common.ExecutionEngine = (*NoOpExecutionEngine)(nil)

var epochSubProcessingByPhase = map[string][]string{
	"phase0": {
		"justification_and_finalization",
		"rewards_and_penalties",
		"registry_updates",
		"slashings",
		"eth1_data_reset",
		"effective_balance_updates",
		"slashings_reset",
		"randao_mixes_reset",
		"historical_roots_update",
		"participation_record_updates",
	},
	"altair": {
		"justification_and_finalization",
		"inactivity_updates",
		"rewards_and_penalties",
		"registry_updates",
		"slashings",
		"eth1_data_reset",
		"effective_balance_updates",
		"slashings_reset",
		"randao_mixes_reset",
		"historical_roots_update",
		"participation_flag_updates",
		"sync_committee_updates",
	},
	"merge": {
		"justification_and_finalization",
		"rewards_and_penalties",
		"registry_updates",
		"slashings",
		"eth1_data_reset",
		"effective_balance_updates",
		"slashings_reset",
		"randao_mixes_reset",
		"historical_roots_update",
		"participation_record_updates",
	},
	"sharding": {
		"pending_shard_confirmations",
		"charge_confirmed_shard_fees",
		"reset_pending_shard_work",
		"justification_and_finalization",
		"rewards_and_penalties",
		"registry_updates",
		"slashings",
		"eth1_data_reset",
		"effective_balance_updates",
		"slashings_reset",
		"randao_mixes_reset",
		"historical_roots_update",
		"participation_record_updates",
		"shard_epoch_increment",
	},
}

var blockOpSubProcessingByPhase = map[string][]string{
	"phase0": {
		"block_header",
		"randao",
		"eth1_data",
		"proposer_slashing",
		"attester_slashing",
		"attestation",
		"deposit",
		"voluntary_exit",
	},
	"altair": {
		"block_header",
		"randao",
		"eth1_data",
		"proposer_slashing",
		"attester_slashing",
		"attestation",
		"deposit",
		"voluntary_exit",
		"sync_aggregate",
	},
	"merge": {
		"block_header",
		"randao",
		"eth1_data",
		"proposer_slashing",
		"attester_slashing",
		"attestation",
		"deposit",
		"voluntary_exit",
		"execution_payload",
	},
	"sharding": {
		"block_header",
		"randao",
		"eth1_data",
		"proposer_slashing",
		"attester_slashing",
		"shard_proposer_slashing",
		"shard_header",
		"attestation",
		"deposit",
		"voluntary_exit",
	},
}
