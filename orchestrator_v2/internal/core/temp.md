

var (
	initState = []int32{0, 0, 0}
	cState    = []int32{1, 0, 0}
	dState    = []int32{2, 0, 0}
	eState    = []int32{3, 1, 0}
	fState    = []int32{3, 2, 0}
)

func (srv *Service) startService(orderId string, clock []int32) error {
	var (
		wg     sync.WaitGroup = sync.WaitGroup{}
		aClock []int32
		aErr   error

		bClock []int32
		bErr   error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		aClock, aErr = srv.transactionChecker.CheckOrder(orderId, clock)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bClock, bErr = srv.transactionChecker.CheckUser(orderId, clock)
	}()

	wg.Wait()

	if aErr != nil {
		return aErr
	}

	if bErr != nil {
		return bErr
	}

	return nil
}



var mapStates = map[*[]int32]func(string, []int32) error{
		&initState: srv.startService,
	}
