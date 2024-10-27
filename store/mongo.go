package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sing3demons/todoapi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	*mongo.Collection
}

func NewMongoStore(db *mongo.Collection) *MongoStore {
	return &MongoStore{db}
}

func (g *MongoStore) Create(todo *model.Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	todo.ID = primitive.NewObjectID().Hex()
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	todo.DeletedAt = nil
	_, err := g.Collection.InsertOne(ctx, todo)

	return err
}

var projection = bson.D{
	{Key: "id", Value: 1},
	{Key: "title", Value: 1},
}

func (g *MongoStore) List() ([]model.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var todos []model.Todo
	filter := bson.D{
		{Key: "deleted_at", Value: nil},
	}
	opts := &options.FindOptions{
		Projection: projection,
	}
	cur, err := g.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &todos); err != nil {
		return nil, err
	}

	for i := range todos {
		todos[i].Href = g.getHref(todos[i].ID)
	}

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

func (g *MongoStore) FindOne(id string) (*model.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var todo model.Todo
	filter := bson.D{
		{Key: "deleted_at", Value: nil},
		{Key: "id", Value: id},
	}
	opts := &options.FindOneOptions{
		Projection: projection,
	}
	err := g.Collection.FindOne(ctx, filter, opts).Decode(&todo)
	if err != nil {
		return nil, err
	}

	todo.Href = g.getHref(todo.ID)
	return &todo, nil
}

func (g *MongoStore) getHref(id string) string {
	return fmt.Sprintf("%s/todo/%s", os.Getenv("HOST"), id)
}
