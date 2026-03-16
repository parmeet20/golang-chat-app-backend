package room

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoomRepo struct {
	collection *mongo.Collection
}

func NewRoomRepo(db *mongo.Database) *RoomRepo {
	return &RoomRepo{
		collection: db.Collection("rooms"),
	}
}

func (r *RoomRepo) Create(room *Room) error {
	room.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(
		context.Background(),
		room,
	)

	return err
}

func (r *RoomRepo) FindByRoomName(name string) (*Room, error) {
	var room Room

	err := r.collection.FindOne(
		context.Background(),
		bson.M{"name": name},
	).Decode(&room)

	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *RoomRepo) FindByID(id primitive.ObjectID) (*Room, error) {
	var room Room

	err := r.collection.FindOne(
		context.Background(),
		bson.M{"_id": id},
	).Decode(&room)

	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *RoomRepo) FindByIDWithMembers(id primitive.ObjectID) (*RoomWithMembers, error) {

	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: id},
			}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "users"},
				{Key: "localField", Value: "members"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "members"},
			}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "users"},
				{Key: "localField", Value: "created_by"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "created_by"},
			}},
		},
		{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$created_by"},
			}},
		},
		{
			{Key: "$project", Value: bson.D{
				{Key: "name", Value: 1},
				{Key: "created_at", Value: 1},
				{Key: "online", Value: 1},

				{Key: "created_by._id", Value: 1},
				{Key: "created_by.username", Value: 1},

				{Key: "members._id", Value: 1},
				{Key: "members.username", Value: 1},
			}},
		},
	}

	cursor, err := r.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var rooms []RoomWithMembers

	if err := cursor.All(context.Background(), &rooms); err != nil {
		return nil, err
	}

	if len(rooms) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &rooms[0], nil
}

func (r *RoomRepo) FindAllRooms() ([]Room, error) {

	limit := int64(50)

	opts := options.Find().SetLimit(limit)

	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{},
		opts,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	var rooms []Room

	if err := cursor.All(context.Background(), &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *RoomRepo) JoinRoom(roomId primitive.ObjectID, userId primitive.ObjectID) error {

	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": roomId},
		bson.M{
			"$addToSet": bson.M{
				"members": userId,
			},
		},
	)

	return err
}

func (r *RoomRepo) LeaveRoom(roomId primitive.ObjectID, userId primitive.ObjectID) error {

	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": roomId},
		bson.M{
			"$pull": bson.M{
				"members": userId,
			},
		},
	)

	return err
}

func (r *RoomRepo) FindByUserID(userId primitive.ObjectID) ([]Room, error) {
	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{"members": userId},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var rooms []Room
	if err := cursor.All(context.Background(), &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}
