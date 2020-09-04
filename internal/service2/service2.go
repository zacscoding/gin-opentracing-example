package service2

import (
	"gin-opentracing-example/internal/response"
	"gin-opentracing-example/pkg/logging"
	"gin-opentracing-example/pkg/middleware"
	"gin-opentracing-example/pkg/remote"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"math/rand"
	"net/http"
	"time"
)

func StartService2Server() {
	gin.SetMode(gin.DebugMode)
	e := gin.Default()

	e.Use(middleware.NewRequestIdMiddleware(), middleware.NewTracingMiddleware(opentracing.GlobalTracer()))
	e.GET("/service2/trace", func(ctx *gin.Context) {
		sleep := rand.Intn(3)
		logger := logging.FromContext(ctx.Request.Context())
		logger.Infow("## Requested /service2/trace parameters", "header", ctx.Request.Header, "sleep", sleep)
		time.Sleep(time.Duration(sleep) * time.Second)

		ret := response.NewResponse(ctx)

		// call service5
		code, body, err := remote.HttpGet(ctx.Request.Context(), "http://service5:3500/service5/trace", "/service5/trace")
		if err != nil {
			ret["service5"] = gin.H{
				"code":  code,
				"error": err.Error(),
			}
		} else {
			ret["service5"] = gin.H{
				"code": code,
				"body": body,
			}
		}
		ctx.JSON(http.StatusOK, ret)
	})

	if err := e.Run(":3200"); err != nil {
		panic(err)
	}
}
