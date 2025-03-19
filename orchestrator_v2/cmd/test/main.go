package main

import (
	"context"
	"log"
	"time"

	commonProtoBus "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/common"
	fraudProtoBus "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/fraud_detection"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	target = "localhost:50051"
)

func initOrder() {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	defer conn.Close()

	c := fraudProtoBus.NewFraudDetectionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	r, err := c.InitOrder(ctx, &commonProtoBus.InitRequest{
		OrderId: "1",
		Order: &commonProtoBus.Order{
			Items: []*commonProtoBus.Item{
				{
					Name:     "Book A",
					Quantity: 1,
				},
			},
			DiscountCode:            "holi",
			GiftMessage:             "Grr",
			GiftWrapping:            true,
			TermsAccepted:           false,
			NotificationPreferences: []string{},
			Device: &commonProtoBus.Device{
				Type:  "",
				Model: "",
				Os:    "",
			},
			Browser: &commonProtoBus.Browser{
				Name:    "",
				Version: "",
			},
			AppVersion:       "",
			ScreenResolution: "",
			Referrer:         "",
			DeviceLanguage:   "",
			CreditCard: &commonProtoBus.CreditCard{
				Number:         "1111111111111111",
				ExpirationDate: "04/27",
				Cvv:            "903",
			},
			ShippingMethod: "Standard",
			BillingAddress: &commonProtoBus.Address{
				Street:  "Tartu",
				City:    "Tartu",
				Country: "Estonia",
				State:   "Tartu",
				Zip:     "35500",
			},
			UserComment: "Hi",
			User: &commonProtoBus.User{
				Name:    "Test",
				Contact: "test@example.com",
			},
		},
	})

	if err != nil {
		log.Panic("Error while calling: ", err.Error())
	}

	log.Println("InitOrder")
	log.Println("IsValid: ", r.IsValid)
	log.Println("ErrMessage: ", r.ErrMessage)
	log.Println("------------------")
}

func checkUser() {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	defer conn.Close()
	c := fraudProtoBus.NewFraudDetectionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	r, err := c.CheckUser(ctx, &commonProtoBus.NextRequest{
		OrderId: "1",
		IncomingVectorClock:   []int32{0, 0, 0},
	})

	if err != nil {
		log.Panic("Error while calling: ", err.Error())
	}
	log.Println("CheckUser")
	log.Println("IsValid: ", r.IsValid) // if IsValid is False, handle error
	log.Println("VectorClock: ", r.VectorClock)
	log.Println("ErrMessage", r.ErrMessage)
	log.Println("------------------")
}

func checkCreditCard() {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	defer conn.Close()
	c := fraudProtoBus.NewFraudDetectionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	r, err := c.CheckCreditCard(ctx, &commonProtoBus.NextRequest{
		OrderId: "1",
		IncomingVectorClock:   []int32{0, 0, 0},
	})

	if err != nil {
		log.Panic("Error while calling: ", err.Error())
	}
	log.Println("CheckCreditCard")
	log.Println("IsValid: ", r.IsValid) // if IsValid is False, handle error
	log.Println("VectorClock: ", r.VectorClock)
	log.Println("ErrMessage: ", r.ErrMessage)
	log.Println("------------------")
}

func main() {
	initOrder()
	checkUser()
	checkCreditCard()
}
