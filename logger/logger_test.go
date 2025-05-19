package logger

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestLogger(t *testing.T) {
	loggerOpt := LoggerOption{
		OtlpEndpoint: "log.otl.signoz.sembuh.ai",
		Environment:  "development",
	}

	lp, err := InitLoggerProvider(loggerOpt)
	assert.NoError(t, err)
	assert.NotNil(t, lp)
	defer func() {
		err := lp.Shutdown(context.Background())
		assert.NoError(t, err)
	}()

	logger := InitZapLogger("usecase/payment", lp)
	defer logger.Sync()

	t.Run("Test Info", func(t *testing.T) {
		logger.Info("This is an info message")
	})

	t.Run("Test Error", func(t *testing.T) {
		logger.Error("This is an error message")
	})

	t.Run("Test Debug", func(t *testing.T) {
		logger.Debug("This is a debug message")
	})

	t.Run("Test Panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("This is a panic message")
			}
		}()
		panic("This is a panic message")
	})

	t.Run("Test With", func(t *testing.T) {
		logger.With(zap.String("key", "value")).Info("This is an info message with key-value pair")
	})

	t.Run("Test With Context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "key", "value")
		logger.With(zap.String("key", ctx.Value("key").(string))).Info("This is an info message with context")
	})

	t.Run("Test With Fields", func(t *testing.T) {
		fields := []zap.Field{
			zap.String("key1", "value1"),
			zap.Int("key2", 2),
		}
		logger.With(fields...).Info("This is an info message with fields")
	})
}
