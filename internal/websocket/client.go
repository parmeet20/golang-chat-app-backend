package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/parmeet20/golang-chatapp/internal/message"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	id   string
	conn *websocket.Conn
	room *Room
	send chan []byte
}

func NewClient(id string, conn *websocket.Conn, room *Room) *Client {
	return &Client{
		id:   id,
		conn: conn,
		room: room,
		send: make(chan []byte, 256),
	}
}

func (c *Client) readPump() {

	defer func() {
		c.room.unRegisterClient <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {

		_, msgBytes, err := c.conn.ReadMessage()

		if err != nil {
			break
		}

		userId, _ := primitive.ObjectIDFromHex(c.id)
		roomId, _ := primitive.ObjectIDFromHex(c.room.id)

		msg := message.Message{
			SenderID:  userId,
			RoomID:    roomId,
			Content:   string(msgBytes),
			CreatedAt: time.Now(),
		}

		if err := c.room.messageService.CreateMessage(&msg); err != nil {
			log.Println(err)
			continue
		}

		data, _ := json.Marshal(msg)

		c.room.broadcast <- data
	}
}
func (c *Client) writePump() {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {

		case msg, ok := <-c.send:

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
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
