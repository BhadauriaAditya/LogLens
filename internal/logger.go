package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"
)

// Logger struct (internal state)
type logger struct {
	mu     sync.Mutex
	logDir string
}

// Singleton instance
var instance *logger
var once sync.Once

// LogLens returns the singleton logger instance (Laravel-like façade)
func LogLens() *logger {
	once.Do(func() {
		instance = &logger{
			logDir: "./logs",
		}
		if err := os.MkdirAll(instance.logDir, 0755); err != nil {
			log.Fatalf("failed to create log directory: %v", err)
		}

		// global panic handler
		go func() {
			if r := recover(); r != nil {
				instance.write("panic", "default", fmt.Sprintf("%v", r), debug.Stack())
			}
		}()
	})
	return instance
}

// ========== PUBLIC METHODS ==========

func (l *logger) Info(channel string, msg string, args ...interface{}) {
	l.write("INFO", channel, fmt.Sprintf(msg, args...), nil)
}

func (l *logger) Warn(channel string, msg string, args ...interface{}) {
	l.write("WARN", channel, fmt.Sprintf(msg, args...), nil)
}

func (l *logger) Error(channel string, err interface{}, args ...interface{}) {
	switch v := err.(type) {
	case error:
		// error object + traceback
		l.write("ERROR", channel, v.Error(), debug.Stack())
	case string:
		// formatted string + traceback
		l.write("ERROR", channel, fmt.Sprintf(v, args...), debug.Stack())
	default:
		// fallback
		l.write("ERROR", channel, fmt.Sprintf("%v", v), debug.Stack())
	}
}

// ========== INTERNAL WRITE ==========

func (l *logger) write(level, channel, msg string, stack []byte) {
	l.mu.Lock()
	defer l.mu.Unlock()

	date := time.Now().Format("2006-01-02")
	logFile := filepath.Join(l.logDir, date+".log")

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("failed to open log file: %v", err)
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	entry := fmt.Sprintf("[%s] [%s] [%s] %s\n", timestamp, level, channel, msg)

	if stack != nil {
		entry += fmt.Sprintf("Traceback:\n%s\n", stack)
	}

	if _, err := f.WriteString(entry); err != nil {
		log.Printf("failed to write log: %v", err)
	}
}
