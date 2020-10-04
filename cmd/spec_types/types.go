package spec_types

import (
	. "github.com/protolambda/zrnt/eth2/beacon"
)

type SpecType struct {
	Name     string
	TypeName string
	Alloc    func() interface{}
}

var SpecTypes = []*SpecType{
	{"state", "BeaconState", func() interface{} { return new(BeaconState) }},
	{"block", "BeaconBlock", func() interface{} { return new(BeaconBlock) }},
	{"block_header", "BlockHeader", func() interface{} { return new(BeaconBlockHeader) }},
	{"signed_block", "SignedBeaconBlock", func() interface{} { return new(SignedBeaconBlock) }},
	{"signed_block_header", "SignedBlockHeader", func() interface{} { return new(SignedBeaconBlockHeader) }},
	{"block_body", "BeaconBlockBody", func() interface{} { return new(BeaconBlockBody) }},
	{"attestation_data", "AttestationData", func() interface{} { return new(AttestationData) }},
	{"attestation", "Attestation", func() interface{} { return new(Attestation) }},
	{"attester_slashing", "AttesterSlashing", func() interface{} { return new(AttesterSlashing) }},
	{"proposer_slashing", "ProposerSlashing", func() interface{} { return new(ProposerSlashing) }},
	{"deposit", "Deposit", func() interface{} { return new(Deposit) }},
	{"deposit_data", "DepositData", func() interface{} { return new(DepositData) }},
	{"deposit_message", "DepositMessage", func() interface{} { return new(DepositMessage) }},
	{"voluntary_exit", "VoluntaryExit", func() interface{} { return new(VoluntaryExit) }},
	{"signed_voluntary_exit", "SignedVoluntaryExit", func() interface{} { return new(SignedVoluntaryExit) }},
	{"eth1_data", "Eth1Data", func() interface{} { return new(Eth1Data) }},
	{"fork_data", "ForkData", func() interface{} { return new(ForkData) }},
}
