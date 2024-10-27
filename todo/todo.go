package todo

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sing3demons/todoapi/model"
	"github.com/sing3demons/todoapi/router"
	"github.com/sing3demons/todoapi/store"
)

type TodoHandler struct {
	store store.Storer
}

func NewTodoHandler(store store.Storer) *TodoHandler {
	return &TodoHandler{store: store}
}

func (t *TodoHandler) NewTask(c router.IContext) {
	cmd := "new task"
	node := "client"
	logger := c.Log("new_task")
	var todo model.Todo
	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
		return
	}
	logger.AddInput(node, cmd, todo)

	if todo.Title == "sleep" {
		logger.AddError(node, cmd, "output", map[string]any{
			"error": "not allowed",
		}, fmt.Errorf("not allowed"))
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

	logger.AddOutput(node, cmd, todo)
	logger.End()

	c.JSON(http.StatusCreated, map[string]any{
		"ID": todo.ID,
	})
}

func (t *TodoHandler) List(c router.IContext) {
	cmd := "list task"
	logger := c.Log("tasks_list")
	logger.AddInput("client", cmd, nil)
	todos, err := t.store.List()
	if err != nil {
		logger.Error(cmd, slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	logger.AddOutput("client", cmd, todos)
	logger.End()
	c.JSON(http.StatusOK, todos)
}

func (t *TodoHandler) Remove(c router.IContext) {
	logger := c.Log("remove_task")
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

func (t *TodoHandler) FindOne(c router.IContext) {
	logger := c.Log("find_task")
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
