package main

import "github.com/connorwalsh/loopchat-daemon/core"

func main() {
	daemon, err := core.New()
	if err != nil {
		panic(err)
	}

	daemon.Start()
}
