package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/todoapi/mlog"
)

func PingHandler(c *gin.Context) {
	logger := log.With(slog.String("session", c.GetString(mlog.Session)))

	input := c.Request.Header
	logger.Info("client.input", slog.Any("header", input))

	data := map[string]interface{}{
		"message": "pong",
	}

	logger.Info("client.output", slog.Any("data", data))
	c.JSON(200, data)
}
