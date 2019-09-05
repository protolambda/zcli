package pretty

import (
	"fmt"
	. "github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2/beacon/attestations"
	"github.com/protolambda/zrnt/eth2/beacon/deposits"
	"github.com/protolambda/zrnt/eth2/beacon/eth1"
	"github.com/protolambda/zrnt/eth2/beacon/exits"
	"github.com/protolambda/zrnt/eth2/beacon/header"
	"github.com/protolambda/zrnt/eth2/beacon/slashings/attslash"
	"github.com/protolambda/zrnt/eth2/beacon/slashings/propslash"
	"github.com/protolambda/zrnt/eth2/beacon/transfers"
	"github.com/protolambda/zrnt/eth2/phase0"
	"github.com/protolambda/zssz"
	"github.com/protolambda/zssz/types"
	"github.com/spf13/cobra"
)

var PrettyCmd *cobra.Command

type CmdType struct {
	Name     string
	TypeName string
	SSZTyp   types.SSZ
	Alloc    func() interface{}
}

func (t *CmdType) MakeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s [input path]", t.Name),
		Short: fmt.Sprintf("Pretty print a %s, if the input path is not specified, input is read from STDIN", t.TypeName),
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var path string
			if len(args) > 0 {
				path = args[0]
			}
			dst := t.Alloc()

			err := LoadSSZInputPath(cmd, path, true, dst, phase0.BeaconStateSSZ)
			if Check(err, cmd.ErrOrStderr(), "cannot load input") {
				return
			}

			zssz.Pretty(cmd.OutOrStdout(), "  ", dst, t.SSZTyp)
		},
	}
}

var CmdTypes = []*CmdType{
	{"state", "BeaconState", phase0.BeaconStateSSZ, func() interface{} { return new(phase0.BeaconState) }},
	{"block", "BeaconBlock", phase0.BeaconBlockSSZ, func() interface{} { return new(phase0.BeaconBlock) }},
	{"block_header", "BlockHeader", header.BeaconBlockHeaderSSZ, func() interface{} { return new(header.BeaconBlockHeader) }},
	{"block_body", "BeaconBlockBody", phase0.BeaconBlockBodySSZ, func() interface{} { return new(phase0.BeaconBlockBody) }},
	{"attestation", "Attestation", attestations.AttestationSSZ, func() interface{} { return new(attestations.Attestation) }},
	{"attester_slashing", "AttesterSlashing", attslash.AttesterSlashingSSZ, func() interface{} { return new(attslash.AttesterSlashing) }},
	{"proposer_slashing", "ProposerSlashing", propslash.ProposerSlashingSSZ, func() interface{} { return new(propslash.ProposerSlashing) }},
	{"deposit", "Deposit", deposits.DepositSSZ, func() interface{} { return new(deposits.Deposit) }},
	{"transfer", "Transfer", transfers.TransferSSZ, func() interface{} { return new(transfers.Transfer) }},
	{"voluntary_exit", "VoluntaryExit", exits.VoluntaryExitSSZ, func() interface{} { return new(exits.VoluntaryExit) }},
	{"deposit_data", "Deposit", deposits.DepositDataSSZ, func() interface{} { return new(deposits.DepositData) }},
	{"eth1_data", "Eth1Data", zssz.GetSSZ((*eth1.Eth1Data)(nil)), func() interface{} { return new(eth1.Eth1Data) }},
}

func init() {
	PrettyCmd = &cobra.Command{
		Use:   "pretty",
		Short: "pretty-print SSZ data",
	}

	for _, t := range CmdTypes {
		PrettyCmd.AddCommand(t.MakeCmd())
	}
}
