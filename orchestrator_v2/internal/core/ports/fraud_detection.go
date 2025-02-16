package ports

import (
	"errors"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
)

var (
	ErrMarshaling = errors.New("err data is not structured to be marshall")
)

type IFraudDetection interface {
	CheckFraud(*domain.Checkout) (string, error)
}
