package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"hikari_sync_player/pkg/logging"
	"time"
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	Send   chan []byte
	RoomId string
}

func NewClient(hub *Hub, conn *websocket.Conn, roomId string) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		Send:   make(chan []byte),
		RoomId: roomId,
	}
}

func (c *Client) ReadLoop() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	websocketJsonMessage := WebsocketJsonMessage{}
	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logging.Error(err)
			}
			break
		}

		err = json.Unmarshal(rawMessage, &websocketJsonMessage)
		if err != nil {
			logging.Error(err)
			continue
		}
		switch websocketJsonMessage.Event {
		case "updatePlayState":
			var newPlayState PlayState
			err := json.Unmarshal(websocketJsonMessage.Data, &newPlayState)
			if err != nil {
				logging.Error("failed to parse updatePlayState command data")
			} else {
				previousPlayState := c.hub.RoomPlayState[c.RoomId]
				if previousPlayState != nil {
					// Store the newPlayState only when newPlayState.Timestamp is older than previousPlayState.Timestamp
					if newPlayState.Timestamp > previousPlayState.Timestamp {
						c.hub.RoomPlayState[c.RoomId] = &newPlayState
					}
				} else {
					c.hub.RoomPlayState[c.RoomId] = &newPlayState
				}
			}
			break
		case "hello", "getPlayState":
			playState := c.hub.RoomPlayState[c.RoomId]
			playStateMessage := WebsocketJsonMessage{
				Event:  "playState",
				UserId: "",
				Data:   nil,
			}
			if playState != nil {
				// Update playTime at current timestamp
				playState = playState.GetCurrentPlayState()
				playStateBytes, err := json.Marshal(playState)
				if err != nil {
					logging.Error(err)
					break
				}
				playStateMessage.Data = playStateBytes
			}
			if websocketJsonMessage.Event == "hello" {
				playStateMessage.Event = "hello"
			}
			playStateMessageBytes, err := json.Marshal(playStateMessage)
			if err != nil {
				logging.Error(err)
				break
			}
			// Only send the play state to the client
			c.Send <- playStateMessageBytes
			break
		default:
			websocketBroadCastMessage := WebsocketBroadcastMessage{
				RoomId: c.RoomId,
				Data:   rawMessage,
			}
			c.hub.Broadcast <- websocketBroadCastMessage
			break
		}
	}
}

func (c *Client) WriteLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(message)
			if err != nil {
				logging.Error(err)
				w.Close()
				return
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
