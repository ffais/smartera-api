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
	collection *mongo.Collection
}

type UserDoc struct {
	ID   string `bson:"_id" json:"id,omitempty"`
	User *User
}

type ProjectOverviews struct {
	ProjectOverviews []ProjectOverview
}

func NewMongoStore(url string, db string, coll string) *MongoStore {

	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("mongo error: %s", err)
		return nil
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil
	}

	collection := client.Database(db).Collection(coll)

	return &MongoStore{collection}
}

func (m MongoStore) Add(name string, user User) error {
	userDoc := &UserDoc{ID: name, User: &user}
	insertResult, err := m.collection.InsertOne(context.TODO(), userDoc)

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
	err := m.collection.FindOne(context.TODO(), bson.M{"_id": name}).Decode(&userDoc)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	return *userDoc.User, nil
}

func (m MongoStore) Update(name string, user User) error {
	userDoc := &UserDoc{ID: name, User: &user}
	updateResult, err := m.collection.ReplaceOne(context.TODO(), bson.M{"_id": name}, userDoc)
	if err != nil {
		return err
	}

	if updateResult.MatchedCount != 0 {
		fmt.Printf("Number of documents replaced: %d\n", updateResult.ModifiedCount)
		return nil
	}
	return ErrNotFound
}

func (m MongoStore) FindAll() (projects ProjectOverviews, err error) {
	type projectOverviewsDoc struct {
		ID   string `bson:"_id" json:"id,omitempty"`
		User struct {
			ProjectOverviews []ProjectOverview `json:"projectoverviews"`
		} `json:"user"`
	}
	filter := bson.D{}
	opts := options.Find().SetProjection(bson.D{{Key: "_id", Value: 0}, {Key: "user.projectoverviews", Value: 1}})
	cursor, err := m.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ProjectOverviews{}, ErrNotFound
		}
		return ProjectOverviews{}, err
	}
	for cursor.Next(context.TODO()) {
		var result projectOverviewsDoc
		if err := cursor.Decode(&result); err != nil {
			return ProjectOverviews{}, err
		}
		projects.ProjectOverviews = append(projects.ProjectOverviews, result.User.ProjectOverviews...)
	}
	if err := cursor.Err(); err != nil {
		return ProjectOverviews{}, err
	}
	defer cursor.Close(context.TODO())
	return projects, nil
}
