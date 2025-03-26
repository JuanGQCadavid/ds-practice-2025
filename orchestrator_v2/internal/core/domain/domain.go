package domain

import "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/common"

type User struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

type CreditCard struct {
	Number         string `json:"number"`
	ExpirationDate string `json:"expirationDate"`
	Cvv            string `json:"cvv"`
}

type Item struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type BillingAddress struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
}
type Device struct {
	Type  string `json:"type"`
	Model string `json:"model"`
	Os    string `json:"os"`
}

type Browser struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Checkout struct {
	User                       User           `json:"user"`
	CreditCard                 CreditCard     `json:"creditCard"`
	UserComment                string         `json:"userComment"`
	Items                      []Item         `json:"items"`
	DiscountCode               string         `json:"discountCode"`
	ShippingMethod             string         `json:"shippingMethod"`
	ClientCard                 string         `json:"clientCard"`
	GiftMessage                string         `json:"giftMessage"`
	BillingAddress             BillingAddress `json:"billingAddress"`
	GiftWrapping               bool           `json:"giftWrapping"`
	TermsAndConditionsAccepted bool           `json:"termsAndConditionsAccepted"`
	NotificationPreferences    []string       `json:"notificationPreferences"`
	Device                     Device         `json:"device"`
	Browser                    Browser        `json:"browser"`
	AppVersion                 string         `json:"appVersion"`
	ScreenResolution           string         `json:"screenResolution"`
	Referrer                   string         `json:"referrer"`
	DeviceLanguage             string         `json:"deviceLanguage"`
}

type SuggestedBook struct {
	BookId string `json:"bookId"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type CheckoutResponse struct {
	OrderId        string           `json:"orderId"`
	Status         string           `json:"status"`
	SuggestedBooks []*SuggestedBook `json:"suggestedBooks"`
}

func FromCheckoutToCommon(checkout *Checkout) *common.Order {

	var (
		items []*common.Item = make([]*common.Item, len(checkout.Items))
	)

	for i := range items {
		items[i] = &common.Item{
			Name:     checkout.Items[i].Name,
			Quantity: int32(checkout.Items[i].Quantity),
		}
	}

	return &common.Order{
		User: &common.User{
			Name:    checkout.User.Name,
			Contact: checkout.User.Contact,
		},
		CreditCard: &common.CreditCard{
			Number:         checkout.CreditCard.Number,
			Cvv:            checkout.CreditCard.Cvv,
			ExpirationDate: checkout.CreditCard.ExpirationDate,
		},
		UserComment:    checkout.UserComment,
		DiscountCode:   checkout.DiscountCode,
		ShippingMethod: checkout.ShippingMethod,
		ClientCard:     checkout.ClientCard,
		GiftMessage:    checkout.GiftMessage,
		BillingAddress: &common.Address{
			Street:  checkout.BillingAddress.Street,
			City:    checkout.BillingAddress.City,
			State:   checkout.BillingAddress.State,
			Zip:     checkout.BillingAddress.Zip,
			Country: checkout.BillingAddress.Country,
		},
		GiftWrapping:  checkout.GiftWrapping,
		TermsAccepted: checkout.TermsAndConditionsAccepted,
		Device: &common.Device{
			Model: checkout.Device.Model,
			Os:    checkout.Device.Os,
			Type:  checkout.Device.Type,
		},
		Browser: &common.Browser{
			Name:    checkout.Browser.Name,
			Version: checkout.Browser.Version,
		},
		Items:            items,
		AppVersion:       checkout.AppVersion,
		ScreenResolution: checkout.ScreenResolution,
		Referrer:         checkout.Referrer,
		DeviceLanguage:   checkout.DeviceLanguage,
	}

}
