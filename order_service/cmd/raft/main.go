package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/JuanGQCadavid/ds-practice-2025/order_service/pb"
	"google.golang.org/grpc"
)

type MessageFromThePlebeTypes string

const (
	CoupDeAaaahType           MessageFromThePlebeTypes = "Coup d'état!!"
	TheSonOfBitchIsStillAlive MessageFromThePlebeTypes = "Damn it! Not yet"
)

type MessageFromThePlebe struct {
	MesssageType MessageFromThePlebeTypes
}

var (
	// gRPC
	listener             net.Listener
	messagesFromThePlebe chan MessageFromThePlebe

	// Raft
	heartBeat       = time.Duration(5 * time.Second)
	electionTimeOut = time.Duration(1 * time.Second)
)

const (
	PROTOCOL     = "tcp"
	DEFAULT_PORT = "50052"
)

type Server struct {
	pb.UnimplementedConsensusServer
	messages chan MessageFromThePlebe
}

func (srv *Server) CoupDeAaaah(ctx context.Context, rq *pb.Empty) (*pb.CoupDEtatResponse, error) {
	log.Println("Oh shit! This is getting bananas!")

	srv.messages <- MessageFromThePlebe{
		MesssageType: CoupDeAaaahType,
	}

	return &pb.CoupDEtatResponse{
		Oks: true,
	}, nil
}

func (srv *Server) YeahImStillAliveBitch(ctx context.Context, rq *pb.Empty) (*pb.Empty, error) {
	log.Println("Yet...")

	srv.messages <- MessageFromThePlebe{
		MesssageType: TheSonOfBitchIsStillAlive,
	}

	return &pb.Empty{}, nil
}

func init() {
	var err error
	listener, err = net.Listen(PROTOCOL, fmt.Sprintf(":%s", DEFAULT_PORT))

	if err != nil {
		log.Panic("Unable to start listener in port ", DEFAULT_PORT, PROTOCOL, " err: ", err.Error())
	}
	messagesFromThePlebe = make(chan MessageFromThePlebe, 100)
}

func startServer(thePlebeSpoke chan MessageFromThePlebe) {

	grpcServer := grpc.NewServer()
	pb.RegisterConsensusServer(grpcServer, &Server{
		messages: thePlebeSpoke,
	})

	log.Println("Listening to all the madafuckers")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {

	go func() {
		startServer(messagesFromThePlebe)
	}()

	ticker := time.NewTicker(heartBeat)

	for {
		select {
		case msg := <-messagesFromThePlebe:
			if msg.MesssageType == TheSonOfBitchIsStillAlive {
				log.Println("Damn it! He is alive...")
				ticker.Reset(heartBeat)
			}

		case time := <-ticker.C:
			log.Println("Heartbit!!", time)
			log.Println("Attack!!!")
		}
	}

}
