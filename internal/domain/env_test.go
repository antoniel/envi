package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValidRow(t *testing.T) {
	tasks := []struct {
		name string
		want []string
		env  EnvString
	}{
		{
			name: "base",
			want: []string{"A=1", "B=2", "C=3"},
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
			want: []string{"B=2", "C=3"},
		},
		{
			name: "malformed",
			env:  "A1\nB=2\nC=3\n",
			want: []string{"B=2", "C=3"},
		},
		{
			name: "emptyLine",
			env:  "A=1\n\nB=2\nC=3\n\n",
			want: []string{"A=1", "B=2", "C=3"},
		},
	}

	for _, task := range tasks {
		t.Run(task.name, func(t *testing.T) {
			got := GetValidRows(task.env)
			// fmt.Print(EQ.Equals[]())
			assert.Equal(t, task.want, got)
		})
	}

}
