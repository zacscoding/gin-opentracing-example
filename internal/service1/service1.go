package service1

import (
	"fmt"
	"gin-opentracing-example/pkg/logging"
	"gin-opentracing-example/pkg/middleware"
	"gin-opentracing-example/pkg/remote"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

func StartService1Server() {
	gin.SetMode(gin.DebugMode)
	e := gin.Default()

	e.Use(middleware.NewRequestIdMiddleware(), middleware.NewTracingMiddleware(opentracing.GlobalTracer()))
	e.GET("/service1/trace", func(ctx *gin.Context) {
		logger := logging.FromContext(ctx.Request.Context())
		logger.Infow("## Requested /service1/trace parameters", "header", ctx.Request.Header)

		tracer := opentracing.GlobalTracer()
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ctx.Request.Header))
		fmt.Println(spanCtx)

		ret := gin.H{
			"service1": gin.H{
				"headers": ctx.Request.Header,
			},
		}
		// 1) call service2
		code, body, err := remote.HttpGet(ctx.Request.Context(), "http://service2:3200/service2/trace", "service2")
		if err != nil {
			ret["service2"] = gin.H{
				"code":  code,
				"error": err.Error(),
			}
		} else {
			ret["service2"] = gin.H{
				"code": code,
				"body": body,
			}
		}
		// 2) call service3
		code, body, err = remote.HttpGet(ctx.Request.Context(), "http://service3:3300/service3/trace", "service3")
		if err != nil {
			ret["service3"] = gin.H{
				"code":  code,
				"error": err.Error(),
			}
		} else {
			ret["service3"] = gin.H{
				"code": code,
				"body": body,
			}
		}
		ctx.JSON(http.StatusOK, ret)
	})

	if err := e.Run(":3100"); err != nil {
		panic(err)
	}
}
