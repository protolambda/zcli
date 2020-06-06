package cmd

import (
	"fmt"
	"github.com/protolambda/zcli/cmd/api_util"
	"github.com/protolambda/zcli/cmd/checker"
	"github.com/protolambda/zcli/cmd/diff"
	"github.com/protolambda/zcli/cmd/genesis"
	"github.com/protolambda/zcli/cmd/info"
	"github.com/protolambda/zcli/cmd/keys"
	"github.com/protolambda/zcli/cmd/meta"
	"github.com/protolambda/zcli/cmd/net"
	"github.com/protolambda/zcli/cmd/pretty"
	"github.com/protolambda/zcli/cmd/roots"
	"github.com/protolambda/zcli/cmd/transition"
	"github.com/protolambda/zcli/util"
	"github.com/protolambda/zrnt/eth2"
	"github.com/protolambda/zrnt/eth2/beacon"
	"github.com/protolambda/zssz"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var RootCmd, VersionCmd *cobra.Command

func init() {
	RootCmd = &cobra.Command{
		Use:   "zcli",
		Short: "ZRNT CLI is a tool for ETH 2 debugging",
		Long:  `A command line tool for ETH 2 debugging, based on ZRNT, the Go exec-spec built by @protolambda.`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print versions of integrated tools",
		Run: func(cmd *cobra.Command, args []string) {
			util.Report(cmd.OutOrStdout(), `
ZCLI: v0.0.29
ZRNT: `+eth2.VERSION+`
ZSSZ: `+zssz.VERSION+`

config: `+beacon.PRESET_NAME+`
`)
		},
	}

	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(genesis.GenesisCmd)
	RootCmd.AddCommand(transition.TransitionCmd)
	RootCmd.AddCommand(pretty.PrettyCmd)
	RootCmd.AddCommand(diff.DiffCmd)
	RootCmd.AddCommand(checker.CheckCmd)
	RootCmd.AddCommand(roots.HashTreeRootCmd)
	RootCmd.AddCommand(keys.KeysCmd)
	RootCmd.AddCommand(meta.MetaCmd)
	RootCmd.AddCommand(api_util.ApiUtilCmd)
	RootCmd.AddCommand(info.InfoCmd)
	RootCmd.AddCommand(net.NetCmd)
}

type writerWrap struct {
	w       io.Writer
	onWrite func()
}

func (ww *writerWrap) Write(p []byte) (n int, err error) {
	n, err = ww.w.Write(p)
	ww.onWrite()
	return
}

func Execute() {
	exitCode := 0
	// If there's any regular error (no panics), then make it exit code 1
	errW := &writerWrap{w: os.Stderr, onWrite: func() {
		exitCode = 1
	}}
	RootCmd.SetErr(errW)
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		exitCode = 1
	}
	os.Exit(exitCode)
}
