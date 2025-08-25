package services

import (
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/utils/logger"
	ws "clusterix-code/internal/websocket"
	"go.uber.org/zap"
)

type SocketService struct {
	Config *SocketServiceConfig
}

type SocketServiceConfig struct {
	Hub *ws.Hub
}

func NewSocketService(config *SocketServiceConfig) *SocketService {
	return &SocketService{
		Config: config,
	}
}

func (s *SocketService) SendMessage(message dto.Message) {
	if _, ok := s.Config.Hub.Channels[message.Channel]; !ok {
		s.Config.Hub.Channels[message.Channel] = make(map[*ws.Client]bool)
	}

	s.Config.Hub.Broadcast <- message
	logger.Info("Message sent to hub", zap.String("channel", message.Channel))
}

func (s *SocketService) HandleMessage(message dto.InputMessage) error {
	//msg := dto.Message{
	//	EventType: message.EventType,
	//	Channel:   message.Channel,
	//	Data:      message.Data,
	//}
	//msgBytes, err := json.Marshal(msg)
	//if err != nil {
	//	return fmt.Errorf("failed to marshal message: %w", err)
	//}

	return nil
}
