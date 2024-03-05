package domain

import (
	"envii/apps/cli/internal/llog"
	"sort"
	"strings"

	A "github.com/IBM/fp-go/array"
	EQ "github.com/IBM/fp-go/eq"
	F "github.com/IBM/fp-go/function"
	S "github.com/IBM/fp-go/string"
	"github.com/charmbracelet/lipgloss"
)

type Diff struct {
	Additions []string
	Deletions []string
}

func uniqueElementsSliceEquals[T comparable](x, y []T) bool {
	if len(x) != len(y) {
		return false
	}
	mapX := make(map[T]bool)
	for _, k := range x {
		mapX[k] = true
	}
	mapY := make(map[T]bool)
	for _, k := range y {
		mapY[k] = true
	}
	for k := range mapX {
		if _, found := mapY[k]; !found {
			return false
		}
	}
	return true
}
func (Diff) Equals(x, y Diff) bool {
	stringSetEq := EQ.FromEquals(uniqueElementsSliceEquals[string])

	return stringSetEq.Equals(x.Additions, y.Additions) &&
		stringSetEq.Equals(x.Deletions, y.Deletions)
}
func (d Diff) HasNoDiff() bool {
	return len(d.Additions) == 0 && len(d.Deletions) == 0
}

var EqDiff = EQ.FromEquals(Diff{}.Equals)

func (d Diff) PrettyPrint() string {
	withAdditionSigh := func(s string) string { return "+ " + s }

	withDeletionSigh := func(s string) string { return "- " + s }
	redForeground := lipgloss.NewStyle().Foreground(lipgloss.Color(llog.Tokens.ErrorColor))
	greenForeground := lipgloss.NewStyle().Foreground(lipgloss.Color(llog.Tokens.SuccessColor))

	additions := F.Pipe3(
		d.Additions,
		A.Map(withAdditionSigh),
		S.Join("\n"),
		func(s string) string {
			if s == "" {
				return ""
			}
			return greenForeground.Render(s)
		},
	)
	deletions := F.Pipe3(
		d.Deletions,
		A.Map(withDeletionSigh),
		S.Join("\n"),
		func(s string) string {
			if s == "" {
				return ""
			}
			return redForeground.Render(s)
		},
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

func getValueFromEnvRow(row string) string {
	return strings.Split(row, "=")[1]
}
func getKeysFromEnvRow(row string) string {
	return strings.Split(row, "=")[0]
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

func MergeEnvsPreservingFirst(a, b EnvString) EnvString {
	validRowsA := GetValidRows(a)
	validRowsB := GetValidRows(b)
	envMap := make(map[string]string)

	for _, row := range validRowsA {
		key := getKeysFromEnvRow(row)
		val := getValueFromEnvRow(row)
		envMap[key] = val
	}

	for _, row := range validRowsB {
		if _, found := envMap[getKeysFromEnvRow(row)]; !found {
			envMap[getKeysFromEnvRow(row)] = getValueFromEnvRow(row)
		}
	}

	envList := []string{}
	for k, v := range envMap {
		envList = append(envList, k+"="+v)
	}

	sort.Strings(envList)
	return sliceToEnvString(envList)
}

func sliceToEnvString(slice []string) EnvString {
	return EnvString(strings.Join(slice, "\n"))
}
