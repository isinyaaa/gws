package main

import (
	"fmt"
	"log"
	"time"

	"github.com/lxzan/gws"
)

const remoteAddr = "127.0.0.1:9001"

func main() {
	const count = 517
	for i := 1; i <= count; i++ {
		testCase(true, i, "gws-client/sync")
	}
	for i := 1; i <= count; i++ {
		testCase(false, i, "gws-client/async")
	}
	updateReports()
}

func testCase(sync bool, id int, agent string) {
	url := fmt.Sprintf("ws://%s/runCase?case=%d&agent=%s", remoteAddr, id, agent)
	handler := &WebSocket{Sync: sync, onexit: make(chan struct{})}
	socket, _, err := gws.NewClient(handler, &gws.ClientOption{
		Addr:             url,
		CheckUtf8Enabled: true,
	})
	if err != nil {
		log.Println(err.Error())
		return
	}
	go socket.ReadLoop()
	<-handler.onexit
}

type WebSocket struct {
	onexit chan struct{}
	Sync   bool
}

func (c *WebSocket) OnOpen(socket *gws.Conn) {
	_ = socket.SetDeadline(time.Now().Add(30 * time.Second))
}

func (c *WebSocket) OnClose(socket *gws.Conn, err error) {
	c.onexit <- struct{}{}
}

func (c *WebSocket) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.WritePong(payload)
}

func (c *WebSocket) OnPong(socket *gws.Conn, payload []byte) {}

func (c *WebSocket) OnMessage(socket *gws.Conn, message *gws.Message) {
	if c.Sync {
		_ = socket.WriteMessage(message.Opcode, message.Bytes())
		_ = message.Close()
	} else {
		socket.WriteAsync(message.Opcode, message.Bytes(), func(err error) { _ = message.Close() })
	}
}

type updateReportsHandler struct {
	gws.BuiltinEventHandler
	onexit chan struct{}
}

func (c *updateReportsHandler) OnOpen(socket *gws.Conn) {
	_ = socket.SetDeadline(time.Now().Add(5 * time.Second))
}

func (c *updateReportsHandler) OnClose(socket *gws.Conn, err error) {
	c.onexit <- struct{}{}
}

func updateReports() {
	url := fmt.Sprintf("ws://%s/updateReports?agent=gws/client", remoteAddr)
	handler := &updateReportsHandler{onexit: make(chan struct{})}
	socket, _, err := gws.NewClient(handler, &gws.ClientOption{
		Addr:             url,
		CheckUtf8Enabled: true,
	})
	if err != nil {
		log.Println(err.Error())
		return
	}
	go socket.ReadLoop()
	<-handler.onexit
}
