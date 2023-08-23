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

		logger, err := NewLogger(zapcore.DebugLevel, "xo/logger", xo.RelativePathOf("./logs/test.debug.log"), nil)
		require.NoError(t, err)
		require.NotNil(t, logger)

		newLogger := logger.With(zap.String("some_test_field", "some_test_value"))
		require.NotNil(t, newLogger)

		logger.Debug("debug message")
		newLogger.Debug("debug message")
	})
	t.Run("Info", func(t *testing.T) {
		t.Parallel()

		logger, err := NewLogger(zapcore.DebugLevel, "xo/logger", xo.RelativePathOf("./logs/test.info.log"), nil)
		require.NoError(t, err)
		require.NotNil(t, logger)

		newLogger := logger.With(zap.String("some_test_field", "some_test_value"))
		require.NotNil(t, newLogger)

		logger.Info("debug message")
		newLogger.Info("debug message")
	})
	t.Run("Warn", func(t *testing.T) {
		t.Parallel()

		logger, err := NewLogger(zapcore.DebugLevel, "xo/logger", xo.RelativePathOf("./logs/test.warn.log"), nil)
		require.NoError(t, err)
		require.NotNil(t, logger)

		newLogger := logger.With(zap.String("some_test_field", "some_test_value"))
		require.NotNil(t, newLogger)

		logger.Warn("debug message")
		newLogger.Warn("debug message")
	})
	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		logger, err := NewLogger(zapcore.DebugLevel, "xo/logger", xo.RelativePathOf("./logs/test.error.log"), nil)
		require.NoError(t, err)
		require.NotNil(t, logger)

		newLogger := logger.With(zap.String("some_test_field", "some_test_value"))
		require.NotNil(t, newLogger)

		logger.Error("debug message")
		newLogger.Error("debug message")
	})
}
