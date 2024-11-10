package store

import (
	"context"
	"time"

	"github.com/sing3demons/todoapi/logger"
	"github.com/sing3demons/todoapi/model"
	"github.com/sing3demons/todoapi/utils"
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

func (g *MongoStore) Create(todo *model.Todo, logger logger.ILogDetail) error {
	todo.ID = primitive.NewObjectID().Hex()
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	todo.DeletedAt = nil

	store := Store{
		mongo:  g.Collection,
		logger: logger,
	}

	return store.Create("create_todo", "todo", "InsertOne", todo)
}

func (g *MongoStore) List(opt FindOption, logger logger.ILogDetail) ([]model.Todo, error) {
	var todos []model.Todo
	cmd := "list_todo"

	store := Store{
		mongo:  g.Collection,
		logger: logger,
	}

	data, err := store.List(cmd, "todos", opt, todos)
	if err != nil {
		return nil, err
	}

	todos = data.([]model.Todo)

	for i := range todos {
		if todos[i].ID != "" {
			todos[i].Href = utils.GenHref(todos[i].ID)
		}
	}

	logger.AddInput("mongo", "list_todo", todos)

	return todos, nil
}

func (g *MongoStore) Delete(id string, logger logger.ILogDetail) error {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "deleted_at", Value: nil},
		{Key: "id", Value: id},
	}

	logger.AddOutput("mongo", "delete_todo", filter).End()

	r, err := g.Collection.DeleteOne(ctx, filter)
	if err != nil {
		logger.AddError("mongo", "delete_todo", "input", nil, err)
		return err
	}

	logger.AddInput("mongo", "delete_todo", r)
	return nil
}

func (g *MongoStore) FindOne(id string, logger logger.ILogDetail) (*model.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var todo model.Todo
	filter := bson.D{
		{Key: "deleted_at", Value: nil},
		{Key: "id", Value: id},
	}
	opts := &options.FindOneOptions{}
	err := g.Collection.FindOne(ctx, filter, opts).Decode(&todo)
	if err != nil {
		return nil, err
	}

	todo.Href = utils.GenHref(todo.ID)
	return &todo, nil
}
