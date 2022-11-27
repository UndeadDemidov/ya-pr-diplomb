package pkg

import (
	"io"

	"github.com/UndeadDemidov/ya-pr-diplomb/config"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

// InitJaeger Jaeger.
func InitJaeger(cfg *config.App) (opentracing.Tracer, io.Closer, error) { //nolint:ireturn
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: cfg.Jaeger.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           cfg.Jaeger.LogSpans,
			LocalAgentHostPort: cfg.Jaeger.Host,
		},
	}

	return jaegerCfgInstance.NewTracer( //nolint:wrapcheck
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
}
