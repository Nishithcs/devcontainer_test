package socket

import (
	"clusterix-code/internal/api/requests"
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/services"
	"clusterix-code/internal/utils/di"
	"clusterix-code/internal/utils/logger"
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type WorkspaceLogHandlerConsumer struct {
	services *services.Services
}

func NewWorkspaceLogHandlerConsumer(c *di.Container) *WorkspaceLogHandlerConsumer {
	services := di.Make[*services.Services](c)
	return &WorkspaceLogHandlerConsumer{
		services: services,
	}
}

func (c *WorkspaceLogHandlerConsumer) HandleMessage(delivery amqp.Delivery) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stackTrace := string(debug.Stack())
			logger.Error("Panic occurred while processing message",
				fmt.Errorf("%v", r),
				zap.String("stack_trace", stackTrace),
				zap.String("message_id", delivery.MessageId))
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	var message dto.Message
	if err := json.Unmarshal(delivery.Body, &message); err != nil {
		logger.Error("Failed to unmarshal socket message", err,
			zap.String("body", string(delivery.Body)),
			zap.String("message_id", delivery.MessageId))
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	var logData dto.LogData
	rawData, _ := json.Marshal(message.Data)
	if err = json.Unmarshal(rawData, &logData); err == nil {
		logRequest := requests.CreateWorkspaceLogRequest{
			WorkspaceID: logData.WorkspaceId,
			Text:        logData.Text,
			Type:        logData.Type,
			Time:        logData.Time,
		}
		if logData.Text != "" {
			if err := c.services.WorkspaceLog.Create(context.Background(), logRequest); err != nil {
				logger.Error("Failed to insert workspace log", err)
				return fmt.Errorf("failed to insert workspace log: %w", err)
			}
		}
	}

	logger.Info("Processing socket message", zap.String("message_id", delivery.MessageId))
	c.services.Socket.SendMessage(message)

	return nil
}
