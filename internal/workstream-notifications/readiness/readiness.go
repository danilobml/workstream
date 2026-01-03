package readiness

import (
	"errors"
	"os"
)

func IsReady() error {
	mongodbUri := os.Getenv("MONGODB_URI")
	redisAddr := os.Getenv("REDIS_ADDR")
	rabbitmqUrl := os.Getenv("RABBITMQ_URL")

	if mongodbUri == "" {
		return errors.New("MongoDB can't be initialized: MONGODB_URI variable could not be read from environment.")
	}

	if redisAddr == "" {
		return errors.New("Redis can't be initialized: REDIS_ADDR variable could not be read from environment.")
	}

	if rabbitmqUrl == "" {
		return errors.New("RabbitMQ can't be initialized: RABBITMQ_URL variable could not be read from environment.")
	}

	return nil
}
