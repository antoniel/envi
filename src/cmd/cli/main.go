package main

import (
	"envi/src/cmd/cli/auth"
	"envi/src/cmd/cli/pull"
	"envi/src/cmd/cli/push"
	"envi/src/cmd/cli/undo"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

var RootCmd = &cobra.Command{
	Use:     "envi [command]",
	Short:   "envi is a tool for synchronizing .env files across teams.",
	Long:    "envi streamlines the management and synchronization of .env files among different environments or team members, facilitating consistent configuration across development workflows.",
	Version: VERSION,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(pull.PullCmd)
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
