package main

import (
	"context"
	"log"
	"time"

	pb "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/fraud_detection"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	target = "localhost:50052"
)

func initOrder() {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	defer conn.Close()

	c := pb.NewFraudDetectionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()


	r, err := c.InitOrder(ctx, &pb.FraudDetectionRequestInit{
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
	log.Println("code: ", r.code) // if code is 400, handle error
	log.Println("------------------")
}

func main() {
	initOrder()
}