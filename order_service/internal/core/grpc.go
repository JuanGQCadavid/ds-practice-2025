package core

import (
	"context"
	"log"

	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/core/domain"
	"github.com/JuanGQCadavid/ds-practice-2025/order_service/pb"
)

type Server struct {
	pb.UnimplementedConsensusServer

	messages         chan domain.MessageFromThePlebe
	democracyUpdates <-chan domain.StatesOfDemocracy

	latestStateOfDemocracy domain.StatesOfDemocracy
}

func NewServer(messages chan domain.MessageFromThePlebe, democracyUpdates <-chan domain.StatesOfDemocracy, latestStateOfDemocracy domain.StatesOfDemocracy) *Server {
	return &Server{
		messages:               messages,
		democracyUpdates:       democracyUpdates,
		latestStateOfDemocracy: latestStateOfDemocracy,
	}
}

func (srv *Server) InitServer() *Server {
	go func() {
		for msg := range srv.democracyUpdates {
			log.Println("New state saved on server", msg)
			srv.latestStateOfDemocracy = msg
		}
		log.Println("Server stop listening democracy updates")
	}()

	return srv
}

func (srv *Server) CoupDeAaaah(ctx context.Context, rq *pb.Empty) (*pb.CoupDEtatResponse, error) {
	log.Println("Oh shit! This is getting bananas!")

	srv.messages <- domain.MessageFromThePlebe{
		MesssageType: domain.CoupDeAaaahType,
	}

	if srv.latestStateOfDemocracy == domain.Submission {
		return &pb.CoupDEtatResponse{
			Oks: true,
		}, nil
	}

	return &pb.CoupDEtatResponse{
		Oks: false,
	}, nil

}

func (srv *Server) YeahImStillAliveBitch(ctx context.Context, rq *pb.Empty) (*pb.Empty, error) {
	srv.messages <- domain.MessageFromThePlebe{
		MesssageType: domain.TheSonOfBitchIsStillAlive,
		Term:         int(rq.Term),
	}

	return &pb.Empty{}, nil
}
