package env_util

import (
	"strings"

	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	S "github.com/IBM/fp-go/string"
	"github.com/charmbracelet/lipgloss"
)

type EnvString string

func GetValidRows(env EnvString) []string {
	return F.Pipe2(
		string(env),
		split("\n"),
		A.Filter(func(s string) bool {
			return !emptyLine(s) &&
				!commentLine(s) &&
				!malformedLine(s)
		}),
	)
}

func emptyLine(s string) bool {
	return len(s) == 0
}
func commentLine(s string) bool {
	return strings.HasPrefix(s, "#")
}
func malformedLine(s string) bool {
	if !strings.Contains(s, "=") {
		return true
	}
	return len(strings.Split(s, "=")) != 2
}

var split = F.Curry2(F.Swap(strings.Split))

type Diff struct {
	Additions []string
	Deletions []string
}

func (d Diff) PrettyPrint() string {
	withAdditionSigh := func(s string) string { return "+ " + s }
	withDeletionSigh := func(s string) string { return "- " + s }

	BgRed := lipgloss.NewStyle().Foreground(lipgloss.Color("#F15C93"))
	BgGreen := lipgloss.NewStyle().Foreground(lipgloss.Color("#C1F3AB"))

	additions := F.Pipe3(
		d.Additions,
		A.Map(withAdditionSigh),
		S.Join("\n"),
		func(s string) string { return BgGreen.Render(s) },
	)
	deletions := F.Pipe3(
		d.Deletions,
		A.Map(withDeletionSigh),
		S.Join("\n"),
		func(s string) string { return BgRed.Render(s) },
	)

	emptyAdditions := len(additions) == 0
	emptyDeletions := len(deletions) == 0

	if emptyAdditions && emptyDeletions {
		return ""
	}

	if emptyAdditions && !emptyDeletions {
		return deletions

	}
	if !emptyAdditions && emptyDeletions {
		return additions
	}

	return strings.Join([]string{additions, deletions}, "\n")
}

var Split = F.Curry2(F.Swap(strings.Split))

func Includes(s string, arr []string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}

func DiffEnvs(local, remote EnvString) Diff {
	localRows := GetValidRows(local)
	remoteRows := GetValidRows(remote)

	localSet := make(map[string]string)
	for _, k := range localRows {
		key := getKeysFromEnvRow(k)
		val := getValueFromEnvRow(k)
		localSet[key] = val
	}
	remoteSet := make(map[string]string)
	for _, k := range remoteRows {
		key := getKeysFromEnvRow(k)
		val := getValueFromEnvRow(k)
		remoteSet[key] = val
	}

	var localMinusRemote = []string{}
	var remoteMinusLocal = []string{}

	for k := range localSet {
		if _, found := remoteSet[k]; !found {
			localMinusRemote = append(localMinusRemote, k)
		}
	}

	for k, v := range remoteSet {
		if lv, found := localSet[k]; !found {
			remoteMinusLocal = append(remoteMinusLocal, k)
		} else if v != lv {
			remoteMinusLocal = append(remoteMinusLocal, k)
		}
	}

	return Diff{
		Additions: remoteMinusLocal,
		Deletions: localMinusRemote,
	}
}

func getValueFromEnvRow(row string) string {
	return strings.Split(row, "=")[1]
}

func getKeysFromEnvRow(row string) string {
	return strings.Split(row, "=")[0]
}
