package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sveniu/rushpath/cmd/rushpath/cmd/uiflows"
)

func init() {
	cmdMFA.AddCommand(cmdMFAU2FAdd)
}

var cmdMFAU2FAdd = &cobra.Command{
	Use:   "u2f-add",
	Short: "Add U2F device for multi-factor authentication",
	Long:  `Rushpath – Add U2F device for multi-factor authentication`,
	RunE:  runMFAU2FAdd,
}

func runMFAU2FAdd(cmd *cobra.Command, args []string) error {
	ui := uiflows.New()
	return ui.MFAU2FAdd()
}
