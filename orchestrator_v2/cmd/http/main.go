package main

import (
	"log"
	"os"
	"time"

	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/core/ports"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/handlers/httphdl"
	"github.com/JuanGQCadavid/ds-practice-2025/orchestrator_v2/internal/repositories/fraud"
	"github.com/gin-gonic/gin"
)

const (
	fraudDNSEnvName string = "fraud_dns"
	SERVICE_PORT    string = ":8000"
)

var (
	fraudService   ports.IFraudDetection
	defaultTimeOut = 4 * time.Second
)

func init() {
	fraudDNS, isThereFraud := os.LookupEnv(fraudDNSEnvName)
	if !isThereFraud {
		log.Panic("Fraud detection system DNS is needed")
	}

	fraudService = fraud.NewFraudDetectionService(fraudDNS, defaultTimeOut)
}

func main() {
	var (
		srv    = core.NewService(fraudService)
		hdl    = httphdl.NewHTTPHandler(srv)
		router = gin.Default()
	)
	hdl.SetRouter(router)
	router.Run(SERVICE_PORT)
}
