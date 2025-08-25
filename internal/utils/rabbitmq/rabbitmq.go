package rabbitmq

import (
	"clusterix-code/internal/config"
	"clusterix-code/internal/utils/di"
	"clusterix-code/internal/utils/logger"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Config  config.RabbitMQConfig
}

func Provider(c *di.Container) (*RabbitMQ, error) {
	cfg := di.Make[*config.Config](c)

	return NewRabbitMQ(cfg.RabbitMQ)
}

func NewRabbitMQ(cfg config.RabbitMQConfig) (*RabbitMQ, error) {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%d/",
		cfg.Protocol, cfg.User, cfg.Password, cfg.Host, cfg.Port)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to open a channel", err)
		return nil, err
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
		Config:  cfg,
	}, nil
}

func (r *RabbitMQ) Close() {
	if err := r.Channel.Close(); err != nil {
		logger.Warn("Failed to close RabbitMQ channel", zap.Error(err))
	}
	if err := r.Conn.Close(); err != nil {
		logger.Warn("Failed to close RabbitMQ connection", zap.Error(err))
	}
}
