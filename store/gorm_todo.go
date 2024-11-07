package store

import (
	"fmt"
	"strings"

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
			todos[i].Href = utils.GenHref(todos[i].ID)
		}
	}

	logger.AddInput("gorm", "list_todo", todos)
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
