package main

import (
	"log"
	"os"
	"time"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/ports"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/handlers/httphdl"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/repositories/fraud"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/repositories/suggestions"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/repositories/transcheck"
	"github.com/gin-gonic/gin"
)

const (
	FRAUD_DNS_ENV_NAME         string = "fraud_dns"
	TRANS_CHECKER_DNS_ENV_NAME string = "transaction_verification_dns"
	SUGGEST_BOOKS_DNS_ENV_NAME string = "suggestions_dns"
	SERVICE_PORT               string = ":8081"
)

var (
	fraudService        ports.IFraudDetection
	transCheckerService ports.ITransactionVerification
	suggestService      ports.ISuggestionsService
	defaultTimeOut      = 4 * time.Second
)

func init() {
	fraudDNS, isThereFraud := os.LookupEnv(FRAUD_DNS_ENV_NAME)
	if !isThereFraud {
		log.Panic("Fraud detection system DNS is needed")
	}

	fraudService = fraud.NewFraudDetectionService(fraudDNS, defaultTimeOut)

	tranServiceDNS, ok := os.LookupEnv(TRANS_CHECKER_DNS_ENV_NAME)
	if !ok {
		log.Panic("transaction dns system DNS is needed")
	}

	transCheckerService = transcheck.NewTransactionVerification(tranServiceDNS, defaultTimeOut)

	suggestionsServiceDNS, ok := os.LookupEnv(SUGGEST_BOOKS_DNS_ENV_NAME)
	if !ok {
		log.Panic("transaction dns system DNS is needed")
	}

	suggestService = suggestions.NewSuggestionService(suggestionsServiceDNS, defaultTimeOut)
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
		)
		hdl    = httphdl.NewHTTPHandler(srv)
		router = gin.Default()
	)

	router.Use(CORSMiddleware())
	hdl.SetRouter(router)
	router.Run(SERVICE_PORT)
}
