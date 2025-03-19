package transcheck

import (
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/repositories/utils"
	pb "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/transaction_verification"
)

type TransactionVerification struct {
	client         pb.TransactionVerificationServiceClient
	defaultTimeOut time.Duration
}

func NewTransactionVerification(target string, defaultTimeOut time.Duration) *TransactionVerification {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	return &TransactionVerification{
		client:         pb.NewTransactionVerificationServiceClient(conn),
		defaultTimeOut: defaultTimeOut,
	}
}

func (srv *TransactionVerification) Init(orderId string, data *domain.Checkout) error {
	return utils.Init(orderId, srv.defaultTimeOut, data, srv.client.InitOrder)
}

func (srv *TransactionVerification) CheckOrder(orderId string, clock []int32) ([]int32, error) {
	return utils.SimpleCall(orderId, clock, srv.defaultTimeOut, srv.client.CheckOrder)
}

func (srv *TransactionVerification) CheckUser(orderId string, clock []int32) ([]int32, error) {
	return utils.SimpleCall(orderId, clock, srv.defaultTimeOut, srv.client.CheckUser)
}

func (srv *TransactionVerification) CheckFormatCreditCard(orderId string, clock []int32) ([]int32, error) {
	return utils.SimpleCall(orderId, clock, srv.defaultTimeOut, srv.client.CheckFormatCreditCard)
}

func (srv *TransactionVerification) CleanOrder(orderId string, clock []int32) ([]int32, error) {
	return utils.SimpleCall(orderId, clock, srv.defaultTimeOut, srv.client.CleanOrder)
}
