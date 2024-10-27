package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sing3demons/todoapi/router"
	"github.com/sing3demons/todoapi/todo"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

var log *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "time" {
			return slog.Attr{
				Key:   "@timestamp",
				Value: a.Value,
			}
		}

		if a.Key == "msg" {
			return slog.Attr{
				Key:   "event",
				Value: a.Value,
			}
		}

		if a.Key == "level" {
			return slog.Attr{}
		}
		return a
	},
}))

func init() {
	slog.SetDefault(log)
}

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Error("Error loading .env file")
	}

	slog.Info("Starting server...")

	r := router.NewFiberRouter(log)

	r.GET("/healthz", Healthz)

	r.GET("/x", X)

	r.GET("/ping", PingHandler)
	r.GET("/transfer/:id", Transfer)

	conn := db{}
	defer conn.Close()

	todoHandler := todo.NewTodoHandler(conn.MongoStore())
	r.POST("/todo", todoHandler.NewTask)
	r.GET("/todo/:id", todoHandler.FindOne)
	r.GET("/todo", todoHandler.List)

	r.Run()
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
