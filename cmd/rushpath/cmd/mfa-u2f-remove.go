package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sveniu/rushpath/cmd/rushpath/cmd/uiflows"
)

func init() {
	cmdMFA.AddCommand(cmdMFAU2FRemove)
}

var cmdMFAU2FRemove = &cobra.Command{
	Use:   "u2f-remove",
	Short: "Remove U2F device",
	Long:  `Rushpath – Remove U2F device`,
	RunE:  runMFAU2FRemove,
}

func runMFAU2FRemove(cmd *cobra.Command, args []string) error {
	ui := uiflows.New()
	return ui.MFAU2FRemove()
}
