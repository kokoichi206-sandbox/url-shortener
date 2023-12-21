package util

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

// By default, jaeger is listening at "127.0.0.1:5775".
func NewJaegerTracer(host, port, service string) (opentracing.Tracer, io.Closer, error) {
	if host == "" || port == "" {
		return nil, nil, fmt.Errorf("host or port is empty")
	}

	addr := net.JoinHostPort(host, port)

	cfg := config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: time.Second,
			LocalAgentHostPort:  addr,
		},
	}

	//nolint:wrapcheck
	return cfg.NewTracer()
}
