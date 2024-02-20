package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/protolambda/ask"
	"github.com/protolambda/zrnt/eth2/beacon"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/bellatrix"
	"github.com/protolambda/zrnt/eth2/beacon/capella"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/deneb"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/zrnt/eth2/configs"
	"github.com/protolambda/zrnt/eth2/execution"

	"github.com/protolambda/zcli/spec_types"
	"github.com/protolambda/zcli/util"
)

type TransitionCmd struct{}

func (c *TransitionCmd) Help() string {
	return "Run state transitions and sub-processes"
}

func (c *TransitionCmd) Cmd(route string) (cmd interface{}, err error) {
	switch route {
	case "phase0", "altair", "bellatrix", "capella", "deneb":
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
	case "epoch":
		return &TransitionEpochCmd{PreFork: c.PreFork}, nil
	case "blocks":
		return &TransitionBlocksCmd{PreFork: c.PreFork}, nil
	case "sub":
		return &TransitionSubRouterCmd{PreFork: c.PreFork}, nil
	}
	return nil, ask.UnrecognizedErr
}

func (c *TransitionSubCmd) Routes() []string {
	return []string{"slots", "epoch", "blocks", "sub"}
}

type TransitionEpochCmd struct {
	PreFork             string
	Timeout             time.Duration `ask:"--timeout" help:"Timeout, e.g. 100ms"`
	configs.SpecOptions `ask:"."`
	Pre                 util.StateInput  `ask:"--pre" help:"Pre-state"`
	Post                util.StateOutput `ask:"--post" help:"Post-state"`
}

func (c *TransitionEpochCmd) Help() string {
	return fmt.Sprintf("Process the epoch transition (%s pre-state), without any slot processing", c.PreFork)
}

func (c *TransitionEpochCmd) Run(ctx context.Context, args ...string) error {
	if c.Timeout != 0 {
		ctx, _ = context.WithTimeout(ctx, c.Timeout)
	}
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	pre, err := c.Pre.Read(spec, c.PreFork)
	if err != nil {
		return err
	}
	state := &beacon.StandardUpgradeableBeaconState{BeaconState: pre}
	epc, err := common.NewEpochsContext(spec, pre)
	if err != nil {
		return err
	}
	if err := state.ProcessEpoch(ctx, spec, epc); err != nil {
		return err
	}
	return c.Post.Write(spec, state)
}

type TransitionSlotsCmd struct {
	PreFork             string
	Slots               uint64        `ask:"<slots>" help:"Number of slots to process"`
	Timeout             time.Duration `ask:"--timeout" help:"Timeout, e.g. 100ms"`
	configs.SpecOptions `ask:"."`
	Pre                 util.StateInput  `ask:"--pre" help:"Pre-state"`
	Post                util.StateOutput `ask:"--post" help:"Post-state"`
	// TODO: maybe fork-override, to transition between forks?
}

func (c *TransitionSlotsCmd) Help() string {
	return fmt.Sprintf("Process empty slots (%s pre-state)", c.PreFork)
}

func (c *TransitionSlotsCmd) Run(ctx context.Context, args ...string) error {
	if c.Timeout != 0 {
		ctx, _ = context.WithTimeout(ctx, c.Timeout)
	}
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	pre, err := c.Pre.Read(spec, c.PreFork)
	if err != nil {
		return err
	}
	state := &beacon.StandardUpgradeableBeaconState{BeaconState: pre}
	epc, err := common.NewEpochsContext(spec, pre)
	if err != nil {
		return err
	}
	slot, err := state.Slot()
	if err != nil {
		return err
	}
	if err := common.ProcessSlots(ctx, spec, epc, state, slot+common.Slot(c.Slots)); err != nil {
		return err
	}
	return c.Post.Write(spec, state)
}

type TransitionBlocksCmd struct {
	PreFork             string
	VerifyStateRoot     bool          `ask:"--verify-state-root" help:"Verify the state root of each block"`
	Timeout             time.Duration `ask:"--timeout" help:"Timeout, e.g. 100ms"`
	configs.SpecOptions `ask:"."`
	Pre                 util.StateInput  `ask:"--pre" help:"Pre-state"`
	Post                util.StateOutput `ask:"--post" help:"Post-state"`
	// TODO: maybe fork-override, to transition between forks?
}

func (c *TransitionBlocksCmd) Help() string {
	return fmt.Sprintf("Process blocks (%s pre-state)", c.PreFork)
}

func (c *TransitionBlocksCmd) Run(ctx context.Context, args ...string) error {
	if c.Timeout != 0 {
		ctx, _ = context.WithTimeout(ctx, c.Timeout)
	}
	spec, err := c.Spec()
	if err != nil {
		return err
	}
	spec.ExecutionEngine = new(execution.NoOpExecutionEngine)
	pre, err := c.Pre.Read(spec, c.PreFork)
	if err != nil {
		return err
	}
	state := &beacon.StandardUpgradeableBeaconState{BeaconState: pre}
	epc, err := common.NewEpochsContext(spec, pre)
	if err != nil {
		return err
	}
	genesisValRoot, err := state.GenesisValidatorsRoot()
	if err != nil {
		return err
	}
	phase := c.PreFork
	for i, arg := range args {
		var obj interface {
			common.EnvelopeBuilder
			common.SpecObj
		}
		var digest common.ForkDigest
		switch phase {
		case "phase0":
			obj = new(phase0.SignedBeaconBlock)
			digest = common.ComputeForkDigest(spec.GENESIS_FORK_VERSION, genesisValRoot)
		case "altair":
			obj = new(altair.SignedBeaconBlock)
			digest = common.ComputeForkDigest(spec.ALTAIR_FORK_VERSION, genesisValRoot)
		case "bellatrix":
			obj = new(bellatrix.SignedBeaconBlock)
			digest = common.ComputeForkDigest(spec.BELLATRIX_FORK_VERSION, genesisValRoot)
		case "capella":
			obj = new(capella.SignedBeaconBlock)
			digest = common.ComputeForkDigest(spec.CAPELLA_FORK_VERSION, genesisValRoot)
		case "deneb":
			obj = new(deneb.SignedBeaconBlock)
			digest = common.ComputeForkDigest(spec.DENEB_FORK_VERSION, genesisValRoot)
		}
		input := util.ObjInput(arg)
		if err := input.Read(spec.Wrap(obj)); err != nil {
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

func (c *TransitionSubRouterCmd) Routes() []string {
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
	PreFork             string
	Transition          string
	Timeout             time.Duration `ask:"--timeout" help:"Timeout, e.g. 100ms"`
	configs.SpecOptions `ask:"."`
	Pre                 util.StateInput  `ask:"--pre" help:"Pre-state"`
	Post                util.StateOutput `ask:"--post" help:"Post-state"`
}

func (c *TransitionEpochSubCmd) Help() string {
	return fmt.Sprintf("Run epoch-sub-process %s (%s pre-state)", c.Transition, c.PreFork)
}

func (c *TransitionEpochSubCmd) Run(ctx context.Context, args ...string) error {
	if c.Timeout != 0 {
		ctx, _ = context.WithTimeout(ctx, c.Timeout)
	}
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
	case "justification_and_finalization", "inactivity_updates", "rewards_and_penalties":
		switch c.PreFork {
		case "phase0":
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
				return maybeOutput(phase0.ProcessEpochRewardsAndPenalties(ctx, spec, epc, attesterData, state.(common.BeaconState)))
			}
		case "altair", "bellatrix", "capella":
			attesterData, err := altair.ComputeEpochAttesterData(ctx, spec, epc, flats, state.(altair.AltairLikeBeaconState))
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
				return maybeOutput(altair.ProcessInactivityUpdates(ctx, spec, attesterData, state.(altair.AltairLikeBeaconState)))
			case "rewards_and_penalties":
				return maybeOutput(altair.ProcessEpochRewardsAndPenalties(ctx, spec, epc, attesterData, state.(altair.AltairLikeBeaconState)))
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
		switch c.PreFork {
		case "phase0", "altair", "bellatrix":
			return maybeOutput(phase0.ProcessHistoricalRootsUpdate(ctx, spec, epc, state))
		default:
			return errors.New("historical_roots_update is only available before Capella")
		}
	case "participation_record_updates":
		if c.PreFork != "phase0" {
			return errors.New("participation_record_updates was removed after Phase0")
		}
		return maybeOutput(phase0.ProcessParticipationRecordUpdates(ctx, spec, epc, state.(phase0.Phase0PendingAttestationsBeaconState)))
	case "participation_flag_updates":
		if c.PreFork == "phase0" {
			return errors.New("participation_flag_updates was introduced after Phase0")
		}
		return maybeOutput(altair.ProcessParticipationFlagUpdates(ctx, spec, state.(*altair.BeaconStateView)))
	case "sync_committee_updates":
		if c.PreFork == "phase0" {
			return errors.New("sync_committee_updates is only available after Phase0")
		}
		return maybeOutput(altair.ProcessSyncCommitteeUpdates(ctx, spec, epc, state.(*altair.BeaconStateView)))
	case "historical_summaries_update":
		switch c.PreFork {
		case "phase0", "altair", "bellatrix":
			return errors.New("historical_summaries_update is only available after Bellatrix")
		default:
			return maybeOutput(capella.ProcessHistoricalSummariesUpdate(ctx, spec, epc, state.(capella.HistoricalSummariesBeaconState)))
		}
	}
	return ask.UnrecognizedErr
}

type TransitionBlockSubCmd struct {
	PreFork             string
	Transition          string
	Timeout             time.Duration `ask:"--timeout" help:"Timeout, e.g. 100ms"`
	configs.SpecOptions `ask:"."`
	Pre                 util.StateInput  `ask:"--pre" help:"Pre-state"`
	Op                  util.ObjInput    `ask:"<op>" help:"Block operation input"`
	Post                util.StateOutput `ask:"--post" help:"Post-state"`
}

func (c *TransitionBlockSubCmd) Help() string {
	return fmt.Sprintf("Run block-sub-process %s (%s pre-state)", c.Transition, c.PreFork)
}

func (c *TransitionBlockSubCmd) Run(ctx context.Context, args ...string) error {
	if c.Timeout != 0 {
		ctx, _ = context.WithTimeout(ctx, c.Timeout)
	}
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
		var attSl phase0.AttesterSlashing
		if err := c.Op.Read(spec.Wrap(&attSl)); err != nil {
			return err
		}
		return maybeOutput(phase0.ProcessAttesterSlashing(spec, epc, state, &attSl))
	case "attestation":
		switch c.PreFork {
		case "phase0":
			var att phase0.Attestation
			if err := c.Op.Read(spec.Wrap(&att)); err != nil {
				return err
			}
			return maybeOutput(phase0.ProcessAttestation(spec, epc, state.(phase0.Phase0PendingAttestationsBeaconState), &att))
		case "altair", "bellatrix", "capella":
			var att phase0.Attestation
			if err := c.Op.Read(spec.Wrap(&att)); err != nil {
				return err
			}
			return maybeOutput(altair.ProcessAttestation(spec, epc, state.(*altair.BeaconStateView), &att))
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
	case "bls_to_execution_change":
		switch c.PreFork {
		case "phase0", "altair", "bellatrix":
			return fmt.Errorf("fork %s does not have bls_to_execution_change processing", c.PreFork)
		case "capella", "deneb":
			var change common.SignedBLSToExecutionChange
			if err := c.Op.Read(&change); err != nil {
				return err
			}
			return maybeOutput(capella.ProcessBLSToExecutionChange(ctx, spec, nil, state, &change))
		}
	case "sync_aggregate":
		if c.PreFork == "phase0" {
			return fmt.Errorf("fork %s does not have sync_aggregate processing", c.PreFork)
		}
		var agg altair.SyncAggregate
		if err := c.Op.Read(spec.Wrap(&agg)); err != nil {
			return err
		}
		return maybeOutput(altair.ProcessSyncAggregate(ctx, spec, epc, state.(*altair.BeaconStateView), &agg))
	case "execution_payload":
		switch c.PreFork {
		case "phase0", "altair":
			return fmt.Errorf("fork %s does not have execution_payload processing", c.PreFork)
		case "bellatrix":
			var body bellatrix.BeaconBlockBody
			if err := c.Op.Read(spec.Wrap(&body)); err != nil {
				return err
			}
			return maybeOutput(bellatrix.ProcessExecutionPayload(ctx, spec, state.(bellatrix.ExecutionTrackingBeaconState),
				&body.ExecutionPayload, new(execution.NoOpExecutionEngine)))
		case "capella":
			var body capella.BeaconBlockBody
			if err := c.Op.Read(spec.Wrap(&body)); err != nil {
				return err
			}
			return maybeOutput(capella.ProcessExecutionPayload(ctx, spec, state.(capella.ExecutionTrackingBeaconState),
				&body.ExecutionPayload, new(execution.NoOpExecutionEngine)))
		case "deneb":
			var body deneb.BeaconBlockBody
			if err := c.Op.Read(spec.Wrap(&body)); err != nil {
				return err
			}
			return maybeOutput(deneb.ProcessExecutionPayload(ctx, spec, state.(deneb.ExecutionTrackingBeaconState),
				&body, new(execution.NoOpExecutionEngine)))
		}
	case "withdrawals":
		switch c.PreFork {
		case "phase0", "altair", "bellatrix":
			return fmt.Errorf("fork %s does not have withdrawals processing", c.PreFork)
		case "capella":
			var payload capella.ExecutionPayload
			if err := c.Op.Read(spec.Wrap(&payload)); err != nil {
				return err
			}
			return maybeOutput(capella.ProcessWithdrawals(ctx, spec, state.(capella.BeaconStateWithWithdrawals), &payload))
		case "deneb":
			var payload deneb.ExecutionPayload
			if err := c.Op.Read(spec.Wrap(&payload)); err != nil {
				return err
			}
			return maybeOutput(capella.ProcessWithdrawals(ctx, spec, state.(capella.BeaconStateWithWithdrawals), &payload))
		}
	}
	return ask.UnrecognizedErr
}

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
	"bellatrix": {
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
		"sync_committee_updates",
	},
	"capella": {
		"justification_and_finalization",
		"rewards_and_penalties",
		"registry_updates",
		"slashings",
		"eth1_data_reset",
		"effective_balance_updates",
		"slashings_reset",
		"randao_mixes_reset",
		"historical_summaries_update",
		"participation_record_updates",
		"sync_committee_updates",
	},
	"deneb": {
		"justification_and_finalization",
		"rewards_and_penalties",
		"registry_updates",
		"slashings",
		"eth1_data_reset",
		"effective_balance_updates",
		"slashings_reset",
		"randao_mixes_reset",
		"historical_summaries_update",
		"participation_record_updates",
		"sync_committee_updates",
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
	"bellatrix": {
		"block_header",
		"randao",
		"eth1_data",
		"proposer_slashing",
		"attester_slashing",
		"attestation",
		"deposit",
		"voluntary_exit",
		"sync_aggregate",
		"execution_payload",
	},
	"capella": {
		"block_header",
		"randao",
		"eth1_data",
		"proposer_slashing",
		"attester_slashing",
		"attestation",
		"deposit",
		"voluntary_exit",
		"bls_to_execution_change",
		"sync_aggregate",
		"execution_payload",
		"withdrawals",
	},
	"deneb": {
		"block_header",
		"randao",
		"eth1_data",
		"proposer_slashing",
		"attester_slashing",
		"attestation",
		"deposit",
		"voluntary_exit",
		"bls_to_execution_change",
		"sync_aggregate",
		"execution_payload",
		"withdrawals",
	},
}
