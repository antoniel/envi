package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffEnvs(t *testing.T) {
	type test struct {
		name   string
		local  EnvString
		remote EnvString
		want   Diff
	}

	tests := []test{
		{
			local:  "A=1\nB=2\nC=3\n",
			remote: "A=1\nB=2\nD=3\n",
			want:   Diff{Additions: []string{"D"}, Deletions: []string{"C"}},
			name:   "Additions and deletions",
		},
		{
			local:  "A=1\nB=2\nC=3\n",
			remote: "A=1\nB=2\nC=3\n",
			want:   Diff{Additions: []string{}, Deletions: []string{}},
			name:   "No Changes",
		},
		{
			local:  "A=1\n",
			remote: "",
			want:   Diff{Additions: []string{}, Deletions: []string{"A"}},
			name:   "Empty remote",
		},
		{
			local:  "A=1\nB=2\nC=3\n",
			remote: "A=3\nB=2\nC=1\n",
			want:   Diff{Additions: []string{"A", "C"}, Deletions: []string{}},
			name:   "Change Values",
		},
	}

	for _, tc := range tests {
		t.Run(string(tc.name), func(t *testing.T) {
			got := DiffEnvs(tc.local, tc.remote)
			eq := EqDiff.Equals(tc.want, got)
			fmt.Print(eq)
			assert.Equal(t, eq, true)
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
			diff: Diff{Additions: []string{"D"}, Deletions: []string{"C"}},
		},
		{
			name: "addition",
			want: `+ D`,
			diff: Diff{Additions: []string{"D"}, Deletions: []string{}},
		},
		{
			name: "deletion",
			want: `- C`,
			diff: Diff{Additions: []string{}, Deletions: []string{"C"}},
		},
		{
			name: "empty",
			want: "",
			diff: Diff{Additions: []string{}, Deletions: []string{}},
		},
	}
	for _, task := range tasks {
		t.Run(task.name, func(t *testing.T) {
			got := task.diff.PrettyPrint()
			assert.Equal(t, task.want, got)
		})
	}
}

func TestEqDiff(t *testing.T) {
	tasks := []struct {
		diffX Diff
		diffY Diff
		want  bool
	}{
		{
			diffX: Diff{Additions: []string{"D"}, Deletions: []string{"C"}},
			diffY: Diff{Additions: []string{"D"}, Deletions: []string{"C"}},
			want:  true,
		},
		{
			diffX: Diff{Additions: []string{"D"}, Deletions: []string{"C"}},
			diffY: Diff{Additions: []string{"D"}, Deletions: []string{}},
			want:  false,
		},
		{
			diffX: Diff{Additions: []string{"D"}, Deletions: []string{"C"}},
			diffY: Diff{Additions: []string{""}, Deletions: []string{"C"}},
			want:  false,
		},
		{
			diffX: Diff{Additions: []string{"D"}, Deletions: []string{"C"}},
			diffY: Diff{Additions: []string{"D"}, Deletions: []string{"C", "D"}},
			want:  false,
		},
		{
			diffX: Diff{Additions: []string{"D"}, Deletions: []string{"C"}},
			diffY: Diff{Additions: []string{"D", "C"}, Deletions: []string{"C"}},
			want:  false,
		},
		{
			diffX: Diff{Additions: []string{"C", "D"}, Deletions: []string{"C"}},
			diffY: Diff{Additions: []string{"D", "C"}, Deletions: []string{"C"}},
			want:  true,
		},
		{
			diffX: Diff{Additions: []string{"A", "B", "C", "D"}, Deletions: []string{"C"}},
			diffY: Diff{Additions: []string{"B", "C", "D", "A"}, Deletions: []string{"C"}},
			want:  true,
		},
	}

	for _, task := range tasks {
		t.Run("EqDiff", func(t *testing.T) {
			got := EqDiff.Equals(task.diffX, task.diffY)
			assert.Equal(t, task.want, got)
		})
	}

}
