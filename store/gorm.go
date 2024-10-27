package store

import (
	"fmt"
	"os"

	"github.com/google/uuid"
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
	todo.ID = uuid.New().String()
	return g.db.Create(todo).Error
}

func (g *GormStore) List() ([]todo.Todo, error) {
	var todos []todo.Todo
	r := g.db.Find(&todos)
	if err := r.Error; err != nil {
		return nil, err
	}

	for i := range todos {
		todos[i].Href = fmt.Sprintf("%s/todo/%s", os.Getenv("HOST"), todos[i].ID)
	}
	return todos, nil
}

func (g *GormStore) Delete(id string) error {
	return g.db.Delete(&todo.Todo{}, id).Error
}

func (g *GormStore) FindOne(id string) (*todo.Todo, error) {
	var todo todo.Todo
	r := g.db.First(&todo, "id = ?", id)
	if err := r.Error; err != nil {
		return nil, err
	}
	todo.Href = fmt.Sprintf("%s/todo/%s", os.Getenv("HOST"), todo.ID)
	return &todo, nil
}
