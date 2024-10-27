package main

import (
	"log/slog"
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
