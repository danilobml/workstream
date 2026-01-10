package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitMongoDB(ctx context.Context, mongoUri, dbName string) (*mongo.Database, *mongo.Client, error) {
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(mongoUri),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("connection error :%w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, fmt.Errorf("mongoDB failed to ping :%w", err)
	}
	log.Println("MongoDb ping success")

	database := client.Database(dbName)

	return database, client, nil
}
