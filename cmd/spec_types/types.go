package spec_types

import (
	"github.com/protolambda/zrnt/eth2/beacon/attestations"
	"github.com/protolambda/zrnt/eth2/beacon/deposits"
	"github.com/protolambda/zrnt/eth2/beacon/eth1"
	"github.com/protolambda/zrnt/eth2/beacon/exits"
	"github.com/protolambda/zrnt/eth2/beacon/header"
	"github.com/protolambda/zrnt/eth2/beacon/slashings/attslash"
	"github.com/protolambda/zrnt/eth2/beacon/slashings/propslash"
	"github.com/protolambda/zrnt/eth2/core"
	"github.com/protolambda/zrnt/eth2/phase0"
	"github.com/protolambda/zssz"
	"github.com/protolambda/zssz/types"
)

type SpecType struct {
	Name     string
	TypeName string
	SSZTyp   types.SSZ
	Alloc    func() interface{}
}

type APIBeaconState struct {
	Root        core.Root
	BeaconState phase0.BeaconState
}

var APIBeaconStateSSZ = zssz.GetSSZ((*APIBeaconState)(nil))

var SpecTypes = []*SpecType{
	{"state_dump", "APIBeaconState", APIBeaconStateSSZ, func() interface{} { return new(APIBeaconState) }},
	{"state", "BeaconState", phase0.BeaconStateSSZ, func() interface{} { return new(phase0.BeaconState) }},
	{"block", "BeaconBlock", phase0.BeaconBlockSSZ, func() interface{} { return new(phase0.BeaconBlock) }},
	{"block_header", "BlockHeader", header.BeaconBlockHeaderSSZ, func() interface{} { return new(header.BeaconBlockHeader) }},
	{"signed_block", "SignedBeaconBlock", phase0.SignedBeaconBlockSSZ, func() interface{} { return new(phase0.SignedBeaconBlock) }},
	{"signed_block_header", "SignedBlockHeader", header.SignedBeaconBlockHeaderSSZ, func() interface{} { return new(header.SignedBeaconBlockHeader) }},
	{"block_body", "BeaconBlockBody", phase0.BeaconBlockBodySSZ, func() interface{} { return new(phase0.BeaconBlockBody) }},
	{"attestation_data", "AttestationData", attestations.AttestationDataSSZ, func() interface{} { return new(attestations.AttestationData) }},
	{"attestation", "Attestation", attestations.AttestationSSZ, func() interface{} { return new(attestations.Attestation) }},
	{"attester_slashing", "AttesterSlashing", attslash.AttesterSlashingSSZ, func() interface{} { return new(attslash.AttesterSlashing) }},
	{"proposer_slashing", "ProposerSlashing", propslash.ProposerSlashingSSZ, func() interface{} { return new(propslash.ProposerSlashing) }},
	{"deposit", "Deposit", deposits.DepositSSZ, func() interface{} { return new(deposits.Deposit) }},
	{"deposit_data", "DepositData", deposits.DepositDataSSZ, func() interface{} { return new(deposits.DepositData) }},
	{"deposit_message", "DepositMessage", deposits.DepositMessageSSZ, func() interface{} { return new(deposits.DepositMessage) }},
	{"voluntary_exit", "VoluntaryExit", exits.VoluntaryExitSSZ, func() interface{} { return new(exits.VoluntaryExit) }},
	{"signed_voluntary_exit", "SignedVoluntaryExit", exits.SignedVoluntaryExitSSZ, func() interface{} { return new(exits.SignedVoluntaryExit) }},
	{"eth1_data", "Eth1Data", zssz.GetSSZ((*eth1.Eth1Data)(nil)), func() interface{} { return new(eth1.Eth1Data) }},
	{"fork_data", "ForkData", core.ForkDataSSZ, func() interface{} { return new(core.ForkData) }},
}
