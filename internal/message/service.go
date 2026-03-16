package message

import "go.mongodb.org/mongo-driver/bson/primitive"

type MessageService struct {
	Repo *MessageRepo
}

func NewMessageService(repo *MessageRepo) *MessageService {
	return &MessageService{Repo: repo}
}

func (s *MessageService) CreateMessage(msg *Message) error {
	return s.Repo.CreateMessage(msg)
}

func (s *MessageService) GetMessagesByRoomId(roomId primitive.ObjectID) ([]Message, error) {
	return s.Repo.GetMessagesByRoomId(roomId)
}
