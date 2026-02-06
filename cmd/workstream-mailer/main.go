package main

import (
	"context"
	"log"
	"os"

	http "github.com/danilobml/workstream/internal/platform/httpserver"
	"github.com/danilobml/workstream/internal/platform/rabbitmq"
	"github.com/danilobml/workstream/internal/workstream-mailer/readiness"
	"github.com/danilobml/workstream/internal/workstream-mailer/services"
)

const (
	serviceName       = "workstream-mailer"
	httpPortName      = "MAILER_HTTP_PORT"
	rabbitmqUrlName   = "RABBITMQ_URL"
	fromEmailName     = "MAILER_FROM_EMAIL"
	fromEmailPassName = "MAILER_FROM_EMAIL_PASS"
	fromEmailSMTPName = "MAILER_FROM_EMAIL_SMTP"
	smtpAddrName      = "MAILER_SMTP_ADDR"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rabbitmqUrl := os.Getenv(rabbitmqUrlName)
	if rabbitmqUrl == "" {
		log.Fatal("unable to read RABBITMQ_URL from env")
	}

	messageClient, err := rabbitmq.NewRabbitMQClient(ctx, rabbitmqUrl, rabbitmq.MailerExchange)
	if err != nil {
		log.Fatal("workstream-mailer - failed to connect to RabbitMQ", err)
	}
	defer messageClient.Close()

	if err := messageClient.DeclareQueues(rabbitmq.MailerQueue, rabbitmq.MailerQueue, rabbitmq.MailerBinding); err != nil {
		log.Fatal("workstream-mailer - failed to declare queues", err)
	}

	mailConfig := getMailConfig()
	mailerService := services.NewLocalMailService(mailConfig)

	messageConsumerService := services.NewRabbitMessageConsumerService(messageClient, mailerService)

	go func() {
		if err := messageConsumerService.Consume(ctx); err != nil {
			log.Printf("workstream-mailer - consumer stopped: %v", err)
			cancel()
		}
	}()

	if err := http.StartServer(
		serviceName,
		httpPortName,
		nil,
		readiness.IsReady,
	); err != nil {
		log.Fatal(err)
	}
}

func getMailConfig() services.LocalMailConfig {
	fromEmail := os.Getenv(fromEmailName)
	if fromEmail == "" {
		log.Fatal("unable to read MAILER_FROM_EMAIL from env")
	}

	fromEmailPass := os.Getenv(fromEmailPassName)
	if fromEmailPass == "" {
		log.Fatal("unable to read MAILER_FROM_EMAIL_PASS from env")
	}

	fromEmailSMTP := os.Getenv(fromEmailSMTPName)
	if fromEmailSMTP == "" {
		log.Fatal("unable to read MAILER_FROM_EMAIL_SMTP from env")
	}

	smtpAddr := os.Getenv(smtpAddrName)
	if smtpAddr == "" {
		log.Fatal("unable to read MAILER_SMTP_ADDR from env")
	}

	return services.LocalMailConfig{
		FromEmail:     fromEmail,
		FromEmailPass: fromEmailPass,
		FromEmailSMTP: fromEmailSMTP,
		SMTPAddr:      smtpAddr,
	}
}
