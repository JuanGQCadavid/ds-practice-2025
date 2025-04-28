package payment

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/JuanGQCadavid/ds-practice-2025/utils/pb/common"
	"github.com/JuanGQCadavid/ds-practice-2025/utils/pb/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentService struct {
	timeout time.Duration
	client  payment.PaymentServiceClient
}

func NewPaymentService(target string) *PaymentService {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}
	return &PaymentService{
		client:  payment.NewPaymentServiceClient(conn),
		timeout: 4 * time.Second,
	}
}

func (srv *PaymentService) Prepare(orderId string, creditCard *common.CreditCard) error {
	ctx, cancel := context.WithTimeout(context.Background(), srv.timeout)
	defer cancel()

	resp, err := srv.client.Prepare(ctx, &payment.PrepareRequest{
		OrderID:    orderId,
		CreditCard: creditCard,
	})

	if err != nil {
		log.Println("Payment - Prepare - error: ", err.Error())
		return err
	}

	if !resp.IsValid {
		return fmt.Errorf("err %s", resp.ErrMessage)
	}

	return nil
}

func (srv *PaymentService) Commit(orderId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), srv.timeout)
	defer cancel()

	resp, err := srv.client.Commit(ctx, &payment.CommitRequest{
		OrderID: orderId,
	})

	if err != nil {
		log.Println("Payment - Commit - error: ", err.Error())
		return err
	}

	if !resp.IsValid {
		return fmt.Errorf("err %s", resp.ErrMessage)
	}

	return nil
}

func (srv *PaymentService) Abort(orderId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), srv.timeout)
	defer cancel()

	resp, err := srv.client.Abort(ctx, &payment.AbortRequest{
		OrderID: orderId,
	})

	if err != nil {
		log.Println("Payment - Abort - error: ", err.Error())
		return err
	}

	if !resp.IsValid {
		return fmt.Errorf("err %s", resp.ErrMessage)
	}

	return nil
}
