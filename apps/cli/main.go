package main

import (
	"envii/apps/cli/cmd/auth"
	"envii/apps/cli/cmd/pull"
	"envii/apps/cli/cmd/push"
	"envii/apps/cli/cmd/undo"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

var RootCmd = &cobra.Command{
	Use:     "envii [command]",
	Short:   "envii is a tool for synchronizing .env files across teams.",
	Long:    "envii streamlines the management and synchronization of .env files among different environments or team members, facilitating consistent configuration across development workflows.",
	Version: VERSION,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	RootCmd.PersistentFlags().StringP("provider", "p", "zipper", "Provider to use to pull the .env file: zipper | k8s")
	RootCmd.PersistentFlags().StringP("k8s-values-path", "k", "", "Path to the k8s values file")
	RootCmd.PersistentFlags().StringP("secrets-declaration", "s", "", "Path or identifier for the secrets declaration")
	RootCmd.AddCommand(pull.PullCmd)
	pull.PullCmd.Flags().Bool("preserve", false, "Preserve any existing value in the .env file")

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
