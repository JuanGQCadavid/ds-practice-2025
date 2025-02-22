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
	suggestedBooks     ports.ISuggestionsService
}

func NewService(
	fraudDetection ports.IFraudDetection,
	transactionChecker ports.ITransactionVerification,
	suggestedBooks ports.ISuggestionsService,
) *Service {
	return &Service{
		fraudDetection:     fraudDetection,
		transactionChecker: transactionChecker,
		suggestedBooks:     suggestedBooks,
	}
}

func (srv *Service) Checkout(checkout *domain.Checkout) (*domain.CheckoutResponse, error, error) {

	var (
		wg                  sync.WaitGroup = sync.WaitGroup{}
		fraudDetectionStatus string
		fraudError          error

		transVerificationErrMessage string
		transError                  error

		suggestedBooks []*domain.SuggestedBook
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		fraudDetectionStatus, fraudError = srv.fraudDetection.CheckFraud(checkout)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		transVerificationErrMessage, transError = srv.transactionChecker.CheckTransaction(checkout)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		suggestedBooks, _ = srv.suggestedBooks.SuggestBooks(checkout.Items)
	}()

	wg.Wait()

	if fraudError != nil {
		return nil, ports.ErrInternalError, nil
	}

	if fraudDetectionStatus != "200" {
		return nil, ports.ErrFraudDetected, nil
	}

	if transError == ports.ErrTransIsNotValid {
		return nil, transError, errors.New(transVerificationErrMessage)
	}

	return &domain.CheckoutResponse{
		Status:         "Order Approved",
		OrderId:        "Pending",
		SuggestedBooks: suggestedBooks,
	}, nil, nil
}
