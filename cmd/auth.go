package cmd

import (
	"github.com/spf13/cobra"
)

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Display commands for authenticating envi with an account",
	RunE:  AuthCmdE,
}

func init() {
	AuthCmd.AddCommand(AuthLoginCmd)
}

func AuthCmdE(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}
