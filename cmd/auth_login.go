package cmd

import (
	"envi/internal/llog"
	"envi/internal/provider"
	"envi/internal/storage"

	E "github.com/IBM/fp-go/either"
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

	if E.IsRight(provider.PersistToken(applicationPath, maybeTokenFromFlag)) {
		showSuccessMessage()
		return nil
	}

	_, err := E.Unwrap(provider.GetOrAskAndPersistToken(applicationPath))
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
		Foreground(lipgloss.Color(llog.Tokens.ForegroundColor)).
		PaddingLeft(2).
		MarginBottom(1)

	title := styleTitle.Render("Authenticated! Quick Commands:")
	cmdPull := styleCommand.PaddingLeft(2).Render("`envi pull`:")
	cmdAuth := styleCommand.PaddingLeft(2).Render("`envi auth`:")
	helpText := styleText.Render("For more, `envi --help`.")

	message := lipgloss.JoinVertical(lipgloss.Left, title,
		lipgloss.JoinHorizontal(lipgloss.Left, cmdPull, styleText.Render("Sync .env files.")),
		lipgloss.JoinHorizontal(lipgloss.Left, cmdAuth, styleText.Render("Manage account.")),
		helpText,
	)

	println(message)
}

func AuthIsLogged() bool {
	persistedToken, err := storage.AccessToken.Get()
	if err != nil {
		return false
	}
	return persistedToken != ""
}
