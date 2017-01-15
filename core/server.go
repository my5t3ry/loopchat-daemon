package core

import "github.com/satori/go.uuid"

type LoopChat struct {
	Clients  map[uuid.UUID]*Client
	Sessions map[uuid.UUID]*Session
}
