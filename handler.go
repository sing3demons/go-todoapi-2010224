package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/sing3demons/todoapi/router"
	"github.com/sing3demons/todoapi/utils"
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

	logger.Debug("parsing...", slog.String("id", id))
	time.Sleep(time.Millisecond * 200)
	logger.Debug("validating...", slog.String("id", id))
	time.Sleep(time.Millisecond * 100)
	logger.Debug("staging...", slog.String("id", id))
	time.Sleep(time.Millisecond * 200)
	logger.Debug("transection starting...", slog.String("id", id))
	time.Sleep(time.Millisecond * 300)
	logger.Debug("drawing...", slog.String("id", id))
	time.Sleep(time.Millisecond * 400)
	logger.Debug("depositing...", slog.String("id", id))
	time.Sleep(time.Millisecond * 400)
	logger.Debug("transection ending...", slog.String("id", id))
	time.Sleep(time.Millisecond * 100)
	logger.Debug("responding...", slog.String("id", id))
	time.Sleep(time.Millisecond * 100)

	data := map[string]any{
		"message": "success" + id,
		"id":      id,
		"href":    utils.GenHref(id),
	}
	logger.AddOutput(node, cmd, data).End()
	c.JSON(http.StatusOK, data)
}
