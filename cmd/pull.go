package cmd

import (
	"envi/internal/domain"
	"envi/internal/llog"
	"envi/internal/provider"
	"envi/internal/storage"
	"envi/internal/ui"
	"envi/internal/utils"
	"errors"
	"fmt"
	"os"
	"time"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	l "github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var PullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls the latest .env file from the remote server and replaces the local .env file with it.",
	RunE:  PullCmdFunc,
}

type EnvSyncState struct {
	RemoteEnvValues string
	LocalEnvValues  string
	DiffRemoteLocal domain.Diff
}

func PullCmdFunc(cmd *cobra.Command, args []string) error {
	err := F.Pipe3(
		SyncEnvState(),
		E.Chain(backupEnvFileIOEither),
		E.Chain(SaveEnvFileIOEither(storage.LocalHistory, os.WriteFile)),
		E.Fold(F.Identity, handleRight),
	)
	if err != nil {
		llog.L.Error(err)
		return nil
	}

	return nil
}

func SyncEnvState() E.Either[error, EnvSyncState] {
	remoteEnvValuesSetter := utils.Setter[string, EnvSyncState]("RemoteEnvValues")
	localEnvValuesSetter := utils.Setter[string, EnvSyncState]("LocalEnvValues")
	diffRemoteLocalSetter := utils.Setter[domain.Diff, EnvSyncState]("DiffRemoteLocal")

	eitherEnvSyncState := F.Pipe3(
		E.Do[error](EnvSyncState{}),
		E.Bind(remoteEnvValuesSetter, fetchRemoteEnvComputation(provider.ZipperFetchRemoteEnvValues)),
		E.Bind(localEnvValuesSetter, getCurrentEnvValuesComputation),
		E.Bind(diffRemoteLocalSetter, diffEnvValuesComputation),
	)

	return eitherEnvSyncState
}

func backupEnvFileIOEither(s EnvSyncState) E.Either[error, EnvSyncState] {
	if s.DiffRemoteLocal.HasNoDiff() {
		// No need to backup
		return E.Right[error](s)
	}
	fileName := fmt.Sprintf(".env.local.backup.%d", time.Now().UnixNano())
	err := os.Rename(".env", fileName)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("No .env file found. Skipping backup.")
		return E.Either[error, EnvSyncState](E.Right[error](s))
	}
	if err != nil {
		return E.Left[EnvSyncState](fmt.Errorf("❌ error while backing up .env file:\n%s", err))
	}
	return E.Right[error](s)
}

type writeFileFn = func(string, []byte, os.FileMode) error

func SaveEnvFileIOEither(storageImp storage.LocalHistorySave, writeFile writeFileFn) func(s EnvSyncState) E.Either[error, EnvSyncState] {
	return func(s EnvSyncState) E.Either[error, EnvSyncState] {
		// Should keep the current .env file in storage.history before saving the new one
		// This is to allow the user to undo the operation
		storageError := storageImp.Save(s.LocalEnvValues)

		if errors.Is(storageError, storage.ErrUnableToPersistLocalHistory) {
			return E.Left[EnvSyncState](storageError)
		}

		err := writeFile(".env", []byte(s.RemoteEnvValues), 0666)
		if err != nil {
			return E.Left[EnvSyncState](fmt.Errorf("❌ error while saving .env file:\n%s", err))
		}
		return E.Right[error](s)
	}
}

func fetchRemoteEnvComputation(fetchRemoteValueImplementation func() (string, error)) func(s EnvSyncState) E.Either[error, string] {
	return func(s EnvSyncState) E.Either[error, string] {
		doneFn := ui.ProgressBar("Fetching remote .env file...")
		defer doneFn()

		// remoteEnvValues, err := fetchRemoteEnvValues(s.CallbackURL, s.AccessToken)
		remoteEnvValues, err := fetchRemoteValueImplementation()
		if err != nil {
			return E.Left[string](err)
		}
		return E.Right[error](remoteEnvValues)
	}
}
func getCurrentEnvValuesComputation(s EnvSyncState) E.Either[error, string] {
	localEnvFile, err := getCurrentEnvValues()
	if err != nil {
		return E.Left[string](err)
	}
	return E.Right[error](localEnvFile)
}
func diffEnvValuesComputation(s EnvSyncState) E.Either[error, domain.Diff] {
	return E.Right[error](diffEnvValues(s.LocalEnvValues, s.RemoteEnvValues))
}
func handleRight(s EnvSyncState) error {
	showEnvUpdateSuccessMessage(s.DiffRemoteLocal.PrettyPrint())
	return nil
}

func showEnvUpdateSuccessMessage(diffPrintStr string) {
	var styleSuccess = l.NewStyle().
		Foreground(l.Color(llog.Tokens.SuccessColor))

	var styleHint = l.NewStyle().
		Foreground(l.Color(llog.Tokens.HintColor)).
		Padding(0, 1).
		Margin(0, 0, 1, 0)

	if diffPrintStr == "" {
		successMessage := styleSuccess.Render(".env file is already up to date.")
		fmt.Println(successMessage)
		return
	}
	successMessage := styleSuccess.
		Bold(true).
		Margin(0, 0, 1, 0).
		Render(".env file updated successfully.")
	undoHintMessage := styleHint.Render(
		"To undo this operation, use",
		llog.StyleCommand().Render("`envi undo`"))

	message := l.JoinVertical(l.Left, successMessage, diffPrintStr, undoHintMessage)

	fmt.Println(message)
}

var ErrEnvFileNotFound = errors.New("- env file not found")
var ErrUnableToCreateEnvFile = errors.New("- unable to create env file")

func getCurrentEnvValues() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	envFile, err := os.ReadFile(currentDir + "/.env")

	if errors.Is(err, os.ErrNotExist) {
		return string(""), nil
	}

	if err != nil {
		return "", err
	}

	return string(envFile), nil
}

func diffEnvValues(local string, remote string) domain.Diff {
	return domain.DiffEnvs(domain.EnvString(local), domain.EnvString(remote))
}
