package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/todoapi/logger"
	"github.com/sing3demons/todoapi/todo"
)

type MyContext struct {
	*gin.Context
}

func NewMyContext(c *gin.Context) *MyContext {
	return &MyContext{c}
}

func (c *MyContext) Bind(v any) error {
	return c.ShouldBindJSON(v)
}

func (c *MyContext) JSON(code int, v any) {
	c.Context.JSON(code, v)
}

func (c *MyContext) Log() *logger.Logger {
	route := c.FullPath()
	method := c.Request.Method
	device := c.Context.Request.UserAgent()
	attribute := map[string]any{
		"route":  route,
		"method": method,
		"device": device,
	}
	switch l := c.Value("logger").(type) {
	case *slog.Logger:
		return logger.New(l, attribute)
	default:
		return logger.New(slog.Default(), attribute)
	}
}

func (c *MyContext) Get(key string) any {
	v, _ := c.Context.Get(key)
	return v
}

func (c *MyContext) TransactionID() string {
	return c.Request.Header.Get("TransactionID")
}

func (c *MyContext) Param(key string) string {
	return c.Context.Param(key)
}

func NewGinHandler(handler func(todo.IContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(NewMyContext(c))
	}
}
