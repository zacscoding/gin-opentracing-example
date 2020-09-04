package service2

import (
	"gin-opentracing-example/pkg/logging"
	"gin-opentracing-example/pkg/middleware"
	"gin-opentracing-example/pkg/remote"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

func StartService2Server() {
	gin.SetMode(gin.DebugMode)
	e := gin.Default()

	e.Use(middleware.NewRequestIdMiddleware(), middleware.NewTracingMiddleware(opentracing.GlobalTracer()))
	e.GET("/service2/trace", func(ctx *gin.Context) {
		logger := logging.FromContext(ctx.Request.Context())
		logger.Infow("## Requested /service2/trace parameters", "header", ctx.Request.Header)

		ret := gin.H{
			"headers": ctx.Request.Header,
		}

		// call service5
		code, body, err := remote.HttpGet(ctx.Request.Context(), "http://service5:3200/service5/trace", "service5")
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
