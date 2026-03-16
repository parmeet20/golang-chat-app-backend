package room

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomService struct {
	roomRepo *RoomRepo
}

func NewRoomService(roomRepo *RoomRepo) *RoomService {
	return &RoomService{roomRepo: roomRepo}
}

func (s *RoomService) CreateRoom(room *Room) error {
	room.CreatedAt = time.Now()
	return s.roomRepo.Create(room)
}

func (s *RoomService) GetRoomByID(id primitive.ObjectID) (*Room, error) {
	return s.roomRepo.FindByID(id)
}

func (s *RoomService) GetRoomByIDWithMembers(id primitive.ObjectID) (*RoomWithMembers, error) {
	return s.roomRepo.FindByIDWithMembers(id)
}

func (s *RoomService) FindAllRooms() ([]Room, error) {
	return s.roomRepo.FindAllRooms()
}

func (s *RoomService) JoinRoom(roomId primitive.ObjectID, userId primitive.ObjectID) error {
	return s.roomRepo.JoinRoom(roomId, userId)
}

func (s *RoomService) LeaveRoom(roomId primitive.ObjectID, userId primitive.ObjectID) error {
	return s.roomRepo.LeaveRoom(roomId, userId)
}

func (s *RoomService) GetRoomsByUserID(userId primitive.ObjectID) ([]Room, error) {
	return s.roomRepo.FindByUserID(userId)
}
