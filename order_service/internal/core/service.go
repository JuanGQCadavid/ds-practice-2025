package core

import (
	"log"
	"sync"
	"time"

	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/core/domain"
	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/repositories/db"
	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/repositories/payment"
	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/repositories/queue"
	"github.com/JuanGQCadavid/ds-practice-2025/utils/pb/common"
)

type Service struct {
	repository     *queue.QueueRepository
	dbRepositoy    *db.DBService
	paymentService *payment.PaymentService

	democracyUpdates <-chan domain.StatesOfDemocracy
	pull             bool
}

func (srv *Service) Shhhhhhhhhhh() {
	srv.pull = true
}

func NewService(
	repository *queue.QueueRepository,
	dbRepositoy *db.DBService,
	paymentService *payment.PaymentService,
	democracyUpdates <-chan domain.StatesOfDemocracy,
) *Service {
	return &Service{
		repository:       repository,
		dbRepositoy:      dbRepositoy,
		paymentService:   paymentService,
		pull:             false,
		democracyUpdates: democracyUpdates,
	}
}

func (srv *Service) broadcast(orderID string, order *common.Order, functions ...func(string, *common.Order) error) bool {
	var (
		noVotes   = make([]error, 0)
		errosLock = sync.Mutex{}
		wg        sync.WaitGroup
	)

	for _, f := range functions {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := f(orderID, order); err != nil {
				log.Println("Transaction ", orderID, " received a NO")
				errosLock.Lock()
				noVotes = append(noVotes, err)
				errosLock.Unlock()
			} else {
				log.Println("Transaction ", orderID, " received a YES")
			}
		}()
	}

	wg.Wait()

	if len(noVotes) > 0 {
		log.Println("There at least one NO vote, transaction not approved.")
		return false
	}

	log.Println("All vote YES!")
	return true

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
		log.Printf("Order %s recieved.\n", pullMessage.OrderId)

		log.Println("----------------")
		log.Println("")
		log.Println("Prepare stage: ", pullMessage.OrderId)

		if !srv.broadcast(
			pullMessage.OrderId,
			pullMessage.Order,
			func(s string, o *common.Order) error {
				return srv.paymentService.Prepare(pullMessage.OrderId, pullMessage.Order.CreditCard)
			},
			func(s string, o *common.Order) error {
				return srv.dbRepositoy.Prepare(pullMessage.OrderId, pullMessage.Order.Items)
			},
		) {
			log.Println("Prepare for commit ", pullMessage.OrderId, " fail,  starting transaction abort.")
			srv.broadcast(
				pullMessage.OrderId,
				pullMessage.Order,
				func(s string, o *common.Order) error {
					return srv.paymentService.Abort(pullMessage.OrderId)
				},
				func(s string, o *common.Order) error {
					return srv.dbRepositoy.Abort(pullMessage.OrderId)
				},
			)
			continue
		}

		log.Println("----------------")
		log.Println("")
		log.Println("Commit stage: ", pullMessage.OrderId)

		if !srv.broadcast(
			pullMessage.OrderId,
			pullMessage.Order,
			func(s string, o *common.Order) error {
				return srv.paymentService.Commit(pullMessage.OrderId)
			},
			func(s string, o *common.Order) error {
				return srv.dbRepositoy.Commit(pullMessage.OrderId)
			},
		) {
			log.Println("commit for id ", pullMessage.OrderId, " fail,  starting transaction abort.")
			srv.broadcast(
				pullMessage.OrderId,
				pullMessage.Order,
				func(s string, o *common.Order) error {
					return srv.paymentService.Abort(pullMessage.OrderId)
				},
				func(s string, o *common.Order) error {
					return srv.dbRepositoy.Abort(pullMessage.OrderId)
				},
			)
			continue
		}
		log.Println("----------------")
		log.Println("Alles ist gut my dear transaction ", pullMessage.OrderId)
		log.Println("It was a pleasure to have you here")
		log.Println("")
		log.Println("And they lived happily ever after")

		consecutivesEmpty = 0
	}

}
