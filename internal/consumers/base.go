package consumers

import (
	"clusterix-code/internal/config"
	"clusterix-code/internal/services"
	"clusterix-code/internal/utils/di"
	"clusterix-code/internal/utils/logger"
	"clusterix-code/internal/utils/rabbitmq"
	"context"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Consumer interface {
	HandleMessage(delivery amqp.Delivery) error
}

// ConsumerManager manages the lifecycle of RabbitMQ
type ConsumerManager struct {
	services  *services.Services
	rabbitMQ  *rabbitmq.RabbitMQ
	cfg       *config.Config
	consumers map[string]Consumer
}

func Provider(c *di.Container) (*ConsumerManager, error) {
	cfg := di.Make[*config.Config](c)
	services := di.Make[*services.Services](c)
	rabbitMQ := di.Make[*rabbitmq.RabbitMQ](c)

	return NewConsumerManager(services, rabbitMQ, cfg), nil
}

// NewConsumerManager creates a new ConsumerManager instance.
func NewConsumerManager(services *services.Services, rabbitMQ *rabbitmq.RabbitMQ, cfg *config.Config) *ConsumerManager {
	return &ConsumerManager{
		services:  services,
		rabbitMQ:  rabbitMQ,
		cfg:       cfg,
		consumers: make(map[string]Consumer),
	}
}

// RegisterConsumer registers a consumer for a specific queue.
func (m *ConsumerManager) RegisterConsumer(queueName string, consumer Consumer) {
	m.consumers[queueName] = consumer
}

// Start starts all registered consumers
func (m *ConsumerManager) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, queueConfig := range m.cfg.RabbitMQ.Queues {
		consumer, ok := m.consumers[queueConfig.Name]
		if !ok {
			logger.Warn("No consumer found for queue", zap.String("queue", queueConfig.Name))
			continue
		}
		// Protect each goroutine with panic recovery
		go func(q config.QueueConfig, c Consumer) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("Panic in consumer goroutine", nil, zap.Any("error", r), zap.String("queue", q.Name))
				}
			}()
			m.startConsumer(ctx, q, c)
		}(queueConfig, consumer)
	}

	// Wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	logger.Info("Shutdown signal received")
}

// startConsumer starts a single consumer for a specific queue.
func (m *ConsumerManager) startConsumer(ctx context.Context, queueConfig config.QueueConfig, consumer Consumer) {
	ch, err := m.rabbitMQ.Conn.Channel()
	if err != nil {
		logger.Fatal("Failed to open a channel", zap.Error(err))
	}
	defer ch.Close()

	// Declare the queue if needed
	if queueConfig.ShouldCreate {
		_, err := ch.QueueDeclare(
			queueConfig.Name,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			logger.Fatal("Failed to declare queue", zap.Error(err), zap.String("queue", queueConfig.Name))
		}
	}

	// Binding the queue if needed
	if len(queueConfig.Bindings) > 0 {
		for _, binding := range queueConfig.Bindings {
			if err := ch.QueueBind(
				queueConfig.Name,
				binding.RoutingKey,
				binding.ExchangeName,
				false,
				nil,
			); err != nil {
				logger.Fatal("Failed to binding queue", zap.Error(err), zap.String("queue", queueConfig.Name))
			}
		}
	}

	// Set QoS/prefetch
	if err := ch.Qos(1, 0, false); err != nil {
		logger.Fatal("Failed to Set QoS/prefetch", zap.Error(err), zap.String("queue", queueConfig.Name))
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		queueConfig.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Fatal("Failed to register a consumer", zap.Error(err), zap.String("queue", queueConfig.Name))
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("Context canceled, stopping consumer", zap.String("queue", queueConfig.Name))
			return
		case msg, ok := <-msgs:
			if !ok {
				logger.Warn("Message channel closed", zap.String("queue", queueConfig.Name))
				return
			}
			err := consumer.HandleMessage(msg)
			if err != nil {
				logger.Error("Failed to process message", err, zap.String("queue", queueConfig.Name))
				if err := msg.Nack(false, true); err != nil {
					logger.Error("Failed to Nack message", err, zap.String("queue", queueConfig.Name))
				}
			} else {
				if err := msg.Ack(false); err != nil {
					logger.Error("Failed to Ack message", err, zap.String("queue", queueConfig.Name))
				}
			}
		}
	}
}
