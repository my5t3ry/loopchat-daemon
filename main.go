package main

import (
	"flag"

	"github.com/connorwalsh/loopchat-daemon/core"
)

var addr = flag.String("addr", ":6666", "http service address")

func main() {
	flag.Parse()

	daemon, err := core.New()
	if err != nil {
		panic(err)
	}

	daemon.Start(addr)
}
