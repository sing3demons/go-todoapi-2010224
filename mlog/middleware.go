package mlog

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Middleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(key, logSessionID(c, logger))
		c.Next()
	}
}

const (
	Session string = "x-session"
)

func logSessionID(c *gin.Context, logger *slog.Logger) *slog.Logger {
	session := c.Request.Header.Get(Session)
	if session == "" {
		uuidV7, err := uuid.NewV7()
		if err != nil {
			session = "unknown" + time.Now().String()
		} else {
			session = uuidV7.String()
		}
		c.Request.Header.Set(Session, session)
	}

	c.Set(Session, session)

	return logger.With(slog.String("session", session))
}
