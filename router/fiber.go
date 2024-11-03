package router

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sing3demons/todoapi/logger"
)

type FiberContext struct {
	*fiber.Ctx
}

func logSessionID(c *fiber.Ctx, logger *slog.Logger) *slog.Logger {
	session := string(c.Request().Header.Peek("x-session"))
	if session == "" {
		uuidV7, err := uuid.NewV7()
		if err != nil {
			session = "unknown" + time.Now().String()
		} else {
			session = uuidV7.String()
		}
		c.Request().Header.Set("x-session", session)
	}

	c.Locals("session", session)

	return logger.With(slog.String("session", session))
}

func NewFiberRouter(logger *slog.Logger) *FiberRouter {
	r := fiber.New()
	r.Use(func(c *fiber.Ctx) error {
		c.Locals("logger", logSessionID(c, logger))
		return c.Next()
	})
	return &FiberRouter{r}
}

func NewFiberContext(c *fiber.Ctx) *FiberContext {
	return &FiberContext{Ctx: c}
}

func (c *FiberContext) Bind(v any) error {
	return c.Ctx.BodyParser(v)
}

func (c *FiberContext) JSON(code int, v any) {
	c.Ctx.JSON(v)
}

func (c *FiberContext) Query(key string) string {
	return c.Ctx.Query(key)
}

func (c *FiberContext) Incoming() map[string]any {
	var data = make(map[string]any)
	body := make(map[string]any)
	c.Ctx.BodyParser(&body)
	params := c.Ctx.AllParams()
	query := c.Ctx.Queries()

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

func (c *FiberContext) Log(name string) logger.ILogDetail {
	route := c.Ctx.Route().Path
	method := c.Ctx.Method()
	device := c.Ctx.Get("User-Agent")

	var instance string
	instance, err := os.Hostname()
	if err != nil {
		instance = c.Ctx.Hostname()
	}

	attribute := map[string]any{
		"route":    route,
		"method":   method,
		"device":   device,
		"instance": instance,
	}

	switch l := c.Ctx.Locals("logger").(type) {
	case *slog.Logger:
		return logger.New(l, name, attribute)
	default:
		return logger.New(slog.Default(), name, attribute)
	}
}

func (c *FiberContext) Get(key string) any {
	return c.Ctx.Locals(key)
}

func (c *FiberContext) TransactionID() string {
	return string(c.Ctx.Request().Header.Peek("TransactionID"))
}

func (c *FiberContext) Param(key string) string {
	return c.Ctx.Params(key)
}

func NewFiberHandler(handler func(IContext)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		handler(NewFiberContext(c))
		return nil
	}
}

type FiberRouter struct {
	*fiber.App
}

// func (r *FiberRouter) Run(addr string) error {
// 	return r.App.Listen(addr)
// }

func (r *FiberRouter) GET(path string, h func(IContext)) {
	r.App.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberContext(c))
		return nil
	})
}

func (r *FiberRouter) POST(path string, h func(IContext)) {
	r.App.Post(path, func(c *fiber.Ctx) error {
		h(NewFiberContext(c))
		return nil
	})
}

func (r *FiberRouter) DELETE(path string, h func(IContext)) {
	r.App.Delete(path, func(c *fiber.Ctx) error {
		h(NewFiberContext(c))
		return nil
	})
}

func (r *FiberRouter) PUT(path string, h func(IContext)) {
	r.App.Put(path, func(c *fiber.Ctx) error {
		h(NewFiberContext(c))
		return nil
	})
}

func (r *FiberRouter) PATCH(path string, h func(IContext)) {
	r.App.Patch(path, func(c *fiber.Ctx) error {
		h(NewFiberContext(c))
		return nil
	})
}

func (r *FiberRouter) Use(h func(IContext)) {
	r.App.Use(func(c *fiber.Ctx) error {
		h(NewFiberContext(c))
		return nil
	})
}

func (r *FiberRouter) Group(path string, h func(IContext)) {
	r.App.Group(path, func(c *fiber.Ctx) error {
		h(NewFiberContext(c))
		return nil
	})
}

func (r *FiberRouter) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := r.Listen(":" + os.Getenv("PORT")); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()

	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	if err := r.Shutdown(); err != nil {
		fmt.Println(err)
	}
}
