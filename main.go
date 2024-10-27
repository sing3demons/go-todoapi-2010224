package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
	"github.com/sing3demons/todoapi/router"
	"github.com/sing3demons/todoapi/store"
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
			// return slog.Attr{
			// 	Key:   a.Key,
			// 	Value: slog.StringValue("DETAIL"),
			// }

			// remove level
			return slog.Attr{}
		}
		return a
	},
}))

func init() {
	slog.SetDefault(log)
}

func connectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&todo.Todo{}); err != nil {
		log.Error("failed to migrate", slog.Any("error", err))
	}
	return db
}

func connectMongo() *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Error("failed to connect mongodb", slog.Any("error", err))
	}
	collection := client.Database("myapp").Collection("todos")
	collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{Key: "created_at", Value: 1}},
		},
	})
	// defer client.Disconnect(ctx)
	return collection
}

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Error("Error loading .env file")
	}

	slog.Info("Starting server...")

	r := router.NewMyRouter(log)

	r.GET("/healthz", func(c todo.IContext) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/x", func(c todo.IContext) {
		c.JSON(200, gin.H{
			"buildcommit": buildcommit,
			"buildtime":   buildtime,
		})
	})

	r.GET("/ping", PingHandler)
	r.GET("/transfer/:id", Transfer)

	// store := store.NewGormStore(connectDB())
	store := store.NewMongoStore(connectMongo())
	todoHandler := todo.NewTodoHandler(store)
	r.POST("/todo", todoHandler.NewTask)
	r.GET("/todo", todoHandler.List)

	r.Run()
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

	data := gin.H{"message": "success" + id,
		"id": id,
	}
	logger.Info("finish", slog.Any("data", data))
	c.JSON(http.StatusOK, data)
}
