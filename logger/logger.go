package logger

import (
	"context"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
)

func InitLoggerProvider(o LoggerOption) (*log.LoggerProvider, error) {
	// Set up the OTLP HTTP exporter
	ctx := context.Background()
	exporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint(o.OtlpEndpoint),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// Create a new logger provider with the exporter
	processor := log.NewBatchProcessor(exporter)

	provider := log.NewLoggerProvider(
		log.WithProcessor(processor),
	)

	return provider, nil
}

func InitZapLogger(scopeName string, provider *log.LoggerProvider) *zap.Logger {
	logger := zap.New(otelzap.NewCore(scopeName, otelzap.WithLoggerProvider(provider)))
	return logger
}
