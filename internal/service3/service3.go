package service3

import (
	"gin-opentracing-example/internal/response"
	"gin-opentracing-example/pkg/logging"
	"gin-opentracing-example/pkg/middleware"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

func StartService3Server() {
	gin.SetMode(gin.DebugMode)
	e := gin.Default()

	e.Use(middleware.NewRequestIdMiddleware(), middleware.NewTracingMiddleware(middleware.MWComponentName("service3")))
	e.GET("/service3/trace", func(ctx *gin.Context) {
		sleep := rand.Intn(3)
		logger := logging.FromContext(ctx.Request.Context())
		logger.Infow("## Requested /service3/trace parameters", "header", ctx.Request.Header, "sleep", sleep)
		time.Sleep(time.Duration(sleep) * time.Second)
		ctx.JSON(http.StatusOK, response.NewResponse(ctx))
	})

	if err := e.Run(":3300"); err != nil {
		panic(err)
	}
}
