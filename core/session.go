package core

import "github.com/satori/go.uuid"

type Session struct {
	ID         uuid.UUID
	Clients    map[uuid.UUID]*Client
	register   chan *Client
	unregister chan *Client
	// client --> [i|n] --> session
	incoming chan []byte
}

func NewSession() *Session {
	return &Session{
		ID:         uuid.NewV4(),
		Clients:    make(map[uuid.UUID]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		incoming:   make(chan []byte),
	}
}

func (s *Session) start() {
	for {
		select {
		case client := <-s.register:
			s.Clients[client.ID] = client
		case client := <-s.unregister:
			if _, ok := s.Clients[client.ID]; ok {
				delete(s.Clients, client.ID)
				close(client.outgoing)
			}
		case msg := <-s.incoming:
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
	var result = []byte{}

	return result
}
