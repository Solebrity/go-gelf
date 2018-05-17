// Package gelf provides a library for logging messages in the GELF format for a Graylog server.
package gelf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Logger is a client to the Graylog server
type Logger struct {
	GraylogHost string
	GraylogPort int
	Version     string
	Level       int
	Hostname    string
}

// NewLogger creates a new Logger
func NewLogger(host string, port int, applicationHost string) *Logger {
	var err error
	defaultLevel := 6
	if os.Getenv("GELF_LOG_LEVEL") != "" {
		defaultLevel, err = strconv.Atoi(os.Getenv("GELF_LOG_LEVEL"))
		if err != nil {
			panic(err)
		}
	}
	return &Logger{host, port, "1.1", defaultLevel, applicationHost}
}

// Debug logs a message with a debug priority level
func (l *Logger) Debug(shortM string, longM string, metadata map[string]interface{}) {
	m := l.basicMessage(shortM, longM, metadata)
	m["level"] = 7
	l.sendMessage(m)
}

// Info logs a message with an info priority level
func (l *Logger) Info(shortM string, longM string, metadata map[string]interface{}) {
	m := l.basicMessage(shortM, longM, metadata)
	l.sendMessage(m)
}

// Notice logs a message with the notice priority level
func (l *Logger) Notice(shortM string, longM string, metadata map[string]interface{}) {
	m := l.basicMessage(shortM, longM, metadata)
	m["level"] = 5
	l.sendMessage(m)
}

// Warning logs a message with the warn priority level
func (l *Logger) Warning(shortM string, longM string, metadata map[string]interface{}) {
	m := l.basicMessage(shortM, longM, metadata)
	m["level"] = 4
	l.sendMessage(m)
}

// Error logs a message with the error priority level
func (l *Logger) Error(shortM string, longM string, metadata map[string]interface{}) {
	m := l.basicMessage(shortM, longM, metadata)
	m["level"] = 3
	l.sendMessage(m)
}

// Critical logs a message with the critical priority level
func (l *Logger) Critical(shortM string, longM string, metadata map[string]interface{}) {
	m := l.basicMessage(shortM, longM, metadata)
	m["level"] = 2
	l.sendMessage(m)
}

// Alert logs a message with the alert priority level
func (l *Logger) Alert(shortM string, longM string, metadata map[string]interface{}) {
	m := l.basicMessage(shortM, longM, metadata)
	m["level"] = 1
	l.sendMessage(m)
}

func (l *Logger) basicMessage(shortM, longM string, metadata map[string]interface{}) map[string]interface{} {
	m := map[string]interface{}{
		"version":       l.Version,
		"host":          l.Hostname,
		"level":         l.Level,
		"timestamp":     time.Now().Unix(),
		"short_message": shortM,
		"full_message":  longM}

	for k, v := range metadata {
		if !strings.HasPrefix(k, "_") {
			k = "_" + k
		}
		m[k] = v
	}
	return m
}

func (l *Logger) sendMessage(m map[string]interface{}) {
	if os.Getenv("GELF_DISABLED") != "" {
		if strings.ToLower(os.Getenv("GELF_DISABLED")) != "false" && os.Getenv("GELF_DISABLED") != "0" {
			return
		}
	}

	data, err := json.Marshal(m)

	if err != nil {
		fmt.Println(err)
		return
	}

	url := fmt.Sprintf("http://%s:%d/gelf", l.GraylogHost, l.GraylogPort)

	_, err = http.Post(url, "application/json", bytes.NewReader(data))

	if err != nil {
		fmt.Println(err)
		return
	}
}
