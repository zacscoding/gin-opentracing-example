package main

import (
	"fmt"
	"gin-opentracing-example/internal/service1"
	"gin-opentracing-example/internal/service2"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
	"github.com/uber/jaeger-client-go/zipkin"
	"os"
)

func main() {
	os.Args = append(os.Args, "service1")
	if len(os.Args) < 2 {
		fmt.Println("empty service name")
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
		fmt.Println("## Start service1")
		service1.StartService1Server()
	case "service2":
		fmt.Println("## Start service2")
		service2.StartService2Server()
	case "service3":
		fmt.Println("## Start service3")
	case "service4":
		fmt.Println("## Start service4")
	case "service5":
		fmt.Println("## Start service5")
	default:
		fmt.Println("Unknown service name:", os.Args[1])
	}
}
