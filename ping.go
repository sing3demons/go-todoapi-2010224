package main

import (
	"log/slog"

	"github.com/sing3demons/todoapi/todo"
)

func PingHandler(c todo.IContext) {
	logger := c.Log()

	data := map[string]interface{}{
		"message": "pong",
	}

	logger.Info("client.output", slog.Any("data", data))
	c.JSON(200, data)
}
