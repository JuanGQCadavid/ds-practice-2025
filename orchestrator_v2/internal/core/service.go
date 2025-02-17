package core

import (
	"sync"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/ports"
)

type Service struct {
	fraudDetection     ports.IFraudDetection
	transactionChecker ports.ITransactionVerification
}

func NewService(fraudDetection ports.IFraudDetection, transactionChecker ports.ITransactionVerification) *Service {
	return &Service{
		fraudDetection:     fraudDetection,
		transactionChecker: transactionChecker,
	}
}

func (srv *Service) Checkout(checkout *domain.Checkout) (*domain.CheckoutResponse, error) {

	var (
		wg                  sync.WaitGroup = sync.WaitGroup{}
		fradDetectionStatus string
		fraudError          error

		transVereficationStatus string
		transError              error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		fradDetectionStatus, fraudError = srv.fraudDetection.CheckFraud(checkout)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		transVereficationStatus, transError = srv.transactionChecker.CheckTransaction(checkout)
	}()

	wg.Wait()

	if transError != nil || fraudError != nil {
		return nil, ports.ErrInternalError
	}

	if transVereficationStatus != "200" {
		return nil, ports.ErrTransIsNotValid
	}

	if fradDetectionStatus != "200" {
		return nil, ports.ErrFraudDetected
	}

	return &domain.CheckoutResponse{
		Status:  "Okey",
		OrderId: "Pending",
		SuggestedBooks: []domain.SuggestedBook{
			{
				BookId: "Pending",
				Title:  "Pending",
				Author: "Pending",
			},
		},
	}, nil
}
