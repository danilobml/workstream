package models

import "time"

type Event struct {
	EventID    string    `json:"event_id"`
	EventType  string    `json:"event_type"`
	OccurredAt time.Time `json:"occurred_at"`
	Producer   string    `json:"producer"`
	TraceID    string    `json:"trace_id,omitempty"`
	Payload    any       `json:"payload"`
}
