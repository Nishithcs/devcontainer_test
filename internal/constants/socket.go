package constants

import "time"

const (
	// Default limit for pagination.
	DefaultLimit = 10

	// Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * 9) / 10

	// Maximum message size allowed from peer.
	MaxMessageSize = 512
)

type Action string

type EventType string

const (
	ActionSubscription   Action = "subscribe"
	ActionUnsubscription Action = "unsubscribe"
	ActionChat           Action = "chat"
)

const (
	WorkspaceCreated EventType = "workspace_created"
	WorkspaceLogs    EventType = "workspace_log"
	WorkspaceStatus  EventType = "workspace_status"
)
