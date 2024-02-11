package cmd

import (
	"envi/internal/llog"
	"envi/internal/storage"
	form "envi/internal/ui"
	"fmt"
	"os"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var AuthLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to envi",
	RunE:  AuthCmdRunE,
}

func init() {
	AuthLoginCmd.Flags().StringP("token", "t", "", "Token")
}

func AuthCmdRunE(cmd *cobra.Command, args []string) error {
	applicationPath := storage.GetApplicationDataPath()
	maybeTokenFromFlag := cmd.Flag("token").Value.String()

	accessTokenE := F.Pipe1(
		GetAccessToken(maybeTokenFromFlag, applicationPath),
		E.Chain(F.Bind1st(persistToken, applicationPath)),
	)

	_, err := E.Unwrap(accessTokenE)
	if err != nil {
		return err
	}
	showSuccessMessage()
	return nil
}

func showSuccessMessage() {
	var styleTitle = llog.StyleTitle()
	var styleCommand = llog.StyleCommand()

	var styleText = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		PaddingLeft(2).
		MarginBottom(1)

	title := styleTitle.Render("Authenticated! Quick Commands:")
	cmdPull := styleCommand.Render("`envi pull`:")
	cmdAuth := styleCommand.Render("`envi auth`:")
	helpText := styleText.Render("For more, `envi --help`.")

	message := lipgloss.JoinVertical(lipgloss.Left, title,
		lipgloss.JoinHorizontal(lipgloss.Left, cmdPull, styleText.Render("Sync .env files.")),
		lipgloss.JoinHorizontal(lipgloss.Left, cmdAuth, styleText.Render("Manage account.")),
		helpText,
	)

	println(message)
}

func requestTokenToUser() string {
	response := form.NewPrompt([]form.Question{
		form.NewQuestion("Enter your token: ").
			WithEchoMode(textinput.EchoPassword),
	})
	return response[0]
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

	tokenFromUserInput := requestTokenToUser()
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
		return "", err
	}
	return string(token), nil
}

func AuthIsLogged() bool {
	persistedToken, err := getPersistedToken(storage.GetApplicationDataPath())
	if err != nil {
		return false
	}
	return persistedToken != ""
}
