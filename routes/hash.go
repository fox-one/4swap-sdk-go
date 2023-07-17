package routes

import (
	"github.com/speps/go-hashids"
)

var h *hashids.HashID

func init() {
	h = newHashID()
}

func newHashID() *hashids.HashID {
	hd := hashids.NewData()
	hd.Salt = "uniswap routes"
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}

	return h
}
