package main

import (
	"chatting_in_golang/trace"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	// forward 는 수신한 메시지를 보관하는 채널이며,
	// 해당 채팅방에 있는 클라이언트들에게 전달되어야 함
	forward chan []byte
	// join 은 방에 들어오려는 클라이언트를 위한 채널
	join chan *client
	// leave 는 방에서 나가려는 클라이언트를 위한 채널
	leave chan *client
	// clients 는 현재 채팅방에 있는 모든 클라이언트를 보유함
	clients map[*client]bool
	// tracer 는 방 안에서 활동 정보를 추적함
	tracer trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// 입장
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.leave:
			// 퇴장
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", string(msg))

			// 모든 클라이언트들에게 메시지 전달
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace(" -- sent to client")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client

	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
