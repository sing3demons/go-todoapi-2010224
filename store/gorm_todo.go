package store

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sing3demons/todoapi/logger"
	"github.com/sing3demons/todoapi/model"
	"github.com/sing3demons/todoapi/utils"
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (g *GormStore) Create(todo *model.Todo, logger logger.ILogDetail) error {
	store := Store{
		sql:    g.db,
		logger: logger,
	}

	todo.ID = uuid.New().String()
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	return store.Create("create_todo", "todo", "create", todo)
}

func (g *GormStore) List(opt FindOption, logger logger.ILogDetail) ([]model.Todo, error) {
	var todos []model.Todo

	store := Store{
		sql:    g.db,
		logger: logger,
	}

	data, err := store.List("list_todo", "todos", opt, todos)
	if err != nil {
		return nil, err
	}

	todos = data.([]model.Todo)

	for i := range todos {
		if todos[i].ID != "" {
			todos[i].Href = utils.GenHref(todos[i].ID)
		}
	}

	return todos, nil
}

func (g *GormStore) Delete(id string, logger logger.ILogDetail) error {
	node := "gorm"
	cmd := "delete_todo"
	query := "id = ?"
	logger.AddOutput(node, cmd, map[string]any{
		"query": strings.Replace(query, "?", id, 1),
	}).End()
	r := g.db.Where(query, id).Delete(&model.Todo{})

	if r.Error != nil {
		logger.AddError(node, cmd, "output", nil, r.Error)
		return r.Error
	}
	logger.AddOutput(node, cmd, r.RowsAffected).End()
	return nil
}

func (g *GormStore) FindOne(id string, logger logger.ILogDetail) (*model.Todo, error) {
	node := "gorm"
	cmd := "find_one_todo"
	logger.AddOutput(node, cmd, id).End()
	var todo model.Todo
	r := g.db.First(&todo, "id = ?", id)
	if err := r.Error; err != nil {
		logger.AddError(node, cmd, "output", nil, err)
		return nil, err
	}
	todo.Href = utils.GenHref(todo.ID)
	logger.AddInput(node, cmd, todo)
	return &todo, nil
}
