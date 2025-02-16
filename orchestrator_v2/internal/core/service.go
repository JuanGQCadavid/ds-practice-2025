package core

import (
	"log"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/ports"
)

type Service struct {
	fraudDetection ports.IFraudDetection
}

func NewService(fraudDetection ports.IFraudDetection) *Service {
	return &Service{
		fraudDetection: fraudDetection,
	}
}
func (srv *Service) Checkout(checkout *domain.Checkout) (*domain.CheckoutResponse, error) {
	fraduCode, err := srv.fraudDetection.CheckFraud(checkout)

	if err != nil {
		return nil, err
	}

	log.Println(fraduCode)

	return &domain.CheckoutResponse{
		Status: fraduCode,
	}, nil
}
