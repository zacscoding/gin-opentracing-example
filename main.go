package main

import (
	"gin-opentracing-example/internal/service1"
	"gin-opentracing-example/internal/service2"
	"gin-opentracing-example/internal/service3"
	"gin-opentracing-example/internal/service4"
	"gin-opentracing-example/internal/service5"
	"gin-opentracing-example/pkg/logging"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
	"github.com/uber/jaeger-client-go/zipkin"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		logging.DefaultLogger().Error("empty service name")
		return
	}

	propagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	t := transport.NewHTTPTransport("http://jaegertracing:14268/api/traces", transport.HTTPBatchSize(1))
	tracer, closer := jaeger.NewTracer(
		os.Args[1],
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(t),
		jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, propagator),
		jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, propagator),
		jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
	)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	switch os.Args[1] {
	case "service1":
		service1.StartService1Server()
	case "service2":
		service2.StartService2Server()
	case "service3":
		service3.StartService3Server()
	case "service4":
		service4.StartService4Server()
	case "service5":
		service5.StartService5Server()
	default:
		logging.DefaultLogger().Error("Unknown service name: " + os.Args[1])
	}
}
