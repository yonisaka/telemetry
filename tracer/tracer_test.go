package tracer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTracer(t *testing.T) {
	tracerOpt := TracerOption{
		CollectorURL: "127.0.0.1:4318",
		InsecureMode: false,
		ServiceName:  "test-app",
		Environment:  "development",
	}
	tp, err := InitTracer(tracerOpt)
	assert.NoError(t, err)
	assert.NotNil(t, tp)
	defer func() {
		err := tp.Shutdown(context.Background())
		assert.NoError(t, err)
	}()

	ctx := context.Background()
	ctx, span := StartSpan(ctx, "test-span")
	assert.NotNil(t, ctx)
	assert.NotNil(t, span)
	defer span.Finish()

	t.Run("Trace ID", func(t *testing.T) {
		assert.NotEmpty(t, span.TraceId())
	})

	t.Run("Trace Events", func(t *testing.T) {
		span.AddEvents("test-event")
	})

	t.Run("Trace Tags String", func(t *testing.T) {
		span.AddTags(SpanTagString("test-key-string", "test-value"))
	})

	t.Run("Trace Tags Int64", func(t *testing.T) {
		span.AddTags(SpanTagInt64("test-key-int", 100))
	})

	t.Run("Trace Tags Bool", func(t *testing.T) {
		span.AddTags(SpanTagBool("test-key-bool", true))
	})

	t.Run("Trace Tags Float64", func(t *testing.T) {
		span.AddTags(SpanTagFloat64("test-key-float", 100.0))
	})

	t.Run("Trace Tags String Slice", func(t *testing.T) {
		span.AddTags(SpanTagStringSlice("test-key-string-slice", []string{"test-value"}))
	})

	t.Run("Trace Tags Int", func(t *testing.T) {
		span.AddTags(SpanTagInt("test-key-int", 100))
	})

	t.Run("Trace Tags Any", func(t *testing.T) {
		span.AddAnyTags(map[string]interface{}{
			"test-key-string": "test-value",
			"test-key-int":    100,
			"test-key-bool":   true,
			"test-key-float":  100.0,
		})
	})

	t.Run("Trace Error", func(t *testing.T) {
		span.AddError(assert.AnError)
	})
}
