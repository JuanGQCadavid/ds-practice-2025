package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/suggestions"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedBookSuggestionsServiceServer
}

var (
	defaultResponse *pb.BookSuggest = &pb.BookSuggest{
		Books: []*pb.BookSuggest_Book{
			{
				BookId: "1",
				Author: "George Orwell",
				Title:  "1984",
			},
			{
				BookId: "2",
				Author: "F. Scott Fitzgerald",
				Title:  "The Great Gatsby",
			},
			{
				BookId: "3",
				Author: "Harper Lee",
				Title:  "To Kill a Mockingbird",
			},
			{
				BookId: "4",
				Author: "J.R.R. Tolkien",
				Title:  "The Hobbit",
			},
			{
				BookId: "5",
				Author: "Mary Shelley",
				Title:  "Frankenstein",
			},
		},
	}
)

func (srv *Server) SuggestBooks(ctx context.Context, rq *pb.ItemsBought) (*pb.BookSuggest, error) {
	log.Println("Request, len of items", len(rq.Items))

	for i, item := range rq.Items {
		log.Println(i, item.Name, item.Quantity)
	}

	return defaultResponse, nil
}

var (
	listener net.Listener
)

const (
	PORT_NUMBER_ENV_NAME = "port_to_listening"
	DEFAULT_PORT         = "50050"
	PROTOCOL             = "tcp"
)

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
	pb.RegisterBookSuggestionsServiceServer(grpcServer, &Server{})

	log.Println("Suggestions server start process")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
