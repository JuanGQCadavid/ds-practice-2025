package transcheck

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/ports"
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

func (srv *TransactionVerification) CheckTransaction(data *domain.Checkout) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), srv.defaultTimeOut)
	defer cancel()

	jsonData, err := json.MarshalIndent(data, "", "	")

	if err != nil {
		log.Println("err while converting struct to JSON, struct: ", *data, " err: ", err.Error())
		return "", ports.ErrMarshaling
	}

	resp, err := srv.client.CheckTransaction(ctx, &pb.TransactionVerificationRequest{
		Json: string(jsonData),
	})

	if err != nil {
		log.Println("Ups! we got an error on fraud: ", err.Error())
		return "", err
	}

	if !resp.IsValid {
		return resp.ErrMessage, ports.ErrTransIsNotValid
	}

	return "", nil
}
