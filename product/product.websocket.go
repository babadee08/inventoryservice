package product

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

type message struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

func productSocket(ws *websocket.Conn) {
	done := make(chan struct{})
	fmt.Println("new websocket connection established")
	go func(c *websocket.Conn) {
		for {
			var msg message
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				log.Println(err)
				break
			}
			fmt.Printf("received message %v\n", msg.Data)
		}
		close(done)
	}(ws)
loop:
	for {
		select {
		case <-done:
			fmt.Println("Connection was closed")
			break loop
		default:
			products, err := GetTopTenProducts()
			if err != nil {
				log.Println(err)
				break
			}
			if err := websocket.JSON.Send(ws, products); err != nil {
				log.Println(err)
				break
			}
			time.Sleep(10 * time.Second)
		}
	}
	fmt.Println("closing the websocket")
	defer ws.Close()
}
