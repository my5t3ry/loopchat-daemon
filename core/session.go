package core

import (
	"fmt"

	"github.com/fatih/color"
)

type Session struct {
	ID         string             `json:"id"`
	FmtID      string             `json:"-"`
	Clients    map[string]*Client `json:"peers"`
	register   chan *Client
	unregister chan *Client
	// client --> [i|n] --> session
	incoming chan []byte
	end      chan *Session
}

func NewSession(id string, end chan *Session) *Session {
	return &Session{
		ID:         id,
		FmtID:      color.RedString("Session " + id),
		Clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		incoming:   make(chan []byte),
		end:        end,
	}
}

func (s *Session) Start() {
	for {
		select {
		case client := <-s.register:
			fmt.Printf("%s registering %s\n",
				s.FmtID,
				client.FmtName)
			s.Clients[client.ID] = client

			// handle registration
			go s.HandleRegistration(client)

		case client := <-s.unregister:
			fmt.Printf("%s unregistering %s\n",
				s.FmtID,
				client.FmtName)
			if _, ok := s.Clients[client.ID]; ok {
				delete(s.Clients, client.ID)
				close(client.outgoing)
			}
			if len(s.Clients) == 0 {
				// end session if there are no more clients
				s.end <- s
			}
		case msg, ok := <-s.incoming:
			if !ok {
				// session ended
				return
			}

			// handle incoming messages
			go s.HandleIncoming(msg)
		}

	}
}

func (s *Session) send(msg []byte, clients ...*Client) {
	for _, client := range clients {
		select {
		case client.outgoing <- msg:
		default:
			// if we can't reach a client-- shut them down
			close(client.outgoing)
			delete(s.Clients, client.ID)
		}
	}
}

func (s *Session) GetClients() []*Client {
	var clients = []*Client{}
	for _, client := range s.Clients {
		clients = append(clients, client)
	}

	return clients
}
