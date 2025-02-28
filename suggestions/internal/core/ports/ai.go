package ports

import "github.com/JuanGQCadavid/ds-practice-2025/suggestions/internal/core/domain"

type ISuggestionsAI interface {
	SuggestBooks(books []string, size int) []*domain.Book
}
