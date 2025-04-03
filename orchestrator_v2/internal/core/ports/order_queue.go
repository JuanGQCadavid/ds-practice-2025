package ports

import "github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"

type IOrderQueue interface {
	Enqueue(string, *domain.Checkout) error
}
