package core

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID         string
	conn       *websocket.Conn
	register   chan *Client
	unregister chan *Client
	// server --> [o|u|t] -> client
	outgoing chan []byte
	// client --> [i|n] --> server
	incoming chan []byte
}

func (c *Client) read() {
	// unregister & close connection upon return
	defer func() {
		c.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// begin reading
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		// pass incoming message to server
		msg = bytes.TrimSpace(bytes.Replace(msg, []byte{"\n"}, []byte{" "}, -1))
		c.incoming <- msg
	}
}
