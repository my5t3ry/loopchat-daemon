package core

import (
	"encoding/json"
	"log"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (s *Session) HandleIncoming(msgIn []byte) {
	msgOut := Message{
		Type:    "chat",
		Payload: msgIn,
	}

	// marshal message
	bytes, err := json.Marshal(msgOut)
	if err != nil {
		log.Fatal(err)
	}

	s.send(bytes, s.GetClients()...)
}

func (s *Session) HandleRegistration(c *Client) {
	// create a Message
	msg := Message{
		Type:    "session",
		Payload: s,
	}

	// marshal message
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	// send to client
	s.send(bytes, c)
}
