package spec_types

import (
	. "github.com/protolambda/zrnt/eth2/beacon"
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
	Root        Root
	BeaconState BeaconState
}

var APIBeaconStateSSZ = zssz.GetSSZ((*APIBeaconState)(nil))

var SpecTypes = []*SpecType{
	{"state_dump", "APIBeaconState", APIBeaconStateSSZ, func() interface{} { return new(APIBeaconState) }},
	{"state", "BeaconState", BeaconStateSSZ, func() interface{} { return new(BeaconState) }},
	{"block", "BeaconBlock", BeaconBlockSSZ, func() interface{} { return new(BeaconBlock) }},
	{"block_header", "BlockHeader", BeaconBlockHeaderSSZ, func() interface{} { return new(BeaconBlockHeader) }},
	{"signed_block", "SignedBeaconBlock", SignedBeaconBlockSSZ, func() interface{} { return new(SignedBeaconBlock) }},
	{"signed_block_header", "SignedBlockHeader", SignedBeaconBlockHeaderSSZ, func() interface{} { return new(SignedBeaconBlockHeader) }},
	{"block_body", "BeaconBlockBody", BeaconBlockBodySSZ, func() interface{} { return new(BeaconBlockBody) }},
	{"attestation_data", "AttestationData", AttestationDataSSZ, func() interface{} { return new(AttestationData) }},
	{"attestation", "Attestation", AttestationSSZ, func() interface{} { return new(Attestation) }},
	{"attester_slashing", "AttesterSlashing", AttesterSlashingSSZ, func() interface{} { return new(AttesterSlashing) }},
	{"proposer_slashing", "ProposerSlashing", ProposerSlashingSSZ, func() interface{} { return new(ProposerSlashing) }},
	{"deposit", "Deposit", DepositSSZ, func() interface{} { return new(Deposit) }},
	{"deposit_data", "DepositData", DepositDataSSZ, func() interface{} { return new(DepositData) }},
	{"deposit_message", "DepositMessage", DepositMessageSSZ, func() interface{} { return new(DepositMessage) }},
	{"voluntary_exit", "VoluntaryExit", VoluntaryExitSSZ, func() interface{} { return new(VoluntaryExit) }},
	{"signed_voluntary_exit", "SignedVoluntaryExit", SignedVoluntaryExitSSZ, func() interface{} { return new(SignedVoluntaryExit) }},
	{"eth1_data", "Eth1Data", zssz.GetSSZ((*Eth1Data)(nil)), func() interface{} { return new(Eth1Data) }},
	{"fork_data", "ForkData", ForkDataSSZ, func() interface{} { return new(ForkData) }},
}
