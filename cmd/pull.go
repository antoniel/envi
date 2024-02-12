package cmd

import (
	"envi/internal/domain"
	"envi/internal/llog"
	"envi/internal/provider"
	"envi/internal/storage"
	"envi/internal/ui"
	"errors"
	"fmt"
	"log"
	"os"

	E "github.com/IBM/fp-go/either"
	l "github.com/charmbracelet/lipgloss"
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

	callbackURL := provider.GetZipperProviderDefaultUrl()
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

	diff := diffEnvValues(localEnvFile, remoteEnvValues)
	diffPrintStr := diff.PrettyPrint()
	showEnvUpdateSuccessMessage(diffPrintStr)

	return nil
}
func showEnvUpdateSuccessMessage(diffPrintStr string) {
	var styleSuccess = l.NewStyle().
		Bold(true).
		Foreground(l.Color("#4CAF50")).
		Padding(0, 1).
		Margin(0, 0, 1, 0)

	var styleHint = l.NewStyle().
		Foreground(l.Color("#6272A4")). // Cor mais escura para a dica
		Padding(0, 1).
		Margin(1, 0, 1, 0)

	successMessage := styleSuccess.Render(".env file updated successfully.")
	undoHint := styleHint.Render("To undo this operation, use", llog.StyleCommand().Render("`envi undo`."))

	message := l.JoinVertical(l.Left, successMessage, diffPrintStr, undoHint)

	fmt.Println(message)
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

func diffEnvValues(local string, remote string) domain.Diff {
	return domain.DiffEnvs(domain.EnvString(local), domain.EnvString(remote))
}
