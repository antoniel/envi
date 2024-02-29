package storage

import (
	"fmt"
	"os"
)

const token_file_name = "token"

type accessToken struct{}

var AccessToken = accessToken{}

func (accessToken) Get() (string, error) {
	path := GetApplicationDataPath()
	token, err := os.ReadFile(path + "/" + token_file_name)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

var ErrUnableToPersistToken = fmt.Errorf("unable to persist token")
var ErrInvalidToken = fmt.Errorf("invalid token")

func (accessToken) Save(token string) error {
	path := GetApplicationDataPath()
	if token == "" {
		// return ErrInvalidToken,
		return ErrInvalidToken
	}
	some := os.WriteFile(path+"/"+token_file_name, []byte(token), 0644)
	if some != nil {
		return (ErrUnableToPersistToken)
	}
	return nil
}

func (accessToken) Clear() error {
	path := GetApplicationDataPath()
	err := os.Remove(path + "/" + token_file_name)
	if err != nil {
		return err
	}
	return nil
}
