package devpod

import (
	"clusterix-code/internal/data/enums"
	"clusterix-code/internal/data/repositories"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type StatusUpdater func(ctx context.Context, workspaceID uint64, newStatus enums.WorkspaceStatus, message string) error

func StartPersistentSupervisor(ctx context.Context, workspaceRepo *repositories.WorkspaceRepository, updateStatus StatusUpdater) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("🛑 Persistent supervisor shutting down\n")
			return
		default:
			syncWorkspacesWithDatabase(ctx, workspaceRepo, updateStatus)
			time.Sleep(10 * time.Second)
		}
	}
}

func syncWorkspacesWithDatabase(ctx context.Context, repo *repositories.WorkspaceRepository, updateStatus StatusUpdater) {
	out, err := exec.Command("devpod", "list").CombinedOutput()
	if err != nil {
		fmt.Printf("devpod list failed: %v\n", err)
		return
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "NAME") || strings.Contains(line, "---") {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 1 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		if name == "" {
			continue
		}
		// 🔥 In your DB, workspace name should match DevPod's --id
		workspaceId, err := strconv.ParseUint(name, 10, 64)
		if err != nil {
			fmt.Printf("Workspace ID '%s' is not a valid uint64\n", name)
			continue
		}

		workspace, err := repo.GetByID(ctx, workspaceId)
		if err != nil {
			// If not found → insert as new workspace or just log
			fmt.Printf("Workspace '%s' found in DevPod but missing in DB\n", name)
			continue
		}

		// Update status if needed
		if workspace.Status != "running" {
			err := updateStatus(ctx, workspace.ID, "running", "Workspace detected running in DevPod")
			if err != nil {
				fmt.Printf("Failed to update workspace '%s' status: %v\n", name, err)
			} else {
				fmt.Printf("✅ Updated workspace '%s' status to running\n", name)
			}
		} else {
			fmt.Println("Something happened for the Workspace")
		}
	}
}
