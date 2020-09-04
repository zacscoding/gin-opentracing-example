// copied from https://github.com/opentracing-contrib/go-gin/blob/master/ginhttp/server.go
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
	"net/url"
)

const defaultComponentName = "net/http"

type mwOptions struct {
	opNameFunc    func(r *http.Request) string
	spanObserver  func(span opentracing.Span, r *http.Request)
	urlTagFunc    func(u *url.URL) string
	componentName string
}

// MWOption controls the behavior of the Middleware.
type MWOption func(*mwOptions)

// OperationNameFunc returns a MWOption that uses given function f
// to generate operation name for each server-side span.
func OperationNameFunc(f func(r *http.Request) string) MWOption {
	return func(options *mwOptions) {
		options.opNameFunc = f
	}
}

// MWComponentName returns a MWOption that sets the component name
// for the server-side span.
func MWComponentName(componentName string) MWOption {
	return func(options *mwOptions) {
		options.componentName = componentName
	}
}

// MWSpanObserver returns a MWOption that observe the span
// for the server-side span.
func MWSpanObserver(f func(span opentracing.Span, r *http.Request)) MWOption {
	return func(options *mwOptions) {
		options.spanObserver = f
	}
}

// MWURLTagFunc returns a MWOption that uses given function f
// to set the span's http.url tag. Can be used to change the default
// http.url tag, eg to redact sensitive information.
func MWURLTagFunc(f func(u *url.URL) string) MWOption {
	return func(options *mwOptions) {
		options.urlTagFunc = f
	}
}

// Middleware is a gin native version of the equivalent middleware in:
//   https://github.com/opentracing-contrib/go-stdlib/
func NewTracingMiddleware(tr opentracing.Tracer, options ...MWOption) gin.HandlerFunc {
	opts := mwOptions{
		opNameFunc: func(r *http.Request) string {
			return "HTTP " + r.Method
		},
		spanObserver: func(span opentracing.Span, r *http.Request) {},
		urlTagFunc: func(u *url.URL) string {
			return u.String()
		},
	}
	for _, opt := range options {
		opt(&opts)
	}

	return func(c *gin.Context) {
		var span opentracing.Span
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		wireContext, err := tr.Extract(opentracing.HTTPHeaders, carrier)
		if err != nil {
			span = opentracing.StartSpan(c.Request.URL.Path)
		} else {
			span = opentracing.StartSpan(c.Request.URL.Path, opentracing.ChildOf(wireContext))
		}
		defer span.Finish()
		span.SetTag("request.id", c.Request.Header.Get("X-Request-Id"))

		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, opts.urlTagFunc(c.Request.URL))
		opts.spanObserver(span, c.Request)
		// set component name, use "net/http" if caller does not specify
		componentName := opts.componentName
		if componentName == "" {
			componentName = defaultComponentName
		}
		ext.Component.Set(span, componentName)
		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))
		c.Next()
		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
	}
}
