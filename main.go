package main

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

	"github.com/joho/godotenv"
	"github.com/sing3demons/todoapi/mlog"
	"github.com/sing3demons/todoapi/todo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

var log *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "time" {
			return slog.Attr{
				Key: "@timestamp",
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

	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&todo.Todo{}); err != nil {
		log.Error("failed to migrate", slog.Any("error", err))
	}
	slog.Info("Starting server...")
	// log json

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(mlog.Middleware(log))

	r.GET("/healthz", func(c *gin.Context) {
		c.Status(200)
	})

	r.GET("/x", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"buildcommit": buildcommit,
			"buildtime":   buildtime,
		})
	})

	r.GET("/ping", PingHandler)
	r.GET("/transfer/:id", Transfer)

	gormStore := todo.NewGormStore(db)
	todoHandler := todo.NewTodoHandler(gormStore)
	r.POST("/todo", todoHandler.NewTask)
	r.GET("/todo", todoHandler.List)

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
			log.Error("listen", slog.Any("error", err))
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

}

func Transfer(c *gin.Context) {
	logger := log.With(slog.String("session", c.GetString(mlog.Session)))
	id := c.Param("id")
	device := c.Request.UserAgent()

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

	data := gin.H{"message": "success" + id,
		"id":     id,
		"device": device}
	logger.Info("finish", slog.Any("data", data))
	c.JSON(http.StatusOK, data)
}
