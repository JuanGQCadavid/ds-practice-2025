package fraud

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
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

func (srv *FraudDetectionService) CheckFraud(data *domain.Checkout) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), srv.defaultTimeOut)
	defer cancel()

	resp, err := srv.client.CheckFraud(ctx, &pb.FraudDetectionRequest{
		CreditCard: &pb.CreditCard{
			Number:         data.CreditCard.Number,
			Cvv:            data.CreditCard.Cvv,
			ExpirationDate: data.CreditCard.ExpirationDate,
		},
	})

	if err != nil {
		log.Println("Ups! we got an error on fraud: ", err.Error())
		return "", err
	}

	return resp.Code, nil
}
