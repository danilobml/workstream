package repositories

import (
	"context"
	"time"

	"github.com/danilobml/workstream/internal/platform/errs"
	"github.com/danilobml/workstream/internal/workstream-notifications/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProcessedEventsRepo interface {
	Insert(ctx context.Context, event models.ProcessedEvent) error
	MarkProcessed(ctx context.Context, eventId string, processedAt time.Time) error
	Find(ctx context.Context, eventId string) (*models.ProcessedEvent, error)
}

type MongoProcessedEventsRepo struct {
	db         *mongo.Database
	collection string
}

func NewMongoProcessedEventsRepo(db *mongo.Database) *MongoProcessedEventsRepo {
	return &MongoProcessedEventsRepo{
		db:         db,
		collection: "processed_events",
	}
}

func (mr *MongoProcessedEventsRepo) Insert(ctx context.Context, event models.ProcessedEvent) error {
	collection := mr.db.Collection(mr.collection)

	_, err := collection.InsertOne(ctx, event)
	if mongo.IsDuplicateKeyError(err) {
		return errs.ErrAlreadyProcessed
	} else if err != nil {
		return err
	}

	return nil
}

func (mr *MongoProcessedEventsRepo) MarkProcessed(ctx context.Context, eventId string, processedAt time.Time) error {
	collection := mr.db.Collection(mr.collection)

	filter := bson.D{
		{Key: "event_id", Value: eventId},
		{Key: "processed_at", Value: nil},
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "processed_at", Value: processedAt}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (mr *MongoProcessedEventsRepo) Find(ctx context.Context, eventId string) (*models.ProcessedEvent, error) {
	collection := mr.db.Collection(mr.collection)

	var event models.ProcessedEvent
	filter := bson.M{"event_id": eventId}
	err := collection.FindOne(ctx, filter).Decode(&event)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &event, nil
}
