package services

import (
	"context"
	"errors"
	"time"

	"github.com/danilobml/workstream/internal/platform/errs"
	plmodel "github.com/danilobml/workstream/internal/platform/models"
	model "github.com/danilobml/workstream/internal/workstream-notifications/models"
	"github.com/danilobml/workstream/internal/workstream-notifications/repositories"
)

type EventsProcessor interface {
	SaveNewEvent(ctx context.Context, event plmodel.Event) error
}

type EventsProcessorService struct {
	repo repositories.ProcessedEventsRepo
}

func NewEventsProcessorService(repo repositories.ProcessedEventsRepo) *EventsProcessorService {
	return &EventsProcessorService{
		repo: repo,
	}
}

func (ns *EventsProcessorService) SaveNewEvent(ctx context.Context, event plmodel.Event) error {
	claim := model.ProcessedEvent{
		EventID:    event.EventID,
		EventType:  event.EventType,
		OccurredAt: event.OccurredAt,
		Producer:   event.Producer,
		TraceID:    event.TraceID,
		Payload:    event.Payload,
	}

	err := ns.repo.Insert(ctx, claim)
	if err != nil && !errors.Is(err, errs.ErrAlreadyProcessed) {
		return err
	}

	if errors.Is(err, errs.ErrAlreadyProcessed) {
		existing, findErr := ns.repo.Find(ctx, event.EventID)
		if findErr != nil {
			return findErr
		}

		if existing.ProcessedAt != nil {
			return errs.ErrAlreadyProcessed
		}

		return errs.ErrInProgress
	}

	// Perform other operations:
	return ns.repo.MarkProcessed(ctx, event.EventID, time.Now())
}
