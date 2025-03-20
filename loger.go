package mylogger

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Graylog configuration
const (
	GraylogHost = "34.93.173.35" // Change to your actual Graylog server IP
	GraylogPort = "12201"        // Default GELF port
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

// Logger struct
type Logger struct {
	logger *logrus.Logger
}

// NewLogger initializes a new logger
func NewLogger() *Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(os.Stdout) // Also log to the console

	return &Logger{logger: l}
}

// Log logs a message and sends it to Graylog
func (l *Logger) Log(level, message string, data LogData) {
	data.Timestamp = time.Now().UTC().Format(time.RFC3339)
	data.Message = maskSensitiveData(message) // Mask sensitive data if needed

	jsonData, _ := json.Marshal(data)
	
	// Log locally
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

	// Send to Graylog
	go sendToGraylog(jsonData) // Run in a goroutine to avoid blocking
}

// Mask sensitive data
func maskSensitiveData(message string) string {
	if len(message) > 2 {
		return string(message[0]) + "****" + string(message[len(message)-1])
	}
	return "****"
}

// sendToGraylog sends log data to Graylog over both UDP and TCP
func sendToGraylog(logData []byte) {
	udpAddr := fmt.Sprintf("%s:%s", GraylogHost, GraylogPort)
	tcpAddr := fmt.Sprintf("%s:%s", GraylogHost, GraylogPort)

	// Send over UDP
	err := sendUDP(udpAddr, logData)
	if err != nil {
		log.Println("Failed to send log over UDP:", err)
	}

	// Send over TCP
	err = sendTCP(tcpAddr, logData)
	if err != nil {
		log.Println("Failed to send log over TCP:", err)
	}
}

// sendUDP sends log data over UDP
func sendUDP(address string, data []byte) error {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(data)
	return err
}

// sendTCP sends log data over TCP
func sendTCP(address string, data []byte) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(append(data, '\n')) // GELF messages should end with a newline
	return err
}
