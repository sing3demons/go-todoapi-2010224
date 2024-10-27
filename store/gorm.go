package store

import (
	"github.com/sing3demons/todoapi/todo"
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (g *GormStore) Create(todo *todo.Todo) error {
	return g.db.Create(todo).Error
}

func (g *GormStore) List() ([]todo.Todo, error) {
	var todos []todo.Todo
	r := g.db.Find(&todos)
	if err := r.Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (g *GormStore) Delete(id int) error {
	return g.db.Delete(&todo.Todo{}, id).Error
}
