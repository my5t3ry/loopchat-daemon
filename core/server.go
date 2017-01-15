package core

import (
	"fmt"

	"github.com/gocraft/web"
	"github.com/satori/go.uuid"
)

type Context struct{}

type LoopChat struct {
	Clients  map[uuid.UUID]*Client
	Sessions map[uuid.UUID]*Session
}

func New() (*LoopChat, error) {
	return &LoopChat{}, nil
}

func (d *LoopChat) Start(addr string) {
	router := web.New(d)
	router.Get("/", (*LoopChat).Hello)
	fmt.Println("starting...")

}
