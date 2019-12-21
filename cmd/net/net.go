package net

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
	. "github.com/protolambda/zcli/util"
	"github.com/spf13/cobra"
	"net"
	"strings"
)

var NetCmd, EnrCmd *cobra.Command

func init() {
	NetCmd = &cobra.Command{
		Use:   "net",
		Short: "Util tools for networking",
	}
	EnrCmd = &cobra.Command{
		Use:   "enr [base64 ENR string]",
		Short: "Decode ENR record. If the ENR string is not specified, input is read from STDIN. Base64 raw-URL-safe type (RFC 4648).",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			var input string
			if len(args) == 1 {
				input = args[0]
			} else {
				var buf bytes.Buffer
				_, err := buf.ReadFrom(cmd.InOrStdin())
				if Check(err, cmd.ErrOrStderr(), "cannot read ENR from input") {
					return
				}
				input = buf.String()
			}
			input = strings.TrimSpace(input)
			Report(cmd.OutOrStdout(), "input: %s", input)

			data, err := base64.RawURLEncoding.DecodeString(input)
			if Check(err, cmd.ErrOrStderr(), "ENR is not valid base64") {
				return
			}
			var record enr.Record
			if Check(rlp.Decode(bytes.NewReader(data), &record), cmd.ErrOrStderr(), "invalid ENR RLP encoding") {
				return
			}
			idSchemeName := record.IdentityScheme()
			switch idSchemeName {
			case "v4":
				var ip enr.IPv4
				_ = record.Load(&ip)
				var tcpPort enr.TCP
				_ = record.Load(&tcpPort)
				var udpPort enr.UDP
				_ = record.Load(&udpPort)
				Report(cmd.OutOrStdout(), "ip4: %s tcp: %d udp: %d", net.IP(ip).String(), tcpPort, udpPort)
			case "v6":
				var ip enr.IPv6
				_ = record.Load(&ip)
				var tcpPort enr.TCP6
				_ = record.Load(&tcpPort)
				var udpPort enr.UDP6
				_ = record.Load(&udpPort)
				Report(cmd.OutOrStdout(), "ip6: %s tcp: %d udp: %d", net.IP(ip).String(), tcpPort, udpPort)
			}

			Report(cmd.OutOrStdout(), "signature bytes: %s", hex.EncodeToString(record.Signature()))

			if idScheme, ok := enode.ValidSchemes[idSchemeName]; ok {
				if err := record.VerifySignature(idScheme); err != nil {
					Report(cmd.OutOrStdout(), "signature failed to verify for id scheme %s: %v", idSchemeName, err)
				} else {
					Report(cmd.OutOrStdout(), "signature verified")
				}
			} else {
				Report(cmd.OutOrStdout(), "identity scheme %s cannot be verified", idSchemeName)
			}
		},
	}

	NetCmd.AddCommand(EnrCmd)
}
