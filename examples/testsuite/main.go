package main

import (
	"context"
	"fmt"
	"github.com/lxzan/gws"
	"net/http"
)

func main() {
	var config = gws.Config{CompressEnabled: true, CheckTextEncoding: true, LogEnabled: true, MaxContentLength: 32 * 1024 * 1024}
	var handler = new(WebSocket)

	http.HandleFunc("/connect", func(writer http.ResponseWriter, request *http.Request) {
		socket, err := gws.Accept(context.Background(), writer, request, handler, config)
		if err != nil {
			return
		}
		socket.Listen()
	})

	_ = http.ListenAndServe(":3000", nil)
}

type WebSocket struct{}

func (c *WebSocket) OnClose(socket *gws.Conn, code uint16, reason []byte) {
	fmt.Printf("onclose: code=%d, payload=%s\n", code, string(reason))
}

func (c *WebSocket) OnError(socket *gws.Conn, err error) {
	fmt.Printf("onerror: err=%s\n", err.Error())
}

func (c *WebSocket) OnOpen(socket *gws.Conn) {
	println("connected")
}

func (c *WebSocket) OnPing(socket *gws.Conn, payload []byte) {
	fmt.Printf("onping: payload=%s\n", string(payload))
	socket.WritePong(payload)
}

func (c *WebSocket) OnPong(socket *gws.Conn, payload []byte) {}

func (c *WebSocket) OnMessage(socket *gws.Conn, message *gws.Message) {
	socket.WriteMessage(message.Typ(), message.Bytes())
	message.Close()
}
