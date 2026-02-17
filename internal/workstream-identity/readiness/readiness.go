package readiness

import (
	"errors"
	"os"
)

func IsReady() error {
	postgresDsn := os.Getenv("POSTGRES_DSN")
	rabbitmqUrl := os.Getenv("RABBITMQ_URL")

	if postgresDsn == "" {
		return errors.New("Postgres DB can't be initialized: POSTGRES_DSN variable could not be read from environment.")
	}

	if rabbitmqUrl == "" {
		return errors.New("RabbitMQ can't be initialized: RABBITMQ_URL variable could not be read from environment.")
	}

	return nil
}
