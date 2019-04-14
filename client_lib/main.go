package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func main() {
	client, _, err := websocket.DefaultDialer.Dial("ws://localhost:8888", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("start sending message")

	go func() {
		for {
			_, message, err := client.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Println(message)
		}
	}()

	client.WriteMessage(websocket.TextMessage, []byte("ab"))

	time.Sleep(1000 * time.Millisecond)
	//client.WriteMessage(websocket.TextMessage, []byte("b"))
	//client.WriteMessage(websocket.TextMessage, []byte("c"))
}
