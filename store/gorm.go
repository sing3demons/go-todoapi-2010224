package store

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/sing3demons/todoapi/model"
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (g *GormStore) Create(todo *model.Todo) error {
	todo.ID = uuid.New().String()
	return g.db.Create(todo).Error
}

func (g *GormStore) List(opt FindOption) ([]model.Todo, error) {
	var todos []model.Todo
	var conds []interface{}
	if opt.SearchItem != nil {
		for k, v := range opt.SearchItem {
			conds = append(conds, fmt.Sprintf("%s = ?", k), v)
		}
	}

	if opt.SortItem != nil {
		conds = append(conds, "order by")
		for k, v := range opt.SortItem {
			conds = append(conds, fmt.Sprintf("%s %s", k, v))
		}
	}

	selectTodo := []string{}
	if opt.SelectItem != nil {
		selectTodo = opt.SelectItem
	}
	r := g.db.Select(selectTodo).Find(&todos, conds...)
	if err := r.Error; err != nil {
		return nil, err
	}

	for i := range todos {
		if todos[i].ID != "" {
			todos[i].Href = g.getHref(todos[i].ID)
		}
	}
	return todos, nil
}

func (g *GormStore) Delete(id string) error {
	return g.db.Delete(&model.Todo{}, id).Error
}

func (g *GormStore) getHref(id string) string {
	return fmt.Sprintf("%s/todo/%s", os.Getenv("HOST"), id)
}

func (g *GormStore) FindOne(id string) (*model.Todo, error) {
	var todo model.Todo
	r := g.db.First(&todo, "id = ?", id)
	if err := r.Error; err != nil {
		return nil, err
	}
	todo.Href = g.getHref(todo.ID)
	return &todo, nil
}
