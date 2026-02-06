package services

import "github.com/danilobml/workstream/internal/platform/models"

type Mailer interface {
	SendMail(mailInput models.MailInput) error
}
