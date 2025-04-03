package main

import (
	"log"
	"os"

	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/core"
	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/repositories/queue"
)

const (
	QUEUE_ENV_NAME = "order_queue_dns"
)

var (
	repo *queue.QueueRepository
)

func init() {
	dns, ok := os.LookupEnv(QUEUE_ENV_NAME)

	if !ok {
		log.Fatalln("Missing QUEUE dns")
	}

	repo = queue.NewQueueRepository(dns)
}
func main() {
	var (
		service = core.NewService(repo, nil)
	)
	service.Listen()
}
