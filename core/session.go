package core

import (
	"fmt"

	"github.com/fatih/color"
)

type Session struct {
	ID         string
	Clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	// client --> [i|n] --> session
	incoming chan []byte
	end      chan *Session
}

func NewSession(id string, end chan *Session) *Session {
	return &Session{
		ID:         id,
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
			fmt.Printf("%s %s registering %s %s\n",
				color.RedString("Session"),
				color.RedString(s.ID),
				color.YellowString("Client"),
				color.YellowString(client.ID))
			s.Clients[client.ID] = client
		case client := <-s.unregister:
			fmt.Printf("%s %s unregistering %s %s\n",
				color.RedString("Session"),
				color.RedString(s.ID),
				color.YellowString("Client"),
				color.YellowString(client.ID))
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
			result := s.handle(msg)

			// send result to all clients within session
			for _, client := range s.Clients {
				select {
				case client.outgoing <- result:
				default:
					// if we can't reach a client-- shut them down
					close(client.outgoing)
					delete(s.Clients, client.ID)
				}
			}
		}

	}
}

func (s *Session) handle(msg []byte) []byte {
	var result = msg

	return result
}
