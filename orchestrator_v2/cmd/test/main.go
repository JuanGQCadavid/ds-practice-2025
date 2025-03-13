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

func main() {
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
			Items:                   []*pb.Item{},
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
				Number:         "1111",
				ExpirationDate: "1234",
				Cvv:            "1234",
			},
			ShippingMethod: "Por domicilio",
			BillingAddress: &pb.Address{
				Street:  "Tartu",
				City:    "Tartu",
				Country: "Estonia",
				State:   "Tartu",
				Zip:     "Holi",
			},
			UserComment: "Hi",
			User: &pb.User{
				Name: "Test",
			},
		},
	})

	// r, err := c.CheckFraud(ctx, &pb.FraudDetectionRequest{
	// 	Json: `{
	// 		"creditCard": {
	// 			"number": "4111111111111111",
	// 			"cvv": "123",
	// 			"expirationDate": "12/26"
	// 		}
	// 	}`,
	// })

	if err != nil {
		log.Panic("Error while calling: ", err.Error())
	}
	log.Println(r)

	log.Println("Response: ", r.IsValid)
	log.Println("Error message: ", r.ErrMessage)

}
