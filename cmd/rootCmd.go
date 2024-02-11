package cmd

import (
	"log"
	"os"

	gap "github.com/muesli/go-app-paths"
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

func initTaskDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o770)
		}
		return err
	}
	return nil
}

func getApplicationDataPath() string {
	// get XDG paths
	scope := gap.NewScope(gap.User, "envi-storage")
	dirs, err := scope.DataDirs()
	if err != nil {
		log.Fatal(err)
	}
	// create the app base dir, if it doesn't exist
	var taskDir string
	if len(dirs) > 0 {
		taskDir = dirs[0]
	} else {
		taskDir, _ = os.UserHomeDir()
	}

	if err := initTaskDir(taskDir); err != nil {
		log.Fatal(err)
	}
	return taskDir
}

func init() {
	RootCmd.AddCommand(PullCmd)
	RootCmd.AddCommand(AuthCmd)
}
