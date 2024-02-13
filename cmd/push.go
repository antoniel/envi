package cmd

import (
	"envi/internal/llog"
	"envi/internal/ui"
	"fmt"
	"time"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var PushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushes the latest .env file to the remote server.",
	RunE:  PushCmdFunc,
}

func PushCmdFunc(cmd *cobra.Command, args []string) error {
	err := F.Pipe3(
		SyncEnvState(),
		E.Chain(validateIfEnvFileHasChanges),
		E.Chain(func(s EnvSyncState) E.Either[error, string] {
			done := ui.ProgressBar("Pushing .env file to remote server")
			eitherPushLocalEnvValues := pushLocalEnvValues(s.AccessToken, s.CallbackURL, s.LocalEnvValues)
			done()
			return eitherPushLocalEnvValues
		}),
		E.Fold(
			F.Identity,
			showMessageSuccess,
		),
	)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}

func validateIfEnvFileHasChanges(s EnvSyncState) E.Either[error, EnvSyncState] {
	if s.DiffRemoteLocal.HasNoDiff() {
		return E.Left[EnvSyncState, error](fmt.Errorf(llog.SuccessStyle().Render("✅ No changes detected in .env file.")))
	}
	return E.Right[error](s)
}

func showMessageSuccess(_ string) error {
	fmt.Println(llog.SuccessStyle().Render("✅ .env file pushed successfully"))
	return nil
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

func pushLocalEnvValues(accessToken string, callbackUrl string, localEnvValues string) E.Either[error, string] {
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
