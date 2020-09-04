package remote

import (
	"context"
	"encoding/json"
	"gin-opentracing-example/pkg/trace"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

type StatusCode int

const (
	AcceptJson = "application/json"
)

var cli *fasthttp.Client

func init() {
	cli = &fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		ReadBufferSize:           4096,
		WriteBufferSize:          4096,
		ReadTimeout:              10 * time.Second,
		WriteTimeout:             10 * time.Second,
		MaxIdleConnDuration:      90 * time.Second,
	}
}

func HttpGet(ctx context.Context, requestURI, operationName string) (StatusCode, map[string]interface{}, error) {
	// setup http client
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	var err error
	header := make(http.Header)

	// TOOD : refactor to trace package
	// setup span
	tracer := opentracing.GlobalTracer()
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := tracer.StartSpan(operationName, opentracing.ChildOf(span.Context()))
		defer span.Finish()
		err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
		if err != nil {
			return 0, nil, err
		}
	}

	req.SetRequestURI(requestURI)
	req.Header.SetMethod("GET")
	req.Header.Add("Accept", AcceptJson)
	req.Header.Add("X-Request-ID",trace.ExtractRequestId(ctx))
	for k, v := range header {
		if len(v) > 0 {
			req.Header.Add(k, v[0])
		}
	}

	// do request
	deadline, ok := ctx.Deadline()
	if ok {
		err = cli.DoDeadline(req, resp, deadline)
	} else {
		err = cli.Do(req, resp)
	}
	code := StatusCode(resp.StatusCode())
	if err != nil {
		return code, nil, err
	}
	var respBody map[string]interface{}
	if err2 := json.Unmarshal(resp.Body(), &respBody); err2 != nil {
		return code, nil, err2
	}
	return code, respBody, nil
}

func (s StatusCode) Is2xxSuccessful() bool {
	return s >= 200 && s < 300
}
