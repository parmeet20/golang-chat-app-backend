package message

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/parmeet20/golang-chatapp/internal/room"
	"github.com/parmeet20/golang-chatapp/pkg/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageController struct {
	Service     *MessageService
	RoomService *room.RoomService
}

func NewMessageController(service *MessageService, roomService *room.RoomService) *MessageController {
	return &MessageController{Service: service, RoomService: roomService}
}

func (c *MessageController) GetMessagesByRoomId(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "id")
	if roomId == "" {
		response.JSON(w, http.StatusBadRequest, "Room ID is required")
		return
	}
	objectID, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		response.JSON(w, http.StatusNotFound, "Invalid room ID")
		return
	}
	if _, err = c.RoomService.GetRoomByID(objectID); err != nil {
		response.JSON(w, http.StatusNotFound, "Room not found")
		return
	}
	rooms, err := c.Service.GetMessagesByRoomId(objectID)
	if err != nil {
		response.JSON(w, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, rooms)
}
