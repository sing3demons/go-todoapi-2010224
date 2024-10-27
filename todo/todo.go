package todo

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/sing3demons/todoapi/logger"
)

type NullTime struct {
	Time time.Time
}

type Todo struct {
	ID        string     `gorm:"primarykey" json:"id" bson:"id"`
	Title     string     `json:"text" binding:"required"`
	Href      string     `json:"href,omitempty"`
	CreatedAt time.Time  `json:"-" bson:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"-" bson:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"-" bson:"deleted_at,omitempty"`
}

func (Todo) TableName() string {
	return "todos"
}

type storer interface {
	Create(*Todo) error
	List() ([]Todo, error)
	Delete(string) error
	FindOne(id string) (*Todo, error)
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

	err := t.store.Delete(idParam)
	if err != nil {
		logger.Error(cmd, slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	data := map[string]any{
		"ID":     idParam,
		"status": "success",
	}

	logger.Info(cmd, slog.Any("data", data))

	c.JSON(http.StatusOK, data)
}

func (t *TodoHandler) FindOne(c IContext) {
	logger := c.Log()
	idParam := c.Param("id")
	cmd := "find task"

	logger.Info(cmd, slog.Group("param", slog.String("id", idParam)))

	todo, err := t.store.FindOne(idParam)
	if err != nil {
		logger.Error(cmd, slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	logger.Info(cmd, slog.Any("data", todo))
	c.JSON(http.StatusOK, todo)
}
