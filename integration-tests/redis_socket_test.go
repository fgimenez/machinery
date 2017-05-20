package integrationtests

import (
	"fmt"
	"os"
	"testing"

	"github.com/RichardKnop/machinery/v1/config"
)

func TestRedisSocket(t *testing.T) {
	redisSocket := os.Getenv("REDIS_SOCKET")
	if redisSocket == "" {
		return
	}

	// Redis broker, Redis result backend
	server := setup(&config.Config{
		Broker:        fmt.Sprintf("redis+socket://%v", redisSocket),
		DefaultQueue:  "test_queue",
		ResultBackend: fmt.Sprintf("redis+socket://%v", redisSocket),
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
