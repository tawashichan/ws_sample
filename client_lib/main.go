package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

func main() {
	client, _, err := websocket.DefaultDialer.Dial("ws://localhost:8888", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("start sending message")
	client.WriteMessage(websocket.TextMessage, []byte("abcdefg"))
}
