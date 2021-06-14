package main

import (
	"context"
	"fmt"
	"github.com/protolambda/ask"
	"github.com/protolambda/zcli/commands"
	"os"
)

type MainCmd struct {
}

func (c *MainCmd) Help() string {
	return "Run ZCLI. See sub-commands."
}

func (c *MainCmd) Cmd(route string) (cmd interface{}, err error) {
	switch route {
	case "convert":
		cmd = &commands.ConvertCmd{}
	case "diff":
		cmd = &commands.DiffCmd{}
	case "meta":
		cmd = &commands.MetaCmd{}
	case "proof":
		cmd = &commands.ProofCmd{}
	case "root":
		cmd = &commands.RootCmd{}
	case "transition":
		cmd = &commands.TransitionCmd{}
	case "tree":
		cmd = &commands.TreeCmd{}
	case "version":
		cmd = &commands.VersionCmd{}
	default:
		return nil, ask.UnrecognizedErr
	}
	return
}

func (c *MainCmd) Routes() []string {
	return []string{"convert", "diff", "meta", "proof", "root", "transition", "tree", "version"}
}

func main() {
	cmd := &MainCmd{}
	descr, err := ask.Load(cmd)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to load main command: %v", err.Error())
		os.Exit(1)
	}

	if cmd, isHelp, err := descr.Execute(context.Background(), os.Args[1:]...); err != nil && err != ask.UnrecognizedErr {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	} else if cmd == nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to load command")
		os.Exit(1)
	} else if isHelp || (err == ask.UnrecognizedErr) {
		_, _ = fmt.Fprintln(os.Stdout, cmd.Usage())
		os.Exit(0)
	}
}
