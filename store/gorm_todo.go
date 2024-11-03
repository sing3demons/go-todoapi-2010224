package store

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/sing3demons/todoapi/logger"
	"github.com/sing3demons/todoapi/model"
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (g *GormStore) Create(todo *model.Todo, logger logger.ILogDetail) error {
	logger.AddOutput("gorm", "create_todo", map[string]interface{}{
		"body": *todo,
	}).End()

	todo.ID = uuid.New().String()
	if err := g.db.Create(todo).Error; err != nil {
		logger.AddError("gorm", "create_todo", "output", todo, err)
		return err
	}

	logger.AddInput("gorm", "create_todo", todo)
	return nil
}

func (g *GormStore) List(opt FindOption, logger logger.ILogDetail) ([]model.Todo, error) {
	var todos []model.Todo
	var conds []interface{}
	if opt.SearchItem != nil {
		for k, v := range opt.SearchItem {
			conds = append(conds, fmt.Sprintf("%s = ?", k), v)
		}
	}

	var order []string
	if opt.SortItem != nil {
		for k, v := range opt.SortItem {
			order = append(order, fmt.Sprintf("%s %s", k, v))
		}
	}

	selectTodo := []string{}
	if opt.SelectItem != nil {
		selectTodo = opt.SelectItem
	}
	logger.AddOutput("gorm", "list_todo", opt).End()

	r := g.db.Select(selectTodo).Order(strings.Join(order, ",")).Find(&todos, conds...)
	if err := r.Error; err != nil {
		logger.AddError("gorm", "list_todo", "output", nil, err)
		return nil, err
	}

	for i := range todos {
		if todos[i].ID != "" {
			todos[i].Href = GenHref(todos[i].ID)
		}
	}

	logger.AddInput("gorm", "list_todo", todos)
	return todos, nil
}

func (g *GormStore) Delete(id string) error {
	return g.db.Delete(&model.Todo{}, id).Error
}

func (g *GormStore) FindOne(id string) (*model.Todo, error) {
	var todo model.Todo
	r := g.db.First(&todo, "id = ?", id)
	if err := r.Error; err != nil {
		return nil, err
	}
	todo.Href = GenHref(todo.ID)
	return &todo, nil
}
