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
	MarkProcessed(ctx context.Context, eventId string, claimedAt time.Time, processedAt time.Time) error
	Find(ctx context.Context, eventId string) (*models.ProcessedEvent, error)
	TryClaim(ctx context.Context, eventId string, now time.Time, lease time.Duration) (bool, error)
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
	}

	return err
}

func (mr *MongoProcessedEventsRepo) MarkProcessed(
	ctx context.Context,
	eventId string,
	claimedAt time.Time,
	processedAt time.Time,
) error {
	collection := mr.db.Collection(mr.collection)

	filter := bson.D{
		{Key: "event_id", Value: eventId},
		{Key: "claimed_at", Value: claimedAt},
		{
			Key: "$or",
			Value: bson.A{
				bson.D{{Key: "processed_at", Value: bson.D{{Key: "$exists", Value: false}}}},
				bson.D{{Key: "processed_at", Value: nil}},
			},
		},
	}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "processed_at", Value: processedAt},
	}}}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errs.ErrInProgress
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

// TryClaim tries to set claimed_at = now if:
// - processed_at is nil (not done)
// - AND (claimed_at is nil/missing OR claimed_at < now-lease)
func (mr *MongoProcessedEventsRepo) TryClaim(
	ctx context.Context,
	eventId string,
	now time.Time,
	lease time.Duration,
) (bool, error) {
	collection := mr.db.Collection(mr.collection)
	expiredBefore := now.Add(-lease)

	filter := bson.D{
		{Key: "event_id", Value: eventId},
		{
			Key: "$or",
			Value: bson.A{
				bson.D{{Key: "processed_at", Value: bson.D{{Key: "$exists", Value: false}}}},
				bson.D{{Key: "processed_at", Value: nil}},
			},
		},
		{
			Key: "$or",
			Value: bson.A{
				bson.D{{Key: "claimed_at", Value: bson.D{{Key: "$exists", Value: false}}}},
				bson.D{{Key: "claimed_at", Value: nil}},
				bson.D{{Key: "claimed_at", Value: bson.D{{Key: "$lt", Value: expiredBefore}}}},
			},
		},
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "claimed_at", Value: now}}}}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return res.MatchedCount == 1 && res.ModifiedCount == 1, nil
}
