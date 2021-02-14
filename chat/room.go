package chat

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		fmt.Printf("\n%s", string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			fmt.Printf("E: Failed %v", err)
			return
		}
	}
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("I: %v", r.Host)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("E: %v", err)
	}

	reader(ws)
}
