package provider

import (
	"errors"
	"os"
	"path/filepath"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"gopkg.in/yaml.v3"
)

func GetK8sProviderDefaultUrl() string {
	return "https://k8s-provider.zipper.run/api"
}

// Loads the remote env value from a path locally
func K8sPullRemoteEnvValuesConstructor(k8sValuesPath string) func() (string, error) {
	return func() (string, error) {
		getK8sFileContent(k8sValuesPath)
		return E.Unwrap(F.Pipe2(
			k8sValuesPath,
			getK8sFileContent,
			E.Chain(parseK8sValues),
		))
	}
}
func getK8sFileContent(k8sValuesPath string) E.Either[error, []byte] {
	joinErr := func(errs ...error) error {
		return errors.Join(
			append([]error{errors.New("GetK8sFileContent ")}, errs...)...)
	}
	absPath, err := getAbsolutePath(k8sValuesPath)
	if err != nil {
		return E.Left[[]byte](joinErr(err))
	}
	file, err := os.Stat(absPath)
	if err != nil {
		return E.Left[[]byte](joinErr(errors.New("❌ unable to get k8s values path"), err))
	}
	if filepath.Ext(absPath) != ".yaml" {
		return E.Left[[]byte](joinErr(errors.New("❌ k8s values path is not a yaml file")))
	}
	if file.IsDir() {
		return E.Left[[]byte](errors.New("❌ k8s values path is a directory"))
	}
	b, err := os.ReadFile(absPath)
	if err != nil {
		return E.Left[[]byte](errors.Join(errors.New("❌ unable to read k8s values path"), err))
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
	Env map[string]interface{} `yaml:"env"`
}

func parseK8sValues(k8sValues []byte) E.Either[error, string] {
	var k8sValuesMap k8sValuesEnv
	err := yaml.Unmarshal((k8sValues), &k8sValuesMap)
	if err != nil {
		return E.Left[string](errors.New("❌ unable to parse k8s values\n ℹ️ expected a yaml file with a 'env' key"))
	}

	var envValues string
	for k, v := range k8sValuesMap.Env {
		envValues += k + "=" + v.(string) + "\n"
	}
	return E.Right[error](envValues)
}
