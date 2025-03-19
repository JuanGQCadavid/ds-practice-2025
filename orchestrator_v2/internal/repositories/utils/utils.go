package utils

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
	"github.com/JuanGQCadavid/ds-practice-2025/utils/pb/common"
	"google.golang.org/grpc"
)

func Init(
	orderId string,
	defaultTimeOut time.Duration,
	data *domain.Checkout,
	initFunction func(ctx context.Context, in *common.InitRequest, opts ...grpc.CallOption) (*common.InitResponse, error),
) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOut)
	defer cancel()

	_, err := initFunction(ctx, &common.InitRequest{
		OrderId: orderId,
		Order:   domain.FromCheckoutToCommon(data),
	})

	if err != nil {
		log.Println("Ups! we got an error on init Trans: ", err.Error())
		return err
	}
	return nil
}

func SimpleCall(
	orderId string,
	vectorClock []int32,
	defaultTimeOut time.Duration,
	f func(ctx context.Context, in *common.NextRequest, opts ...grpc.CallOption) (*common.NextResponse, error),
) ([]int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOut)
	defer cancel()

	resp, err := f(ctx, &common.NextRequest{
		OrderId:             orderId,
		IncomingVectorClock: vectorClock,
	})

	if err != nil {
		log.Println("Ups! we got an error on grpc call:", err.Error())
		return nil, err
	}

	if !resp.IsValid {
		log.Println("Err grpc:", resp.ErrMessage)
		return nil, errors.New(resp.ErrMessage)
	}

	return resp.VectorClock, nil
}
