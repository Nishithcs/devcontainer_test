package devpod

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

type logLevel int

const (
	InfoLog logLevel = iota
	DoneLog
	FatalLog
)

type DevpodService struct{}

func NewDevpodService() *DevpodService {
	_ = ensureAWSProvider()
	if err := setDevpodIdleTimeout(); err != nil {
		fmt.Printf("[DevpodService] WARNING: failed to set idle timeout: %v\n", err)
	}
	return &DevpodService{}
}

func forwardPortWithSocat(workspaceID, userID uint64, internalPort string) {
	externalPort, err := getAvailablePort(20000)
	if err != nil {
		processLog(workspaceID, userID, "LOG", false, "Could not find available external port: %v", err)
		return
	}

	processLog(workspaceID, userID, "LOG", false, "Detected devpod internal port: %s; assigning external port: %d", internalPort, externalPort)

	if err := waitForPort(internalPort, 30*time.Second); err != nil {
		processLog(workspaceID, userID, "LOG", false, "Timeout waiting for devpod port %s: %v", internalPort, err)
		return
	}

	cmd := exec.Command("socat",
		fmt.Sprintf("TCP-LISTEN:%d,fork,reuseaddr", externalPort),
		fmt.Sprintf("TCP:127.0.0.1:%s", internalPort))

	stderr, err := cmd.StderrPipe()
	if err != nil {
		processLog(workspaceID, userID, "LOG", false, "Failed to get socat stderr for port %d: %v", externalPort, err)
	}

	if err := cmd.Start(); err != nil {
		processLog(workspaceID, userID, "LOG", false, "Failed to start socat forwarding %d->%s: %v", externalPort, internalPort, err)
		return
	}

	processLog(workspaceID, userID, "LOG", false, "Socat started forwarding 0.0.0.0:%d -> 127.0.0.1:%s", externalPort, internalPort)

	if stderr != nil {
		go func(port int) {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				processLog(workspaceID, userID, "LOG", false, "[socat-%d] %s", port, scanner.Text())
			}
		}(externalPort)
	}

	processLog(workspaceID, userID, "LOG", false, "Workspace available at http://localhost:%d/?folder=...", externalPort)
}

func handleCallback(callback func(string, string) error, workspaceID, userID uint64, line string, logType string) {
	if callback != nil {
		if err := callback(line, logType); err != nil {
			processLog(workspaceID, userID, "LOG", true, "onRunning callback error: %v", err)
		}
	}
}

func processLog(workspaceID uint64, userID uint64, logType string, isUserLog bool, format string, args ...interface{}) {
	//message := fmt.Sprintf(format, args...)
	//log.Printf("[devpod-%d] [%s] %s", workspaceID, logType, message)
}

func processLogWithLevel(level logLevel, workspaceID, userID uint64, message string) {
	var styled, logType string
	switch level {
	case DoneLog:
		logType = "DONE"
		styled = fmt.Sprintf("✅ \033[1;32m[%s] %s\033[0m", logType, message)
	case FatalLog:
		logType = "FATAL"
		styled = fmt.Sprintf("🚨 \033[1;31m[%s] %s\033[0m", logType, message)
	case InfoLog:
		logType = "INFO"
		styled = fmt.Sprintf("🔸 [%s] %s", logType, message)
	default:
		logType = "LOG"
		styled = message
	}
	processLog(workspaceID, userID, logType, true, styled)
}

func waitForPort(port string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.Dial("tcp", "127.0.0.1:"+port)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for port %s", port)
}

func getAvailablePort(start int) (int, error) {
	for port := start; port < start+1000; port++ {
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		l, err := net.Listen("tcp", addr)
		if err == nil {
			_ = l.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available ports")
}

func ensureAWSProvider() error {
	cmd := exec.Command("devpod", "provider", "add", "aws",
		"--option", "AWS_REGION="+os.Getenv("AWS_REGION"),
		"--option", "AWS_ACCESS_KEY_ID="+os.Getenv("AWS_ACCESS_KEY_ID"),
		"--option", "AWS_SECRET_ACCESS_KEY="+os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"--option", "subnetId="+os.Getenv("AWS_SUBNET_ID"),
		"--option", "securityGroupIds="+os.Getenv("AWS_SECURITY_GROUP_IDS"),
		"--option", "INACTIVITY_TIMEOUT=10m",
		"--option", "AWS_INSTANCE_TYPE=t2.nano",
	)
	// cmd := exec.Command("devpod", "provider", "add", "docker")
	cmd.Env = os.Environ()
	_, _ = cmd.CombinedOutput()
	return nil
}

func setDevpodIdleTimeout() error {
	cmd := exec.Command("devpod", "context", "set-options", "-o", "EXIT_AFTER_TIMEOUT=false")
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set devpod idle timeout: %w\noutput: %s", err, string(output))
	}
	processLog(0, 0, "LOG", false, "Configured idle timeout: %s", string(output))
	return nil
}

func isDoneLine(line string) bool {
	return strings.Contains(line, "done")
}

func isInfoLine(line string) bool {
	return strings.Contains(line, "info")
}

func isFatalLine(line string) bool {
	return strings.Contains(line, "fatal")
}
