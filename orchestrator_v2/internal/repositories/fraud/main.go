package fraud

import (
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/repositories/utils"
	pb "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/fraud_detection"
)

type FraudDetectionService struct {
	client         pb.FraudDetectionServiceClient
	defaultTimeOut time.Duration
}

func NewFraudDetectionService(target string, defaultTimeOut time.Duration) *FraudDetectionService {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	return &FraudDetectionService{
		client:         pb.NewFraudDetectionServiceClient(conn),
		defaultTimeOut: defaultTimeOut,
	}
}

func (srv *FraudDetectionService) CheckUser(orderId string, clock []int32) ([]int32, error) {
	return utils.SimpleCall(orderId, clock, srv.defaultTimeOut, srv.client.CheckUser)
}

func (srv *FraudDetectionService) CheckCreditCard(orderId string, clock []int32) ([]int32, error) {
	return utils.SimpleCall(orderId, clock, srv.defaultTimeOut, srv.client.CheckCreditCard)
}
func (srv *FraudDetectionService) CleanOrder(orderId string, clock []int32) ([]int32, error) {
	return utils.SimpleCall(orderId, clock, srv.defaultTimeOut, srv.client.CleanOrder)
}

func (srv *FraudDetectionService) Init(orderId string, data *domain.Checkout) error {
	return utils.Init(orderId, srv.defaultTimeOut, data, srv.client.InitOrder)
}
