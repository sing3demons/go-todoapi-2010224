package todo

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type GinContext struct {
	*gin.Context
}

// func NewGinContext(c *gin.Context) *GinContext {
// 	return &GinContext{c}
// }

func (c *GinContext) Bind(v any) error {
	return c.ShouldBindJSON(v)
}

func (c *GinContext) JSON(code int, v any) {
	c.Context.JSON(code, v)
}

func (c *GinContext) Log() *slog.Logger {
	switch logger := c.Value("logger").(type) {
	case *slog.Logger:
		return logger
	default:
		return slog.Default()
	}
}
