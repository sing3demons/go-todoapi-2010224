package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sing3demons/todoapi/logger"
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

func (g *MongoStore) Create(todo *model.Todo, logger logger.ILogDetail) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	logger.AddOutput("mongo", "create_todo", map[string]interface{}{
		"body": *todo,
	})
	todo.ID = primitive.NewObjectID().Hex()
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	todo.DeletedAt = nil
	_, err := g.Collection.InsertOne(ctx, todo)
	logger.AddInput("mongo", "create_todo", todo)

	return err
}

var projection = bson.D{
	{Key: "id", Value: 1},
	{Key: "title", Value: 1},
}

func (g *MongoStore) List(opt FindOption, logger logger.ILogDetail) ([]model.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var todos []model.Todo
	filter := bson.D{
		{Key: "deleted_at", Value: nil},
	}

	if opt.SearchItem != nil {
		for k, v := range opt.SearchItem {
			filter = append(filter, bson.E{Key: k, Value: v})
		}
	}

	opts := &options.FindOptions{}
	opts.Sort = bson.D{{Key: "created_at", Value: -1}}

	if opt.SelectItem != nil {
		projection = bson.D{}
		for _, v := range opt.SelectItem {
			projection = append(projection, bson.E{Key: v, Value: 1})
		}

		opts.Projection = projection
	}

	if opt.SortItem != nil {
		opts.Sort = opt.SortItem
	}

	if opt.SortItem != nil {
		for k, v := range opt.SortItem {
			if v == "asc" {
				opts.Sort = bson.D{{Key: k, Value: 1}}
			} else {
				opts.Sort = bson.D{{Key: k, Value: -1}}
			}
		}
	}

	logger.AddOutput("mongo", "list_todo", opt).End()

	cur, err := g.Collection.Find(ctx, filter, opts)
	if err != nil {
		logger.AddError("mongo", "list_todo", "output", nil, err)
		return nil, err
	}

	if err := cur.All(ctx, &todos); err != nil {
		return nil, err
	}

	for i := range todos {
		if todos[i].ID != "" {
			todos[i].Href = GenHref(todos[i].ID)
		}
	}

	logger.AddInput("mongo", "list_todo", todos)

	return todos, nil
}

func (g *MongoStore) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "deleted_at", Value: nil},
		{Key: "id", Value: id},
	}

	_, err := g.Collection.DeleteOne(ctx, filter)
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

	todo.Href = GenHref(todo.ID)
	return &todo, nil
}

func GenHref(id string) string {
	if os.Getenv("HOST") == "" {
		return fmt.Sprintf("%s/todo/%s", "{{HOST}}", id)
	}
	return fmt.Sprintf("%s/todo/%s", os.Getenv("HOST"), id)
}