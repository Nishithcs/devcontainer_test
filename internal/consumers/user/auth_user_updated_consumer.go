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

type AuthUserUpdatedConsumer struct {
	services *services.Services
}

func NewAuthUserUpdatedConsumer(c *di.Container) *AuthUserUpdatedConsumer {
	services := di.Make[*services.Services](c)

	return &AuthUserUpdatedConsumer{
		services: services,
	}
}

func (c *AuthUserUpdatedConsumer) HandleMessage(delivery amqp.Delivery) error {
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
		logger.Error("Failed to unmarshal auth user updated message", err,
			zap.String("body", string(delivery.Body)),
			zap.String("message_id", delivery.MessageId))
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	logger.Info("Processing auth user updated message", zap.String("message_id", delivery.MessageId), zap.Uint64("user_id", message.ID))
	err := c.services.User.SyncUpdatedUser(context.Background(), message.ID)
	if err != nil {
		logger.Error("Failed to process auth user updated message", err,
			zap.String("message_id", delivery.MessageId))
		return fmt.Errorf("failed to process auth user updated message: %w", err)
	}

	return nil
}
