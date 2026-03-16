package websocket

import (
	"log"

	"github.com/parmeet20/golang-chatapp/internal/message"
	"github.com/parmeet20/golang-chatapp/internal/room"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hub struct {
	rooms map[string]*Room

	registerRoom   chan registerRoomRequest
	unRegisterRoom chan *Room
	getRoom        chan roomRequest
	quit           chan struct{}
}

type roomRequest struct {
	roomID string
	resp   chan *Room
}

type registerRoomRequest struct {
	room *Room
	resp chan *Room
}

func NewHub() *Hub {
	return &Hub{
		rooms:          make(map[string]*Room),
		registerRoom:   make(chan registerRoomRequest, 32),
		unRegisterRoom: make(chan *Room, 32),
		getRoom:        make(chan roomRequest),
		quit:           make(chan struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case req := <-h.registerRoom:
			if existing, ok := h.rooms[req.room.id]; ok {
				req.resp <- existing
			} else {
				h.rooms[req.room.id] = req.room
				go req.room.Run()
				req.resp <- req.room
			}

		case room := <-h.unRegisterRoom:
			delete(h.rooms, room.id)

		case req := <-h.getRoom:
			req.resp <- h.rooms[req.roomID]
		
		case <-h.quit:
			for _, room := range h.rooms {
				close(room.quit)
			}
			return
		}
	}
}

func (h *Hub) Stop() {
	close(h.quit)
}

func (h *Hub) GetOrCreateRoom(
	roomId string,
	roomService *room.RoomService,
	messageService *message.MessageService,
) *Room {

	resp := make(chan *Room)

	h.getRoom <- roomRequest{
		roomID: roomId,
		resp:   resp,
	}

	existing := <-resp
	if existing != nil {
		return existing
	}

	objectRoomId, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		log.Println("invalid room id:", err)
		return nil
	}

	dbRoom, err := roomService.GetRoomByID(objectRoomId)
	if err != nil {
		dbRoom = &room.Room{ID: objectRoomId}

		if err := roomService.CreateRoom(dbRoom); err != nil {
			log.Println("error creating room:", err)
			return nil
		}
	}

	newRoom := NewRoom(roomId, messageService)

	req := registerRoomRequest{
		room: newRoom,
		resp: make(chan *Room),
	}
	h.registerRoom <- req

	return <-req.resp
}
