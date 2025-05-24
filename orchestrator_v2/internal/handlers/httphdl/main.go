package httphdl

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/domain"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/ports"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("orchestactor-server")

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
	logger := log.
		With().
		Str("path", context.Request.URL.Path).
		Str("trace_ID", uuid.NewString()).
		Logger()

	ctx := logger.WithContext(context)

	// body, _ := context.Request.GetBody()
	// bodyBytes, err := io.ReadAll(body)
	// if err != nil {
	// 	log.Println("We fail to read the body")
	// }
	// bodyString := string(bodyBytes)

	_, span := tracer.Start(context.Request.Context(), "checkout", trace.WithAttributes(attribute.String("request", "Dude")))
	defer span.End()

	var checkoutRequest *domain.Checkout = &domain.Checkout{}

	if err := context.BindJSON(&checkoutRequest); err != nil {
		logger.Error().Msg("error while casting request")
		context.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			ErrorType: ErrorType{
				Code:    "NO IDEA",
				Message: "error while casting request",
			},
		})
		return
	}

	response, genErr := hdl.service.Checkout(ctx, checkoutRequest)

	switch genErr {
	case ports.ErrInternalError:
		context.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
			ErrorType: ErrorType{
				Code:    "Internal error",
				Message: genErr.Error(),
			},
		})
		return

	case ports.ErrTransIsNotValid:
		context.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			ErrorType: ErrorType{
				Code: "Bad request, trans is not valid",
			},
		})
		return
	case ports.ErrFraudDetected:
		context.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
			ErrorType: ErrorType{
				Code:    genErr.Error(),
				Message: "It seems someone is trying to commit a crime...",
			},
		})
		return
	}

	if genErr != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			ErrorType: ErrorType{
				Code:    genErr.Error(),
				Message: genErr.Error(),
			},
		})
		return
	}

	logger.Info().Msgf("%+v\n", response)

	context.JSON(http.StatusOK, response)
}
