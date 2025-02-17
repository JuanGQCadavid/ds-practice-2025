package ports

import (
	"errors"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
)

var (
	ErrTransIsNotValid = errors.New("err trans is not valid")
)

type ITransactionVerification interface {
	CheckTransaction(*domain.Checkout) (string, error)
}
