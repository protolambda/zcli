package transition

import (
	"bytes"
	"fmt"
	"github.com/protolambda/zrnt/eth2/core"
	"github.com/protolambda/zrnt/eth2/phase0"
	"github.com/protolambda/zssz"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
)

var TransitionCmd, BlocksCmd, SlotsCmd *cobra.Command

func init() {
	TransitionCmd = &cobra.Command{
		Use:   "transition",
		Short: "Run a state-transition",
	}
	TransitionCmd.PersistentFlags().String("pre", "", "Pre (Input) path. If none is specified, input is read from STDOUT")
	TransitionCmd.PersistentFlags().String("post", "", "Post (Output) path. If none is specified, output is written to STDOUT")

	SlotsCmd = &cobra.Command{
		Use:   "slots <number>",
		Short: "Process empty slots on the pre-state to get a post-state",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("expected one argument: <number>")
			}
			_, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("argument %v is a not a valid number", args[0])
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			isDelta, err := cmd.Flags().GetBool("delta")
			if check(err, cmd.ErrOrStderr(), "delta flag could not be parsed") {
				return
			}

			slots, _ := strconv.ParseUint(args[0], 10, 64)

			state, err := loadPreFull(cmd)
			if check(err, cmd.ErrOrStderr(), "pre state could not be loaded") {
				return
			}

			to := core.Slot(slots)
			if isDelta {
				to += state.Slot
			} else if to < state.Slot {
				report(cmd.ErrOrStderr(), "to slot is lower than pre-state slot")
				return
			}

			state.ProcessSlots(to)
			err = writePost(cmd, state.BeaconState)
			if check(err, cmd.ErrOrStderr(), "could not write post-state") {
				return
			}
		},
	}
	SlotsCmd.Flags().Bool("delta", false, "to interpret the slot number as a delta from the pre-state")
	TransitionCmd.AddCommand(SlotsCmd)

	BlocksCmd = &cobra.Command{
		Use:   "blocks",
		Short: "Process blocks on the pre-state to get a post-state",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			verifyStateRoot, err := cmd.Flags().GetBool("verify-state-root")
			if check(err, cmd.ErrOrStderr(), "verify-state-root could not be parsed") {
				return
			}

			state, err := loadPreFull(cmd)
			if check(err, cmd.ErrOrStderr(), "could not load pre-state") {
				return
			}

			for i := 0; i < len(args); i++ {
				b, err := loadBlock(args[i])
				if check(err, cmd.ErrOrStderr(), "could not load block: %s", args[i]) {
					return
				}

				blockProc := &phase0.BlockProcessFeature{Block: b, Meta: state}

				err = state.StateTransition(blockProc, verifyStateRoot)
				if check(err, cmd.ErrOrStderr(), "failed block transition to block %s", args[i]) {
					return
				}
			}

			err = writePost(cmd, state.BeaconState)
			if check(err, cmd.ErrOrStderr(), "could not write post-state") {
				return
			}
		},
	}

	BlocksCmd.Flags().Bool("verify-state-root", true, "change the state-root verification step")

	TransitionCmd.AddCommand(BlocksCmd)
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

func loadBlock(blockPath string) (*phase0.BeaconBlock, error) {
	r, err := os.Open(blockPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read block from input path: %s\n%v", blockPath, err)
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read block into buffer: %s\n%v", blockPath, err)
	}

	var block phase0.BeaconBlock
	err = zssz.Decode(&buf, uint64(buf.Len()), &block, phase0.BeaconBlockSSZ)
	if err != nil {
		return nil, fmt.Errorf("cannot decode block SSZ: %s\n%v", blockPath, err)
	}

	return &block, nil
}

func loadPreFull(cmd *cobra.Command) (*phase0.FullFeaturedState, error) {
	inPath, err := cmd.Flags().GetString("pre")
	if err != nil {
		return nil, fmt.Errorf("pre path could not be parsed")
	}

	var r io.Reader
	if inPath == "" {
		r = os.Stdin
	} else {
		r, err = os.Open(inPath)
		if err != nil {
			return nil, fmt.Errorf("cannot read pre from input path: %v", err)
		}
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read pre-state into buffer: %v", err)
	}

	var pre phase0.BeaconState
	err = zssz.Decode(&buf, uint64(buf.Len()), &pre, phase0.BeaconStateSSZ)
	if err != nil {
		return nil, fmt.Errorf("cannot decode pre-state: %v", err)
	}

	preFull := phase0.NewFullFeaturedState(&pre)
	preFull.LoadPrecomputedData()

	return preFull, nil
}

func writePost(cmd *cobra.Command, state *phase0.BeaconState) error {
	outPath, err := cmd.Flags().GetString("post")
	if err != nil {
		return fmt.Errorf("post path could not be parsed: %v", err)
	}

	var w io.Writer
	if outPath == "" {
		w = os.Stdout
	} else {
		w, err = os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	}

	_, err = zssz.Encode(w, state, phase0.BeaconStateSSZ)
	if err != nil {
		return fmt.Errorf("cannot encode post-state: %v", err)
	}
	return nil
}
