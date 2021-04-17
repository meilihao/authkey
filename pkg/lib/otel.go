// from https://github.com/open-telemetry/opentelemetry-go/blob/master/example/otel-collector/main.go
// see [opentelemetry-java/QUICKSTART.md](https://github.com/open-telemetry/opentelemetry-java/blob/master/QUICKSTART.md)
// [Documentation / Go / Getting Started](https://opentelemetry.io/docs/go/getting-started/)
package lib

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type (
	LoggerKey struct{}
)

var (
	_spanLogger *zap.Logger
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func InitOTEL(endpoint, serviceName string, logger, spanLogger *zap.Logger) (func(), error) {
	_spanLogger = spanLogger
	if endpoint == "" {
		return func() {}, nil
	}

	ctx := context.Background()

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns
	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(endpoint),
		otlpgrpc.WithDialOption(grpc.WithBlock()), // useful for testing
	)
	exp, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create exporter")
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create resource")
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	cont := controller.New(
		processor.New(
			simple.NewWithExactDistribution(),
			exp,
		),
		controller.WithPusher(exp),
		controller.WithCollectPeriod(2*time.Second),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})
	// set global TracerProvider (the default is noopTracerProvider).
	otel.SetTracerProvider(tracerProvider)
	global.SetMeterProvider(cont.MeterProvider())
	if err = cont.Start(context.Background()); err != nil {
		return nil, errors.Wrap(err, "failed to start controller")
	}

	return func() {
		// Shutdown will flush any remaining spans.
		if err := tracerProvider.Shutdown(ctx); err != nil {
			logger.Error(err.Error(), zap.String("reason", "failed to shutdown provider"))
		}

		// Push any last metric events to the exporter.
		if err := cont.Stop(context.Background()); err != nil {
			logger.Error(err.Error(), zap.String("reason", "failed to stop exporter"))
		}
	}, nil
}

func SpanLog(ctx context.Context, span trace.Span, l zapcore.Level, msg string, kv ...attribute.KeyValue) {
	//var logger *zap.Logger
	// if tmp := ctx.Value(LoggerKey{}); tmp == nil { // 不使用该方式, 因为代码实现不美观且会导致在gin handler context中注入LoggerKey{}前的gin middleware DebugReq()无法使用SpanLog()
	if _spanLogger == nil {
		span.AddEvent(msg, trace.WithAttributes(kv...))

		return
	}

	if ce := _spanLogger.Check(l, msg); ce != nil {
		sctx := span.SpanContext()

		var fs []zap.Field
		if sctx.IsValid() {
			fs = make([]zap.Field, 0, len(kv)+2)
			fs = append(fs, zap.String("trace_id", sctx.TraceID.String()))
			fs = append(fs, zap.String("span_id", sctx.SpanID.String()))
		} else {
			fs = make([]zap.Field, 0, len(kv))
		}

		if len(kv) > 0 {
			for _, attr := range kv {
				switch attr.Value.Type() {
				case attribute.STRING:
					fs = append(fs, zap.String(string(attr.Key), attr.Value.AsString()))
				case attribute.INT64:
					fs = append(fs, zap.Int64(string(attr.Key), attr.Value.AsInt64()))
				case attribute.BOOL:
					fs = append(fs, zap.Bool(string(attr.Key), attr.Value.AsBool()))
				case attribute.FLOAT64:
					fs = append(fs, zap.Float64(string(attr.Key), attr.Value.AsFloat64()))
				default:
					fs = append(fs, zap.Any(string(attr.Key), attr.Value))
				}
			}
		}

		ce.Write(fs...)

		kv = append(kv, attribute.String("level", l.String()))
		span.AddEvent(msg, trace.WithAttributes(kv...))
	}
}
