package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	cmdRoot.AddCommand(cmdMFA)
}

var cmdMFA = &cobra.Command{
	Use:   "mfa",
	Short: "Manage multi-factor authentication",
	Long:  `Rushpath – multi-factor authentication management`,
}
