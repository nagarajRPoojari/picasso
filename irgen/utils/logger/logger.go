package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

const (
	colorReset = "\033[0m"
	colorDebug = "\033[36m" // Cyan
	colorInfo  = "\033[32m" // Green
	colorWarn  = "\033[33m" // Yellow
	colorError = "\033[31m" // Red
	colorFatal = "\033[35m" // Magenta
)

// getColor returns the ANSI color code for a log level
func getColor(level LogLevel) string {
	switch level {
	case DEBUG:
		return colorDebug
	case INFO:
		return colorInfo
	case WARN:
		return colorWarn
	case ERROR:
		return colorError
	case FATAL:
		return colorFatal
	default:
		return colorReset
	}
}

// Logger represents a logger instance
type Logger struct {
	level      LogLevel
	mu         sync.Mutex
	out        io.Writer
	errOut     io.Writer
	timeFormat string
	useColor   bool
}

// Global default logger
var defaultLogger *Logger

func init() {
	defaultLogger = New()
}

// New creates a new Logger with default settings
func New() *Logger {
	level := INFO
	if envLevel := os.Getenv("PICASSO_LOG_LEVEL"); envLevel != "" {
		switch envLevel {
		case "DEBUG":
			level = DEBUG
		case "INFO":
			level = INFO
		case "WARN":
			level = WARN
		case "ERROR":
			level = ERROR
		case "FATAL":
			level = FATAL
		}
	}

	return &Logger{
		level:      level,
		out:        os.Stdout,
		errOut:     os.Stderr,
		timeFormat: "2006-01-02 15:04:05",
		useColor:   true,
	}
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetOutput sets the output writer for INFO and DEBUG messages
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

// SetErrorOutput sets the output writer for ERROR and FATAL messages
func (l *Logger) SetErrorOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.errOut = w
}

// SetTimeFormat sets the time format for log messages
func (l *Logger) SetTimeFormat(format string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timeFormat = format
}

// SetUseColor enables or disables colored output
func (l *Logger) SetUseColor(useColor bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.useColor = useColor
}

// log writes a log message at the specified level
func (l *Logger) log(module string, level LogLevel, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	timestamp := time.Now().Format(l.timeFormat)

	out := l.out
	if level >= ERROR {
		out = l.errOut
	}

	message := fmt.Sprintf(format, args...)

	if l.useColor {
		// Build module part if provided
		moduleStr := ""
		if module != "" {
			moduleStr = fmt.Sprintf(" [%s]", module)
		}

		fmt.Fprintf(out, "%s[%s] [%s]%s%s%s\n",
			getColor(level),
			timestamp,
			level.String(),
			moduleStr,
			colorReset,
			message)
	} else {
		// Build module part if provided
		moduleStr := ""
		if module != "" {
			moduleStr = fmt.Sprintf(" [%s]", module)
		}

		fmt.Fprintf(out, "[%s] [%s]%s %s\n",
			timestamp,
			level.String(),
			moduleStr,
			message)
	}
}

func (l *Logger) Debug(module string, format string, args ...interface{}) {
	l.log(module, DEBUG, format, args...)
}

func (l *Logger) Info(module string, format string, args ...interface{}) {
	l.log(module, INFO, format, args...)
}

func (l *Logger) Warn(module string, format string, args ...interface{}) {
	l.log(module, WARN, format, args...)
}

func (l *Logger) Error(module string, format string, args ...interface{}) {
	l.log(module, ERROR, format, args...)
}

func (l *Logger) Fatal(module string, format string, args ...interface{}) {
	l.log(module, FATAL, format, args...)
	os.Exit(1)
}

func (l *Logger) Log(module string, level LogLevel, format string, args ...interface{}) {
	l.log(module, level, format, args...)
}

func SetLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

func SetOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

func SetErrorOutput(w io.Writer) {
	defaultLogger.SetErrorOutput(w)
}

func SetTimeFormat(format string) {
	defaultLogger.SetTimeFormat(format)
}

func SetUseColor(useColor bool) {
	defaultLogger.SetUseColor(useColor)
}

func Debug(module string, format string, args ...interface{}) {
	defaultLogger.Debug(module, format, args...)
}

func Info(module string, format string, args ...interface{}) {
	defaultLogger.Info(module, format, args...)
}

func Warn(module string, format string, args ...interface{}) {
	defaultLogger.Warn(module, format, args...)
}

func Error(module string, format string, args ...interface{}) {
	defaultLogger.Error(module, format, args...)
}

func Fatal(module string, format string, args ...interface{}) {
	defaultLogger.Fatal(module, format, args...)
}

func Log(module string, level LogLevel, format string, args ...interface{}) {
	defaultLogger.Log(module, level, format, args...)
}
