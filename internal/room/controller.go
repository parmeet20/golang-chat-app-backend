package room

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/parmeet20/golang-chatapp/internal/auth"
	"github.com/parmeet20/golang-chatapp/pkg/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomController struct {
	Service   *RoomService
	validator *validator.Validate
}

func NewRoomController(service *RoomService) *RoomController {
	return &RoomController{
		Service:   service,
		validator: validator.New(),
	}
}

func (c *RoomController) CreateRoom(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	var room Room
	if err := json.NewDecoder(body).Decode(&room); err != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.validator.Struct(room); err != nil {
		response.JSON(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}
	claims, err := auth.GetClaims(r)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, err.Error())
		return
	}
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	room.CreatedBy = userID
	room.CreatedAt = time.Now()
	room.Members = append(room.Members, userID)

	if err := c.Service.CreateRoom(&room); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func (c *RoomController) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid Room ID")
		return
	}

	room, err := c.Service.GetRoomByIDWithMembers(objID)
	if err != nil {
		response.JSON(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, room)
}
func (c *RoomController) FindAllRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := c.Service.FindAllRooms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.JSON(w, http.StatusOK, rooms)
}

func (c *RoomController) JoinRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid Room ID")
		return
	}
	claims, err := auth.GetClaims(r)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, err.Error())
		return
	}
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	if err := c.Service.JoinRoom(objID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *RoomController) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid Room ID")
		return
	}
	claims, err := auth.GetClaims(r)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, err.Error())
		return
	}
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	if err := c.Service.LeaveRoom(objID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *RoomController) GetRoomsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	rooms, err := c.Service.GetRoomsByUserID(objID)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, rooms)
}
