package core

import (
	"log"
	"time"

	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/repositories/queue"
)

type Service struct {
	repository *queue.QueueRepository
}

func NewService(repository *queue.QueueRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (srv *Service) Listen() {
	log.Println("Listening queue.")

	var (
		consecutivesEmpty uint16 = 0
	)

	for true {
		pullMessage := srv.repository.Pull()

		if pullMessage == nil {
			consecutivesEmpty++

			if consecutivesEmpty < 10 {
				time.Sleep(1 * time.Second)
				continue
			} else if consecutivesEmpty < 20 {
				time.Sleep(2 * time.Second)
				continue
			} else if consecutivesEmpty < 100 {
				time.Sleep(3 * time.Second)
				continue
			}
			consecutivesEmpty = 0
			continue
		}

		consecutivesEmpty = 0
		log.Printf("Order %s recieved.\n", pullMessage.OrderId)
	}

}
