package logger

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"time"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
)

// LogPrettyFormatter defines the format for log file.
type LogPrettyFormatter struct {
	logrus.TextFormatter
	MinimumCallerDepth int
}

// NewLogFileFormatter return the log format for log file.
//
// eg: 2023-06-01T12:00:00 [info] [controllers/some_controller/code_file.go:99] foo key=value
func NewLogPrettyFormatter() *LogPrettyFormatter {
	return &LogPrettyFormatter{
		TextFormatter: logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			FullTimestamp:   true,
		},
		MinimumCallerDepth: 0,
	}
}

// Format renders a single log entry for log file
//
// the original file log format is defined here: github.com/sirupsen/logrus/text_formatter.TextFormatter{}.Format().
func (f *LogPrettyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(map[string]any)
	for k, v := range entry.Data {
		data[k] = v
	}

	keys := make([]string, 0, len(data))

	for k := range data {
		if k == "caller_file" {
			continue
		}

		keys = append(keys, k)
	}

	if !f.DisableSorting {
		if nil != f.SortingFunc {
			f.SortingFunc(keys)
		} else {
			sort.Strings(keys)
		}
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	prefixStr := entry.Time.Format(timestampFormat) + " "
	var renderFunc func(a ...any) string

	switch entry.Level {
	case logrus.TraceLevel:
		renderFunc = color.FgGray.Render
	case logrus.DebugLevel:
		renderFunc = color.FgGreen.Render
	case logrus.InfoLevel:
		renderFunc = color.FgCyan.Render
	case logrus.WarnLevel:
		renderFunc = color.FgYellow.Render
	case logrus.ErrorLevel:
		renderFunc = color.FgRed.Render
	case logrus.FatalLevel:
		renderFunc = color.FgMagenta.Render
	case logrus.PanicLevel:
		renderFunc = color.FgMagenta.Render
	default:
		renderFunc = color.FgGray.Render
	}

	prefixStr += renderFunc("[", entry.Level.String(), "]")

	b.WriteString(prefixStr)

	if data["caller_file"] != nil {
		fmt.Fprintf(b, " [%s]", data["caller_file"])
		delete(data, "file")
	} else if entry.Context != nil {
		caller, _ := entry.Context.Value(runtimeCaller).(*runtime.Frame)
		if caller != nil {
			fmt.Fprintf(b, " [%s:%d]", caller.File, caller.Line)
		}
	}

	if entry.Message != "" {
		b.WriteString(" " + entry.Message)
	}

	for _, key := range keys {
		value := data[key]
		appendKeyValue(b, key, value, f.QuoteEmptyFields)
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

// appendKeyValue append value with key to data that to be appended to log file.
func appendKeyValue(b *bytes.Buffer, key string, value interface{}, QuoteEmptyFields bool) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}

	b.WriteString(key)
	b.WriteByte('=')
	appendValue(b, value, QuoteEmptyFields)
}

// appendValue append value to data used for method appendKeyValue.
func appendValue(b *bytes.Buffer, value interface{}, QuoteEmptyFields bool) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !needsQuoting(stringVal, QuoteEmptyFields) {
		b.WriteString(stringVal)
	} else {
		fmt.Fprintf(b, "%q", stringVal)
	}
}

// needsQuoting check where text needs to be quoted.
func needsQuoting(text string, quoteEmptyFields bool) bool {
	if quoteEmptyFields && len(text) == 0 {
		return true
	}

	for _, ch := range text {
		if (ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+' {
			continue
		}

		return true
	}

	return false
}
