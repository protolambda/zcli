package cmd

import (
	"fmt"
	"github.com/protolambda/zcli/cmd/genesis"
	"github.com/protolambda/zcli/cmd/transition"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var RootCmd *cobra.Command

func init() {
	RootCmd = &cobra.Command{
		Use:   "zli",
		Short: "ZRNT CLI is a tool for ETH 2 debugging",
		Long: `A command line tool for ETH 2 debugging, based on ZRNT, the Go exec-spec built by @protolambda.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	RootCmd.AddCommand(GenesisCmd)
	RootCmd.AddCommand(TransitionCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func report(out io.Writer, msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(out, msg, args...)
}

func check(err error, out io.Writer, msg string, args ...interface{}) bool {
	if err != nil {
		report(out, msg, args...)
		report(out, "%v", err)
		return true
	} else {
		return false
	}
}
