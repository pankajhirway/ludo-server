package players

import (
	"log"

	"github.com/gorilla/websocket"
)

type Player struct {
	c    *websocket.Conn
	Id   string `json:"userId"`
	Name string `json:"name"`
}

func (p *Player) Send(message string) {
	if p.c != nil {
		p.c.WriteMessage(websocket.TextMessage, []byte(message))
	}
}

func (p *Player) StartListening() {
	defer p.c.Close()
	for {
		_, message, err := p.c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
	}
}

func (p *Player) Set(conn *websocket.Conn, id string) {
	p.c = conn
	p.Id = id
}
