package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type InputOutputLog struct {
	Invoke   string      `json:"Invoke"`
	Event    string      `json:"Event"`
	Protocol *string     `json:"Protocol,omitempty"`
	Type     string      `json:"Type"`
	RawData  interface{} `json:"RawData,omitempty"`
	Data     interface{} `json:"Data"`
	ResTime  *string     `json:"ResTime,omitempty"`
}

type DetailLog struct {
	LogType         string           `json:"LogType"`
	Host            string           `json:"Host"`
	AppName         string           `json:"AppName"`
	Instance        *string          `json:"Instance,omitempty"`
	Session         string           `json:"Session"`
	Scenario        string           `json:"Scenario"`
	InputTimeStamp  *string          `json:"InputTimeStamp,omitempty"`
	Input           []InputOutputLog `json:"Input"`
	OutputTimeStamp *string          `json:"OutputTimeStamp,omitempty"`
	Output          []InputOutputLog `json:"Output"`
	ProcessingTime  *string          `json:"ProcessingTime,omitempty"`
	startTimeDate   time.Time        `json:"-"`
	inputTime       *time.Time
	outputTime      *time.Time
	timeCounter     map[string]time.Time   `json:"-"`
	Attributes      map[string]interface{} `json:"attributes,omitempty"`
}

func NewDetailLog(c *gin.Context) DetailLog {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{
					Key: "@timestamp",
				}
			}

			if a.Key == "msg" {
				return slog.Attr{
					Key: "log_name",
				}
			}
			return a
		},
	})))

	// slog.With(slog.String("session", c.GetString("session")))
	host := c.Request.Host
	startTimeDate := time.Now()

	return DetailLog{
		LogType:       "DetailLog",
		startTimeDate: startTimeDate,
		Host:          host,
		Session:       c.GetString("session"),
		AppName:       "AppName",
		Scenario:      "Scenario",
		timeCounter:   make(map[string]time.Time),
		Attributes: map[string]interface{}{
			"device":        c.Request.UserAgent(),
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"query":         c.Request.URL.RawQuery,
			"proto":         c.Request.Proto,
			"remote":        c.Request.RemoteAddr,
			"uri":           c.Request.RequestURI,
			"url":           c.Request.URL,
			"auth":          c.Request.Header.Get("Authorization"),
			"size":          c.Writer.Size(),
			"error":         c.Errors,
			"keys":          c.Keys,
			"fullPath":      c.FullPath(),
			"accepted":      c.Accepted,
			"clientIP":      c.ClientIP(),
			"contentType":   c.ContentType(),
			"contentLength": c.Request.ContentLength,
		},
		Input:  []InputOutputLog{},
		Output: []InputOutputLog{},
	}
}

func (d *DetailLog) AddInputLog(invoke, event, protocol, logType string, rawData, data interface{}) {
	now := time.Now()
	if d.startTimeDate.IsZero() {
		d.startTimeDate = now
	}

	d.Input = append(d.Input, InputOutputLog{
		Invoke:   invoke,
		Event:    event,
		Protocol: &protocol,
		Type:     logType,
		RawData:  rawData,
		Data:     data,
	})
}

func (d *DetailLog) AddOutputLog(invoke, event, protocol, logType string, rawData, data interface{}) {
	d.Output = append(d.Output, InputOutputLog{
		Invoke:   invoke,
		Event:    event,
		Protocol: &protocol,
		Type:     logType,
		RawData:  rawData,
		Data:     data,
	})
}

func (d *DetailLog) GetInputTimeStamp() *string {
	inputTime := time.Now()
	d.inputTime = &inputTime
	inputTimeStr := inputTime.Format(time.RFC3339)
	return &inputTimeStr
}

func (d *DetailLog) End() {
	ProcessingTime := time.Now().Sub(d.startTimeDate).String()
	d.ProcessingTime = &ProcessingTime

	outputTime := time.Now()
	d.outputTime = &outputTime

	inputTimeStamp := d.inputTime.Format(time.RFC3339)
	d.InputTimeStamp = &inputTimeStamp

	outputTimeStamp := d.outputTime.Format(time.RFC3339)
	d.OutputTimeStamp = &outputTimeStamp

	slog.Info("DetailLog", d)
	d.inputTime = nil
	d.outputTime = nil
}
