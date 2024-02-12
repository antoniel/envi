package storage

import (
	"fmt"
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

const TOKEN_FILE_NAME = "token"

func GetApplicationDataPath() string {

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

type accessToken struct{}

var AccessToken = accessToken{}

func (accessToken) Get() (string, error) {
	path := GetApplicationDataPath()
	token, err := os.ReadFile(path + "/" + TOKEN_FILE_NAME)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

var ErrUnableToPersistToken = fmt.Errorf("unable to persist token")
var ErrInvalidToken = fmt.Errorf("invalid token")

func (accessToken) Set(token string) error {
	path := GetApplicationDataPath()
	if token == "" {
		// return ErrInvalidToken,
		return ErrInvalidToken
	}
	some := os.WriteFile(path+"/"+TOKEN_FILE_NAME, []byte(token), 0644)
	if some != nil {
		return (ErrUnableToPersistToken)
	}
	return nil
}

func (accessToken) Clear() error {
	path := GetApplicationDataPath()
	err := os.Remove(path + "/" + TOKEN_FILE_NAME)
	if err != nil {
		return err
	}
	return nil
}
