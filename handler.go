package main

import (
	"log/slog"
	"net/http"
	"time"

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

func X(c todo.IContext) {
	c.JSON(200, map[string]any{
		"buildcommit": buildcommit,
		"buildtime":   buildtime,
	})
}

func Healthz(c todo.IContext) {
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

func Transfer(c todo.IContext) {
	logger := c.Log()
	id := c.Param("id")

	logger.Info("parsing...", slog.String("id", id))
	time.Sleep(time.Millisecond * 200)
	logger.Info("validating...", slog.String("id", id))
	time.Sleep(time.Millisecond * 100)
	logger.Info("staging...", slog.String("id", id))
	time.Sleep(time.Millisecond * 200)
	logger.Info("transection starting...", slog.String("id", id))
	time.Sleep(time.Millisecond * 300)
	logger.Info("drawing...", slog.String("id", id))
	time.Sleep(time.Millisecond * 400)
	logger.Info("depositing...", slog.String("id", id))
	time.Sleep(time.Millisecond * 400)
	logger.Info("transection ending...", slog.String("id", id))
	time.Sleep(time.Millisecond * 100)
	logger.Info("responding...", slog.String("id", id))
	time.Sleep(time.Millisecond * 100)

	data := map[string]any{"message": "success" + id,
		"id": id,
	}
	logger.Info("finish", slog.Any("data", data))
	c.JSON(http.StatusOK, data)
}
