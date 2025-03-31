package main

import (
	"context"
	"log"
	"time"

	commonProtoBus "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/common"
	// 	transactionProtoBus "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/transaction_verification"
	//fraudProtoBus "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/fraud_detection"
	queueProtoBus "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/order_queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	target = "localhost:50054"
)

func enqueue() {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	defer conn.Close()
	c := queueProtoBus.NewOrderQueueServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	defer cancel()

	r, err := c.Enqueue(ctx, &queueProtoBus.EnqueueRequest{
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
				Number:         "1289789012345678",
				ExpirationDate: "04/27",
				Cvv:            "903",
			},
			ShippingMethod: "Standard",
			ClientCard:     "None",
			BillingAddress: &commonProtoBus.Address{
				Street:  "Tartu",
				City:    "Tartu",
				Country: "Estonia",
				State:   "Tartu",
				Zip:     "35500",
			},
			UserComment: "Hi",
			User: &commonProtoBus.User{
				Name:    "Maria Perez",
				Contact: "mp@example.com",
			},
		},
	})
	if err != nil {
		log.Panic("Error while calling: ", err.Error())
	}
	log.Println("Enqueue")
	log.Println("IsValid: ", r.IsValid) // if IsValid is False, handle error
	log.Println("ErrMessage", r.ErrMessage)
	log.Println("------------------")
}

func dequeue() {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	defer conn.Close()
	c := queueProtoBus.NewOrderQueueServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	r, err := c.Dequeue(ctx, &queueProtoBus.EmptyRequest{})

	if err != nil {
		log.Panic("Error while calling: ", err.Error())
	}

	log.Println("Dequeue")
	log.Println("IsValid: ", r.IsValid) // if IsValid is False, handle error
	log.Println("ErrMessage", r.ErrMessage)
	log.Println("Order: ", r.OrderId)
	log.Println("Order: ", r.Order)
	log.Println("------------------")
}

func clean() {

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Panic("Error while creating conn: ", err.Error())
	}

	defer conn.Close()
	c := queueProtoBus.NewOrderQueueServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	r, err := c.Clean(ctx, &queueProtoBus.EmptyRequest{})

	if err != nil {
		log.Panic("Error while calling: ", err.Error())
	}

	log.Println("Clean")
	log.Println("IsValid: ", r.IsValid) // if IsValid is False, handle error
	log.Println("------------------")

}

func main() {
	enqueue()
	dequeue()
	clean()
}
