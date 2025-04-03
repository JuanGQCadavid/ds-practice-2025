package orderqueue

import (
	"context"
	"log"
	"time"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
	pb "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/order_queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderQueue struct {
	client         pb.OrderQueueServiceClient
	defaultTimeOut time.Duration
}

func NewOrderQueue(target string, defaultTimeOut time.Duration) *OrderQueue {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	return &OrderQueue{
		client:         pb.NewOrderQueueServiceClient(conn),
		defaultTimeOut: defaultTimeOut,
	}
}

func (ord *OrderQueue) Enqueue(orderId string, data *domain.Checkout) error {
	// TODO Maybe we could have a queue of pending orders that fail.
	ctx, cancel := context.WithTimeout(context.Background(), ord.defaultTimeOut)
	defer cancel()

	_, err := ord.client.Enqueue(ctx, &pb.EnqueueRequest{
		OrderId: orderId,
		Order:   domain.FromCheckoutToCommon(data),
	})

	if err != nil {
		log.Printf("Ups! we got an error on sending the order %s to the queue, err %s: \n", orderId, err.Error())
		return err
	}

	log.Printf("Order id %s is now on the queue \n", orderId)
	return nil
}
