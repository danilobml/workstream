package readiness

import (
	"errors"
	"os"
)

func IsReady() error {
	rabbitmqUrl := os.Getenv("RABBITMQ_URL")

	if rabbitmqUrl == "" {
		return errors.New("RabbitMQ can't be initialized: RABBITMQ_URL variable could not be read from environment.")
	}

	return nil
}
