package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/danilobml/workstream/internal/platform/errs"
	"github.com/danilobml/workstream/internal/platform/models"
	plmodel "github.com/danilobml/workstream/internal/platform/models"
	model "github.com/danilobml/workstream/internal/workstream-notifications/models"
	"github.com/danilobml/workstream/internal/workstream-notifications/repositories"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventsProcessor interface {
	ProcessEvent(ctx context.Context, event plmodel.Event) error
}

type EventsProcessorService struct {
	repo            repositories.ProcessedEventsRepo
	messageProducer EventsService
}

func NewEventsProcessorService(repo repositories.ProcessedEventsRepo, messsageProducer EventsService) *EventsProcessorService {
	return &EventsProcessorService{
		repo:            repo,
		messageProducer: messsageProducer,
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
		user := models.User{
			Email: "danilobml@hotmail.com",
		}

		mailBody := fmt.Sprintf("Payload: %s", event.Payload)

		log.Printf("publishing mail event for task event_id=%s", event.EventID)
		if err := ns.sendMailMessage(ctx, user, "Event Procesed: Task", mailBody); err != nil {
			log.Printf("createMail failed: %v", err)
			return err
		}
		log.Printf("mail event published")

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

	// TODO Dummy User - Will be replaced, when user service is implemented:
	user := models.User{
		Email: "danilobml@hotmail.com",
	}

	mailBody := fmt.Sprintf("Payload: %s", event.Payload)

	log.Printf("publishing mail event for task event_id=%s", event.EventID)
	if err := ns.sendMailMessage(ctx, user, "Event Procesed: Task", mailBody); err != nil {
		log.Printf("createMail failed: %v", err)
		return err
	}
	log.Printf("mail event published")

	return ns.repo.MarkProcessed(ctx, event.EventID, claimTime, time.Now())
}

func (ns *EventsProcessorService) sendMailMessage(ctx context.Context, user models.User, subject, body string) error {
	mailInput := models.MailInput{
		To:      []string{user.Email},
		Subject: subject,
		Body:    body,
	}
	payloadBytes, err := json.Marshal(mailInput)
	if err != nil {
		return status.Error(codes.Internal, "failed to marshal mail payload")
	}

	event := plmodel.Event{
		EventID:    uuid.NewString(),
		EventType:  "workstream.mail.sent.v1",
		OccurredAt: time.Now(),
		Producer:   "workstream-notifications",
		TraceID:    uuid.NewString(),
		Payload:    json.RawMessage(payloadBytes),
	}

	err = ns.messageProducer.Publish(ctx, event)
	if err != nil {
		return status.Error(codes.Internal, "failed to send event")
	}

	return nil
}
