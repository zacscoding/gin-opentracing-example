package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-opentracing-example/pkg/logging"
	"gin-opentracing-example/pkg/middleware"
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

	header := make(http.Header)

	// setup span
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	v := span.BaggageItem(middleware.RequestIdKey)
	fmt.Println("## Check request id :", v)
	tracer := opentracing.GlobalTracer()
	err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
	logging.FromContext(ctx).Infow(">> Inject http header", "headers", header)

	if err != nil {
		return 0, nil, err
	}

	req.SetRequestURI(requestURI)
	req.Header.SetMethod("GET")
	req.Header.Add("Accept", AcceptJson)
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
