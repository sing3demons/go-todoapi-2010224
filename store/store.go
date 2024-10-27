package store

import "github.com/sing3demons/todoapi/model"

type Storer interface {
	Create(*model.Todo) error
	List() ([]model.Todo, error)
	Delete(string) error
	FindOne(id string) (*model.Todo, error)
}
