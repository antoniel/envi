package cmd

import (
	"errors"
	"net/url"
	"strings"

	O "github.com/IBM/fp-go/option"
)

var ErrEnviConfigFileNotFound = errors.New("ErrEnviConfigFileNotFound - envii config file not found")

func getEnviConfigFile(readFile func(string) ([]byte, error)) ([]byte, error) {
	envFile, err := readFile(".envi")
	if err != nil {
		return nil, ErrEnviConfigFileNotFound
	}
	return envFile, nil
}

var ErrMissingAccessToken = errors.New("ErrMissingAccessToken - missing ACCESS_TOKEN")

var ErrInvalidCallbackURL = errors.New("invalid callback URL")

func GetProviderURL() string {
	return "https://envii.zipper.run/api"
}

func isValidURL(urlString string) bool {
	u, err := url.ParseRequestURI(urlString)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func predicateAccessToken(x string) func(line string) bool {
	return func(line string) bool {
		return strings.HasPrefix(line, x)
	}
}
func getEnviConfigValueByKey(key string) func(line string) O.Option[string] {
	return func(line string) O.Option[string] {
		tokens := strings.Split(line, "=")
		if len(tokens) == 2 && tokens[0] == key {
			return O.Some(tokens[1])
		}
		return O.None[string]()
	}
}
