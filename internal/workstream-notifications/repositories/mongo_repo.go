package repositories

import (
	"context"

	"github.com/danilobml/workstream/internal/platform/errs"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationsRepo interface {
	Create(ctx context.Context, data any, store string) error
}

type MongoRepo struct {
	db *mongo.Database
}

func NewMongoRepo(db *mongo.Database) *MongoRepo {
	return &MongoRepo{
		db: db,
	}
}

func (mr *MongoRepo) Create(ctx context.Context, data any, store string) error {
	collection := mr.db.Collection(store)

	_, err := collection.InsertOne(ctx, data)
	if mongo.IsDuplicateKeyError(err) {
		return errs.ErrAlreadyProcessed
	} else if err != nil {
		return err
	}

	return nil
}
