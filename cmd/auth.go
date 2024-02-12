package cmd

import (
	"envi/internal/llog"
	"fmt"

	l "github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Display commands for authenticating envi with an account",
	RunE:  AuthCmdE,
}

func init() {
	AuthCmd.AddCommand(AuthLoginCmd)
	AuthCmd.AddCommand(AuthLogoutCmd)
	AuthCmd.PersistentFlags().StringP("token", "t", "", "Token")
}

func AuthCmdE(cmd *cobra.Command, args []string) error {
	if !AuthIsLogged() {
		fmt.Println(
			llog.StyleTitle().
				MarginBottom(1).
				Render("Session not found"))
		AuthLoginCmd.RunE(cmd, args)
		return nil
	}
	showLoginMessage()
	return nil
}

func showLoginMessage() {

	// Estilos
	var styleTitle = llog.StyleTitle()

	var hint = func() l.Style {
		return l.NewStyle().
			Foreground(l.Color(llog.Tokens.HintColor)).
			PaddingLeft(2).
			MarginBottom(1)
	}
	cmd := l.NewStyle().Foreground(l.Color(llog.Tokens.CommandForegroundColor))

	title := styleTitle.Render("Authenticated with Zipper:")
	loginStatus := hint().Render("You are currently logged using the Zipper provider.")
	logoutMessage := hint().Italic(true).Render("To log out, use", cmd.Render("`envi auth logout`."))

	message := l.JoinVertical(
		l.Left,
		title,
		loginStatus,
		logoutMessage,
	)

	fmt.Println(message)
}
