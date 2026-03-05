package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Level represents log severity levels
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var levelNames = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

// Logger is a structured logger with level-based output
type Logger struct {
	logger     *log.Logger
	mu         sync.Mutex
	level      Level
	output     io.Writer
	withCaller bool
	debugMode  bool // when true, stack traces are included on Error/Fatal
	prefix     string
}

// Global logger instance
var defaultLogger *Logger
var once sync.Once

// DefaultConfig returns default logger configuration
func DefaultConfig() *Logger {
	return &Logger{
		logger:     log.New(os.Stdout, "", 0),
		level:      InfoLevel,
		output:     os.Stdout,
		withCaller: false,
		debugMode:  false,
	}
}

// New creates a new logger with custom configuration
func New(output io.Writer, level Level, withCaller bool) *Logger {
	return &Logger{
		logger:     log.New(output, "", 0),
		level:      level,
		output:     output,
		withCaller: withCaller,
	}
}

// Default returns the global logger instance (singleton)
func Default() *Logger {
	once.Do(func() {
		defaultLogger = &Logger{
			logger:     log.New(os.Stdout, "", 0),
			level:      InfoLevel,
			output:     os.Stdout,
			withCaller: false,
			debugMode:  false,
		}
	})
	return defaultLogger
}

// Initialize sets up the global logger with the given level string and debug flag.
// This should be called once at application startup.
func Initialize(levelStr string, debug bool) *Logger {
	l := Default()
	l.mu.Lock()
	defer l.mu.Unlock()

	l.level = ParseLevel(levelStr)
	l.debugMode = debug

	if debug {
		l.level = DebugLevel
		l.withCaller = true
	}

	return l
}

// ParseLevel converts a string level name to a Level constant.
// Returns InfoLevel for unrecognized strings.
func ParseLevel(s string) Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetDebugMode enables or disables debug mode (stack traces on errors)
func (l *Logger) SetDebugMode(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.debugMode = enabled
}

// IsDebug returns true if the logger is in debug mode
func (l *Logger) IsDebug() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.debugMode
}

// SetOutput changes the output writer for the logger
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
	l.logger = log.New(w, "", 0)
}

// WithPrefix returns a new Logger that prepends a fixed prefix to all messages.
// Useful for request-scoped logging (e.g. request ID).
func (l *Logger) WithPrefix(prefix string) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	return &Logger{
		logger:     l.logger,
		level:      l.level,
		output:     l.output,
		withCaller: l.withCaller,
		debugMode:  l.debugMode,
		prefix:     prefix,
	}
}

// WithRequestID returns a prefixed logger with the given request ID
func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithPrefix(fmt.Sprintf("[req:%s]", requestID))
}

// logWithSkip does the actual logging and accepts an explicit caller skip
func (l *Logger) logWithSkip(level Level, skip int, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var levelName string
	if int(level) >= 0 && int(level) < len(levelNames) {
		levelName = levelNames[level]
	} else {
		levelName = "UNKNOWN"
	}

	message := fmt.Sprintf(format, args...)

	// Prepend prefix if set (e.g. request ID)
	if l.prefix != "" {
		message = l.prefix + " " + message
	}

	if l.withCaller {
		_, file, line, ok := runtime.Caller(skip)
		if ok {
			file = filepath.Base(file)
			l.logger.Printf("[%s] %s [%s:%d] %s", timestamp, levelName, file, line, message)
		} else {
			l.logger.Printf("[%s] %s %s", timestamp, levelName, message)
		}
	} else {
		l.logger.Printf("[%s] %s %s", timestamp, levelName, message)
	}

	// In debug mode, append stack trace for Error and Fatal levels
	if l.debugMode && level >= ErrorLevel {
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, false)
		l.logger.Printf("[%s] %s STACK TRACE:\n%s", timestamp, levelName, string(buf[:n]))
	}
}

// log is preserved for compatibility and calls logWithSkip with a default skip
func (l *Logger) log(level Level, format string, args ...interface{}) {
	l.logWithSkip(level, 2, format, args...)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DebugLevel, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(InfoLevel, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WarnLevel, format, args...)
}

// Error logs an error message (with stack trace in debug mode)
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

// Fatal logs a fatal message (with stack trace in debug mode) and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FatalLevel, format, args...)
	os.Exit(1)
}

// WithCaller enables caller information in log output
func (l *Logger) WithCaller() *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.withCaller = true
	return l
}

// --- Package-level convenience functions using default logger ---

func Debug(format string, args ...interface{}) {
	Default().logWithSkip(DebugLevel, 3, format, args...)
}

func Info(format string, args ...interface{}) {
	Default().logWithSkip(InfoLevel, 3, format, args...)
}

func Warn(format string, args ...interface{}) {
	Default().logWithSkip(WarnLevel, 3, format, args...)
}

func Error(format string, args ...interface{}) {
	Default().logWithSkip(ErrorLevel, 3, format, args...)
}

func Fatal(format string, args ...interface{}) {
	Default().logWithSkip(FatalLevel, 3, format, args...)
	os.Exit(1)
}

// IsDebugMode returns whether the global logger is in debug mode
func IsDebugMode() bool {
	return Default().IsDebug()
}
