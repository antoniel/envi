package main

import (
	"engov/apps/cli/cli/auth"
	"engov/apps/cli/cli/pull"
	"engov/apps/cli/cli/push"
	"engov/apps/cli/cli/undo"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

var RootCmd = &cobra.Command{
	Use:     "engov [command]",
	Short:   "engov is a tool for synchronizing .env files across teams.",
	Long:    "engov streamlines the management and synchronization of .env files among different environments or team members, facilitating consistent configuration across development workflows.",
	Version: VERSION,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(pull.PullCmd)
	RootCmd.PersistentFlags().StringP("provider", "p", "zipper", "Provider to use to pull the .env file: zipper | k8s")
	RootCmd.PersistentFlags().StringP("k8s-values-path", "k", "", "Path to the k8s values file")
	RootCmd.PersistentFlags().StringP("secrets-declaration", "s", "", "Path or identifier for the secrets declaration")

	RootCmd.AddCommand(auth.AuthCmd)
	RootCmd.AddCommand(push.PushCmd)
	RootCmd.AddCommand(undo.UndoCmd)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}