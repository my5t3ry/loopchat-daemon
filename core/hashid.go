package core

import (
	"github.com/satori/go.uuid"
	"github.com/speps/go-hashids"
)

func getHashID() string {
	hd := hashids.NewData()
	hd.Salt = uuid.NewV4().String()
	hd.MinLength = 5
	h := hashids.NewWithData(hd)
	id, _ := h.Encode([]int{666})

	return id
}
