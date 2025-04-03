package core

import (
	"log"
	"time"

	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/core/domain"
	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/repositories/queue"
)

type Service struct {
	repository       *queue.QueueRepository
	democracyUpdates <-chan domain.StatesOfDemocracy
	pull             bool
}

func NewService(repository *queue.QueueRepository, democracyUpdates <-chan domain.StatesOfDemocracy) *Service {
	return &Service{
		repository:       repository,
		pull:             false,
		democracyUpdates: democracyUpdates,
	}
}

func (srv *Service) Init() {
	log.Println("Puller: listening to democracy updates")
	go func() {
		for msg := range srv.democracyUpdates {
			log.Println("New state saved on the puller: ", msg)

			if msg == domain.Reich {
				srv.pull = true
				log.Println("Oh shit, time to work....")
				continue
			}

			srv.pull = false
		}
		log.Println("Puller stop listening the democracy updates")
	}()
}

func (srv *Service) Listen() {
	log.Println("Listening queue.")

	var (
		consecutivesEmpty uint16 = 0
	)

	for true {

		if !srv.pull {
			continue
		}

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
