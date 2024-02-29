package auth

import (
	"engov/apps/cli/internal/llog"
	"engov/apps/cli/internal/provider"
	"engov/apps/cli/internal/storage"

	E "github.com/IBM/fp-go/either"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var AuthLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to engov",
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
	cmdPull := styleCommand.PaddingLeft(2).Render("`engov pull`:")
	cmdAuth := styleCommand.PaddingLeft(2).Render("`engov auth`:")
	helpText := styleText.Render("For more, `engov --help`.")

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
