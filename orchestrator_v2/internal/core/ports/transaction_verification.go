package ports

import (
	"errors"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
)

var (
	ErrTransIsNotValid = errors.New("err trans is not valid")
	ErrUserNotValid    = errors.New("err user check fail")
)

type ITransactionVerification interface {
	Init(orderId string, data *domain.Checkout) error
	CheckOrder(orderId string, clock []int32) ([]int32, error)
	CheckUser(orderId string, clock []int32) ([]int32, error)
	CheckFormatCreditCard(orderId string, clock []int32) ([]int32, error)
	CleanOrder(orderId string, clock []int32) ([]int32, error)
}
