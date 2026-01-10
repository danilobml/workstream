package repositories

import (
	"context"

	"github.com/danilobml/workstream/internal/platform/errs"
	"github.com/danilobml/workstream/internal/workstream-notifications/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProcesssedEventsRepo interface {
	Insert(ctx context.Context, event models.ProcessedEvent) error
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
