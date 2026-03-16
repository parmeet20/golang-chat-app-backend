package message

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepo struct {
	Collection *mongo.Collection
}

func NewMessageRepo(db *mongo.Database) *MessageRepo {
	return &MessageRepo{Collection: db.Collection("messages")}
}

func (r *MessageRepo) CreateMessage(msg *Message) error {
	if _, err := r.Collection.InsertOne(context.Background(), msg); err != nil {
		return err
	}
	return nil
}
func (r *MessageRepo) GetMessagesByRoomId(roomId primitive.ObjectID) ([]Message, error) {
	var messages []Message

	limit := int64(50)
	cursor, err := r.Collection.Find(context.Background(), bson.M{"room_id": roomId}, options.Find().SetLimit(limit))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
