package core

import "fmt"

type LoopChatDaemon struct {
	Clients  []Client
	Sessions map[string][]Connection
}

func New() (*LoopChatDaemon, error) {
	return &LoopChatDaemon{}, nil
}

func (d *LoopChatDaemon) Start() {
	fmt.Println("starting...")
}
