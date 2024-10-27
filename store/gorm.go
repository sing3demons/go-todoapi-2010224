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

var selectTodo = []string{"id", "title"}

func (g *GormStore) List() ([]todo.Todo, error) {
	var todos []todo.Todo
	r := g.db.Select(selectTodo).Find(&todos)
	if err := r.Error; err != nil {
		return nil, err
	}

	for i := range todos {
		todos[i].Href = g.getHref(todos[i].ID)
	}
	return todos, nil
}

func (g *GormStore) Delete(id string) error {
	return g.db.Delete(&todo.Todo{}, id).Error
}

func (g *GormStore) getHref(id string) string {
	return fmt.Sprintf("%s/todo/%s", os.Getenv("HOST"), id)
}

func (g *GormStore) FindOne(id string) (*todo.Todo, error) {
	var todo todo.Todo
	r := g.db.Select(selectTodo).First(&todo, "id = ?", id)
	if err := r.Error; err != nil {
		return nil, err
	}
	todo.Href = g.getHref(todo.ID)
	return &todo, nil
}
