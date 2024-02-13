package cmd

import (
	"envi/internal/env_util"
	"strings"
	"testing"

	F "github.com/IBM/fp-go/function"
	"github.com/stretchr/testify/assert"
)

type Diff struct {
	additions []string
	deletions []string
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

func DiffEnvs(local, remote env_util.EnvString) Diff {
	localKeys := env_util.GetValidKeys(local)
	remoteKeys := env_util.GetValidKeys(remote)

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

func TestDiffEnvs(t *testing.T) {
	type test struct {
		local  env_util.EnvString
		remote env_util.EnvString
		want   Diff
	}

	tests := []test{
		{
			local:  "A=1\nB=2\nC=3\n",
			remote: "A=1\nB=2\nD=4\n",
			want:   Diff{additions: []string{"D"}, deletions: []string{"C"}},
		},
		{
			local:  "A=1\nB=2\nC=3\n",
			remote: "A=1\nB=2\nD=4\n",
			want:   Diff{additions: []string{"D"}, deletions: []string{"C"}},
		},
		{
			local:  "A=1\n",
			remote: "",
			want:   Diff{additions: []string{}, deletions: []string{"A"}},
		},
	}

	for _, tc := range tests {
		t.Run(string(tc.local), func(t *testing.T) {
			got := DiffEnvs(tc.local, tc.remote)
			assert.Equal(t, tc.want, got)
		})
	}
}
