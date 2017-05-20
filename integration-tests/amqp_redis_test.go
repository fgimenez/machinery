package integrationtests

import (
	"fmt"
	"os"
	"testing"

	"github.com/RichardKnop/machinery/v1/config"
)

func TestAmqpRedis(t *testing.T) {
	amqpURL := os.Getenv("AMQP_URL")
	redisURL := os.Getenv("REDIS_URL")
	if amqpURL == "" || redisURL == "" {
		return
	}

	// AMQP broker, Redis result backend
	server := setup(&config.Config{
		Broker:        amqpURL,
		DefaultQueue:  "test_queue",
		ResultBackend: fmt.Sprintf("redis://%v", redisURL),
		AMQP: &config.AMQPConfig{
			Exchange:      "test_exchange",
			ExchangeType:  "direct",
			BindingKey:    "test_task",
			PrefetchCount: 1,
		},
	})
	worker := server.NewWorker("test_worker")
	go worker.Launch()
	testSendTask(server, t)
	testSendGroup(server, t)
	testSendChord(server, t)
	testSendChain(server, t)
	testReturnJustError(server, t)
	testReturnMultipleValues(server, t)
	testPanic(server, t)
	worker.Quit()
}
