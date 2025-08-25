package jobs

import (
	"clusterix-code/internal/constants"
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/services"
	"encoding/json"
	"fmt"
)

// PublishLogMessage publishes a log message to the message queue
func PublishLogMessage(
	publisherSvc *services.PublisherService,
	workspaceID uint64,
	message string,
) error {
	messageTime, messageType, messageText := getClearLog(message)

	body := dto.Message{
		EventType: constants.WorkspaceLogs,
		Channel:   fmt.Sprintf("workspace_%d_logs", workspaceID),
		Data: map[string]interface{}{
			"workspace_id": workspaceID,
			"time":         messageTime,
			"type":         messageType,
			"text":         messageText,
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	publisherSvc.Publish(constants.CLUSTERIX_CODE_V1_EXCHANGE, constants.WORKSPACE_LOG_HANDLER_QUEUE, payload)
	return nil
}
