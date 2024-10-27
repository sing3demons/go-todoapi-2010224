package logger

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type Logger struct {
	*slog.Logger
	event       string
	attributes  []LogEvent
	attribute   map[string]any
	startTime   time.Time
	ProcessTime time.Duration
}

type LogEvent struct {
	Name       string      `json:"name"`
	Timestamp  string      `json:"timestamp"`
	Attributes interface{} `json:"attributes,omitempty"`
}

func New(s *slog.Logger, attribute map[string]any) *Logger {
	// name, route, method, device string

	return &Logger{Logger: s, attribute: attribute, startTime: time.Now()}
}

func (l *Logger) AddEvent(node, cmd string, data interface{}) {
	l.event = fmt.Sprintf("%s.%s", node, cmd)

	attribute := LogEvent{
		Name:       cmd,
		Timestamp:  time.Now().Format(time.RFC3339),
		Attributes: data,
	}

	l.attributes = append(l.attributes, attribute)

}

func (l *Logger) End() {
	l.Logger.Info(strings.ReplaceAll(l.event, " ", "_"),
	slog.String("log_name", "DETAIL"),
		slog.Any("startTime", l.startTime),
		slog.Any("endTime", time.Now()),
		slog.Any("processTime", time.Since(l.startTime)),
		slog.Any("attribute", l.attribute),
		slog.Any("events", l.attributes))
}
