package websocket

import "github.com/gorilla/websocket"

type Client struct {
	Hub      *Hub
	Id       string
	Socket   *websocket.Conn
	Outbound chan []byte
}

func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		Hub:      hub,
		Socket:   socket,
		Outbound: make(chan []byte),
	}
}

func (c *Client) Write() {
	for {
		select {
		case message, ok := <-c.Outbound:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
