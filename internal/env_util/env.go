package env_util

import (
	"strings"

	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	S "github.com/IBM/fp-go/string"
	"github.com/charmbracelet/lipgloss"
)

type EnvString string

func GetValidKeys(env EnvString) []string {
	return F.Pipe4(
		string(env),
		split("\n"),
		A.Filter(func(s string) bool { return !emptyLine(s) && !commentLine(s) && !malformedLine(s) }),
		A.Map(split("=")),
		A.Map(func(a []string) string { return a[0] }),
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
	additions []string
	deletions []string
}

func (d Diff) PrettyPrint() string {
	withAdditionSigh := func(s string) string { return "+ " + s }
	withDeletionSigh := func(s string) string { return "- " + s }

	BgRed := lipgloss.NewStyle().Foreground(lipgloss.Color("#F15C93"))
	BgGreen := lipgloss.NewStyle().Foreground(lipgloss.Color("#C1F3AB"))

	additions := F.Pipe3(
		d.additions,
		A.Map(withAdditionSigh),
		S.Join("\n"),
		func(s string) string { return BgGreen.Render(s) },
	)
	deletions := F.Pipe3(
		d.deletions,
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
	localKeys := GetValidKeys(local)
	remoteKeys := GetValidKeys(remote)

	var localMinusRemote = []string{}
	for _, k := range localKeys {
		if !Includes(k, remoteKeys) {
			localMinusRemote = append(localMinusRemote, k)
		}
	}
	var remoteMinusLocal = []string{}
	for _, k := range remoteKeys {
		if !Includes(k, localKeys) {
			remoteMinusLocal = append(remoteMinusLocal, k)
		}
	}

	return Diff{
		additions: remoteMinusLocal,
		deletions: localMinusRemote,
	}
}
