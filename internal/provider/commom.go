package provider

import (
	"envi/internal/storage"
	"envi/internal/ui"
	"errors"
	"fmt"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"

	"github.com/charmbracelet/bubbles/textinput"
)

func getAccessToken(applicationDataPath string) E.Either[error, string] {
	persistedToken, err := storage.AccessToken.Get()
	if err == nil {
		return E.Right[error](persistedToken)
	}

	tokenFromUserInput := requestTokenToUser()
	if tokenFromUserInput != "" {
		return E.Right[error](tokenFromUserInput)
	}

	return E.Left[string](ErrTokenNotProvided)
}

func requestTokenToUser() string {
	response := ui.NewPrompt([]ui.Question{
		ui.NewQuestion("Enter your token: ").
			WithEchoMode(textinput.EchoPassword),
	})
	return response[0]
}

var ErrTokenNotProvided = fmt.Errorf("token not provided")

func PersistToken(path string, token string) E.Either[error, string] {
	return E.FromError(storage.AccessToken.Save)(token)
}
func GetOrAskAndPersistToken(applicationDataPath string) E.Either[error, string] {
	path := storage.GetApplicationDataPath()
	persistTokenFn := F.Bind1st(PersistToken, path)
	return F.Pipe1(
		getAccessToken(storage.GetApplicationDataPath()),
		E.Chain(persistTokenFn),
	)
}

var ErrUnableToFetchRemoteEnvValues = errors.New("- Unable to fetch remote env values, check your connection and credentials, then try again")