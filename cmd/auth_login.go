package cmd

import (
	"envi/internal/storage"
	form "envi/internal/ui"
	"fmt"
	"os"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/spf13/cobra"
)

var AuthLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to envi",
	RunE:  AuthCmdRunE,
}

func AuthCmdRunE(cmd *cobra.Command, args []string) error {
	applicationPath := storage.GetApplicationDataPath()
	c2PersistToken := F.Curry2(persistToken)

	maybeTokenFromFlag := cmd.Flag("token").Value.String()

	accessTokenE := F.Pipe1(
		GetAccessToken(maybeTokenFromFlag, applicationPath),
		E.Chain(c2PersistToken(applicationPath)),
	)

	_, err := E.Unwrap(accessTokenE)
	if err != nil {
		return err
	}

	return nil
}

func requestTokenToUser(tokenFromFlag string) string {
	var token string
	token = tokenFromFlag

	if token == "" {
		response := form.NewPrompt([]form.Question{
			form.NewQuestion("Enter your token: ").
				WithEchoMode(textinput.EchoPassword),
		})
		token = response[0]
	}

	return token
}

var ErrUnableToPersistToken = fmt.Errorf("unable to persist token")
var ErrInvalidToken = fmt.Errorf("invalid token")

const TOKEN_FILE_NAME = "token"

func GetAccessToken(maybeTokenFromFlag, applicationDataPath string) E.Either[error, string] {
	if maybeTokenFromFlag != "" {
		return E.Right[error](maybeTokenFromFlag)
	}

	persistedToken, err := getPersistedToken(applicationDataPath)
	if err == nil {
		return E.Right[error](persistedToken)
	}

	tokenFromUserInput := requestTokenToUser(maybeTokenFromFlag)
	if tokenFromUserInput != "" {
		return E.Right[error](tokenFromUserInput)
	}

	return E.Left[string](ErrInvalidToken)
}

func persistToken(path string, token string) E.Either[error, string] {
	if token == "" {
		// return ErrInvalidToken,
		return E.Left[string](ErrInvalidToken)
	}
	some := os.WriteFile(path+"/"+TOKEN_FILE_NAME, []byte(token), 0644)
	if some != nil {
		return E.Left[string](ErrUnableToPersistToken)
	}
	return E.Right[error](token)
}

func getPersistedToken(path string) (string, error) {
	token, err := os.ReadFile(path + "/" + TOKEN_FILE_NAME)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(token), nil
}
