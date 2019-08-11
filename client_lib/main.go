package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8888", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("start sending message")

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Println("new message")
			fmt.Println(string(message))
		}
	}()

	go conn.WriteMessage(websocket.TextMessage, []byte("ab"))
	go conn.WriteMessage(websocket.TextMessage, []byte("cd"))

	time.Sleep(5000 * time.Millisecond)

	//client.WriteMessage(websocket.TextMessage, []byte("b"))
	//client.WriteMessage(websocket.TextMessage, []byte("c"))
}
