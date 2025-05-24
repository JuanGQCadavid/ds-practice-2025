package main

import (
	"context"
	"os"
	"time"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/ports"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/handlers/httphdl"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/repositories/fraud"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/repositories/orderqueue"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/repositories/suggestions"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/repositories/transcheck"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	FRAUD_DNS_ENV_NAME         string = "fraud_dns"
	TRANS_CHECKER_DNS_ENV_NAME string = "transaction_verification_dns"
	SUGGEST_BOOKS_DNS_ENV_NAME string = "suggestions_dns"
	ORDER_QUEUE_DNS_ENV_NAME   string = "order_queue_dns"
	SERVICE_PORT               string = ":8081"
)

var (
	fraudService        ports.IFraudDetection
	transCheckerService ports.ITransactionVerification
	suggestService      ports.ISuggestionsService
	orderQueue          ports.IOrderQueue
	defaultTimeOut      = 4 * time.Second
)

func init() {
	fraudDNS, isThereFraud := os.LookupEnv(FRAUD_DNS_ENV_NAME)
	if !isThereFraud {
		log.Fatal().Msg("Fraud detection system DNS is needed")
	}

	fraudService = fraud.NewFraudDetectionService(fraudDNS, defaultTimeOut)

	tranServiceDNS, ok := os.LookupEnv(TRANS_CHECKER_DNS_ENV_NAME)
	if !ok {
		log.Fatal().Msg("Transaction dns system DNS is needed")
	}

	transCheckerService = transcheck.NewTransactionVerification(tranServiceDNS, defaultTimeOut)

	suggestionsServiceDNS, ok := os.LookupEnv(SUGGEST_BOOKS_DNS_ENV_NAME)
	if !ok {
		log.Fatal().Msg("Suggestions dns system DNS is needed")
	}

	suggestService = suggestions.NewSuggestionService(suggestionsServiceDNS, defaultTimeOut)

	orderDNS, ok := os.LookupEnv(ORDER_QUEUE_DNS_ENV_NAME)
	if !ok {
		log.Fatal().Msg("Order queue dns system DNS is needed")
	}

	orderQueue = orderqueue.NewOrderQueue(orderDNS, defaultTimeOut)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	var (
		srv = core.NewService(
			fraudService,
			transCheckerService,
			suggestService,
			orderQueue,
		)
		hdl    = httphdl.NewHTTPHandler(srv)
		router = gin.Default()
	)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	shutdown, err := SetupOTelSDK(context.Background())

	if err != nil {
		log.Fatal().Msg("Imposible to setup OTEL")
	}

	defer shutdown(context.Background())

	router.Use(CORSMiddleware())
	router.Use(otelgin.Middleware("orchestactor"))

	router.Use(
		func(c *gin.Context) {
			start := time.Now()
			c.Next()
			log.Info().
				Int("status", c.Writer.Status()).
				Dur("latency", time.Since(start)).
				Str("client_ip", c.ClientIP()).
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Send()
		},
	)

	hdl.SetRouter(router)
	router.Run(SERVICE_PORT)
}
