package logger

import (
	"testing"

	"github.com/nekomeowww/xo"
	"github.com/stretchr/testify/require"
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
