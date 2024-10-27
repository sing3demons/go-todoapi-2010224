package store

import (
	"context"
	"fmt"
	"time"

	"github.com/sing3demons/todoapi/todo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStore struct {
	*mongo.Collection
}

func NewMongoStore(db *mongo.Collection) *MongoStore {
	return &MongoStore{db}
}

func (g *MongoStore) Create(todo *todo.Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	todo.ID = primitive.NewObjectID().Hex()
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	todo.DeletedAt = nil
	_, err := g.Collection.InsertOne(ctx, todo)

	return err
}

func (g *MongoStore) List() ([]todo.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var todos []todo.Todo
	filter := bson.D{
		{Key: "deleted_at", Value: nil},
	}
	cur, err := g.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &todos); err != nil {
		return nil, err
	}

	fmt.Println("todos", todos)

	return todos, nil
}

func (g *MongoStore) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := g.Collection.DeleteOne(ctx, bson.D{
		{Key: "deleted_at", Value: nil},
		{Key: "id", Value: id},
	})
	return err
}
