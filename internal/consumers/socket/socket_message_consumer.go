package socket

import (
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/services"
	"clusterix-code/internal/utils/di"
	"clusterix-code/internal/utils/logger"
	"encoding/json"
	"fmt"
	"runtime/debug"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type SocketMessageConsumer struct {
	services *services.Services
}

func NewSocketMessageConsumer(c *di.Container) *SocketMessageConsumer {
	services := di.Make[*services.Services](c)
	return &SocketMessageConsumer{
		services: services,
	}
}

func (c *SocketMessageConsumer) HandleMessage(delivery amqp.Delivery) (err error) {
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

	logger.Info("Processing socket message", zap.String("message_id", delivery.MessageId))
	c.services.Socket.SendMessage(message)

	return nil
}
