package models

import (
	"encoding/json"
	"time"
)

type Event struct {
	EventID    string          `json:"event_id"`
	EventType  string          `json:"event_type"`
	OccurredAt time.Time       `json:"occurred_at"`
	Producer   string          `json:"producer"`
	TraceID    string          `json:"trace_id,omitempty"`
	Payload    json.RawMessage `json:"payload"`
}
