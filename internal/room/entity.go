package room

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name      string               `json:"name" bson:"name" validate:"required,min=3,max=50"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	CreatedBy primitive.ObjectID   `json:"created_by" bson:"created_by"`
	Members   []primitive.ObjectID `json:"members" bson:"members"`
	Online    int                  `json:"online" bson:"online"`
}

type RoomWithMembers struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"created_at"`

	CreatedBy struct {
		ID       primitive.ObjectID `bson:"_id"`
		Username string             `bson:"username"`
	} `bson:"created_by"`

	Members []struct {
		ID       primitive.ObjectID `bson:"_id"`
		Username string             `bson:"username"`
	} `bson:"members"`

	Online int `bson:"online"`
}
