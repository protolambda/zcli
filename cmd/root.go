package cmd

import (
	"fmt"
	"github.com/protolambda/zcli/cmd/genesis"
	"github.com/protolambda/zcli/cmd/transition"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "zli",
	Short: "ZRNT CLI is a tool for ETH 2 debugging",
	Long: `A command line tool for ETH 2 debugging, based on ZRNT, the Go exec-spec built by @protolambda.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	rootCmd.AddCommand(genesis.GenesisCmd)
	rootCmd.AddCommand(transition.TransitionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
