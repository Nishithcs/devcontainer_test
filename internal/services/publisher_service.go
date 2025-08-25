package services

import (
	"clusterix-code/internal/utils/logger"
	"clusterix-code/internal/utils/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PublisherServiceConfig struct {
	RabbitMQ *rabbitmq.RabbitMQ
}

type PublisherService struct {
	rabbitMQ *rabbitmq.RabbitMQ
}

func NewPublisherService(config *PublisherServiceConfig) *PublisherService {
	return &PublisherService{
		rabbitMQ: config.RabbitMQ,
	}
}

func (s *PublisherService) Publish(exchange, routingKey string, body []byte) error {
	ch, err := s.rabbitMQ.Conn.Channel()
	if err != nil {
		logger.Error("Failed to open a channel", err)
		return err
	}
	defer ch.Close()

	err = ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		logger.Error("Failed to publish message", err)
		return err
	}
	return nil
}
