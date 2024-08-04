package main

import (
	"gateway/application/routing/delivery"
	"gateway/application/routing/delivery/service"
	"gateway/application/routing/usecase"
	"gateway/config"
	_ "gateway/docs"
	"gateway/pkg/logger"
	"gateway/pkg/middleware"
	"gateway/pkg/redis"
	"gateway/pkg/utils"
	"gateway/pkg/wrapper"
	"github.com/gin-contrib/cors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.elastic.co/apm/module/apmgin/v2"
)

// @title Sample PUBLIC Gateway
// @version 1.0
// @description REST -> GRPC API Gateway / Permission and More for all public services

// @contact.name TuanNguyen
// @contact.url http://www.swagger.io/support
// @contact.email tuannadldev@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3010
// @BasePath /
// @query.collection.format multi
func main() {
	configPath := utils.GetConfigPath(os.Getenv("ENVIRONMENT"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalln("Failed at config", err)
	}
	// Init AppLog
	logger.Newlogger(logger.ConfigLogger{})
	appLogger := logger.GetLogger()
	rdb, err := redis.InitConnection(cfg)
	if err != nil {
		appLogger.Panic("Can't connect Redis ", err)
	}
	defer rdb.Close()

	if os.Getenv("ENVIRONMENT") != "local" {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}
	r := gin.Default()
	r.RedirectTrailingSlash = false
	// if you don't use any proxy, you can disable this feature by using
	// r.SetTrustedProxies(nil)
	middleware.UseErrorHandling(r)
	r.Use(middleware.JSON())
	r.Use(apmgin.Middleware(r))
	r.Use(middleware.BucketRateLimiter(rdb))

	// - No origin allowed by default
	// - GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	cors_config := cors.DefaultConfig()
	if len(cfg.Cors.AllowOrigins) > 0 {
		cors_config.AllowOrigins = cfg.Cors.AllowOrigins
	} else {
		cors_config.AllowAllOrigins = true
	}
	if len(cfg.Cors.AllowHeaders) > 0 {
		cors_config.AllowHeaders = cfg.Cors.AllowHeaders
	} else {
		cors_config.AllowHeaders = []string{"*"}
	}

	r.Use(cors.New(cors_config))

	serviceClient := service.NewServiceClient(cfg)

	if os.Getenv("ENVIRONMENT") != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	pingHandler := wrapper.WithContext(func(ctx *wrapper.Context) {
		ip := ctx.Request.Header.Get("X-Forwarded-For")
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"version": cfg.Server.AppVersion,
			"ip":      ip,
		})
	})

	r.GET("/health-check", pingHandler)
	r.HEAD("/health-check", pingHandler)

	routingUC := usecase.NewRoutingUseCase(serviceClient)
	delivery.RegisterRoutes(r, cfg, rdb, routingUC)
	r.Run(cfg.Server.Port)
}
