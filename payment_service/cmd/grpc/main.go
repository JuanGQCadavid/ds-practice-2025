package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/JuanGQCadavid/ds-practice-2025/utils/pb/common"
	pb "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/payment"
	"google.golang.org/grpc"
)

var (
	listener net.Listener
)

const (
	PORT_NUMBER_ENV_NAME = "port_to_listening"
	DEFAULT_PORT         = "50055"
	PROTOCOL             = "tcp"
)

type Server struct {
	pb.UnimplementedPaymentServiceServer
	Pending map[string]*common.CreditCard
}

func (srv *Server) Prepare(ctx context.Context, rq *pb.PrepareRequest) (*common.NextResponse, error) {
	log.Println("Prepare, Id", rq.OrderID)
	log.Printf("Credit card %+v \n", rq.CreditCard)

	if rq.CreditCard == nil {
		return &common.NextResponse{
			ErrMessage: "Not credit card!",
			IsValid:    false,
		}, nil
	}

	srv.Pending[rq.OrderID] = rq.CreditCard

	return &common.NextResponse{
		IsValid: true,
	}, nil
}

func (srv *Server) Commit(ctx context.Context, rq *pb.CommitRequest) (*common.NextResponse, error) {
	log.Println("Commit, Id", rq.OrderID)

	if srv.Pending[rq.OrderID] == nil {
		log.Println("Hmmmm, i don't know that transaction id...", rq.OrderID)
		return &common.NextResponse{
			ErrMessage: "Hmmmm, i don't know that transaction id...",
			IsValid:    false,
		}, nil
	}

	log.Println("Processing payment of Order ID: ", rq.OrderID)
	log.Println(`
		Computing the most elaborate encryption algorithm
		Requesting more CPU units from the NASA servers
		...
		...
		NASA remove the international space station compute power for us
		...
		DONE.
	`)

	delete(srv.Pending, rq.OrderID)

	return &common.NextResponse{
		IsValid: true,
	}, nil
}

func (srv *Server) Abort(ctx context.Context, rq *pb.AbortRequest) (*common.NextResponse, error) {
	log.Println("Abort, Id", rq.OrderID)
	if srv.Pending[rq.OrderID] == nil {
		log.Println("Hmmmm, i don't know that transaction id...", rq.OrderID)
		return &common.NextResponse{
			ErrMessage: "Hmmmm, i don't know that transaction id...",
			IsValid:    false,
		}, nil
	}

	log.Println("Aborting ", rq.OrderID)
	log.Println(` 
		It is not you.... It is me..... 
		and the executor btw us....
		`)

	delete(srv.Pending, rq.OrderID)

	return &common.NextResponse{
		ErrMessage: "",
		IsValid:    true,
	}, nil
}

func init() {

	portTo, ok := os.LookupEnv(PORT_NUMBER_ENV_NAME)

	if !ok {
		portTo = DEFAULT_PORT
	}

	var err error
	listener, err = net.Listen(PROTOCOL, fmt.Sprintf(":%s", portTo))

	if err != nil {
		log.Panic("Unable to start listener in port ", portTo, PROTOCOL, " err: ", err.Error())
	}

}

func main() {
	grpcServer := grpc.NewServer()

	pb.RegisterPaymentServiceServer(grpcServer, &Server{
		Pending: make(map[string]*common.CreditCard),
	})

	log.Println("Payment server start process")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
