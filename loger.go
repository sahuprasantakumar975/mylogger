package mylogger

import (
	"encoding/json"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// LogData represents the structured log format
type LogData struct {
	Timestamp    string `json:"timestamp"`
	Level        string `json:"level"`
	Message      string `json:"message,omitempty"`
	IPAddress    string `json:"ip_address,omitempty"`
	AppName      string `json:"appname"`
	Hostname     string `json:"hostname,omitempty"`
	TransactionID string `json:"tr_id,omitempty"`
	Channel      string `json:"channel,omitempty"`
	BankCode     string `json:"bank_code,omitempty"`
	ReferenceID  string `json:"reference_id,omitempty"`
	RRN          string `json:"rrn,omitempty"`
	PublishID    string `json:"publish_id,omitempty"`
	CFTrID       string `json:"cf_trid,omitempty"`
	DeviceInfo   string `json:"device_info,omitempty"`
	ParamA       string `json:"param_a,omitempty"`
	ParamB       string `json:"param_b,omitempty"`
	ParamC       string `json:"param_c,omitempty"`
}

// Logger wraps logrus
type Logger struct {
	logger *logrus.Logger
}

// NewLogger initializes a new logger
func NewLogger() *Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(os.Stdout) // Can be changed to a file
	return &Logger{logger: l}
}

// Log method to log messages
func (l *Logger) Log(level, message string, data LogData) {
	data.Timestamp = time.Now().UTC().Format(time.RFC3339)
	data.Message = maskSensitiveData(message) // Mask if needed

	jsonData, _ := json.Marshal(data)
	switch level {
	case "INFO":
		l.logger.Info(string(jsonData))
	case "ERROR":
		l.logger.Error(string(jsonData))
	case "DEBUG":
		l.logger.Debug(string(jsonData))
	default:
		l.logger.Warn(string(jsonData))
	}
}

// Mask sensitive data
func maskSensitiveData(message string) string {
	// Example: mask everything except first and last character
	if len(message) > 2 {
		return string(message[0]) + "****" + string(message[len(message)-1])
	}
	return "****"
}
