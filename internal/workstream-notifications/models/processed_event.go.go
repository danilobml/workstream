package models

import "time"

type ProcessedEvent struct {
	EventID     string     `bson:"event_id"`
	EventType   string     `bson:"event_type"`
	OccurredAt  time.Time  `bson:"occurred_at"`
	Producer    string     `bson:"producer"`
	TraceID     string     `bson:"trace_id"`
	Payload     []byte     `bson:"payload"`
	ProcessedAt *time.Time `bson:"processed_at,omitempty"`
	ClaimedAt   *time.Time `bson:"claimed_at,omitempty"`
}
