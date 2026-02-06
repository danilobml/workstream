package services

import (
	"bytes"
	"log"
	"net/smtp"
	"strings"

	"github.com/danilobml/workstream/internal/platform/errs"
	"github.com/danilobml/workstream/internal/platform/models"
)

type LocalMailConfig struct {
	FromEmail     string
	FromEmailPass string
	FromEmailSMTP string
	SMTPAddr      string
}

type LocalMailService struct {
	fromEmail     string
	fromEmailPass string
	fromEmailSMTP string
	smtpAddr      string
}

func NewLocalMailService(cfg LocalMailConfig) *LocalMailService {
	return &LocalMailService{
		fromEmail:     cfg.FromEmail,
		fromEmailPass: cfg.FromEmailPass,
		fromEmailSMTP: cfg.FromEmailSMTP,
		smtpAddr:      cfg.SMTPAddr,
	}
}

func (ms *LocalMailService) SendMail(mailInput models.MailInput) error {
	if ms.fromEmail == "" || ms.fromEmailPass == "" || ms.fromEmailSMTP == "" || ms.smtpAddr == "" {
		return errs.ErrMailServiceDisabled
	}

	var msg bytes.Buffer
	msg.WriteString("From: " + ms.fromEmail + "\r\n")
	msg.WriteString("To: " + strings.Join(mailInput.To, ", ") + "\r\n")
	msg.WriteString("Subject: " + mailInput.Subject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	msg.WriteString("Content-Transfer-Encoding: 7bit\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(mailInput.Body + "\r\n")

	auth := smtp.PlainAuth("", ms.fromEmail, ms.fromEmailPass, ms.fromEmailSMTP)

	log.Printf("Sending email to: %s", strings.Join(mailInput.To, ", "))

	err := smtp.SendMail(ms.smtpAddr, auth, ms.fromEmail, mailInput.To, msg.Bytes())
	if err != nil {
		return err
	}

	log.Printf("Email sent successfully to: %s", strings.Join(mailInput.To, ", "))
	return nil
}

