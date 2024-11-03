package main

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sing3demons/todoapi/router"
	"github.com/sing3demons/todoapi/todo"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

var f = &lumberjack.Logger{
	LocalTime:  true,
	Compress:   true,
	Filename:   "log/details/todoapi.log",
	MaxSize:    10, // megabytes
	MaxBackups: 5,
	MaxAge:     30, // days
}

var log *slog.Logger = slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, f), &slog.HandlerOptions{
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
	todoHandler := todo.NewTodoHandler(conn.GormStore())
	r.POST("/todo", todoHandler.NewTask)
	r.GET("/todo/:id", todoHandler.FindOne)
	r.GET("/todo", todoHandler.List)

	r.Run()
}
