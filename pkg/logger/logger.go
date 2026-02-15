package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
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
}

// Global logger instance
var defaultLogger *Logger
var once sync.Once

// DefaultConfig returns default logger configuration
func DefaultConfig() *Logger {
	l := &Logger{
		logger:     log.New(os.Stdout, "", 0),
		level:      InfoLevel,
		output:     os.Stdout,
		withCaller: false,
	}
	return l
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
		}
	})
	return defaultLogger
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// log writes a log entry with the given level and message
// logWithSkip does the actual logging and accepts an explicit caller skip
func (l *Logger) logWithSkip(level Level, skip int, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// check level under lock to avoid races with SetLevel
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// bounds-check level name
	var levelName string
	if int(level) >= 0 && int(level) < len(levelNames) {
		levelName = levelNames[level]
	} else {
		levelName = "UNKNOWN"
	}

	message := fmt.Sprintf(format, args...)

	if l.withCaller {
		_, file, line, ok := runtime.Caller(skip)
		if ok {
			file = filepath.Base(file)
			l.logger.Printf("[%s] %s [%s:%d] %s", timestamp, levelName, file, line, message)
			return
		}
	}

	l.logger.Printf("[%s] %s %s", timestamp, levelName, message)
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

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

// Fatal logs a fatal message and exits
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

// Package-level convenience functions using default logger
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
}
