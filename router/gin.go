package router

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sing3demons/todoapi/logger"
	"github.com/sing3demons/todoapi/mlog"
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

func (c *MyContext) Log(name string) logger.ILogDetail {
	route := c.FullPath()
	method := c.Request.Method
	device := c.Context.Request.UserAgent()

	var instance string
	instance, err := os.Hostname()
	if err != nil {
		instance = c.Context.Request.Host
	}
	attribute := map[string]any{
		"route":    route,
		"method":   method,
		"device":   device,
		"instants": instance,
	}
	switch l := c.Value("logger").(type) {
	case *slog.Logger:
		return logger.New(l, name, attribute)
	default:
		return logger.New(slog.Default(), name, attribute)
	}
}

func (c *MyContext) Incoming() map[string]any {
	var data = make(map[string]any)
	body := make(map[string]any)
	c.Context.BindJSON(&body)
	params := c.Context.Keys
	query := c.Context.Request.URL.Query()

	if len(query) != 0 {
		data["query"] = query
	}

	if len(params) != 0 {
		data["params"] = params
	}

	if len(body) != 0 {
		data["body"] = body
	}

	return data
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

func (c *MyContext) Query(key string) string {
	return c.Context.Query(key)
}

func NewGinHandler(handler func(IContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(NewMyContext(c))
	}
}

type MyRouter struct {
	*gin.Engine
}

func NewMyRouter(logger *slog.Logger) *MyRouter {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(mlog.Middleware(logger))
	return &MyRouter{r}
}

func (r *MyRouter) GET(path string, handler func(IContext)) {
	r.Engine.GET(path, NewGinHandler(handler))
}

func (r *MyRouter) POST(path string, handler func(IContext)) {
	r.Engine.POST(path, NewGinHandler(handler))
}

func (r *MyRouter) DELETE(path string, handler func(IContext)) {
	r.Engine.DELETE(path, NewGinHandler(handler))
}

func (r *MyRouter) PUT(path string, handler func(IContext)) {
	r.Engine.PUT(path, NewGinHandler(handler))
}

func (r *MyRouter) PATCH(path string, handler func(IContext)) {
	r.Engine.PATCH(path, NewGinHandler(handler))
}

func (r *MyRouter) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		fmt.Println("server started at", s.Addr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}

	return nil
}
