package env_util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValidKeys(t *testing.T) {
	tasks := []struct {
		name string
		want []string
		env  EnvString
	}{
		{
			name: "base",
			want: []string{"A", "B", "C"},
			env:  "A=1\nB=2\nC=3\n",
		},
		{
			name: "empty",
			want: []string{},
			env:  "",
		},
		{
			name: "comment",
			env:  "#A=1\nB=2\nC=3\n",
			want: []string{"B", "C"},
		},
		{
			name: "malformed",
			env:  "A1\nB=2\nC=3\n",
			want: []string{"B", "C"},
		},
		{
			name: "emptyLine",
			env:  "A=1\n\nB=2\nC=3\n\n",
			want: []string{"A", "B", "C"},
		},
	}

	for _, task := range tasks {
		t.Run(task.name, func(t *testing.T) {
			got := GetValidKeys(task.env)
			assert.Equal(t, task.want, got)
		})
	}

}

func TestDiffEnvs(t *testing.T) {
	type test struct {
		local  EnvString
		remote EnvString
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

func TestDiffPrettyPrint(t *testing.T) {

	tasks := []struct {
		name string
		want string
		diff Diff
	}{
		{
			name: "base",
			want: `+ D
- C`,
			diff: Diff{additions: []string{"D"}, deletions: []string{"C"}},
		},
		{
			name: "addition",
			want: `+ D`,
			diff: Diff{additions: []string{"D"}, deletions: []string{}},
		},
		{
			name: "deletion",
			want: `- C`,
			diff: Diff{additions: []string{}, deletions: []string{"C"}},
		},
		{
			name: "empty",
			want: "",
			diff: Diff{additions: []string{}, deletions: []string{}},
		},
	}
	for _, task := range tasks {
		t.Run(task.name, func(t *testing.T) {
			got := task.diff.PrettyPrint()
			assert.Equal(t, task.want, got)
		})
	}
}
