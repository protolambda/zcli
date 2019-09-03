package genesis

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var GenesisCmd, RandomCmd *cobra.Command

func init() {
	check := func(msg string, err error) {
		if err != nil {
			fmt.Printf("%s: %v", msg, err)
			os.Exit(1)
		}
	}
	GenesisCmd = &cobra.Command{
		Use:   "genesis",
		Short: "Generate a genesis state",
	}

	RandomCmd = &cobra.Command{
		Use:   "random",
		Short: "Generate a random genesis state",
		Run: func(cmd *cobra.Command, args []string) {
			count, err := cmd.Flags().GetUint32("count")
			check("count is invalid", err)
			outPath, err := cmd.Flags().GetString("out")
			check("out path is invalid", err)

			fmt.Printf("count: %d out: %s\n", count, outPath)
		},
	}
	RandomCmd.Flags().Uint32("count", 64, "Number of random validators")
	RandomCmd.Flags().String("out", "", "Output path. If none is specified, output is written to STDOUT")

	GenesisCmd.AddCommand(RandomCmd)
}
