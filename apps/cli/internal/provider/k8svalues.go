package provider

import (
	"context"
	"engov/apps/cli/internal/domain"
	"engov/apps/cli/internal/llog"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	"gopkg.in/yaml.v3"
)

func GetK8sProviderDefaultUrl() string {
	return "https://k8s-provider.zipper.run/api"
}

type K8sPullOption interface {
	Apply(s *k8sPullOpts)
}

// Loads the remote env value from a path locally
type k8sPullOpts struct {
	PullSecrets bool
}

type WithPullSecrets struct {
	Enabled bool
}

func (w WithPullSecrets) Apply(s *k8sPullOpts) {
	s.PullSecrets = w.Enabled
}

func K8sPullRemoteEnvValuesConstructor(
	k8sValuesPath string,
	secretValuesPath string,
	opts ...K8sPullOption) func() (string, error) {
	return func() (string, error) {
		var allOpts k8sPullOpts
		for _, opt := range opts {
			opt.Apply(&allOpts)
		}

		maybeSecretsEnv := maybeGetSecrets(secretValuesPath, allOpts.PullSecrets)
		localK8sFile := F.Pipe2(
			k8sValuesPath,
			getFileContent,
			E.Chain(parseK8sValues),
		)
		mergedEnvs, err := E.Unwrap(mergeEnvs(localK8sFile, maybeSecretsEnv))
		return string(mergedEnvs), err
	}
}

func mergeEnvs(
	envStrFromLocalK8s E.Either[error, domain.EnvString],
	maybeEnvStrFromSecrets O.Option[domain.EnvString]) E.Either[error, domain.EnvString] {

	k8sEnvStr, err := E.Unwrap(envStrFromLocalK8s)
	if err != nil {
		return envStrFromLocalK8s
	}
	secretsEnvStr, ok := O.Unwrap(maybeEnvStrFromSecrets)
	if !ok {
		return envStrFromLocalK8s
	}
	mergedEnvs := string(k8sEnvStr) + "\n# Secrets 🔪\n" + string(secretsEnvStr)

	return E.Right[error](domain.EnvString(mergedEnvs))
}

func getFileContent(filePath string) E.Either[error, []byte] {
	joinErr := func(errs ...error) error {
		return errors.Join(
			append([]error{errors.New("😭 GetK8sFileContent - Unable to resolve path ")}, errs...)...)
	}
	absPath, err := getAbsolutePath(filePath)
	if err != nil {
		return E.Left[[]byte](joinErr(err))
	}
	file, err := os.Stat(absPath)
	if err != nil {
		return E.Left[[]byte](joinErr(err))
	}
	if filepath.Ext(absPath) != ".yaml" {
		return E.Left[[]byte](joinErr(err))
	}
	if file.IsDir() {
		return E.Left[[]byte](joinErr(err))
	}
	b, err := os.ReadFile(absPath)
	if err != nil {
		return E.Left[[]byte](errors.Join(err))
	}
	return E.Right[error](b)
}

func getAbsolutePath(k8sValuesPath string) (string, error) {
	var absPath string
	if filepath.IsAbs(k8sValuesPath) {
		return k8sValuesPath, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	absPath = filepath.Join(wd, k8sValuesPath)
	return absPath, nil
}

type k8sValuesEnv struct {
	Env map[string]string `yaml:"env"`
}

func parseK8sValues(k8sValues []byte) E.Either[error, domain.EnvString] {
	var k8sValuesMap k8sValuesEnv
	err := yaml.Unmarshal((k8sValues), &k8sValuesMap)
	if err != nil {
		return E.Left[domain.EnvString](errors.New("❌ unable to parse k8s values\n ℹ️ expected a yaml file with a 'env' key"))
	}

	var envValues string
	for k, v := range k8sValuesMap.Env {
		envValues += k + "=" + v + "\n"
	}
	return E.Right[error](domain.EnvString(envValues))
}

func maybeGetSecrets(secretValuesPath string, enabled bool) O.Option[domain.EnvString] {
	if !enabled {
		return O.None[domain.EnvString]()
	}

	secretsDeclarationMap, err := E.Unwrap(F.Pipe2(
		getFileContent(secretValuesPath),
		E.Chain(parseSecretsDeclarationFile),
		E.Map[error](getEnvSecretsFromDeclarationMap),
	))

	if err != nil {
		return O.None[domain.EnvString]()
	}

	return secretsDeclarationMap
}

// [EnvName] : [SecretVersionName]
type secretsDeclarationMap map[string]string

func parseSecretsDeclarationFile(file []byte) E.Either[error, secretsDeclarationMap] {
	var secretsEnv struct {
		Secrets []struct {
			VersionName string `yaml:"versionName"`
			Env         string `yaml:"env"`
		}
	}
	err := yaml.Unmarshal(file, &secretsEnv)
	if err != nil {
		return E.Left[secretsDeclarationMap](errors.New("❌ unable to parse secrets declaration file\n ℹ️ expected a yaml file with a 'secrets' key"))
	}

	var secretsDeclarationMap secretsDeclarationMap = make(secretsDeclarationMap)
	for _, v := range secretsEnv.Secrets {
		envName := v.Env
		versionName := v.VersionName
		secretsDeclarationMap[envName] = versionName
	}
	return E.Right[error](secretsDeclarationMap)
}

func getEnvSecretsFromDeclarationMap(secretsDeclarationMap secretsDeclarationMap) O.Option[domain.EnvString] {
	var envString string
	var envStringMutex sync.Mutex
	wg := sync.WaitGroup{}
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		fmt.Println(err)
		return O.None[domain.EnvString]()
	}
	defer client.Close()
	for envName, secretVersion := range secretsDeclarationMap {
		wg.Add(1)
		go func(envName string, secretVersion string) {
			defer wg.Done()

			result, err := client.AccessSecretVersion(
				ctx, &secretmanagerpb.AccessSecretVersionRequest{Name: secretVersion})

			if err != nil {
				llog.L.Error("Unable to get secret value", "secret", envName)
				return
			}
			envStringMutex.Lock()
			envString += envName + "=" + string(result.Payload.Data) + "\n"
			envStringMutex.Unlock()
		}(envName, secretVersion)
	}

	wg.Wait()

	return O.Some(domain.EnvString(envString))
}
