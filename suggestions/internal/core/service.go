package core

import (
	"fmt"
	"github.com/JuanGQCadavid/ds-practice-2025/suggestions/internal/core/ports"
	pb "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/suggestions"
)

type SuggestionSrv struct {
	aiRepository    ports.ISuggestionsAI
	suggestionsSize int
}

func NewSuggestionSrv(aiRepository ports.ISuggestionsAI, suggestionsSize int) *SuggestionSrv {
	return &SuggestionSrv{
		aiRepository:    aiRepository,
		suggestionsSize: suggestionsSize,
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

func (svc *SuggestionSrv) BooksSuggestions(items *pb.ItemsBought) *pb.BookSuggest {
	var (
		books = make([]string, len(items.Items))
	)

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
