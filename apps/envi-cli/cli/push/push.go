package push

import (
	"envi/apps/envi-cli/cli/pull"
	"envi/apps/envi-cli/internal/llog"
	"envi/apps/envi-cli/internal/provider"
	"envi/apps/envi-cli/internal/ui"
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
	pullFn, errPullFn := pull.GetPullFn(cmd)

	if errPullFn != nil {
		return errPullFn
	}
	err := F.Pipe3(
		pull.SyncEnvState(pullFn),
		E.Chain(validateIfEnvFileHasChanges),
		E.Chain(func(s pull.EnvSyncState) E.Either[error, string] {
			done := ui.ProgressBar("Pushing .env file to remote server")
			eitherPushLocalEnvValues := provider.ZipperPushLocalEnvsToRemote(s.LocalEnvValues)
			done()
			return eitherPushLocalEnvValues
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
