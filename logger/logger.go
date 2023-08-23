// Package logger 日志包，用于日志输出和打印
package logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

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
}

// Debug 打印 debug 级别日志。
func (l *Logger) Debug(msg string, fields ...zapcore.Field) {
	l.ZapLogger.Debug(msg, fields...)

	data := make(map[string]any)
	for k, v := range l.LogrusLogger.Data {
		data[k] = v
	}

	entry := logrus.NewEntry(l.LogrusLogger.Logger)
	SetCallFrame(entry, l.namespace, 1)

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
	SetCallFrame(entry, l.namespace, 1)

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
	SetCallFrame(entry, l.namespace, 1)

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
	SetCallFrame(entry, l.namespace, 1)

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
	SetCallFrame(entry, l.namespace, 1)

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
	SetCallFrame(entry, l.namespace, 1)

	for k, v := range data {
		entry = entry.WithField(k, v)
	}

	for _, v := range fields {
		entry = entry.WithField(v.Key, ZapField(v).MatchValue())
	}

	return &Logger{
		ZapLogger:    newZapLogger,
		LogrusLogger: entry,
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

// NewLogger 按需创建 logger 实例。
func NewLogger(level zapcore.Level, namespace string, logFilePath string, hook []logrus.Hook) (*Logger, error) {
	var err error
	if logFilePath != "" {
		err = autoCreateLogFile(logFilePath)
		if err != nil {
			return nil, err
		}
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.InitialFields = map[string]interface{}{
		"app_name": namespace,
	}

	if logFilePath != "" {
		config.OutputPaths = []string{logFilePath}
		config.ErrorOutputPaths = []string{logFilePath}
	}

	zapLogger, err := config.Build(zap.WithCaller(true))
	if err != nil {
		return nil, err
	}

	logrusLogger := logrus.New()
	if len(hook) > 0 {
		for _, h := range hook {
			logrusLogger.Hooks.Add(h)
		}
	}

	logrusLogger.SetFormatter(NewLogFileFormatter())
	logrusLogger.SetReportCaller(true)
	logrusLogger.Level = zapCoreLevelToLogrusLevel(level)

	l := &Logger{
		LogrusLogger: logrus.NewEntry(logrusLogger),
		ZapLogger:    zapLogger,
		namespace:    namespace,
	}

	l.Debug("logger init successfully for both logrus and zap",
		zap.String("log_file_path", logFilePath),
		zap.String("log_level", level.String()),
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
