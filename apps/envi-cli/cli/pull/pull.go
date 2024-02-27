package pull

import (
	"envi/apps/envi-cli/internal/domain"
	"envi/apps/envi-cli/internal/llog"
	"envi/apps/envi-cli/internal/provider"
	"envi/apps/envi-cli/internal/storage"
	"envi/apps/envi-cli/internal/ui"
	"envi/apps/envi-cli/internal/utils"
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

func GetPullFn(cmd *cobra.Command) (provider.PullFn, error) {
	noop := func() (string, error) {
		return "", nil
	}
	validProviders := []string{"zipper", "k8s"}
	providerType := cmd.Flag("provider").Value.String()
	k8sValuesPath := cmd.Flag("k8s-values-path").Value.String()

	if !utils.Contains(validProviders, providerType) {
		return nil, errors.New("❌ invalid provider type")
	}

	if providerType == "k8s" {
		if k8sValuesPath == "" {
			return noop, errors.New("❌ k8s-values-path flag is required when using k8s provider")
		}
		return provider.K8sPullRemoteEnvValuesConstructor(k8sValuesPath), nil
	}

	return provider.ZipperPullRemoteEnvValues, nil
}

func PullCmdFunc(cmd *cobra.Command, args []string) error {
	pullFn, errPullFn := GetPullFn(cmd)

	if errPullFn != nil {
		return errPullFn
	}

	err := F.Pipe3(
		SyncEnvState(pullFn),
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

func SyncEnvState(pullFn provider.PullFn) E.Either[error, EnvSyncState] {
	remoteEnvSetter := utils.Setter[string, EnvSyncState]("RemoteEnvValues")
	localEnvSetter := utils.Setter[string, EnvSyncState]("LocalEnvValues")
	diffEnvsSetter := utils.Setter[domain.Diff, EnvSyncState]("DiffRemoteLocal")

	eitherEnvSyncState := F.Pipe3(
		E.Do[error](EnvSyncState{}),
		E.Bind(remoteEnvSetter, fetchRemoteEnvComputation(pullFn)),
		E.Bind(localEnvSetter, getLocalEnvComputation),
		E.Bind(diffEnvsSetter, diffEnvsComputation),
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

		remoteEnvValues, err := fetchRemoteValueImplementation()
		if err != nil {
			return E.Left[string](err)
		}
		return E.Right[error](remoteEnvValues)
	}
}
func getLocalEnvComputation(s EnvSyncState) E.Either[error, string] {
	localEnvFile, err := getCurrentEnvValues()
	if err != nil {
		return E.Left[string](err)
	}
	return E.Right[error](localEnvFile)
}
func diffEnvsComputation(s EnvSyncState) E.Either[error, domain.Diff] {
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
