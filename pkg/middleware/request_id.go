package middleware

import (
	"gin-opentracing-example/pkg/logging"
	"gin-opentracing-example/pkg/trace"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RequestIdHeaderKey = "X-Request-ID" // request id header key
	RequestIdKey       = "requestId"    // request id context key
)

// NewRequestIdMiddleware creates a request id middleware
// (1) put requestId to context parsed from x-request-id in header or created
// (2) put a new logger with requestId into context
// (3) write requestId to header in response
func NewRequestIdMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logging.DefaultLogger().Infow("## Check headers", "uri", ctx.Request.RequestURI, "headers", ctx.Request.Header)
		requestId := ctx.Request.Header.Get(RequestIdHeaderKey)
		if requestId == "" {
			requestId = uuid.New().String()
			logging.DefaultLogger().Info("empty request id > " + requestId)
		} else {
			logging.DefaultLogger().Info("use exist request id > " + requestId)
		}
		// attach request id
		ctx.Request = ctx.Request.WithContext(trace.WithRequestId(ctx.Request.Context(), requestId))
		// attach logger
		ctx.Request = ctx.Request.WithContext(logging.WithLogger(ctx.Request.Context(),
			logging.DefaultLogger().With("requestId", requestId)))
		ctx.Writer.Header().Add(RequestIdHeaderKey, requestId)
	}
}
