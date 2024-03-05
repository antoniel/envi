package undo

import (
	"envii/apps/cli/cmd/pull"
	"envii/apps/cli/internal/llog"
	"envii/apps/cli/internal/storage"
	"errors"
	"fmt"
	"os"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/spf13/cobra"
)

var UndoCmd = &cobra.Command{
	Use:   "undo",
	Short: "undoes the last `envii pull` command",
	RunE:  UndoCmdFunc,
}

func UndoCmdFunc(cmd *cobra.Command, args []string) error {
	prevEnv, err := storage.LocalHistory.Get()
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println(llog.ErrorStyle().Render("❌ Unable to undo, no history found"))
		return nil
	}
	if err != nil {
		return err
	}
	F.Pipe2(
		pull.EnvSyncState{RemoteEnvValues: prevEnv},
		pull.SaveEnvResultIOEither(storage.LocalHistory, os.WriteFile),
		E.Fold(
			F.Identity,
			func(_ pull.EnvSyncState) error {
				showMessageUndoSuccessfully()
				return nil
			},
		),
	)

	return nil
}

func showMessageUndoSuccessfully() {
	fmt.Println(llog.SuccessStyle().Render("✅ Undo successfully"))
}
