package fraud

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/ports"
	pb "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/fraud_detection"
)

type FraudDetectionService struct {
	client         pb.FraudDetectionServiceClient
	defaultTimeOut time.Duration
}

func NewFraudDetectionService(target string, defaultTimeOut time.Duration) *FraudDetectionService {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Erro while creating conn: ", err.Error())
	}

	return &FraudDetectionService{
		client:         pb.NewFraudDetectionServiceClient(conn),
		defaultTimeOut: defaultTimeOut,
	}
}

func (srv *FraudDetectionService) CheckFraud(data *domain.Checkout) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), srv.defaultTimeOut)
	defer cancel()

	// TODO - This should not be a JSON!
	jsonData, err := json.MarshalIndent(data, "", "	")

	if err != nil {
		log.Println("err while converting struct to JSON, struct: ", *data, " err: ", err.Error())
		return "", ports.ErrMarshaling
	}

	resp, err := srv.client.CheckFraud(ctx, &pb.FraudDetectionRequest{
		Json: string(jsonData),
	})

	if err != nil {
		log.Println("Ups! we got an error on fraud: ", err.Error())
		return "", err
	}

	return resp.Code, nil
}
