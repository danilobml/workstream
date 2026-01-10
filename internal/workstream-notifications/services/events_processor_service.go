package services

import (
	"context"
	"time"

	plmodel "github.com/danilobml/workstream/internal/platform/models"
	model "github.com/danilobml/workstream/internal/workstream-notifications/models"
	"github.com/danilobml/workstream/internal/workstream-notifications/repositories"
)

type EventsProcessor interface {
	SaveNewEvent(ctx context.Context, event plmodel.Event) error
}

type EventsProcessorService struct {
	repo repositories.ProcesssedEventsRepo
}

func NewEventsProcessorService(repo repositories.ProcesssedEventsRepo) *EventsProcessorService {
	return &EventsProcessorService{
		repo: repo,
	}
}

func (ns *EventsProcessorService) SaveNewEvent(ctx context.Context, event plmodel.Event) error {
	newNotificationProcessedEvent := model.ProcessedEvent{
		EventID:     event.EventID,
		EventType:   event.EventType,
		OccurredAt:  event.OccurredAt,
		Producer:    event.Producer,
		TraceID:     event.TraceID,
		Payload:     event.Payload,
		ProcessedAt: time.Now(),
	}
	err := ns.repo.Insert(ctx, newNotificationProcessedEvent)
	if err != nil {
		return err
	}

	return nil
}
