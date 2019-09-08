package genesis

import (
	"encoding/hex"
	. "github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2/core"
	"github.com/protolambda/zrnt/eth2/phase0"
	"github.com/protolambda/zrnt/eth2/util/hashing"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
)

var GenesisCmd, MockedCmd *cobra.Command

type KeyPair struct {
	Priv string `yaml:"privkey"`
	Pub  string `yaml:"pubkey"`
}
type KeyPairs []KeyPair

func init() {
	GenesisCmd = &cobra.Command{
		Use:   "genesis",
		Short: "Generate a genesis state",
	}

	MockedCmd = &cobra.Command{
		Use:   "mock",
		Short: "Generate a genesis state from a predefined set of keys",
		Run: func(cmd *cobra.Command, args []string) {
			eth1RootStr, err := cmd.Flags().GetString("eth1-root")
			if Check(err, cmd.ErrOrStderr(), "cannot parse eth1-root") {
				return
			}
			var eth1Root [32]byte
			if Check(decodeRoot(eth1RootStr, &eth1Root), cmd.ErrOrStderr(), "could not decode eth1-root") {
				return
			}
			genesisTime, err := cmd.Flags().GetUint64("genesis-time")
			if Check(err, cmd.ErrOrStderr(), "could not parse genesis time") {
				return
			}

			count, err := cmd.Flags().GetUint32("count")
			if Check(err, cmd.ErrOrStderr(), "count is invalid") {
				return
			}
			keysPath, err := cmd.Flags().GetString("keys")
			if Check(err, cmd.ErrOrStderr(), "keys path is invalid") {
				return
			}

			var r io.Reader
			if keysPath == "" {
				r = cmd.InOrStdin()
			} else {
				r, err = os.Open(keysPath)
				if err != nil {
					Report(cmd.ErrOrStderr(), "cannot open key pairs file: %s\n%v", keysPath, err)
					return
				}
			}
			var keys KeyPairs
			dec := yaml.NewDecoder(r)
			if Check(dec.Decode(&keys), cmd.ErrOrStderr(), "cannot read key pairs from YAML file: %s", keysPath) {
				return
			}

			if count > uint32(len(keys)) {
				Report(cmd.ErrOrStderr(), "not enough keys available, expected at least %d, got %d", count, len(keys))
				return
			}

			var validators []phase0.KickstartValidatorData
			var privKeys [][32]byte
			for i := uint32(0); i < count; i++ {
				k := &keys[i]
				var pub core.BLSPubkey
				if strings.HasPrefix(k.Pub, "0x") {
					k.Pub = k.Pub[2:]
				}
				k.Pub = strings.Repeat("0", (48*2)-len(k.Pub)) + k.Pub
				if _, err := hex.Decode(pub[:], []byte(k.Pub[:])); Check(err, cmd.ErrOrStderr(), "could not decode pubkey for %d", i) {
					return
				}
				var priv [32]byte
				if Check(decodeRoot(k.Priv, &priv), cmd.ErrOrStderr(), "cannot parse priv key") {
					return
				}
				withdrawal := hashing.Hash(pub[:])
				withdrawal[0] = core.BLS_WITHDRAWAL_PREFIX
				validators = append(validators, phase0.KickstartValidatorData{
					Pubkey:                pub,
					WithdrawalCredentials: withdrawal,
					Balance:               core.MAX_EFFECTIVE_BALANCE,
				})
				privKeys = append(privKeys, priv)
			}

			state, err := phase0.KickStartStateWithSignatures(eth1Root, core.Timestamp(genesisTime), validators, privKeys)
			if Check(err, cmd.ErrOrStderr(), "cannot create beacon state") {
				return
			}

			if Check(WriteStateOutput(cmd, "out", state.BeaconState), cmd.ErrOrStderr(), "cannot output state") {
				return
			}
		},
	}
	MockedCmd.Flags().Uint64("genesis-time", 0, "Genesis time, decimal base")
	MockedCmd.Flags().String("eth1-root", "0x4242424242424242424242424242424242424242424242424242424242424242", "Eth1 root, hex encoded")
	MockedCmd.Flags().Uint32("count", 64, "Number of validators")
	MockedCmd.Flags().String("keys", "", "YAML keys path. If none is specified, keys are read from STDIN")
	MockedCmd.Flags().String("out", "", "Output path. If none is specified, output is written to STDOUT")

	GenesisCmd.AddCommand(MockedCmd)
}

func decodeRoot(inputHex string, out *[32]byte) (err error) {
	if strings.HasPrefix(inputHex, "0x") {
		inputHex = inputHex[2:]
	}
	inputHex = strings.Repeat("0", (32*2)-len(inputHex))+inputHex
	_, err = hex.Decode(out[:], []byte(inputHex[:]))
	return
}
