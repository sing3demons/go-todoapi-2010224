package mlog

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

const key = "logger"

func L(c *gin.Context) *slog.Logger {
	switch logger := c.Value(key).(type) {
	case *slog.Logger:
		return logger
	default:
		return slog.Default()
	}
}
