package cmd

import (
	"envi/internal/domain"
	"envi/internal/llog"
	"envi/internal/provider"
	"envi/internal/storage"
	"envi/internal/ui"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	l "github.com/charmbracelet/lipgloss"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var PullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls the latest .env file from the remote server and replaces the local .env file with it.",
	RunE:  PullCmdFunc,
}

func GeneralSetter[T any, S any](fieldName string, fieldValue T, state S) S {
	val := reflect.ValueOf(state)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Make sure we're dealing with a struct
	if val.Kind() != reflect.Struct {
		panic("❌ State must be a struct or a pointer to struct")
	}

	// Make a copy of the struct to avoid mutating the original
	newState := reflect.New(val.Type()).Elem()
	newState.Set(val)

	// Set the field value
	fld := newState.FieldByName(fieldName)
	if fld.IsValid() && fld.CanSet() {
		fld.Set(reflect.ValueOf(fieldValue))
	} else {
		panic(fmt.Sprintf("🚨 Field `%s` not found or not settable, error at `GeneralSetter`", fieldName))
	}

	return newState.Interface().(S)
}

func setter[T, S any](fieldName string) func(T) func(S) S {
	type TypeOfGeneralSetter = func(string, T, S) S
	generalSetterBound := F.Bind1of3[TypeOfGeneralSetter](GeneralSetter)(fieldName)
	type TypeOfGeneralSetterBound = func(T, S) S
	return F.Curry2[TypeOfGeneralSetterBound](generalSetterBound)
}

type EnvResultState struct {
	AccessToken     string
	CallbackURL     string
	RemoteEnvValues string
	LocalEnvValues  string
	DiffRemoteLocal domain.Diff
}

func PullCmdFunc(cmd *cobra.Command, args []string) error {
	accessTokenSetter := setter[string, EnvResultState]("AccessToken")
	callbackURLSetter := setter[string, EnvResultState]("CallbackURL")
	remoteEnvValuesSetter := setter[string, EnvResultState]("RemoteEnvValues")
	localEnvValuesSetter := setter[string, EnvResultState]("LocalEnvValues")
	diffRemoteLocalSetter := setter[domain.Diff, EnvResultState]("DiffRemoteLocal")

	err := F.Pipe8(
		E.Do[error](EnvResultState{}),
		E.Bind(accessTokenSetter, getAccessTokenComputation),
		E.Bind(callbackURLSetter, getCallBackUrlComputation),
		E.Bind(remoteEnvValuesSetter, fetchRemoteEnvValuesComputation),
		E.Bind(localEnvValuesSetter, getCurrentEnvValuesComputation),
		E.Bind(diffRemoteLocalSetter, diffEnvValuesComputation),
		E.Chain(backupEnvFileIOEither),
		E.Chain(saveEnvFileIOEither),
		E.Fold(F.Identity, handleRight),
	)
	if err != nil {
		llog.L.Error(err)
		return nil
	}

	return nil
}

func backupEnvFileIOEither(s EnvResultState) E.Either[error, EnvResultState] {
	if s.DiffRemoteLocal.HasNoDiff() {
		// No need to backup
		return E.Right[error](s)
	}
	fileName := fmt.Sprintf(".env.local.backup.%d", time.Now().UnixNano())
	err := os.Rename(".env", fileName)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("No .env file found. Skipping backup.")
		return E.Either[error, EnvResultState](E.Right[error](s))
	}
	if err != nil {
		return E.Left[EnvResultState](fmt.Errorf("❌ error while backing up .env file:\n%s", err))
	}
	return E.Right[error](s)
}
func saveEnvFileIOEither(s EnvResultState) E.Either[error, EnvResultState] {
	err := os.WriteFile(".env", []byte(s.RemoteEnvValues), 0666)
	if err != nil {
		return E.Left[EnvResultState](fmt.Errorf("❌ error while saving .env file:\n%s", err))
	}
	return E.Right[error](s)
}
func getAccessTokenComputation(s EnvResultState) E.Either[error, string] {
	path := storage.GetApplicationDataPath()
	persistTokenFn := F.Bind1st(PersistToken, path)
	return F.Pipe1(
		GetAccessToken("", storage.GetApplicationDataPath()),
		E.Chain(persistTokenFn),
	)
}
func getCallBackUrlComputation(s EnvResultState) E.Either[error, string] {
	return E.Right[error](provider.GetZipperProviderDefaultUrl())
}
func fetchRemoteEnvValuesComputation(s EnvResultState) E.Either[error, string] {
	doneFn := ui.ProgressBar("Fetching remote .env file...")
	defer doneFn()
	remoteEnvValues, err := fetchRemoteEnvValues(s.CallbackURL, s.AccessToken)
	if err != nil {
		return E.Left[string](err)
	}
	return E.Right[error](remoteEnvValues)
}
func getCurrentEnvValuesComputation(s EnvResultState) E.Either[error, string] {
	localEnvFile, err := getCurrentEnvValues()
	if err != nil {
		return E.Left[string](err)
	}
	return E.Right[error](localEnvFile)
}
func diffEnvValuesComputation(s EnvResultState) E.Either[error, domain.Diff] {
	return E.Right[error](diffEnvValues(s.LocalEnvValues, s.RemoteEnvValues))
}
func handleRight(s EnvResultState) error {
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

var ErrUnableToFetchRemoteEnvValues = errors.New("- Unable to fetch remote env values, check your connection and credentials, then try again")

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

		return "", ErrUnableToFetchRemoteEnvValues
	}
	responseAsObject := *response.Result().(*map[string]interface{})
	return responseAsObject["data"].(string), nil
}

func diffEnvValues(local string, remote string) domain.Diff {
	return domain.DiffEnvs(domain.EnvString(local), domain.EnvString(remote))
}
