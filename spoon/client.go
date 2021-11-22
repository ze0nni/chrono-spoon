package spoon

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clientIDCounter int = 0

func clientHandle(app *app) func(http.ResponseWriter, *http.Request) {
	var clientIDCounter = 0

	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)

		if nil != ws {
			defer ws.Close()
		}

		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
			return
		} else if err != nil {
			return
		}

		clientIDCounter++

		client := &clientConnection{
			id:  clientID(clientIDCounter),
			app: app,
			ws:  ws,

			responseCh: make(chan interface{}),
		}

		client.listen()
	}

}

type clientID int

type clientConnection struct {
	id  clientID
	app *app
	ws  *websocket.Conn

	responseCh chan interface{}
}

func (c *clientConnection) listen() {
	c.app.clientConnected(c)
	defer c.app.clientDisconnected(c)

	done := make(chan struct{})
	go c.write(done)
	c.read(done)
}

func (client *clientConnection) read(done chan struct{}) {
	log.Printf("[%d] Connected", client.id)
	for {
		client.ws.ReadMessage()
	}
	done <- struct{}{}
	log.Printf("[%d] Disconnected", client.id)
}

func (client *clientConnection) write(done chan struct{}) {
WriteLoop:
	for {
		select {
		case <-done:
			break WriteLoop
		case msg := <-client.responseCh:
			err := client.ws.WriteJSON(msg)
			if err != nil {
				log.Print("[%d] Message sending error: %s", client.id, err.Error())
				break WriteLoop
			}
			log.Printf("[%d] Message sended", client.id)
		}
	}
}
