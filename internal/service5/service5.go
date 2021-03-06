package service5

import (
	"gin-opentracing-example/internal/response"
	"gin-opentracing-example/pkg/logging"
	"gin-opentracing-example/pkg/middleware"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

func StartService5Server() {
	gin.SetMode(gin.DebugMode)
	e := gin.Default()

	e.Use(middleware.NewRequestIdMiddleware(), middleware.NewTracingMiddleware(middleware.MWComponentName("service5")))
	e.GET("/service5/trace", func(ctx *gin.Context) {
		sleep := rand.Intn(3)
		logger := logging.FromContext(ctx.Request.Context())
		logger.Infow("## Requested /service5/trace parameters", "header", ctx.Request.Header, "sleep", sleep)
		time.Sleep(time.Duration(sleep) * time.Second)
		ctx.JSON(http.StatusOK, response.NewResponse(ctx))
	})

	if err := e.Run(":3500"); err != nil {
		panic(err)
	}
}
