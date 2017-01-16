package core

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gocraft/web"
)

type LoopChat struct {
	Sessions map[string]*Session
	begin    chan *Session
	end      chan *Session
}

func New() *LoopChat {
	return &LoopChat{
		Sessions: make(map[string]*Session),
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
			fmt.Printf("Beginning %s %s...\n",
				color.RedString("Session"),
				color.RedString(session.ID))
		case session := <-l.end:
			delete(l.Sessions, session.ID)
			close(session.incoming)
			fmt.Printf("Ending %s %s...\n",
				color.RedString("Session"),
				color.RedString(session.ID))
		}
	}
}

func (l *LoopChat) CreateSession(rw web.ResponseWriter, req *web.Request) {
	sessionID := getHashID()
	session := NewSession(sessionID, l.end)

	// register session
	l.begin <- session

	clientID := getHashID()
	ServeClient(clientID, session, rw, req.Request)
}
