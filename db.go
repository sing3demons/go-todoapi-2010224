package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/sing3demons/todoapi/model"
	"github.com/sing3demons/todoapi/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_CONN")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&model.Todo{}); err != nil {
		log.Error("failed to migrate", slog.Any("error", err))
	}

	return db
}

func connectMongo() *mongo.Client {
	loggerOptions := options.
		Logger().
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetLoggerOptions(loggerOptions))
	if err != nil {
		log.Error("failed to connect mongodb", slog.Any("error", err))
	}

	// defer client.Disconnect(ctx)
	return client
}

type db struct {
	client *mongo.Client
}

func (d *db) GormStore() *store.GormStore {
	return store.NewGormStore(connectDB())
}

func (d *db) MongoStore() *store.MongoStore {
	client := connectMongo()
	collection := client.Database("myapp").Collection("todos")
	collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{Key: "created_at", Value: 1}},
		},
	})

	d.client = client

	return store.NewMongoStore(collection)
}

func (d *db) Close() {
	if d.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := d.client.Disconnect(ctx)
		if err != nil {
			log.Error("failed to disconnect mongodb", slog.Any("error", err))
		}

		log.Info("mongodb disconnected")
	}
}
