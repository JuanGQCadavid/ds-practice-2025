package main

import (
	"log"

	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/repositories/queue"
)

const (
	QUEUE_TARGET = "localhost:50054"
)

var (
	repo *queue.QueueRepository
)

func init() {
	repo = queue.NewQueueRepository(QUEUE_TARGET)
}
func main() {
	pullMessage := repo.Pull()

	if pullMessage == nil {
		log.Println("Empty")
		return
	}

	log.Printf("Order %s recieved.\n", pullMessage.OrderId)
}
