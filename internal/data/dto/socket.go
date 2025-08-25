package dto

import (
	"clusterix-code/internal/constants"
)

type MessageHandler interface {
	HandleMessage(message InputMessage) error
}

type InputMessage struct {
	Action    string              `json:"action"`
	EventType constants.EventType `json:"event_type"`
	Channel   string              `json:"channel"`
	Data      interface{}         `json:"data"`
}

type Message struct {
	EventType constants.EventType `json:"event_type"`
	Channel   string              `json:"channel"`
	Data      interface{}         `json:"data"`
}

type LogData struct {
	WorkspaceId uint64 `json:"workspace_id"`
	Text        string `json:"text"`
	Type        string `json:"type"`
	Time        string `json:"time"`
}
