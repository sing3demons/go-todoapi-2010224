package todo

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/sing3demons/todoapi/logger"
)

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

type Todo struct {
	Title     string `json:"text" binding:"required"`
	ID        uint   `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt NullTime `gorm:"index"`
}

func (Todo) TableName() string {
	return "todos"
}

type storer interface {
	Create(*Todo) error
	List() ([]Todo, error)
	Delete(int) error
}

type TodoHandler struct {
	store storer
}

type IContext interface {
	Bind(interface{}) error
	JSON(int, interface{})
	Log() *logger.Logger
	Get(string) interface{}
	TransactionID() string
	Param(string) string
}

func NewTodoHandler(store storer) *TodoHandler {
	return &TodoHandler{store: store}
}

func (t *TodoHandler) NewTask(c IContext) {
	cmd := "new task"
	node := "client"
	logger := c.Log()
	var todo Todo
	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
		return
	}
	logger.AddEvent(node, cmd, todo)

	if todo.Title == "sleep" {
		logger.Error(cmd, slog.Any("error", "not allowed"))
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": "not allowed",
		})
		return
	}

	err := t.store.Create(&todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	logger.Info(cmd, slog.Any("data", todo))
	logger.End()

	c.JSON(http.StatusCreated, map[string]any{
		"ID": todo.ID,
	})
}

func (t *TodoHandler) List(c IContext) {
	cmd := "list task"
	logger := c.Log()
	todos, err := t.store.List()
	if err != nil {
		logger.Error(cmd, slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	logger.AddEvent("client", cmd, todos)
	logger.End()
	c.JSON(http.StatusOK, todos)
}

func (t *TodoHandler) Remove(c IContext) {
	logger := c.Log()
	idParam := c.Param("id")
	cmd := "remove task"

	logger.Info(cmd, slog.Group("param", slog.String("id", idParam)))

	id, err := strconv.Atoi(idParam)
	if err != nil {
		logger.Error(cmd, slog.Any("error", err))
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
		return
	}

	err = t.store.Delete(id)
	if err != nil {
		logger.Error(cmd, slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	data := map[string]any{
		"ID":     id,
		"status": "success",
	}

	logger.Info(cmd, slog.Any("data", data))

	c.JSON(http.StatusOK, data)
}
