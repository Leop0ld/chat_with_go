package main

import (
	"github.com/gorilla/websocket"
)

// client 는 사용하 1명을 나타냄
type client struct {
	// 해당 클라이언트의 웹소켓
	socket *websocket.Conn
	// 메시지가 전송되는 채널
	send chan []byte
	// 채팅방
	room *room
}

func (c *client) read() {
	// defer 키워드는 함수가 종료될 떄 해당 라인이 실행되게 함
	// 함수가 어디서 끝날지 확실하지 않을 때 유요함
	defer c.socket.Close()

	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
