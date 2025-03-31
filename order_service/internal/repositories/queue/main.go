package queue

import (
	"context"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	queueProtoBus "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/order_queue"
)

type QueueRepository struct {
	client     queueProtoBus.OrderQueueServiceClient
	timeout    time.Duration
	emptyQuery *queueProtoBus.EmptyRequest
}

func NewQueueRepository(target string) *QueueRepository {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}
	return &QueueRepository{
		client:     queueProtoBus.NewOrderQueueServiceClient(conn),
		timeout:    4 * time.Second,
		emptyQuery: &queueProtoBus.EmptyRequest{},
	}
}

func (repo *QueueRepository) Pull() *queueProtoBus.DequeueResponse {
	ctx, cancel := context.WithTimeout(context.Background(), repo.timeout)
	defer cancel()

	r, err := repo.client.Dequeue(ctx, repo.emptyQuery)

	if err != nil {
		log.Println("Error while calling: ", err.Error())
		return nil
	}

	if strings.ContainsAny(r.ErrMessage, "empty") {
		// log.Println("No data from queue")
		return nil
	}

	return r
}
