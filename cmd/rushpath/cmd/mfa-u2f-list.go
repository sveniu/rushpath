package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sveniu/rushpath/cmd/rushpath/cmd/uiflows"
)

func init() {
	cmdMFA.AddCommand(cmdMFAU2FList)
}

var cmdMFAU2FList = &cobra.Command{
	Use:   "u2f-list",
	Short: "List registered U2F devices",
	Long:  `Rushpath – List U2F devices`,
	RunE:  runMFAU2FList,
}

func runMFAU2FList(cmd *cobra.Command, args []string) error {
	ui := uiflows.New()
	return ui.MFAU2FList()
}
