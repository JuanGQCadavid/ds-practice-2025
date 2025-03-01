package suggestions

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
	pb "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/suggestions"
)

type SuggestionService struct {
	client         pb.BookSuggestionsServiceClient
	defaultTimeOut time.Duration
}

func NewSuggestionService(target string, defaultTimeOut time.Duration) *SuggestionService {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	return &SuggestionService{
		client:         pb.NewBookSuggestionsServiceClient(conn),
		defaultTimeOut: defaultTimeOut,
	}
}

func (srv *SuggestionService) SuggestBooks(data []domain.Item) ([]*domain.SuggestedBook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), srv.defaultTimeOut)
	defer cancel()

	request := &pb.ItemsBought{
		Items: make([]*pb.ItemsBought_Item, len(data)),
	}

	for i, val := range data {
		request.Items[i] = &pb.ItemsBought_Item{
			Name:     val.Name,
			Quantity: string(val.Quantity), // TODO Set as int in proto
		}
	}

	resp, err := srv.client.SuggestBooks(ctx, request)

	if err != nil {
		log.Println("Ups! we got an error on suggestions: ", err.Error())
		return nil, err
	}

	result := make([]*domain.SuggestedBook, len(resp.Books))

	for i, b := range resp.Books {
		result[i] = &domain.SuggestedBook{
			BookId: b.BookId,
			Author: b.Author,
			Title:  b.Title,
		}
	}

	return result, nil
}
