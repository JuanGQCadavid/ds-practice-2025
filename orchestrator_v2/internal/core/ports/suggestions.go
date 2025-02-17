package ports

import (
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
)

type ISuggestionsService interface {
	SuggestBooks(data []domain.Item) ([]*domain.SuggestedBook, error)
}
