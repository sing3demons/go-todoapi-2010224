package store

import (
	"github.com/sing3demons/todoapi/logger"
	"github.com/sing3demons/todoapi/model"
)

type Storer interface {
	Create(*model.Todo, logger.ILogDetail) error
	List(opt FindOption, logger logger.ILogDetail) ([]model.Todo, error)
	Delete(id string, logger logger.ILogDetail) error
	FindOne(id string, logger logger.ILogDetail) (*model.Todo, error)
}

type FindOption struct {
	SearchItem  map[string]interface{}
	CommandName string
	SortItem    map[string]interface{}
	SelectItem  []string
}
