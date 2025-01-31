package models

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	client *mongo.Client
}

type UserDoc struct {
	ID   string `bson:"_id" json:"id,omitempty"`
	User *User
}

func NewMongoStore(url string) *MongoStore {
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil
	}

	return &MongoStore{client}
}

func (m MongoStore) Add(name string, user User) error {

	collection := m.client.Database("smartera").Collection("users")
	userDoc := &UserDoc{ID: name, User: &user}
	insertResult, err := collection.InsertOne(context.TODO(), userDoc)

	if err != mongo.ErrNilCursor {
		return err
	}

	if _, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		return nil
	} else {
		return err
	}

}

func (m MongoStore) Get(name string) (User, error) {
	userDoc := &UserDoc{}
	collection := m.client.Database("smartera").Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"_id": name}).Decode(&userDoc)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, ErrNotFound
		}
		panic(err)
	}
	return *userDoc.User, nil
}

func (m MongoStore) Update(name string, user User) error {
	collection := m.client.Database("smartera").Collection("users")
	userDoc := &UserDoc{ID: name, User: &user}
	updateResult, err := collection.ReplaceOne(context.TODO(), bson.M{"_id": name}, userDoc)
	if err != nil {
		panic(err)
	}

	if updateResult.MatchedCount != 0 {
		fmt.Printf("Number of documents replaced: %d\n", updateResult.ModifiedCount)
		return nil
	}
	return ErrNotFound
}
