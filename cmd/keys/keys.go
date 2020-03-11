package keys

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	hbls "github.com/herumi/bls-eth-go-binary/bls"
	. "github.com/protolambda/zcli/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io"
	"math/big"
	"os"
	"path"
	"strconv"
)

var KeysCmd, GenerateCmd, ShardCmd *cobra.Command

type KeyPair struct {
	Priv string `yaml:"privkey"`
	Pub  string `yaml:"pubkey"`
}

type KeyPairs []KeyPair

func init() {
	KeysCmd = &cobra.Command{
		Use:   "keys",
		Short: "Generate and process keys",
	}

	GenerateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a list of keys",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			from, err := cmd.Flags().GetUint("from")
			if Check(err, cmd.ErrOrStderr(), "cannot parse 'from'") {
				return
			}
			to, err := cmd.Flags().GetUint("to")
			if Check(err, cmd.ErrOrStderr(), "cannot parse 'to'") {
				return
			}
			keys := make(KeyPairs, 0)
			for i := uint64(from); i < uint64(to); i++ {
				privKey := GenerateInteropKey(i)
				var secKey hbls.SecretKey
				if Check(secKey.Deserialize(privKey[:]), cmd.ErrOrStderr(), "cannot deserialize secret key") {
					return
				}
				pubKey := secKey.GetPublicKey().Serialize()
				keys = append(keys, KeyPair{
					Priv: fmt.Sprintf("0x%x", privKey),
					Pub:  fmt.Sprintf("0x%x", pubKey),
				})
			}
			keysPath, err := cmd.Flags().GetString("keys")
			if Check(err, cmd.ErrOrStderr(), "keys path is invalid") {
				return
			}

			var w io.Writer
			if keysPath == "" {
				w = cmd.OutOrStdout()
			} else {
				w, err = os.OpenFile(keysPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
				if err != nil {
					Report(cmd.ErrOrStderr(), "cannot open key pairs file: %s\n%v", keysPath, err)
					return
				}
			}
			enc := yaml.NewEncoder(w)
			if Check(enc.Encode(&keys), cmd.ErrOrStderr(), "cannot write keys to output") {
				return
			}
		},
	}
	GenerateCmd.Flags().Uint("from", 0, "Index to start key-generation from (incl.)")
	GenerateCmd.Flags().Uint("to", 16, "Index to end key-generation at (excl.)")
	GenerateCmd.Flags().String("keys", "", "YAML keys path. If none is specified, keys are written to STDOUT")

	ShardCmd = &cobra.Command{
		Use:     "shard <size range 0> [<size range 1> [<size range 2> [...]]]",
		Aliases: []string{"split"},
		Short:   "Shard (split) a YAML list of keys into ranges. Specify sizes as arguments.",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
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

			outPath, err := cmd.Flags().GetString("out")
			if Check(err, cmd.ErrOrStderr(), "out path is invalid") {
				return
			}

			start := uint64(0)
			for i := 0; i < len(args); i++ {
				w, err := os.OpenFile(path.Join(outPath, fmt.Sprintf("key_batch_%d.yaml", i)), os.O_CREATE|os.O_WRONLY, os.ModePerm)
				if err != nil {
					Report(cmd.ErrOrStderr(), "cannot open key pairs file: %s\n%v", keysPath, err)
					return
				}
				size, err := strconv.ParseUint(args[i], 10, 64)
				end := start + size
				if end > uint64(len(keys)) {
					Report(cmd.ErrOrStderr(), "ran out of keys to make key ranges with. Arg: #%d, size: %d, range start: %d, range end: %d, keys len: %d", i, size, start, end, len(keys))
					return
				}
				enc := yaml.NewEncoder(w)
				keyRange := KeyPairs(keys[start:end])
				start = end
				if Check(enc.Encode(&keyRange), cmd.ErrOrStderr(), "cannot write keys to output") {
					return
				}
			}
		},
	}
	ShardCmd.Flags().String("keys", "", "YAML keys path. If none is specified, keys are read from STDIN")
	ShardCmd.Flags().String("out", "", "YAML output directory path. Files are written as 'key_batch_$i.yaml', e.g.: key_batch_0.yaml, key_batch_1.yaml, etc.")
	KeysCmd.AddCommand(GenerateCmd, ShardCmd)
}

var CurveOrder, _ = new(big.Int).SetString("52435875175126190479447740508185965837690552500527637822603658699938581184513", 10)

func GenerateInteropKey(index uint64) (out [32]byte) {
	input := [32]byte{}
	binary.LittleEndian.PutUint64(input[:8], index)
	h := sha256.New()
	h.Write(input[:])
	hOut := h.Sum(nil)
	for i := 0; i < 32; i++ {
		input[31-i] = hOut[i]
	}
	privKey := new(big.Int).SetBytes(input[:])
	privKey.Mod(privKey, CurveOrder)
	pBytes := privKey.Bytes()
	copy(out[32-len(pBytes):], pBytes)
	return
}
