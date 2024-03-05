package pull_test

import (
	"envii/apps/cli/cmd/pull"
	"os"
	"testing"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/stretchr/testify/assert"
)

type mockStorage struct {
	mockSave func(s string) error
}

func (m mockStorage) Save(s string) error {
	return m.mockSave(s)
}

func newMockStorage(mockSave func(s string) error) mockStorage {
	return mockStorage{mockSave: mockSave}
}

func TestSaveEnvFileIOEither(t *testing.T) {

	t.Run("Should persist the current .env in the local history", func(t *testing.T) {
		didSave := false
		mockStorage := newMockStorage(func(s string) error {
			didSave = true
			return nil
		})
		sut := pull.SaveEnvResultIOEither(mockStorage, func(string, []byte, os.FileMode) error {
			return nil
		})

		F.Pipe1(
			sut(pull.EnvSyncState{}),
			E.Fold(func(e error) error {
				t.Fail() //
				return nil
			}, func(s pull.EnvSyncState) error {
				assert.True(t, didSave)
				return nil
			}),
		)
	})
}
