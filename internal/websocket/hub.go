package websocket

import (
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/utils/logger"
	"encoding/json"
	"fmt"
	"log"

	"go.uber.org/zap"
)

type Hub struct {
	Clients       map[*Client]bool
	Channels      map[string]map[*Client]bool
	Broadcast     chan dto.Message
	Register      chan *Client
	Unregister    chan *Client
	Shutdown      chan struct{}
	MessageRouter dto.MessageHandler
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Channels:   make(map[string]map[*Client]bool),
		Broadcast:  make(chan dto.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Shutdown:   make(chan struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			logger.Info("Client registered")
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				for channel, clients := range h.Channels {
					if _, exists := clients[client]; exists {
						delete(clients, client)
						logger.Info("Client removed from channel", zap.String("channel", channel))
						if len(clients) == 0 {
							delete(h.Channels, channel)
							logger.Info("Channel deleted", zap.String("channel", channel))
						}
					}
				}
				close(client.Send)
			}
		case message := <-h.Broadcast:
			if clients, ok := h.Channels[message.Channel]; ok {
				for client := range clients {
					select {
					case client.Send <- serialize(message):
					default:
						close(client.Send)
						delete(h.Clients, client)
					}
				}
			}
		case <-h.Shutdown:
			for client := range h.Clients {
				close(client.Send)
				delete(h.Clients, client)
			}
			close(h.Broadcast)
			close(h.Register)
		}
	}
}

func (h *Hub) Stop() {
	close(h.Shutdown)
}

func (h *Hub) SubscribeClientToChannel(c *Client, channel string) {
	if h.Channels[channel] == nil {
		h.Channels[channel] = make(map[*Client]bool)
		logger.Info("Channel created", zap.String("channel", channel))
	}

	h.Channels[channel][c] = true
	c.Channels[channel] = true
	logger.Info("Client subscribed to channel", zap.String("channel", channel))
}

func (h *Hub) UnsubscribeClientFromChannel(c *Client, channel string) {
	if clients, ok := h.Channels[channel]; ok {
		if _, exists := clients[c]; exists {
			delete(clients, c)
			delete(c.Channels, channel)
			logger.Info("Client unsubscribed from channel", zap.String("channel", channel))
			// Optionally, remove the channel if no clients remain.
			if len(clients) == 0 {
				delete(h.Channels, channel)
				logger.Info("Channel deleted", zap.String("channel", channel))
			}
		}
	}
}

func (h *Hub) HandleMessage(c *Client, message dto.InputMessage) error {
	if h.MessageRouter != nil {
		return h.MessageRouter.HandleMessage(message)
	}
	return fmt.Errorf("no message handler configured")
}

func serialize(message dto.Message) []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error serializing message: %v", err)
	}
	return json
}
