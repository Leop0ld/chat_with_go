package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// client 는 사용하 1명을 나타냄
type client struct {
	// socket 은 해당 클라이언트의 웹소켓
	socket *websocket.Conn
	// send 는 메시지가 전송되는 채널
	send chan *message
	// room 은 클라이언트가 채팅하는 방
	room *room
	// userData 는 사용자에 대한 정보를 보유함
	userData map[string]interface{}
}

func (c *client) read() {
	// defer 키워드는 함수가 종료될 떄 해당 라인이 실행되게 함
	// 함수가 어디서 끝날지 확실하지 않을 때 유용함
	defer c.socket.Close()

	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
