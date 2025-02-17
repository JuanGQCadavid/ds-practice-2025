package ports

import (
	"errors"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
)

var (
	ErrMarshaling    = errors.New("err data is not structured to be marshall")
	ErrInternalError = errors.New("err internal")
	ErrFraudDetected = errors.New("err there is a high chance of fraud")
)

type IFraudDetection interface {
	CheckFraud(*domain.Checkout) (string, error)
}
