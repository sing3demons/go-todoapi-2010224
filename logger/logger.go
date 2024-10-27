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
	Name       string            `json:"name"`
	Timestamp  string            `json:"timestamp"`
	Attributes interface{}       `json:"attributes,omitempty"`
	Msg        map[string]string `json:"msg,omitempty"`
}

func New(s *slog.Logger, name string, attribute map[string]any) *Logger {
	return &Logger{Logger: s, attribute: attribute, startTime: time.Now()}
}

func (l *Logger) addEvent(node, cmd, name string, data interface{}) {
	l.event = fmt.Sprintf("%s.%s", node, cmd)
	attribute := LogEvent{
		Name:       l.name(node, cmd, name),
		Timestamp:  time.Now().Format(time.RFC3339),
		Attributes: data,
	}

	l.attributes = append(l.attributes, attribute)
}

func (l *Logger) name(node, cmd, name string) string {
	return strings.ReplaceAll(fmt.Sprintf("(%s)%s.%s", name, node, cmd), " ", "_")
}

func (l *Logger) AddInput(node, cmd string, data interface{}) {
	l.addEvent(node, cmd, "input", data)
}

func (l *Logger) AddOutput(node, cmd string, custom interface{}) {
	l.addEvent(node, cmd, "output", custom)
	l.End()
}

func (l *Logger) AddError(node, cmd, inOut string, data interface{}, err error) {
	l.event = fmt.Sprintf("%s.%s", node, cmd)

	attribute := LogEvent{
		Name:      l.name(node, cmd, inOut),
		Timestamp:  time.Now().Format(time.RFC3339),
		Attributes: data,
		Msg:        map[string]string{"error": err.Error()},
	}

	l.attributes = append(l.attributes, attribute)
	l.End()
}

func (l *Logger) End() {
	if len(l.attributes) > 0 {
		l.Logger.Info(strings.ReplaceAll(l.event, " ", "_"),
			slog.String("log_name", "DETAIL"),
			slog.Any("startTime", l.startTime),
			slog.Any("endTime", time.Now()),
			slog.Any("processTime", time.Since(l.startTime)),
			slog.Any("attribute", l.attribute),
			slog.Any("events", l.attributes))
	}
	l.attributes = nil
}
