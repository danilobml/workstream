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
	ProcessEvent(ctx context.Context, event plmodel.Event) error
}

type EventsProcessorService struct {
	repo repositories.ProcessedEventsRepo
}

func NewEventsProcessorService(repo repositories.ProcessedEventsRepo) *EventsProcessorService {
	return &EventsProcessorService{
		repo: repo,
	}
}

func (ns *EventsProcessorService) ProcessEvent(ctx context.Context, event plmodel.Event) error {
	claimTime := time.Now()
	lease := 30 * time.Second

	claim := model.ProcessedEvent{
		EventID:     event.EventID,
		EventType:   event.EventType,
		OccurredAt:  event.OccurredAt,
		Producer:    event.Producer,
		TraceID:     event.TraceID,
		Payload:     event.Payload,
		ProcessedAt: nil,
		ClaimedAt:   &claimTime,
	}

	err := ns.repo.Insert(ctx, claim)
	if err == nil {
		// TODO side effects
		return ns.repo.MarkProcessed(ctx, event.EventID, claimTime, time.Now())
	}

	if !errors.Is(err, errs.ErrAlreadyProcessed) {
		return err
	}

	existing, findErr := ns.repo.Find(ctx, event.EventID)
	if findErr != nil {
		return findErr
	}
	if existing != nil && existing.ProcessedAt != nil {
		return errs.ErrAlreadyProcessed
	}

	claimTime = time.Now()
	claimed, claimErr := ns.repo.TryClaim(ctx, event.EventID, claimTime, lease)
	if claimErr != nil {
		return claimErr
	}
	if !claimed {
		return errs.ErrInProgress
	}

	// TODO side effects
	return ns.repo.MarkProcessed(ctx, event.EventID, claimTime, time.Now())
}
