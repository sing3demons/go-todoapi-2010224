package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/sing3demons/todoapi/router"
)

func PingHandler(c router.IContext) {
	logger := c.Log("ping")

	data := map[string]interface{}{
		"message": "pong",
	}

	logger.Info("client.output", slog.Any("data", data))
	c.JSON(200, data)
}

func X(c router.IContext) {
	c.JSON(200, map[string]any{
		"buildcommit": buildcommit,
		"buildtime":   buildtime,
	})
}

func Healthz(c router.IContext) {
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

func Transfer(c router.IContext) {
	logger := c.Log("transfer")
	id := c.Param("id")
	node := "client"
	cmd := "transfer"

	logger.AddInput(node, cmd, c.Incoming())

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
	logger.AddOutput(node, cmd, data).End()
	c.JSON(http.StatusOK, data)
}
