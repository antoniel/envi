package storage

import (
	"fmt"
	"os"
)

type LocalHistorySave interface {
	Save(token string) error
}

const local_history_file_name = "local_history"

type localHistory struct{}

// Local history will persist the last .env that was overwritten by the user
// using the `envii pull` command.
var LocalHistory = localHistory{}

func (localHistory) Get() (string, error) {
	path := GetApplicationDataPath()
	token, err := os.ReadFile(path + "/" + local_history_file_name)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

var ErrUnableToPersistLocalHistory = fmt.Errorf("unable to persist local history")

func (localHistory) Save(envFile string) error {
	path := GetApplicationDataPath()
	some := os.WriteFile(path+"/"+local_history_file_name, []byte(envFile), 0644)
	if some != nil {
		return ErrUnableToPersistLocalHistory
	}
	return nil
}
