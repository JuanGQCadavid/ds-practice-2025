package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/JuanGQCadavid/ds-practice-2025/utils/pb/common"
	"github.com/JuanGQCadavid/ds-practice-2025/utils/pb/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	PREFERED    int64 = 1
	DB_REPLICAS       = map[int64]string{
		1: "db_replica1",
		2: "db_replica2",
		3: "db_replica3",
	}
)

type DBService struct {
	targets      map[int64]string
	actualTarget int64

	timeout time.Duration
	client  database.DatabaseServiceClient
}

func NewDBService() *DBService {
	return NewDBServiceWithTargets(PREFERED, DB_REPLICAS)
}

func NewDBServiceWithTargets(prefered int64, targets map[int64]string) *DBService {
	conn, err := grpc.NewClient(targets[prefered], grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	return &DBService{
		actualTarget: prefered,
		targets:      targets,
		client:       database.NewDatabaseServiceClient(conn),
		timeout:      4 * time.Second,
	}
}

func (srv *DBService) retry() error {
	attemps := 3

	for attemps >= 0 {
		if err := srv.checkService(); err == nil {
			return nil
		}

		log.Println("Sleeping ... ")
		time.Sleep(1 * time.Second)
		attemps -= 1
	}

	log.Println("PANIC! there is not way to reach the db")
	return fmt.Errorf("err there is not way to reach the db")
}

func (srv *DBService) checkService() error {
	ctx, cancel := context.WithTimeout(context.Background(), srv.timeout)
	defer cancel()

	resp, err := srv.client.GetStatus(ctx, &database.EmptyRequest{})

	if err != nil {
		log.Println("Shit!!, db is down, who else can i bother?")
		return fmt.Errorf("err db is donw")
	}

	if srv.actualTarget == resp.LeaderID {
		return nil
	}

	// TODO - should we lock this ?

	log.Println("knock knock! ", srv.targets[resp.LeaderID])

	conn, err := grpc.NewClient(srv.targets[resp.LeaderID], grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("Error while creating conn: ", err.Error())
		return err
	}

	newClient := database.NewDatabaseServiceClient(conn)

	log.Println("Updating leader ID and conenction")
	srv.actualTarget = resp.LeaderID
	srv.client = newClient
	return nil
}

func (srv *DBService) Prepare(orderId string, items []*common.Item) error {
	// Checking leader
	if err := srv.retry(); err != nil {
		log.Println("Database - Prepare - no possible, db down")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), srv.timeout)
	defer cancel()

	booksRequest := make([]*database.BookRequestPrepare, len(items))

	for i, book := range items {
		booksRequest[i] = &database.BookRequestPrepare{
			BookID:   book.Name,
			Quantity: book.Quantity,
		}
	}

	resp, err := srv.client.Prepare(ctx, &database.PrepareRequest{
		OrderID:      orderId,
		BookRequests: booksRequest,
	})

	if err != nil {
		log.Println("Database - Prepare - error: ", err.Error())
		return err
	}

	if !resp.IsValid {
		log.Println("Books ran out of inventory")
		for _, book := range resp.BookRequests {
			log.Println(book.BookID, " only ", book.Quantity)
		}

		return fmt.Errorf("err No stock")
	}

	return nil
}

func (srv *DBService) Commit(orderId string) error {
	// Checking leader
	if err := srv.retry(); err != nil {
		log.Println("Database - Commit - no possible, db down")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), srv.timeout)
	defer cancel()

	resp, err := srv.client.Commit(ctx, &database.CommitRequest{
		OrderID: orderId,
	})

	if err != nil {
		log.Println("Database - Commit - error: ", err.Error())
		return err
	}

	if !resp.IsValid {
		return fmt.Errorf("err %s", resp.ErrMessage)
	}

	return nil
}

func (srv *DBService) Abort(orderId string) error {
	if err := srv.retry(); err != nil {
		log.Println("Database - Abort - no possible, db down")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), srv.timeout)
	defer cancel()

	resp, err := srv.client.Abort(ctx, &database.AbortRequest{
		OrderID: orderId,
	})

	if err != nil {
		log.Println("Database - Abort - error: ", err.Error())
		return err
	}

	if !resp.IsValid {
		return fmt.Errorf("err %s", resp.ErrMessage)
	}

	return nil
}
