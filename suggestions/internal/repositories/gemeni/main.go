package gemeni

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

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
		log.Fatal("Error creating the geminar client, error: ", err.Error())
	}

	return &GemeniAI{
		model: client.GenerativeModel("gemini-1.5-flash"),
	}
}

func (svc *GemeniAI) SuggestBooks(books []string, size int) []*domain.Book {
	ctx := context.Background()
	log.Println(randomBooksQuery)
	resp, err := svc.model.GenerateContent(ctx, genai.Text(randomBooksQuery))
	if err != nil {
		log.Println("Error on calling Geminar AI, ", err.Error())
	}

	return svc.castBooks(resp)

}

func (svc *GemeniAI) castBooks(resp *genai.GenerateContentResponse) []*domain.Book {
	log.Println("Geminar -> Candidates - ", len(resp.Candidates))
	for _, cand := range resp.Candidates {
		msg := ""
		if cand.Content != nil {
			log.Println("Geminar -> Candidates - Content - Parts", len(cand.Content.Parts))
			for _, part := range cand.Content.Parts {
				msg += fmt.Sprintf("%s", part)
			}
		}
		log.Println(msg)
		var books []*domain.Book = make([]*domain.Book, 0)
		if err := json.Unmarshal([]byte(msg), &books); err != nil {
			log.Println("Error on casting Geminar AI, ", err.Error())
			return nil
		}
		return books
	}

	return nil
}
