package main

import (
	"context"
	"log"
	"time"

	pb "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/transaction_verification"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	target = "localhost:50052"
)

func mainV3() {
	mainV2() // InitOrder
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	defer conn.Close()

	c := pb.NewTransactionVerificationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	r, err := c.CheckOrder(ctx, &pb.TransactionVerificationRequestClock{
		OrderId: "1",
		Clock:   []int32{0, 0, 0},
	})

	if err != nil {
		log.Panic("Error while calling: ", err.Error())
	}
    log.Println("CheckOrder")
	log.Println("isValid: ", r.Response.IsValid)
	log.Println("Error message: ", r.Response.ErrMessage)
	log.Println("clock: ", r.Clock)
	log.Println("------------------")
}

func mainV4() {
	mainV3() // InitOrder
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	defer conn.Close()

	c := pb.NewTransactionVerificationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	r, err := c.CheckUser(ctx, &pb.TransactionVerificationRequestClock{
		OrderId: "1",
		Clock:   []int32{0, 0, 0},
	})

	if err != nil {
		log.Panic("Error while calling: ", err.Error())
	}
    log.Println("CheckUser")
	log.Println("isValid: ", r.Response.IsValid)
	log.Println("Error message: ", r.Response.ErrMessage)
	log.Println("clock: ", r.Clock)
	log.Println("------------------")
}

func main() {
	mainV4() // InitOrder
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	defer conn.Close()

	c := pb.NewTransactionVerificationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	r, err := c.CheckFormatCreditCard(ctx, &pb.TransactionVerificationRequestClock{
		OrderId: "1",
		Clock:   []int32{0, 0, 0},
	})

	if err != nil {
		log.Panic("Error while calling: ", err.Error())
	}
    log.Println("CheckFormatCreditCard")
	log.Println("isValid: ", r.Response.IsValid)
	log.Println("Error message: ", r.Response.ErrMessage)
	log.Println("clock: ", r.Clock)
	log.Println("------------------")
}
func mainV2() {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	defer conn.Close()

	c := pb.NewTransactionVerificationServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	r, err := c.InitOrder(ctx, &pb.TransactionVerificationRequestInit{
		OrderId: "1",
		Order: &pb.Order{
			Items:  []*pb.Item{
			    {
			        Name: "Book A",
			        Quantity: 1,
			    },
			},
			DiscountCode:            "holi",
			GiftMessage:             "Grr",
			GiftWrapping:            true,
			TermsAccepted:           false,
			NotificationPreferences: []string{},
			Device: &pb.Device{
				Type:  "",
				Model: "",
				Os:    "",
			},
			Browser: &pb.Browser{
				Name:    "",
				Version: "",
			},
			AppVersion:       "",
			ScreenResolution: "",
			Referrer:         "",
			DeviceLanguage:   "",
			CreditCard: &pb.CreditCard{
				Number:         "1111111111111111",
				ExpirationDate: "04/27",
				Cvv:            "903",
			},
			ShippingMethod: "Standard",
			BillingAddress: &pb.Address{
				Street:  "Tartu",
				City:    "Tartu",
				Country: "Estonia",
				State:   "Tartu",
				Zip:     "35500",
			},
			UserComment: "Hi",
			User: &pb.User{
				Name: "Test",
				Contact: "test@example.com",
			},
		},
	})

	if err != nil {
		log.Panic("Error while calling: ", err.Error())
	}

    log.Println("InitOrder")
	log.Println("isValid: ", r.IsValid) // if isValid is false, handle ID repetition
	log.Println("Error message: ", r.ErrMessage)
	log.Println("------------------")
}
