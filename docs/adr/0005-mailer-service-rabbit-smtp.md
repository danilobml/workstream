# ADR-005: workstream-mailer

- Date: 2026-01-25
- Status: Accepted

## Decision

Implement **workstream-mailer** as a **RabbitMQ consumer** that reads from `workstream.mailer` (MailerQueue) and sends emails via **SMTP** using Go `net/smtp`.

- Message format: `models.Event` where `event.payload` is a JSON-encoded `models.MailInput` (`to`, `subject`, `body`).
- Transport: RabbitMQ (asynchronous; decouples producers from email delivery).
- Delivery mechanism: SMTP using configured credentials (local/dev friendly).

## Error handling

- Invalid envelope (`Event`) or invalid `MailInput` payload → **Nack(requeue=false)** (goes to DLQ if configured).
- Transient send failure → **Nack(requeue=true)** (retry by requeue).
- Success → **Ack**.
- If SMTP config is missing → return `errs.ErrMailServiceDisabled` (treated as failure by consumer).

## Consequences

- Producers publish mail requests as events; they do not send emails directly.
- Email delivery is best-effort with retries; ordering is not guaranteed.
