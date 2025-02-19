package provider

import (
	"envii/apps/cli/internal/domain"
	"envii/apps/cli/internal/storage"
	"fmt"
	"time"

	E "github.com/IBM/fp-go/either"
	"github.com/go-resty/resty/v2"
)

// var ZipperProvider = domain.Provider{Name: "Zipper"}
// var ZipperProvider = zipperProvider{Name: "Zipper"}

type zipperProvider struct {
	Name domain.ProviderName
}

func NewZipperProvider() domain.PushPullProvider {
	return zipperProvider{Name: "Zipper"}
}

func (d zipperProvider) GetName() domain.ProviderName {
	return d.Name
}

func (zipperProvider) PullRemoteEnvValues() (domain.EnvString, error) {
	return zipperPullRemoteEnvValues()
}

func (zipperProvider) PushLocalEnvValues(localEnvValues domain.EnvString) error {
	_, err := E.Unwrap(zipperPushLocalEnvsToRemote(localEnvValues))
	return err
}

func getZipperProviderDefaultUrl() string {
	return "https://envii.zipper.run/api"
}

func zipperPullRemoteEnvValues() (domain.EnvString, error) {
	path := storage.GetApplicationDataPath()
	accessToken, err := E.Unwrap(GetOrAskAndPersistToken(path))
	if err != nil {
		return "", err
	}
	callbackUrl := getZipperProviderDefaultUrl()
	response, err := resty.New().
		R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		SetResult(&map[string]interface{}{}).
		SetBody(map[string]interface{}{"cmd": "pull"}).
		Post(callbackUrl)

	if err != nil {
		return "", err
	}
	if response.IsError() {
		return "", ErrUnableToFetchRemoteEnvValues
	}
	responseAsObject := *response.Result().(*map[string]interface{})
	return responseAsObject["data"].(domain.EnvString), nil
}

var errUnableToPushRemoteEnvValues = fmt.Errorf("❌ Unable to push remote environment values")

type envPushResponse struct {
	Ok   bool `json:"ok"`
	Data struct {
		History []struct {
			Date    time.Time `json:"date"`
			EnvFile string    `json:"envFile"`
		} `json:"history"`
	} `json:"data"`
	Meta struct {
		Request struct {
			ExecutionID string `json:"executionId"`
			Timing      string `json:"timing"`
		} `json:"request"`
	} `json:"__meta"`
}

func zipperPushLocalEnvsToRemote(localEnvValues domain.EnvString) E.Either[error, string] {
	path := storage.GetApplicationDataPath()
	accessToken, err := E.Unwrap(GetOrAskAndPersistToken(path))
	if err != nil {
		return E.Left[string](err)
	}
	callbackUrl := getZipperProviderDefaultUrl()

	response, err := resty.New().
		R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		SetBody(map[string]interface{}{"env": localEnvValues, "cmd": "push"}).
		SetResult(&envPushResponse{}).
		Post(callbackUrl)

	if err != nil {
		return E.Left[string](err)
	}
	if response.IsError() {
		return E.Left[string](errUnableToPushRemoteEnvValues)
	}
	responseAsObject := *response.Result().(*envPushResponse)
	return E.Right[error](responseAsObject.Data.History[len(responseAsObject.Data.History)-1].EnvFile)
}
