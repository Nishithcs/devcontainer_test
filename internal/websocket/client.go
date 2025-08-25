package websocket

import (
	"clusterix-code/internal/constants"
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/utils/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Send     chan []byte
	Channels map[string]bool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins (adjust for production)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Handle error upgrading connection
		return
	}

	client := &Client{
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Channels: make(map[string]bool),
	}

	// Register client
	hub.Register <- client

	// Start client read and write pumps
	go client.WritePump()
	go client.ReadPump(hub)
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(constants.MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(constants.PongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(constants.PongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("Client connection closed", fmt.Errorf("unexpected close error: %w", err))
			}
			break
		}

		var inputMessage dto.InputMessage
		if err := json.Unmarshal(message, &inputMessage); err == nil {
			switch inputMessage.Action {
			case string(constants.ActionSubscription):
				hub.SubscribeClientToChannel(c, inputMessage.Channel)
			case string(constants.ActionUnsubscription):
				hub.UnsubscribeClientFromChannel(c, inputMessage.Channel)
			case string(constants.ActionChat):
				hub.HandleMessage(c, inputMessage)
			default:
				logger.Error("Unknown action", fmt.Errorf("unknown action: %s", inputMessage.Action))
			}
		} else {
			logger.Error("Received message from client", fmt.Errorf("received message from client: %w", err))
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(constants.PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(constants.WriteWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				logger.Error("Client connection closed", fmt.Errorf("hub closed channel"))
				return
			}
			// Write the message.
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logger.Error("Failed to write message to client", fmt.Errorf("failed to write message to client: %w", err))
				return
			}
		case <-ticker.C:
			// Send a ping message to keep the connection alive.
			c.Conn.SetWriteDeadline(time.Now().Add(constants.WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Error("Failed to send ping message to client", fmt.Errorf("failed to send ping message to client: %w", err))
				return
			}
		}
	}
}
