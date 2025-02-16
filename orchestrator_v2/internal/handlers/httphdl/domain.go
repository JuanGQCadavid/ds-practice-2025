package httphdl

type ErrorType struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	ErrorType ErrorType `json:"error"`
}
