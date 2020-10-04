package info

import (
	"fmt"
	. "github.com/protolambda/zcli/util"
	. "github.com/protolambda/zrnt/eth2/beacon"
	"github.com/spf13/cobra"
)

var InfoCmd, RegistryStatusCmd *cobra.Command

func RegistryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   fmt.Sprintf("registry [input (BeaconState) path]"),
		Short: fmt.Sprintf("Print a summary of the validator registry. If the input path is not specified, input is read from STDIN."),
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			spec, err := LoadSpec(cmd)
			if Check(err, cmd.ErrOrStderr(), "cannot load spec") {
				return
			}
			var inPath string
			if len(args) > 0 {
				inPath = args[0]
			}
			state, err := LoadStateInputPath(cmd, inPath, true, spec)
			if Check(err, cmd.ErrOrStderr(), "cannot verify input") {
				return
			}
			verbose, err := cmd.Flags().GetBool("verbose")
			if Check(err, cmd.ErrOrStderr(), "cannot parse verbose flag") {
				return
			}
			{
				currentEpoch := spec.SlotToEpoch(state.Slot)
				fmtStatus := func(v *Validator) string {
					if v.Slashed {
						return "💀"
					}
					if v.WithdrawableEpoch <= currentEpoch {
						return "👋"
					}
					if v.ExitEpoch <= currentEpoch {
						return "❎"
					}
					if v.ActivationEpoch <= currentEpoch {
						return "▶️"
					}
					if v.ActivationEligibilityEpoch <= currentEpoch {
						return "🔜"
					} else {
						return "📦"
					}
				}
				fmtEpoch := func(epoch Epoch) string {
					if epoch == ^Epoch(0) {
						return "♾️"
					} else {
						return fmt.Sprintf("%10d", epoch)
					}
				}
				out := cmd.OutOrStdout()
				_, err := fmt.Fprintf(out, "%8s: %14s %s %10s %10s [%10s %10s %10s %10s]",
					"index", "pubkey", "ℹ️",
					"eff.bal.", "balance",
					"eligible", "activation", "exit", "withdrawal")
				if Check(err, cmd.ErrOrStderr(), "cannot write header") {
					return
				}
				if verbose {
					_, err := fmt.Fprintf(out, " %64s %96s", "pubkey[:7]", "withdrawal-credentials")
					if Check(err, cmd.ErrOrStderr(), "cannot write header") {
						return
					}
				}
				if _, err := fmt.Fprintln(out); Check(err, cmd.ErrOrStderr(), "cannot write header") {
					return
				}
				for i, v := range state.Validators {
					_, err := fmt.Fprintf(out, "%8d: %x %s %10d %10d [%10s %10s %10s %10s]", i, v.Pubkey[:7],
						fmtStatus(v),
						v.EffectiveBalance, state.Balances[i],
						fmtEpoch(v.ActivationEligibilityEpoch), fmtEpoch(v.ActivationEpoch),
						fmtEpoch(v.ExitEpoch), fmtEpoch(v.WithdrawableEpoch))
					if Check(err, cmd.ErrOrStderr(), "cannot write output for validator %d", i) {
						return
					}
					if verbose {
						_, err := fmt.Fprintf(out, " %x %x", v.Pubkey, v.WithdrawalCredentials)
						if Check(err, cmd.ErrOrStderr(), "cannot write output for validator %d", i) {
							return
						}
					}
					if _, err := fmt.Fprintln(out); Check(err, cmd.ErrOrStderr(), "cannot write output for validator %d", i) {
						return
					}
				}
			}
		},
	}
	c.Flags().StringP("spec", "s", "mainnet", "The spec configuration to use. Can also be a path to a yaml config file. E.g. 'mainnet', 'minimal', or 'my_yaml_path.yml")
	c.Flags().BoolP("verbose", "v", false, "verbose=bool")
	return c
}

func init() {
	InfoCmd = &cobra.Command{
		Use:   "info",
		Short: "Information about eth2 data.",
	}

	RegistryStatusCmd = RegistryCmd()

	InfoCmd.AddCommand(RegistryStatusCmd)
}
