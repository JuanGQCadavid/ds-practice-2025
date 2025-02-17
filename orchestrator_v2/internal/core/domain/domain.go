package domain

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
