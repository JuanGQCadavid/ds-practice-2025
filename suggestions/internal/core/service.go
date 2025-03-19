package core

import (
	"fmt"
	"log"

	"github.com/JuanGQCadavid/ds-practice-2025/suggestions/internal/core/ports"
	common "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/common"
	pb "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/suggestions"
)

type SuggestionSrv struct {
	aiRepository    ports.ISuggestionsAI
	suggestionsSize int
	state           map[string]*common.Order
	svcId           int
}

func NewSuggestionSrv(aiRepository ports.ISuggestionsAI, suggestionsSize int) *SuggestionSrv {
	return &SuggestionSrv{
		aiRepository:    aiRepository,
		suggestionsSize: suggestionsSize,
		svcId:           2,
		state:           make(map[string]*common.Order),
	}
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

func (svc *SuggestionSrv) BooksSuggestions(next *common.NextRequest) *pb.BookSuggest {
	if svc.state[next.OrderId] == nil {
		log.Println("Error, the order id ", next.OrderId, " is not on the state.")
		return defaultResponse
	}

	var (
		items = svc.state[next.OrderId]
		books = make([]string, len(items.Items))
	)

	for i := range books {
		books[i] = items.Items[i].Name
	}

	if booksToSuggest := svc.aiRepository.SuggestBooks(books, svc.suggestionsSize); booksToSuggest != nil {
		var (
			resp = &pb.BookSuggest{
				Books: make([]*pb.BookSuggest_Book, len(booksToSuggest)),
			}
		)

		for i, b := range booksToSuggest {
			resp.Books[i] = &pb.BookSuggest_Book{
				BookId: fmt.Sprintf("%d", i),
				Title:  b.Title,
				Author: b.Author,
			}
		}

		return resp
	}

	return defaultResponse

}

func (svc *SuggestionSrv) Init(request *common.InitRequest) *common.InitResponse {
	svc.state[request.OrderId] = request.Order
	return &common.InitResponse{
		IsValid: true,
	}
}
