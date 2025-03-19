package ports

import (
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
)

type ISuggestionsService interface {
	SuggestBooks(orderId string, clock []int32) ([]*domain.SuggestedBook, error)
	Init(orderId string, data *domain.Checkout) error
}
