package cmd

import (
	"envi/internal/storage"
	"fmt"

	l "github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var AuthLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout to envi",
	RunE:  AuthLogOutCmdRunE,
}

func init() {
	AuthLogoutCmd.Flags().StringP("token", "t", "", "Token")
}

func AuthLogOutCmdRunE(cmd *cobra.Command, args []string) error {
	storage.GetApplicationDataPath()
	showLogoutSuccessMessage()
	return nil
}

func showLogoutSuccessMessage() {
	_ = storage.AccessToken.Clear()

	var styleSuccess = l.NewStyle().
		Bold(true).
		Foreground(l.Color("#4CAF50"))

	message := styleSuccess.Render("Logout successful. You have been disconnected from the Zipper provider.")

	fmt.Println(message)
}