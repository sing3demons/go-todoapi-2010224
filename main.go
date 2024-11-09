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

func NewLogger() *slog.Logger {
	serviceName := "todoapi"
	time := time.Now().Format("2006-01-02")

	fileName := "logs/details/" + serviceName + "_" + time + ".log"
	f := &lumberjack.Logger{
		LocalTime:  true,
		Compress:   true,
		Filename:   fileName,
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     1, // days
	}
	return slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, f), &slog.HandlerOptions{
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
}

var log *slog.Logger = NewLogger()

func init() {
	slog.SetDefault(log)

}

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Error("Error loading .env file")
	}

	slog.Debug("Starting server...")

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
	r.DELETE("/todo/:id", todoHandler.Delete)

	r.Run()
}
