package services

import (
	"context"
	"time"

	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/danilobml/workstream/internal/workstream-notifications/mongodb"
	"github.com/danilobml/workstream/internal/workstream-notifications/repositories"
)

type NotificationService interface {
	CreateNewNotification(ctx context.Context, event models.Event) error
}

type MongoService struct {
	repo repositories.NotificationsRepo
}

func NewMongoService(repo repositories.NotificationsRepo) *MongoService {
	return &MongoService{
		repo: repo,
	}
}

func (ms *MongoService) CreateNewNotification(ctx context.Context, event models.Event) error {
	payloadBytes := []byte(event.Payload)

	newNotification := mongodb.ProcessedEvent{
		EventID:     event.EventID,
		EventType:   event.EventType,
		OccurredAt:  event.OccurredAt,
		Producer:    event.Producer,
		TraceID:     event.TraceID,
		Payload:     payloadBytes,
		ProcessedAt: time.Now(),
	}

	err := ms.repo.Create(ctx, newNotification, "processed_events")
	if err != nil {
		return err
	}

	return nil
}
