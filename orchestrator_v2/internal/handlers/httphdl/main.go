package httphdl

import (
	"log"
	"net/http"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	service *core.Service
}

func NewHTTPHandler(service *core.Service) *HTTPHandler {
	return &HTTPHandler{
		service: service,
	}
}

func (hdl *HTTPHandler) SetRouter(router *gin.Engine) {
	router.POST("/checkout", hdl.CheckOut) // OK
}

func (hdl *HTTPHandler) CheckOut(context *gin.Context) {
	var checkoutRequest *domain.Checkout = &domain.Checkout{}

	if err := context.BindJSON(&checkoutRequest); err != nil {
		log.Println("error while casting request")
		context.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			ErrorType: ErrorType{
				Code:    "NO IDEA",
				Message: "error while casting request",
			},
		})
		return
	}

	response, err := hdl.service.Checkout(checkoutRequest)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
			ErrorType: ErrorType{
				Code:    "Internal",
				Message: err.Error(),
			},
		})
		return
	}

	context.JSON(http.StatusOK, response)
}
