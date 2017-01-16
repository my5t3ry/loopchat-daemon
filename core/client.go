package core

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// time limit for writing to peer
	writeWait = 10 * time.Second

	// time limit for reading pong msg from peer
	pongWait = 60 * time.Second

	// periodicity for sending pings to peer. Must be < pongwiat
	pingPeriod = (pongWait * 9) / 10

	// msg limit in bytes
	maxMessageSize = 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Client struct {
	ID      string
	conn    *websocket.Conn
	session *Session
	// session --> [o|u|t] -> client
	outgoing chan []byte
}

func ServeClient(id string, session *Session, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		ID:       id,
		conn:     conn,
		session:  session,
		outgoing: make(chan []byte, 256),
	}

	// register client to session
	client.session.register <- client

	// start writing
	go client.write()

	// start reading
	go client.read()
}

func (c *Client) read() {
	// unregister & close connection upon return
	defer func() {
		c.session.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// begin reading from peer
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		// pass incoming message to server
		msg = bytes.TrimSpace(bytes.Replace(msg, []byte("\n"), []byte(" "), -1))
		c.session.incoming <- msg
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	// begin writing to peer
	for {
		select {
		case msg, ok := <-c.outgoing:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// server is closing the connection
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			// add queued messages
			n := len(c.outgoing)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.outgoing)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
