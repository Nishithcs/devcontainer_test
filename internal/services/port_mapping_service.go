package services

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"sync"
)

const (
	internalStartPort = 10800
	externalStartPort = 20000
	maxMappings       = 999
)

type Mapping struct {
	WorkspaceID  string
	InternalPort int
	ExternalPort int
	Cmd          *exec.Cmd
}

var (
	mu        sync.Mutex
	mappings  = make(map[string]*Mapping) // workspaceID -> Mapping
	usedPorts = make(map[int]bool)
)

func isPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

func findFreeExternalPort() (int, error) {
	for i := 0; i < maxMappings; i++ {
		port := externalStartPort + i
		if !usedPorts[port] && isPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available external ports")
}

func MapWorkspace(workspaceID string, internalPort int) (*Mapping, error) {
	mu.Lock()
	defer mu.Unlock()

	// Check if already mapped
	if m, ok := mappings[workspaceID]; ok {
		return m, nil
	}

	externalPort, err := findFreeExternalPort()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("socat",
		fmt.Sprintf("TCP-LISTEN:%d,fork,reuseaddr", externalPort),
		fmt.Sprintf("TCP:127.0.0.1:%d", internalPort),
	)

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start socat: %w", err)
	}

	mapping := &Mapping{
		WorkspaceID:  workspaceID,
		InternalPort: internalPort,
		ExternalPort: externalPort,
		Cmd:          cmd,
	}

	mappings[workspaceID] = mapping
	usedPorts[externalPort] = true

	log.Printf("Mapped workspace %s: 127.0.0.1:%d → 0.0.0.0:%d\n", workspaceID, internalPort, externalPort)

	return mapping, nil
}
