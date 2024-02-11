package cmd

import (
	"envi/internal/domain"
	"envi/internal/storage"
	"envi/internal/ui"
	"errors"
	"fmt"
	"log"
	"os"

	E "github.com/IBM/fp-go/either"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var PullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls the latest .env file from the remote server and replaces the local .env file with it.",
	RunE:  PullCmdFunc,
}

func PullCmdFunc(cmd *cobra.Command, args []string) error {

	accessToken, err := E.Unwrap(GetAccessToken("", storage.GetApplicationDataPath()))
	if err != nil {
		log.Fatal(err)
	}

	callbackURL := GetProviderURL()
	if err != nil {
		log.Fatal(err)
	}

	doneFn := ui.ProgressBar("Fetching remote .env file...")
	remoteEnvValues, err := fetchRemoteEnvValues(callbackURL, accessToken)
	doneFn()

	if err != nil {
		log.Fatal(err)
	}

	localEnvFile, err := getCurrentEnvValues()
	if err != nil {
		log.Fatal(err)
	}

	diffEnvValues(localEnvFile, remoteEnvValues)

	return nil
}

var ErrEnvFileNotFound = errors.New("ErrEnvFileNotFound - env file not found")

func getCurrentEnvValues() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	envFile, err := os.ReadFile(currentDir + "/.env")

	if errors.Is(err, os.ErrNotExist) {
		return "", ErrEnvFileNotFound
	}

	if err != nil {
		return "", err
	}

	return string(envFile), nil
}

func fetchRemoteEnvValues(callbackUrl, accessToken string) (string, error) {
	response, err := resty.New().
		R().
		EnableTrace().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		SetResult(&map[string]interface{}{}).
		Get(callbackUrl)

	if err != nil {
		return "", err
	}
	if response.IsError() {
		return "", fmt.Errorf("fetch not ok: %s", response.String())
	}
	responseAsObject := *response.Result().(*map[string]interface{})
	return responseAsObject["data"].(string), nil
}

func diffEnvValues(local string, remote string) {
	str := domain.DiffEnvs(domain.EnvString(local), domain.EnvString(remote))
	some := str.PrettyPrint()
	println(some)
}
