package main

import (
	"log"
	"os"

	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/core"
	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/repositories/db"
	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/repositories/payment"
	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/repositories/queue"
)

const (
	QUEUE_ENV_NAME   = "order_queue_dns"
	PAYMENT_ENV_NAME = "payment_dns"
)

var (
	repo        *queue.QueueRepository
	dbRepo      *db.DBService
	paymentRepo *payment.PaymentService
)

func init() {
	dns, ok := os.LookupEnv(QUEUE_ENV_NAME)

	if !ok {
		log.Fatalln("Missing QUEUE dns")
	}

	repo = queue.NewQueueRepository(dns)

	paymentDns, ok := os.LookupEnv(PAYMENT_ENV_NAME)

	if !ok {
		log.Fatalln("Missing PAYMENT dns")
	}

	paymentRepo = payment.NewPaymentService(paymentDns)

	dbRepo = db.NewDBServiceWithTargets(1, map[int64]string{
		1: "localhost:50061",
		2: "localhost:50062",
		3: "localhost:50063",
	})

}
func main() {
	var (
		service = core.NewService(repo, dbRepo, paymentRepo, nil)
	)
	service.Shhhhhhhhhhh()
	service.Listen()
}
