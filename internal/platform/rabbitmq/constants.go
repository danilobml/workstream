package rabbitmq

const (
	NotificationsQueue = "workstream.notifications"
	NotificationsExchange = "workstream.events"
	NotificationsBinding  = "workstream.task.#"

	MailerQueue = "workstream.mailer"
	MailerExchange = "workstream.mailer"
	MailerBinding  = "workstream.mail.#"
)
