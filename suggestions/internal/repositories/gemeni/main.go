package gemeni

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/JuanGQCadavid/ds-practice-2025/suggestions/internal/core/domain"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	//go:embed queries/random.txt
	randomBooksQuery string
)

type GemeniAI struct {
	model *genai.GenerativeModel
}

func NewGemeniAI(apiKey string) *GemeniAI {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		log.Fatal("Error creating the Gemini AI client, error: ", err.Error())
	}

	return &GemeniAI{
		model: client.GenerativeModel("gemini-1.5-flash"),
	}
}

func (svc *GemeniAI) SuggestBooks(books []string, size int) []*domain.Book {
	ctx := context.Background()

	customQuery := fmt.Sprintf(randomBooksQuery, books, size, size)
	resp, err := svc.model.GenerateContent(ctx, genai.Text(customQuery))
	if err != nil {
		log.Println("Error on calling Gemini AI, ", err.Error())
	}

	return svc.castBooks(resp)
}

func (svc *GemeniAI) castBooks(resp *genai.GenerateContentResponse) []*domain.Book {
	for _, cand := range resp.Candidates {
		msg := ""
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				msg += fmt.Sprintf("%s", part)
			}
		}
		var books []*domain.Book = make([]*domain.Book, 0)
		msg = strings.ReplaceAll(msg, "```json", "")
		msg = strings.ReplaceAll(msg, "```", "")

		if err := json.Unmarshal([]byte(msg), &books); err != nil {
			log.Println("Error on casting Gemini AI, ", err.Error())
			log.Println("Using static recommendations")
			return nil
		}

		return books
	}

	return nil
}
