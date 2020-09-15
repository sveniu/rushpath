package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sveniu/rushpath/cmd/rushpath/cmd/uiflows"
)

func init() {
	cmdMFA.AddCommand(cmdMFATOTPEnable)
}

var cmdMFATOTPEnable = &cobra.Command{
	Use:   "totp-enable",
	Short: "Enable TOTP multi-factor authentication",
	Long:  `Rushpath – Enable TOTP multi-factor authentication`,
	RunE:  runMFATOTPEnable,
}

func runMFATOTPEnable(cmd *cobra.Command, args []string) error {
	ui := uiflows.New()
	return ui.MFATOTPEnable()
}
