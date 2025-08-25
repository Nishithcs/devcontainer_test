// rebuild_workspace.go

package devpod

import (
	"clusterix-code/internal/data/dto"
	"context"
	"fmt"
)

func (s *DevpodService) RebuildWorkspace(
	ctx context.Context,
	devpodWorkspaceDTO dto.DevpodWorkspace,
	onSuccess func(string, string) error,
	onFailure func(string, string) error,
) error {
	if err := s.TerminateWorkspace(ctx, devpodWorkspaceDTO, onSuccess, onFailure); err != nil {
		return fmt.Errorf("failed to terminate workspace for rebuild: %w", err)
	}

	return s.StartWorkspace(ctx, devpodWorkspaceDTO, onSuccess, onFailure)
}
