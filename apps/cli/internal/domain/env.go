package domain

import (
	"strings"

	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
)

type _EnvString string
type EnvString = _EnvString

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
