package core

import (
	"github.com/speps/go-hashids"
)

func getHashID() string {
	hd := hashids.NewData()
	hd.Salt = "sdferf"
	hd.MinLength = 5
	h, _ := hashids.NewWithData(hd)
	id, _ := h.Encode([]int{666})

	return id
}
