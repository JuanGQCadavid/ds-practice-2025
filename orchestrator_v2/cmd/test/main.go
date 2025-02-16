package main

import (
	"context"
	"log"
	"time"

	pb "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/fraud_detection"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	target = "localhost:50051"
)

func main() {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Erro while creating conn: ", err.Error())
	}

	defer conn.Close()

	c := pb.NewFraudDetectionServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	r, err := c.CheckFraud(ctx, &pb.FraudDetectionRequest{
		Json: `{
			"creditCard": {
				"number": "4111111111111111",
				"cvv": "123",
				"expirationDate": "12/26"
			}
		}`,
	})

	if err != nil {
		log.Panic("Erro while calling: ", err.Error())
	}
	log.Println(r)

	log.Println("Response: ", r.Code)

}
