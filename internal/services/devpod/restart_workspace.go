// restart_workspace.go

package devpod

import (
	"clusterix-code/internal/data/dto"
	"context"
	"fmt"
)

func (s *DevpodService) RestartWorkspace(
	ctx context.Context,
	devpodWorkspaceDTO dto.DevpodWorkspace,
	onSuccess func(string, string) error,
	onFailure func(string, string) error,
) error {
	if err := s.StopWorkspace(ctx, devpodWorkspaceDTO, onSuccess, onFailure); err != nil {
		return fmt.Errorf("failed to stop workspace for restart: %w", err)
	}
	return s.StartWorkspace(ctx, devpodWorkspaceDTO, onSuccess, onFailure)
}
