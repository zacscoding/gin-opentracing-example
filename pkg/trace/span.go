package trace

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func ExtractSpan(ctx *gin.Context) opentracing.Span {
	carrier := opentracing.HTTPHeadersCarrier(ctx.Request.Header)
	tracer := opentracing.GlobalTracer()
	wireContext, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
	if err != nil {
		return opentracing.StartSpan(ctx.Request.URL.Path)
	}
	return opentracing.StartSpan(ctx.Request.URL.Path, opentracing.ChildOf(wireContext))
}
