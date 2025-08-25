package devpod

import (
	"bufio"
	"clusterix-code/internal/data/dto"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func (s *DevpodService) TerminateWorkspace(
	ctx context.Context,
	devpodWorkspaceDTO dto.DevpodWorkspace,
	onSuccess func(string, string) error,
	onFailure func(string, string) error,
) error {
	cmd := exec.CommandContext(ctx, "devpod", "delete", fmt.Sprintf("%d", devpodWorkspaceDTO.DevpodWorkspaceId), "--force")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if ctx.Err() != nil {
		log.Printf("[devpod-%d] [%s] %s", devpodWorkspaceDTO.DevpodWorkspaceId, "LOG", fmt.Sprintf("context already canceled: %v", ctx.Err()))
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to terminate devpod: %w", err)
	}

	doneChan := make(chan struct{})
	fatalOccurred := new(bool)

	go func() {
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		handleTerminateWorkspaceLogs(scanner, devpodWorkspaceDTO.DevpodWorkspaceId, devpodWorkspaceDTO.UserId, onSuccess, onFailure, doneChan, fatalOccurred)
	}()

	err = cmd.Wait()
	<-doneChan

	if err != nil {
		if *fatalOccurred {
			log.Printf("[devpod-%d] [%s] %s", devpodWorkspaceDTO.DevpodWorkspaceId, "LOG", fmt.Sprintf("devpod exited with fatal, treating as success: %v", err))
			return nil
		}
		log.Printf("[devpod-%d] [%s] %s", devpodWorkspaceDTO.DevpodWorkspaceId, "LOG", fmt.Sprintf("devpod exited with error: %v", err))
		return fmt.Errorf("devpod process exited with error: %w", err)
	}

	log.Printf("[devpod-%d] [%s] %s", devpodWorkspaceDTO.DevpodWorkspaceId, "LOG", "devpod exited successfully")
	return nil
}

func handleTerminateWorkspaceLogs(
	scanner *bufio.Scanner,
	workspaceID, userID uint64,
	onSuccess, onFailure func(string, string) error,
	doneChan chan struct{},
	fatalOccurred *bool,
) {
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case isDoneLine(line):
			processLogWithLevel(DoneLog, workspaceID, userID, line)
			handleCallback(onSuccess, workspaceID, userID, line, "DONE")

		case isFatalLine(line):
			*fatalOccurred = true
			processLogWithLevel(FatalLog, workspaceID, userID, line)
			handleCallback(onFailure, workspaceID, userID, line, "FATAL")
			close(doneChan)
			return

		case isInfoLine(line):
			processLogWithLevel(InfoLog, workspaceID, userID, line)
			handleCallback(onSuccess, workspaceID, userID, line, "INFO")

		default:
			processLog(workspaceID, userID, "LOG", true, line)
			handleCallback(onSuccess, workspaceID, userID, line, "INFO")
		}
	}

	if err := scanner.Err(); err != nil {
		processLog(workspaceID, userID, "LOG", true, "scanner error: %v", err)
	}
	close(doneChan)
}
