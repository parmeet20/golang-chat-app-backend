package websocket

import "github.com/parmeet20/golang-chatapp/internal/message"

type Room struct {
	id               string
	clients          map[*Client]bool
	registerClient   chan *Client
	unRegisterClient chan *Client
	broadcast        chan []byte
	messageService   *message.MessageService
	quit             chan struct{}
}

func NewRoom(id string, msgService *message.MessageService) *Room {
	return &Room{
		id:               id,
		clients:          make(map[*Client]bool),
		registerClient:   make(chan *Client),
		unRegisterClient: make(chan *Client),
		broadcast:        make(chan []byte, 256),
		messageService:   msgService,
		quit:             make(chan struct{}),
	}
}

func (r *Room) Run() {

	for {
		select {

		case client := <-r.registerClient:
			r.clients[client] = true

		case client := <-r.unRegisterClient:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}

		case message := <-r.broadcast:

			for client := range r.clients {

				select {
				case client.send <- message:

				default:
					close(client.send)
					delete(r.clients, client)
				}
			}

		case <-r.quit:
			for client := range r.clients {
				close(client.send)
			}
			return
		}
	}
}
