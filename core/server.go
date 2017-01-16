package core

import (
	"fmt"

	"github.com/gocraft/web"
	"github.com/satori/go.uuid"
)

type LoopChat struct {
	Clients  map[uuid.UUID]*Client
	Sessions map[uuid.UUID]*Session
	begin    chan *Session
	end      chan *Session
}

func New() *LoopChat {
	return &LoopChat{
		Clients:  make(map[uuid.UUID]*Client),
		Sessions: make(map[uuid.UUID]*Session),
		begin:    make(chan *Session),
		end:      make(chan *Session),
	}
}

func (l *LoopChat) Run() {
	for {
		select {
		case session := <-l.begin:
			l.Sessions[session.ID] = session
			go session.Start()
			fmt.Printf("Beginning Session %s...\n", session.ID.String())
		case session := <-l.end:
			delete(l.Sessions, session.ID)
			close(session.incoming)
			fmt.Printf("Ending Session %s...\n", session.ID.String())
		}
	}
}

func (l *LoopChat) CreateSession(rw web.ResponseWriter, req *web.Request) {
	session := NewSession(l.end)

	// register session
	l.begin <- session

	ServeClient(session, rw, req.Request)
}
