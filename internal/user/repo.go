package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{collection: db.Collection("users")}
}

func (r *UserRepo) Create(user *User) error {

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	_, err = r.collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) FindByUsername(username string) (*User, error) {
	var user User

	err := r.collection.
		FindOne(context.Background(), bson.M{"username": username}).
		Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) FindById(id primitive.ObjectID) (*User, error) {
	var user User
	if err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByEmail(email string) (*User, error) {
	var user User
	if err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Login(username, password string) (*User, error) {

	user, err := r.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
