package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sveniu/rushpath/cmd/rushpath/cmd/uiflows"
)

func init() {
	cmdMFA.AddCommand(cmdMFATOTPDisable)
}

var cmdMFATOTPDisable = &cobra.Command{
	Use:   "totp-disable",
	Short: "Disable TOTP multi-factor authentication",
	Long:  `Rushpath – Disable TOTP multi-factor authentication`,
	RunE:  runMFATOTPDisable,
}

func runMFATOTPDisable(cmd *cobra.Command, args []string) error {
	ui := uiflows.New()
	return ui.MFATOTPDisable()
}
