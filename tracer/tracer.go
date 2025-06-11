package tracer

import "C"
import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// InitTracer configures an OpenTelemetry exporter and trace provider
func InitTracer(o TracerOption) (*sdktrace.TracerProvider, error) {
	headers := map[string]string{
		"signoz-access-token": o.SignozToken,
	}

	//secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")) // config can be passed to configure TLS
	//if o.InsecureMode {
	//	secureOption = otlptracegrpc.WithInsecure()
	//}

	otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(o.CollectorURL),
		otlptracehttp.WithHeaders(headers),
	)
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(o.CollectorURL),
			otlptracehttp.WithHeaders(headers),
		),
		//otlptracegrpc.NewClient(
		//	secureOption,
		//	otlptracegrpc.WithEndpoint(o.CollectorURL),
		//	otlptracegrpc.WithHeaders(headers),
		//),
	)
	if err != nil {
		return nil, err
	}

	// Set up the resource attributes
	resourceAttributes := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(o.ServiceName),
		semconv.DeploymentEnvironmentKey.String(o.Environment),
		semconv.TelemetrySDKLanguageGo,
		semconv.TelemetrySDKVersionKey.String("0.1.0"),
	)

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
		sdktrace.WithResource(resourceAttributes),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return traceProvider, nil
}
