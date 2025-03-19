package core

import (
	"fmt"
	"log"
	"sync"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/ports"
	"github.com/google/uuid"
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

func (srv *Service) Bradcast(
	orderId string,
	data *domain.Checkout,
	functions ...func(string, *domain.Checkout) error,
	// payload interface{},
	// functions ...func(interface{}) interface{},
) {
	wg := sync.WaitGroup{}
	for i, f := range functions {
		log.Printf("Broadcast: %d/%d \n", i+1, len(functions))
		wg.Add(1)
		go func() {
			defer wg.Done()
			f(orderId, data)
		}()
	}

	log.Println("Broadcast: Wating")
	wg.Wait()
}

func (srv *Service) updateClock(tick []int32, tack []int32) []int32 {
	for i := range tick {
		if tick[i] > tack[i] {
			tack[i] = tick[i]
		}
	}
	return tack
}

var (
	// initState = []int32{0, 0, 0}
	// cState    = []int32{1, 0, 0}
	// dState    = []int32{2, 0, 0}
	eState = []int32{3, 1, 0}
	// fState = []int32{3, 2, 0}
)

func (srv *Service) checkState(clock, state []int32) bool {
	for i := range clock {
		if clock[i] != state[i] {
			return false
		}
	}

	return true
}

func (srv *Service) Checkout(checkout *domain.Checkout) (*domain.CheckoutResponse, error, error) {
	var (
		orderId string  = uuid.NewString()
		clock   []int32 = []int32{
			0, 0, 0,
		}

		// Sync
		clockMutex = sync.Mutex{}
		genErr     error
		wg         = sync.WaitGroup{}
	)

	srv.Bradcast(
		orderId, checkout,
		srv.fraudDetection.Init,
		srv.transactionChecker.Init,
	)

	wg.Add(1)
	// A - C
	go func() {
		defer wg.Done()

		tack, err := srv.transactionChecker.CheckOrder(orderId, clock)

		if err != nil {
			genErr = err
			return
		}

		tack, err = srv.transactionChecker.CheckFormatCreditCard(orderId, tack)

		if err != nil {
			genErr = err
			return
		}

		clockMutex.Lock()
		clock = srv.updateClock(clock, tack)
		clockMutex.Unlock()
	}()

	wg.Add(1)
	// B - D
	go func() {
		defer wg.Done()

		tack, err := srv.transactionChecker.CheckUser(orderId, clock)

		if err != nil {
			genErr = err
			return
		}

		tack, err = srv.fraudDetection.CheckUser(orderId, tack)

		if err != nil {
			genErr = err
			return
		}

		clockMutex.Lock()
		clock = srv.updateClock(clock, tack)
		clockMutex.Unlock()
	}()

	wg.Wait()

	if genErr != nil {
		return nil, genErr, genErr
	}

	log.Println(clock)

	if !srv.checkState(clock, eState) {
		stateErr := fmt.Errorf("err Clock no on state, clock %v", clock)
		return nil, stateErr, stateErr
	}

	tack, err := srv.fraudDetection.CheckCreditCard(orderId, clock)

	if err != nil {
		return nil, err, err
	}

	clock = srv.updateClock(clock, tack)

	suggestedBooks, _ := srv.suggestedBooks.SuggestBooks(checkout.Items)

	return &domain.CheckoutResponse{
		Status:         "Order Approved",
		OrderId:        orderId,
		SuggestedBooks: suggestedBooks,
	}, nil, nil
}
