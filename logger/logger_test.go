package logger

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nekomeowww/xo"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestWith(t *testing.T) {
	t.Parallel()
	t.Run("Debug", func(t *testing.T) {
		t.Parallel()

		logger, err := NewLogger(
			WithLevel(zapcore.DebugLevel),
			WithNamespace("xo/logger"),
			WithLogFilePath(xo.RelativePathOf("./logs/test.debug.log")),
			WithCallFrameSkip(1),
		)
		require.NoError(t, err)
		require.NotNil(t, logger)

		newLogger := logger.With(zap.String("some_test_field", "some_test_value"))
		require.NotNil(t, newLogger)

		logger.Debug("debug message")
		newLogger.Debug("debug message with with")
	})
	t.Run("Info", func(t *testing.T) {
		t.Parallel()

		logger, err := NewLogger(
			WithLevel(zapcore.DebugLevel),
			WithNamespace("xo/logger"),
			WithLogFilePath(xo.RelativePathOf("./logs/test.info.log")),
			WithCallFrameSkip(1),
		)
		require.NoError(t, err)
		require.NotNil(t, logger)

		newLogger := logger.With(zap.String("some_test_field", "some_test_value"))
		require.NotNil(t, newLogger)

		logger.Info("info message")
		newLogger.Info("info message with with")
	})
	t.Run("Warn", func(t *testing.T) {
		t.Parallel()

		logger, err := NewLogger(
			WithLevel(zapcore.DebugLevel),
			WithNamespace("xo/logger"),
			WithLogFilePath(xo.RelativePathOf("./logs/test.warn.log")),
			WithCallFrameSkip(1),
		)
		require.NoError(t, err)
		require.NotNil(t, logger)

		newLogger := logger.With(zap.String("some_test_field", "some_test_value"))
		require.NotNil(t, newLogger)

		logger.Warn("warn message")
		newLogger.Warn("warn message with with")
	})
	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		logger, err := NewLogger(
			WithLevel(zapcore.DebugLevel),
			WithNamespace("xo/logger"),
			WithLogFilePath(xo.RelativePathOf("./logs/test.error.log")),
			WithCallFrameSkip(1),
		)
		require.NoError(t, err)
		require.NotNil(t, logger)

		newLogger := logger.With(zap.String("some_test_field", "some_test_value"))
		require.NotNil(t, newLogger)

		logger.Error("error message")
		newLogger.Error("error message with with")
	})
}

func TestFormat(t *testing.T) {
	t.Run("WithLogFilePath", func(t *testing.T) {
		logger, err := NewLogger(
			WithFormat(FormatJSON),
			WithLogFilePath(xo.RelativePathOf("./logs/test.json.log")),
		)
		require.NoError(t, err)
		require.NotNil(t, logger)

		logger.Info("info message", zap.String("some_test_field", "some_test_value"))
	})

	t.Run("WithoutLogFilePath", func(t *testing.T) {
		logger, err := NewLogger(
			WithFormat(FormatJSON),
		)
		require.NoError(t, err)
		require.NotNil(t, logger)

		logger.Info("info message", zap.String("some_test_field", "some_test_value"))
	})
}

func TestSimpleSpan(t *testing.T) {
	spanExporter, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	tp := trace.NewTracerProvider(
		trace.WithBatcher(spanExporter),
		trace.WithSampler(trace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)

	ctx, span := tp.Tracer("test").Start(context.Background(), "test-span")
	span.SetAttributes(attribute.String("key", "value"))
	span.End()

	tp.ForceFlush(ctx)
}

func TestOpenTelemetryAndContextual(t *testing.T) {
	spanRecorder := tracetest.NewSpanRecorder()
	tracerProvider := trace.NewTracerProvider(
		trace.WithSpanProcessor(spanRecorder),
		trace.WithSampler(trace.AlwaysSample()),
	)
	otel.SetTracerProvider(tracerProvider)

	err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	require.NoError(t, err)

	t.Parallel()
	t.Run("Debug", func(t *testing.T) {
		t.Parallel()

		tracer := otel.Tracer("test-tracer")

		ctx, span := tracer.Start(context.Background(), "test-span")
		defer span.End()

		// Add some attributes to the span to make it easier to identify
		span.SetAttributes(attribute.String("test", "OpenTelemetryAndContextual"))

		logger, err := NewLogger(
			WithLevel(zapcore.DebugLevel),
			WithNamespace("xo/logger"),
			WithLogFilePath(xo.RelativePathOf("./logs/test.debug.log")),
			WithCallFrameSkip(1),
		)
		require.NoError(t, err)
		require.NotNil(t, logger)

		newLogger := logger.With(zap.String("some_test_field", "some_test_value"))
		require.NotNil(t, newLogger)

		var uuid *uuid.UUID

		logger.DebugContext(ctx, "debug message", zap.Stringer("uuid", uuid))
		newLogger.DebugContext(ctx, "debug message with with")

		// End the span and force flush
		span.End()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = tracerProvider.ForceFlush(ctx)
		require.NoError(t, err, "Failed to flush trace provider")

		// Retrieve and inspect the recorded spans
		spans := spanRecorder.Ended()
		require.NotEmpty(t, spans, "Expected at least one span")

		for _, recordedSpan := range spans {
			t.Logf("Span Name: %s", recordedSpan.Name())
			t.Logf("Trace ID: %s", recordedSpan.SpanContext().TraceID())
			t.Logf("Span ID: %s", recordedSpan.SpanContext().SpanID())

			for _, event := range recordedSpan.Events() {
				t.Logf("Event: %s", event.Name)

				for _, attr := range event.Attributes {
					t.Logf("  Attribute: %s = %v", attr.Key, attr.Value)
				}
			}
		}
	})
}
