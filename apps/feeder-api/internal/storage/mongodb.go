package storage

import (
	"context"
	"fmt"
	"promptu/apps/feeder-api/internal/config"
	"promptu/apps/feeder-api/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func (s *PromptuMongoClient) GetFeed(ctx context.Context, date time.Time) ([]model.Post, error) {
	currDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayAfter := time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, date.Location())

	col := s.MongoClient.Database(s.DatabaseName).Collection("posts")

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel()

	opts := options.Find().SetSort(bson.D{{"utcCreatedAt", -1}})

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"utcCreatedAt", bson.D{{"$gte", currDay}}}},
				bson.D{{"utcCreatedAt", bson.D{{"$lt", dayAfter}}}},
			},
		},
	}

	result, err := col.Find(dbCtx, filter, opts)

	if err != nil {
		panic(err)
	}

	feed := []model.Post{}

	resultCtx, resultCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer resultCancel()

	if err = result.All(resultCtx, &feed); err != nil {
		panic(err)
	}

	return feed, nil
}
