package core

import (
	"errors"
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

func (srv *Service) Checkout(checkout *domain.Checkout) (*domain.CheckoutResponse, error, error) {

	var (
		wg                  sync.WaitGroup = sync.WaitGroup{}
		fradDetectionStatus string
		fraudError          error

		transVereficationErrMessage string
		transError                  error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		fradDetectionStatus, fraudError = srv.fraudDetection.CheckFraud(checkout)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		transVereficationErrMessage, transError = srv.transactionChecker.CheckTransaction(checkout)
	}()

	wg.Wait()

	if fraudError != nil {
		return nil, ports.ErrInternalError, nil
	}

	if fradDetectionStatus != "200" {
		return nil, ports.ErrFraudDetected, nil
	}

	if transError == ports.ErrTransIsNotValid {
		return nil, transError, errors.New(transVereficationErrMessage)
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
	}, nil, nil
}
