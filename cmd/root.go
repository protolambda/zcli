package cmd

import (
	"fmt"
	"github.com/protolambda/zcli/cmd/diff"
	"github.com/protolambda/zcli/cmd/genesis"
	"github.com/protolambda/zcli/cmd/keys"
	"github.com/protolambda/zcli/cmd/meta"
	"github.com/protolambda/zcli/cmd/pretty"
	"github.com/protolambda/zcli/cmd/roots"
	"github.com/protolambda/zcli/cmd/transition"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd *cobra.Command

func init() {
	RootCmd = &cobra.Command{
		Use:   "zcli",
		Short: "ZRNT CLI is a tool for ETH 2 debugging",
		Long:  `A command line tool for ETH 2 debugging, based on ZRNT, the Go exec-spec built by @protolambda.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	RootCmd.AddCommand(genesis.GenesisCmd)
	RootCmd.AddCommand(transition.TransitionCmd)
	RootCmd.AddCommand(pretty.PrettyCmd)
	RootCmd.AddCommand(diff.DiffCmd)
	RootCmd.AddCommand(roots.HashTreeRootCmd, roots.SigningRootCmd)
	RootCmd.AddCommand(keys.KeysCmd)
	RootCmd.AddCommand(meta.MetaCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
