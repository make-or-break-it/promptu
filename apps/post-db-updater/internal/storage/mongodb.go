package storage

import (
	"context"
	"fmt"
	"promptu/apps/post-db-updater/internal/config"
	"promptu/apps/post-db-updater/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PromptuMongoClient struct {
	MongoClient  *mongo.Client
	DatabaseName string
}

func NewMongoDbStore(dbName string) *PromptuMongoClient {
	cfg := config.AppConfig()

	// Read more here: https://www.mongodb.com/docs/drivers/go/current/fundamentals/connection/#std-label-golang-connection-guide
	clientOptions := options.
		Client().
		ApplyURI(fmt.Sprintf("%s/?%s", cfg.MongoDbUrl, cfg.MongoDbConnParams))

	connectCtx, connectCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer connectCancel()

	client, err := mongo.Connect(connectCtx, clientOptions)

	if err != nil {
		panic(err)
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		panic(err)
	}

	return &PromptuMongoClient{
		MongoClient:  client,
		DatabaseName: dbName,
	}
}

func (s *PromptuMongoClient) PostMessage(ctx context.Context, post model.Post) error {
	col := s.MongoClient.Database(s.DatabaseName).Collection("posts")

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel()

	_, err := col.InsertOne(dbCtx, post)

	if err != nil {
		return err
	}

	return nil
}
