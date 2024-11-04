// Package logger
package logger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/nekomeowww/xo/logger/loki"
	"github.com/nekomeowww/xo/logger/otelzap"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapField zapcore.Field

func (f ZapField) MatchValue() any {
	switch f.Type {
	case zapcore.UnknownType: // checked
		return f.Interface
	case zapcore.ArrayMarshalerType: // checked
		return f.Interface
	case zapcore.ObjectMarshalerType: // checked
		return f.Interface
	case zapcore.BinaryType: // checked
		return f.Interface
	case zapcore.BoolType: // checked
		return f.Integer != 0
	case zapcore.ByteStringType: // checked
		return f.Interface
	case zapcore.Complex128Type: // checked
		return f.Interface
	case zapcore.Complex64Type: // checked
		return f.Interface
	case zapcore.DurationType: // checked
		return time.Duration(f.Integer)
	case zapcore.Float64Type: // checked
		return f.Integer
	case zapcore.Float32Type: // checked
		return f.Integer
	case zapcore.Int64Type: // checked
		return f.Integer
	case zapcore.Int32Type: // checked
		return f.Integer
	case zapcore.Int16Type: // checked
		return f.Integer
	case zapcore.Int8Type: // checked
		return f.Integer
	case zapcore.StringType: // checked
		return f.String
	case zapcore.TimeType: // checked
		return time.Unix(0, f.Integer)
	case zapcore.TimeFullType: // checked
		return f.Interface
	case zapcore.Uint64Type: // checked
		return f.Integer
	case zapcore.Uint32Type: // checked
		return f.Integer
	case zapcore.Uint16Type: // checked
		return f.Integer
	case zapcore.Uint8Type: // checked
		return f.Integer
	case zapcore.UintptrType: // checked
		return f.Integer
	case zapcore.ReflectType: // checked
		return f.Interface
	case zapcore.NamespaceType: // checked
		return ""
	case zapcore.StringerType: // checked
		return f.Interface
	case zapcore.ErrorType: // checked
		return f.Interface
	case zapcore.SkipType: // checked
		return ""
	case zapcore.InlineMarshalerType: // checked
		return f.Interface
	}

	return f.Interface
}

type Logger struct {
	LogrusLogger *logrus.Entry
	ZapLogger    *zap.Logger
	otelTracer   trace.Tracer

	withAppendedFields    []zap.Field
	openTelemetryDisabled bool
	namespace             string
	skip                  int
	errorStatusLevel      zapcore.Level
	caller                bool
	stackTrace            bool
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Debug(msg string, fields ...zapcore.Field) {
	l.ZapLogger.Debug(msg, fields...)

	data := make(map[string]any)
	for k, v := range l.LogrusLogger.Data {
		data[k] = v
	}

	entry := logrus.NewEntry(l.LogrusLogger.Logger)
	SetCallFrame(entry, l.namespace, l.skip)

	for k, v := range data {
		entry = entry.WithField(k, v)
	}

	for _, v := range fields {
		entry = entry.WithField(v.Key, ZapField(v).MatchValue())
	}

	entry.Debug(msg)
}

// DebugContext logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger. Besides that, it
// also logs the message to the OpenTelemetry span.
func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	if !l.openTelemetryDisabled {
		l.span(ctx, zapcore.DebugLevel, msg, fields...)
	}

	l.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	l.ZapLogger.Info(msg, fields...)

	data := make(map[string]any)
	for k, v := range l.LogrusLogger.Data {
		data[k] = v
	}

	entry := logrus.NewEntry(l.LogrusLogger.Logger)
	SetCallFrame(entry, l.namespace, l.skip)

	for k, v := range data {
		entry = entry.WithField(k, v)
	}

	for _, v := range fields {
		entry = entry.WithField(v.Key, ZapField(v).MatchValue())
	}

	entry.Info(msg)
}

// InfoContext logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger. Besides that, it
// also logs the message to the OpenTelemetry span.
func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	if !l.openTelemetryDisabled {
		l.span(ctx, zapcore.InfoLevel, msg, fields...)
	}

	l.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger. Besides that, it
// also logs the message to the OpenTelemetry span.
func (l *Logger) Warn(msg string, fields ...zapcore.Field) {
	l.ZapLogger.Warn(msg, fields...)

	data := make(map[string]any)
	for k, v := range l.LogrusLogger.Data {
		data[k] = v
	}

	entry := logrus.NewEntry(l.LogrusLogger.Logger)
	SetCallFrame(entry, l.namespace, l.skip)

	for k, v := range data {
		entry = entry.WithField(k, v)
	}

	for _, v := range fields {
		entry = entry.WithField(v.Key, ZapField(v).MatchValue())
	}

	entry.Warn(msg)
}

// WarnContext logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger. Besides that, it
// also logs the message to the OpenTelemetry span.
func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	if !l.openTelemetryDisabled {
		l.span(ctx, zapcore.WarnLevel, msg, fields...)
	}

	l.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	l.ZapLogger.Error(msg, fields...)

	data := make(map[string]any)
	for k, v := range l.LogrusLogger.Data {
		data[k] = v
	}

	entry := logrus.NewEntry(l.LogrusLogger.Logger)
	SetCallFrame(entry, l.namespace, l.skip)

	for k, v := range data {
		entry = entry.WithField(k, v)
	}

	for _, v := range fields {
		entry = entry.WithField(v.Key, ZapField(v).MatchValue())
	}

	entry.Error(msg)
}

// ErrorContext logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger. Besides that, it
// also logs the message to the OpenTelemetry span.
func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.span(ctx, zapcore.ErrorLevel, msg, fields...)
	l.Error(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// NOTICE: This method calls os.Exit(1) to exit the program. It also prioritizes the execution of logrus' Fatal method over zap's Fatal method.
func (l *Logger) Fatal(msg string, fields ...zapcore.Field) {
	data := make(map[string]any)
	for k, v := range l.LogrusLogger.Data {
		data[k] = v
	}

	entry := logrus.NewEntry(l.LogrusLogger.Logger)
	SetCallFrame(entry, l.namespace, l.skip)

	for k, v := range data {
		entry = entry.WithField(k, v)
	}

	for _, v := range fields {
		entry = entry.WithField(v.Key, ZapField(v).MatchValue())
	}

	entry.Fatal(msg)
	l.ZapLogger.Fatal(msg, fields...)
}

// FatalContext logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger. Besides that, it
// also logs the message to the OpenTelemetry span.
func (l *Logger) FatalContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	if !l.openTelemetryDisabled {
		l.span(ctx, zapcore.FatalLevel, msg, fields...)
	}

	l.Fatal(msg, fields...)
}

// With creates a new logger instance that inherits the context information from the current logger.
// Fields added to the new logger instance do not affect the current logger instance.
func (l *Logger) With(fields ...zapcore.Field) *Logger {
	newZapLogger := l.ZapLogger.With(fields...)

	data := make(map[string]any)
	for k, v := range l.LogrusLogger.Data {
		data[k] = v
	}

	entry := logrus.NewEntry(l.LogrusLogger.Logger)
	SetCallFrame(entry, l.namespace, l.skip)

	for k, v := range data {
		entry = entry.WithField(k, v)
	}

	for _, v := range fields {
		entry = entry.WithField(v.Key, ZapField(v).MatchValue())
	}

	return &Logger{
		ZapLogger:             newZapLogger,
		withAppendedFields:    append(l.withAppendedFields, fields...),
		LogrusLogger:          entry,
		namespace:             l.namespace,
		skip:                  l.skip,
		openTelemetryDisabled: l.openTelemetryDisabled,
	}
}

// WithAndSkip creates a new logger instance that inherits the context information from the current logger.
// Fields added to the new logger instance do not affect the current logger instance.
// The skip parameter is used to determine the number of stack frames to skip when retrieving the caller information.
func (l *Logger) WithAndSkip(skip int, fields ...zapcore.Field) *Logger {
	newZapLogger := l.ZapLogger.WithOptions(zap.AddCallerSkip(1)).With(fields...)

	data := make(map[string]any)
	for k, v := range l.LogrusLogger.Data {
		data[k] = v
	}

	entry := logrus.NewEntry(l.LogrusLogger.Logger)
	SetCallFrame(entry, l.namespace, skip)

	for k, v := range data {
		entry = entry.WithField(k, v)
	}

	for _, v := range fields {
		entry = entry.WithField(v.Key, ZapField(v).MatchValue())
	}

	return &Logger{
		ZapLogger:             newZapLogger,
		LogrusLogger:          entry,
		withAppendedFields:    append(l.withAppendedFields, fields...),
		namespace:             l.namespace,
		skip:                  skip,
		openTelemetryDisabled: l.openTelemetryDisabled,
	}
}

func (l *Logger) span(ctx context.Context, lvl zapcore.Level, msg string, fields ...zap.Field) {
	span := trace.SpanFromContext(ctx)

	if lvl >= l.errorStatusLevel && span.IsRecording() {
		span.SetStatus(codes.Error, msg)
	}

	attrs := make([]attribute.KeyValue, 0, len(fields)+3)
	attrs = append(attrs, semconv.OtelLibraryName("github.com/nekomeowww/xo/logger"))
	attrs = append(attrs, semconv.OtelLibraryVersion("1.0.0"))
	attrs = append(attrs, attribute.String("log.severity", otelzap.LogSeverityFromZapLevel(lvl).String()))
	attrs = append(attrs, attribute.String("log.message", msg))

	for _, field := range l.withAppendedFields {
		attrs = append(attrs, otelzap.AttributesFromZapField(field)...)
	}

	for _, field := range fields {
		attrs = append(attrs, otelzap.AttributesFromZapField(field)...)
	}

	if l.caller {
		if fn, file, line, ok := runtime.Caller(l.skip + 1); ok {
			fn := runtime.FuncForPC(fn).Name()
			if fn != "" {
				attrs = append(attrs, attribute.String("code.function", fn))
			}
			if file != "" {
				attrs = append(attrs, attribute.String("code.filepath", file))
				attrs = append(attrs, attribute.Int("code.lineno", line))
			}
		}
	}

	if l.stackTrace {
		stackTrace := make([]byte, 2048)
		n := runtime.Stack(stackTrace, false)
		attrs = append(attrs, attribute.String("exception.stacktrace", string(stackTrace[:n])))
	}

	span.AddEvent("log", trace.WithAttributes(attrs...))
}

// SetCallFrame set the caller information for the log entry.
func SetCallFrame(entry *logrus.Entry, namespace string, skip int) {
	_, file, line, _ := runtime.Caller(skip + 1)
	pc, _, _, _ := runtime.Caller(skip + 2)
	funcDetail := runtime.FuncForPC(pc)

	var funcName string
	if funcDetail != nil {
		funcName = funcDetail.Name()
	}

	SetCallerFrameWithFileAndLine(entry, namespace, funcName, file, line)
}

type contextKey string

const (
	runtimeCaller contextKey = "ContextKeyRuntimeCaller"
)

// SetCallerFrameWithFileAndLine set the caller information for the log entry.
func SetCallerFrameWithFileAndLine(entry *logrus.Entry, namespace, functionName, file string, line int) {
	splitTarget := filepath.FromSlash("/" + namespace + "/")

	filename := strings.SplitN(file, splitTarget, 2)
	if len(filename) < 2 {
		filename = []string{"", file}
	}

	entry.Context = context.WithValue(context.Background(), runtimeCaller, &runtime.Frame{
		File:     filename[1],
		Line:     line,
		Function: functionName,
	})
}

func zapCoreLevelToLogrusLevel(level zapcore.Level) logrus.Level {
	switch level {
	case zapcore.DebugLevel:
		return logrus.DebugLevel
	case zapcore.InfoLevel:
		return logrus.InfoLevel
	case zapcore.WarnLevel:
		return logrus.WarnLevel
	case zapcore.ErrorLevel:
		return logrus.ErrorLevel
	case zapcore.FatalLevel:
		return logrus.FatalLevel
	case zapcore.PanicLevel, zapcore.DPanicLevel:
		return logrus.PanicLevel
	case zapcore.InvalidLevel:
		return logrus.InfoLevel
	default:
		return logrus.InfoLevel
	}
}

func ReadLogLevelFromEnv() (zapcore.Level, error) {
	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr == "" {
		return zapcore.InfoLevel, nil
	}

	logLevel, err := zapcore.ParseLevel(logLevelStr)
	if err != nil {
		logLevel = zapcore.InfoLevel
		return logLevel, errors.New("log level " + logLevelStr + " in environment variable LOG_LEVEL is invalid, fallbacks to default level: info")
	}
	if logLevel == zapcore.FatalLevel {
		logLevel = zapcore.InfoLevel
		return logLevel, fmt.Errorf("log level fatal in environment variable LOG_LEVEL is invalid, fallbacks to default level: info")
	}

	return logLevel, nil
}

func ReadLogFormatFromEnv() (Format, error) {
	logFormatStr := os.Getenv("LOG_FORMAT")
	if logFormatStr == "" {
		return FormatPretty, nil
	}

	switch logFormatStr {
	case "json":
		return FormatJSON, nil
	case "pretty":
		return FormatPretty, nil
	default:
		return FormatPretty, fmt.Errorf("log format %s in environment variable LOG_FORMAT is invalid, fallbacks to default format: pretty", logFormatStr)
	}
}

type newLoggerOptions struct {
	level                 zapcore.Level
	logFilePath           string
	hook                  []logrus.Hook
	appName               string
	namespace             string
	initialFields         map[string]any
	callFrameSkip         int
	format                Format
	lokiRemoteConfig      *loki.Config
	openTelemetryDisabled bool
}

type NewLoggerCallOption func(*newLoggerOptions)

func WithLevel(level zapcore.Level) NewLoggerCallOption {
	return func(o *newLoggerOptions) {
		o.level = level
	}
}

type Format string

const (
	FormatJSON   Format = "json"
	FormatPretty Format = "pretty"
)

func WithFormat(format Format) NewLoggerCallOption {
	return func(o *newLoggerOptions) {
		o.format = format
	}
}

func WithLogFilePath(logFilePath string) NewLoggerCallOption {
	return func(o *newLoggerOptions) {
		o.logFilePath = logFilePath
	}
}

func WithHook(hook logrus.Hook) NewLoggerCallOption {
	return func(o *newLoggerOptions) {
		o.hook = append(o.hook, hook)
	}
}

func WithAppName(appName string) NewLoggerCallOption {
	return func(o *newLoggerOptions) {
		o.appName = appName
	}
}

func WithNamespace(namespace string) NewLoggerCallOption {
	return func(o *newLoggerOptions) {
		o.namespace = namespace
	}
}

func WithInitialFields(fields map[string]any) NewLoggerCallOption {
	return func(o *newLoggerOptions) {
		o.initialFields = fields
	}
}

func WithCallFrameSkip(skip int) NewLoggerCallOption {
	return func(o *newLoggerOptions) {
		o.callFrameSkip = skip
	}
}

func WithLokiRemoteConfig(config *loki.Config) NewLoggerCallOption {
	return func(o *newLoggerOptions) {
		o.lokiRemoteConfig = config
	}
}

func WithOpenTelemetryDisabled() NewLoggerCallOption {
	return func(o *newLoggerOptions) {
		o.openTelemetryDisabled = true
	}
}

// NewLogger 按需创建 logger 实例。
func NewLogger(callOpts ...NewLoggerCallOption) (*Logger, error) {
	opts := new(newLoggerOptions)
	opts.callFrameSkip = 2
	opts.format = FormatPretty

	for _, opt := range callOpts {
		opt(opts)
	}

	var err error
	if opts.logFilePath != "" {
		err = autoCreateLogFile(opts.logFilePath)
		if err != nil {
			return nil, err
		}
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(opts.level)
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		MessageKey:     "message",
		CallerKey:      "caller",
		FunctionKey:    "function",
		StacktraceKey:  "stack",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		SkipLineEnding: false,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
	}

	config.InitialFields = make(map[string]any)
	if opts.appName != "" {
		config.InitialFields["app_name"] = opts.appName
	}
	if opts.namespace != "" {
		config.InitialFields["namespace"] = opts.namespace
	}
	if len(opts.initialFields) > 0 {
		for k, v := range opts.initialFields {
			config.InitialFields[k] = v
		}
	}
	if opts.logFilePath != "" {
		config.OutputPaths = []string{opts.logFilePath}
		config.ErrorOutputPaths = []string{opts.logFilePath}

		if opts.format == FormatJSON {
			config.OutputPaths = append(config.OutputPaths, "stdout")
			config.ErrorOutputPaths = append(config.ErrorOutputPaths, "stderr")
		}
	} else {
		config.OutputPaths = []string{}
		config.ErrorOutputPaths = []string{}

		if opts.format == FormatJSON {
			config.OutputPaths = append(config.OutputPaths, "stdout")
			config.ErrorOutputPaths = append(config.ErrorOutputPaths, "stderr")
		}
	}
	if opts.lokiRemoteConfig != nil {
		if opts.appName != "" {
			opts.lokiRemoteConfig.Labels["app_name"] = opts.appName
		}
		if opts.namespace != "" {
			opts.lokiRemoteConfig.Labels["namespace"] = opts.namespace
		}

		loki := loki.New(context.Background(), *opts.lokiRemoteConfig)
		config = loki.ApplyConfig(config)
	}

	zapLogger, err := config.Build(zap.WithCaller(true))
	if err != nil {
		return nil, err
	}

	logrusLogger := logrus.New()

	if len(opts.hook) > 0 {
		for _, h := range opts.hook {
			logrusLogger.Hooks.Add(h)
		}
	}

	if opts.format == FormatPretty {
		logrusLogger.SetFormatter(NewLogPrettyFormatter())
		logrusLogger.SetOutput(os.Stdout)
		logrusLogger.SetReportCaller(true)
		logrusLogger.Level = zapCoreLevelToLogrusLevel(opts.level)
	} else {
		logrusLogger.SetOutput(io.Discard)
	}

	l := &Logger{
		LogrusLogger:          logrus.NewEntry(logrusLogger),
		ZapLogger:             zapLogger.WithOptions(zap.AddCallerSkip(opts.callFrameSkip - 2)),
		namespace:             opts.namespace,
		skip:                  opts.callFrameSkip,
		errorStatusLevel:      zapcore.ErrorLevel,
		caller:                true,
		stackTrace:            false,
		openTelemetryDisabled: opts.openTelemetryDisabled,
	}
	if !opts.openTelemetryDisabled {
		l.otelTracer = otel.Tracer("github.com/nekomeowww/xo/logger")
	}

	l.Debug("logger init successfully for both logrus and zap",
		zap.String("log_file_path", opts.logFilePath),
		zap.String("log_level", opts.level.String()),
	)

	return l, nil
}

func autoCreateLogFile(logFilePathStr string) error {
	if logFilePathStr == "" {
		return nil
	}

	logDir := filepath.Dir(logFilePathStr)

	_, err := os.Stat(logDir)
	if err != nil {
		if os.IsNotExist(err) {
			err2 := os.MkdirAll(logDir, 0755)
			if err2 != nil {
				return fmt.Errorf("failed to create %s log directory: %w", logDir, err)
			}
		} else {
			return err
		}
	}

	stat, err := os.Stat(logFilePathStr)
	if err != nil {
		if os.IsNotExist(err) {
			_, err2 := os.Create(logFilePathStr)
			if err2 != nil {
				return fmt.Errorf("failed to create %s log file: %w", logFilePathStr, err)
			}
		} else {
			return err
		}
	}
	if stat != nil && stat.IsDir() {
		return errors.New("path exists but it is a directory")
	}

	return nil
}
