package storage

import (
	"log"
	"os"

	gap "github.com/muesli/go-app-paths"
)

func initTaskDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o770)
		}
		return err
	}
	return nil
}

func GetApplicationDataPath() string {

	// get XDG paths
	scope := gap.NewScope(gap.User, "envii-storage")
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
