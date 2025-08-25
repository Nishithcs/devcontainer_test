package user

import (
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

type AuthUserDeletedConsumer struct {
	services *services.Services
}

func NewAuthUserDeletedConsumer(c *di.Container) *AuthUserDeletedConsumer {
	services := di.Make[*services.Services](c)

	return &AuthUserDeletedConsumer{
		services: services,
	}
}

func (c *AuthUserDeletedConsumer) HandleMessage(delivery amqp.Delivery) error {
	var stackTrace string
	defer func() {
		if r := recover(); r != nil {
			stackTrace = string(debug.Stack())
			logger.Error("Panic occurred while processing message",
				fmt.Errorf("%v", r),
				zap.String("stack_trace", stackTrace),
				zap.String("message_id", delivery.MessageId))
		}

		if ackErr := delivery.Ack(false); ackErr != nil {
			logger.Error("Failed to acknowledge message", ackErr)
		}
	}()

	var message ProcessAuthUserSyncMessage
	if err := json.Unmarshal(delivery.Body, &message); err != nil {
		logger.Error("Failed to unmarshal auth user deleted message", err,
			zap.String("body", string(delivery.Body)),
			zap.String("message_id", delivery.MessageId))
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	logger.Info("Processing auth user deleted message", zap.String("message_id", delivery.MessageId), zap.Uint64("user_id", message.ID))
	err := c.services.User.SyncDeletedUser(context.Background(), message.ID)
	if err != nil {
		logger.Error("Failed to process auth user deleted message", err,
			zap.String("message_id", delivery.MessageId))
		return fmt.Errorf("failed to process auth user deleted message: %w", err)
	}

	return nil
}
