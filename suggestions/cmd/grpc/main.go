package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/JuanGQCadavid/ds-practice-2025/suggestions/internal/core"
	"github.com/JuanGQCadavid/ds-practice-2025/suggestions/internal/repositories/gemeni"
	common "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/common"
	pb "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/suggestions"
	"google.golang.org/grpc"
)

var (
	listener net.Listener
	coreSrv  *core.SuggestionSrv
)

const (
	PORT_NUMBER_ENV_NAME = "port_to_listening"
	GEMINI_API_ENV_NAME  = "gemini_api_key"
	DEFAULT_PORT         = "50053"
	PROTOCOL             = "tcp"
	SUGGESTIONS_SIZE     = 5
)

type Server struct {
	pb.UnimplementedBookSuggestionsServiceServer
	coreService *core.SuggestionSrv
}

func (srv *Server) SuggestBooks(ctx context.Context, rq *common.NextRequest) (*pb.BookSuggest, error) {
	log.Println("SUGGEST, Id", rq.OrderId) // TODO: print order ID and vector clock
	return srv.coreService.BooksSuggestions(rq), nil
}

func (srv *Server) InitOrder(ctx context.Context, rq *common.InitRequest) (*common.InitResponse, error) {
	log.Println("Received order ID", rq.OrderId)
	return srv.coreService.Init(rq), nil
}

func (srv *Server) CleanOrder(ctx context.Context, rq *common.NextRequest) (*common.NextResponse, error) {
	log.Println("Received order ID", rq.OrderId)
	return srv.coreService.Clean(rq), nil
}

func init() {

	portTo, ok := os.LookupEnv(PORT_NUMBER_ENV_NAME)

	if !ok {
		portTo = DEFAULT_PORT
	}

	gemAPIKey, ok := os.LookupEnv(GEMINI_API_ENV_NAME)

	if !ok {
		log.Fatal("Missing gem key env")

	}

	gem := gemeni.NewGemeniAI(gemAPIKey)
	coreSrv = core.NewSuggestionSrv(gem, SUGGESTIONS_SIZE)

	var err error
	listener, err = net.Listen(PROTOCOL, fmt.Sprintf(":%s", portTo))

	if err != nil {
		log.Panic("Unable to start listener in port ", portTo, PROTOCOL, " err: ", err.Error())
	}

}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterBookSuggestionsServiceServer(grpcServer, &Server{
		coreService: coreSrv,
	})

	log.Println("Suggestions server start process")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
