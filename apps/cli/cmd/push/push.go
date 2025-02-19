package push

import (
	"envii/apps/cli/cmd/pull"
	"envii/apps/cli/internal/llog"
	"envii/apps/cli/internal/provider"
	"envii/apps/cli/internal/ui"
	"fmt"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/spf13/cobra"
)

var PushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushes the latest .env file to the remote server.",
	RunE:  PushCmdFunc,
}

func PushCmdFunc(cmd *cobra.Command, args []string) error {
	pullProvider, errPullFn := provider.GetPullProvider(cmd)
	pushProvider, errPushFn := provider.GetPushProvider(cmd)

	if errPullFn != nil {
		return errPullFn
	}
	if errPushFn != nil {
		return errPushFn
	}
	err := F.Pipe3(
		pull.SyncEnvState(pullProvider.PullRemoteEnvValues, pull.SyncEnvStateOptions{
			Preserve: false,
			Provider: pullProvider.GetName(),
		}),
		E.Chain(validateIfEnvFileHasChanges),
		E.Chain(func(s pull.EnvSyncState) E.Either[error, string] {
			done := ui.ProgressBar("Pushing .env file to remote server", pullProvider.GetName())
			eitherPushLocalEnvValues := pushProvider.PushLocalEnvValues(s.LocalEnvValues)
			if eitherPushLocalEnvValues != nil {
				return E.Left[string](eitherPushLocalEnvValues)
			}
			done()
			return E.Right[error]("")
		}),
		E.Fold(
			F.Identity,
			showMessageSuccess,
		),
	)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}

func validateIfEnvFileHasChanges(s pull.EnvSyncState) E.Either[error, pull.EnvSyncState] {
	if s.DiffRemoteLocal.HasNoDiff() {
		return E.Left[pull.EnvSyncState, error](fmt.Errorf(llog.SuccessStyle().Render("✅ No changes detected in .env file.")))
	}
	return E.Right[error](s)
}

func showMessageSuccess(_ string) error {
	fmt.Println(llog.SuccessStyle().Render("✅ .env file pushed successfully"))
	return nil
}
