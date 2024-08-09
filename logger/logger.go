// Package logger 日志包，用于日志输出和打印
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
	"github.com/sirupsen/logrus"
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

	namespace string
	skip      int
}

// Debug 打印 debug 级别日志。
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

// Info 打印 info 级别日志。
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

// Warn 打印 warn 级别日志。
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

// Error 打印错误日志。
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

// Fatal 打印致命错误日志，打印后立即退出程序。
//
// NOTICE: 该方法会调用 os.Exit(1) 退出程序。而且会优先执行 logrus 的 Fatal 方法，然后再执行 zap 的 Fatal 方法。
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

// With 创建一个新的 logger 实例，该实例会继承当前 logger 的上下文信息。
// 添加到新 logger 实例的字段，不会影响当前 logger 实例。
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
		ZapLogger:    newZapLogger,
		LogrusLogger: entry,
		namespace:    l.namespace,
		skip:         l.skip,
	}
}

// With 创建一个新的 logger 实例，该实例会继承当前 logger 的上下文信息。
// 添加到新 logger 实例的字段，不会影响当前 logger 实例。
func (l *Logger) WithAndSkip(skip int, fields ...zapcore.Field) *Logger {
	newZapLogger := l.ZapLogger.With(fields...)

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
		ZapLogger:    newZapLogger,
		LogrusLogger: entry,
		namespace:    l.namespace,
		skip:         skip,
	}
}

// SetCallFrame 设定调用栈。
func SetCallFrame(entry *logrus.Entry, namespace string, skip int) {
	// 获取调用栈的 文件、行号
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

// SetCallerFrameWithFileAndLine 设定调用栈。
func SetCallerFrameWithFileAndLine(entry *logrus.Entry, namespace, functionName, file string, line int) {
	splitTarget := filepath.FromSlash("/" + namespace + "/")
	// 拆解文件名，移除项目所在路径和项目名称，只保留到项目内的文件路径
	filename := strings.SplitN(file, splitTarget, 2)
	// 如果拆解后出现问题，回退到完整路径
	if len(filename) < 2 {
		filename = []string{"", file}
	}

	// 设定 logrus.Entry 的上下文信息
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
		return logLevel, fmt.Errorf("log level " + logLevelStr + " in environment variable LOG_LEVEL is invalid, fallbacks to default level: info")
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
	level            zapcore.Level
	logFilePath      string
	hook             []logrus.Hook
	appName          string
	namespace        string
	initialFields    map[string]any
	callFrameSkip    int
	format           Format
	lokiRemoteConfig *loki.Config
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
		EncodeTime:     zapcore.ISO8601TimeEncoder,
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
		LogrusLogger: logrus.NewEntry(logrusLogger),
		ZapLogger:    zapLogger,
		namespace:    opts.namespace,
		skip:         opts.callFrameSkip,
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
