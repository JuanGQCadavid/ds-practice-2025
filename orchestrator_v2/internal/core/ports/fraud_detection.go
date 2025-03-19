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
	Init(orderId string, data *domain.Checkout) error
	CheckUser(orderId string, clock []int32) ([]int32, error)
	CheckCreditCard(orderId string, clock []int32) ([]int32, error)
	CleanOrder(orderId string, clock []int32) ([]int32, error)
}
