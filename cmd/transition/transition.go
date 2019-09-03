package transition

import (
	"bytes"
	"fmt"
	"github.com/protolambda/zrnt/eth2/phase0"
	"github.com/protolambda/zssz"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var TransitionCmd, BlocksCmd *cobra.Command

func init() {
	TransitionCmd = &cobra.Command{
		Use:   "transition",
		Short: "Run a state-transition",
	}
	TransitionCmd.PersistentFlags().String("pre", "", "Pre (Input) path. If none is specified, input is read from STDOUT")
	TransitionCmd.PersistentFlags().String("post", "", "Post (Output) path. If none is specified, output is written to STDOUT")

	BlocksCmd = &cobra.Command{
		Use:   "blocks",
		Short: "Process blocks on top of the state transition",
		Run: func(cmd *cobra.Command, args []string) {
			verifyStateRoot, err := cmd.Flags().GetBool("verify-state-root")
			check(err, "verify-state-root could not be parsed")

			state := loadPreFull(cmd)

			for i := 0; i < len(args); i++ {
				b := loadBlock(args[i])

				blockProc := &phase0.BlockProcessFeature{Block: b, Meta: state}

				err := state.StateTransition(blockProc, verifyStateRoot)
				check(err, "failed block transition to block %s", args[i])
			}

			writePost(cmd, state.BeaconState)
		},
	}

	BlocksCmd.Flags().Bool("verify-state-root", true, "change the state-root verification step")

	TransitionCmd.AddCommand(BlocksCmd)
}

func check(err error, msg string, args ...interface{}) {
	if err != nil {
		fmt.Printf(msg, args...)
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}

func loadBlock(blockPath string) *phase0.BeaconBlock {
	r, err := os.Open(blockPath)
	check(err, "cannot read block from input path: %s", blockPath)

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	check(err, "cannot read block into buffer: %s", blockPath)

	var block phase0.BeaconBlock
	err = zssz.Decode(&buf, uint64(buf.Len()), &block, phase0.BeaconBlockSSZ)
	check(err, "cannot decode block SSZ: %s", blockPath)

	return &block
}

func loadPreFull(cmd *cobra.Command) *phase0.FullFeaturedState {
	inPath, err := cmd.Flags().GetString("pre")
	check(err, "pre path is invalid")

	var r io.Reader
	if inPath == "" {
		r = os.Stdin
	} else {
		r, err = os.Open(inPath)
		check(err, "cannot read pre from input path")
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	check(err, "cannot read pre-state into buffer")

	var pre phase0.BeaconState
	err = zssz.Decode(&buf, uint64(buf.Len()), &pre, phase0.BeaconStateSSZ)
	check(err, "cannot decode pre-state")

	preFull := phase0.NewFullFeaturedState(&pre)
	preFull.LoadPrecomputedData()
	return preFull
}

func writePost(cmd *cobra.Command, state *phase0.BeaconState) {
	outPath, err := cmd.Flags().GetString("post")
	check(err, "post path is invalid")

	var w io.Writer
	if outPath == "" {
		w = os.Stdout
	} else {
		w, err = os.OpenFile(outPath, os.O_CREATE | os.O_WRONLY, os.ModePerm)
	}

	_, err = zssz.Encode(w, state, phase0.BeaconStateSSZ)
	check(err, "cannot encode post-state")
}
