package devpod

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"clusterix-code/internal/data/dto"
)

// StartWorkspace is your existing CreateWorkspace, just renamed.
func (s *DevpodService) StartWorkspace(
	ctx context.Context,
	devpodWorkspaceDTO dto.DevpodWorkspace,
	onSuccess func(string, string) error,
	onFailure func(string, string) error,
) error {
	repoURL := devpodWorkspaceDTO.RepositoryUrl

	if strings.HasPrefix(repoURL, "https://") {
		repoURL = strings.Replace(repoURL, "https://", fmt.Sprintf("https://git:%s@", devpodWorkspaceDTO.AccessToken), 1)
	} else {
		repoURL = fmt.Sprintf("git:%s@%s", devpodWorkspaceDTO.AccessToken, repoURL)
	}

	cmd := exec.CommandContext(ctx,
		"./binaries/devpod-cli-linux-amd64", "up",
		repoURL,
		"--id", fmt.Sprintf("%d", devpodWorkspaceDTO.DevpodWorkspaceId),
		"--ide", "openvscode",
		"--provider-option", fmt.Sprintf("AWS_INSTANCE_TYPE=%s", devpodWorkspaceDTO.AWSInstanceType),
	)

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
		return fmt.Errorf("failed to start devpod: %w", err)
	}

	doneChan := make(chan struct{})
	fatalOccurred := new(bool)

	go func() {
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		handleStartWorkspaceLogs(scanner, devpodWorkspaceDTO.DevpodWorkspaceId, devpodWorkspaceDTO.UserId, onSuccess, onFailure, doneChan, fatalOccurred)
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

func handleStartWorkspaceLogs(
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
