package response

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// NewResponse returns a gin.H from request header if satisfy below conditions.
// 1) starts with x-b3
// 2) equals to x-request-id
func NewResponse(ctx *gin.Context) gin.H {
	ret := make(gin.H)
	for k, v := range ctx.Request.Header {
		lower := strings.ToLower(k)
		if strings.HasPrefix(lower, "x-b3") || strings.HasPrefix(lower, "x-request-id") {
			ret[k] = v
		}
	}
	return ret
}
