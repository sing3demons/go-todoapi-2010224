package store

import (
	"github.com/sing3demons/todoapi/model"
)

type Storer interface {
	Create(*model.Todo) error
	List(opt FindOption) ([]model.Todo, error)
	Delete(string) error
	FindOne(id string) (*model.Todo, error)
}

type FindOption struct {
	SearchItem  map[string]interface{}
	CommandName string
	SortItem    map[string]interface{}
}
