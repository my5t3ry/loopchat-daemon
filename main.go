package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/connorwalsh/loopchat-daemon/core"
	"github.com/gocraft/web"
)

var (
	// CLI flag for daemon port
	addr = flag.String("addr", "6666", "http service address")

	// global context for app
	loopchat = core.New()
)

type Context struct{}

func (c *Context) CreateSession(rw web.ResponseWriter, req *web.Request) {
	// create new session and add new client
	loopchat.CreateSession(rw, req)
}

func (c *Context) JoinSession(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("ATTEMPTING TO JOIN SESSION")

	// verify whether supplied session ID corresponds to a real session

	// add new client to this existing session
	//loopchat.JoinSession(sessionID)
}

// The LoopChat daemon does not serve the html to the browser currently, but
// establishes the websockets connections with a page which has been loaded
// and requests a websockets connection with this daemon.
func main() {
	go loopchat.Run()

	flag.Parse()

	router := web.New(Context{})
	router.Get("/ws", (*Context).CreateSession)
	router.Get("/ws/:sessionID", (*Context).JoinSession)

	fmt.Println("LoopChat Daemon listening on 127.0.0.1:" + *addr + "...")

	err := http.ListenAndServe(":"+*addr, router)
	if err != nil {
		log.Fatal(err)
	}
}
