package main

import (
	"fmt"
	"log"
	"strconv"
)

type ServiceLogger struct {
	*log.Logger
	serviceName string
}

func NewServiceLogger(base *log.Logger, serviceName string) *ServiceLogger {
	return &ServiceLogger{
		Logger:      base,
		serviceName: serviceName,
	}
}

func (l *ServiceLogger) Debug(message string, fields ...map[string]string) {
	l.log("debug", message, fields...)
}

func (l *ServiceLogger) Info(message string, fields ...map[string]string) {
	l.log("info", message, fields...)
}

func (l *ServiceLogger) Warn(message string, fields ...map[string]string) {
	l.log("warn", message, fields...)
}

func (l *ServiceLogger) Error(message string, fields ...map[string]string) {
	l.log("error", message, fields...)
}

func (l *ServiceLogger) log(level string, message string, fields ...map[string]string) {
	fieldString := ""

	if len(fields) > 0 {
		for key, value := range fields[0] {
			fieldString += fmt.Sprintf(" %s=%q", key, value)
		}
	}

	l.Logger.Printf("service=%s level=%s %smsg=%s",
		l.serviceName,
		level,
		fieldString,
		strconv.Quote(message))
}
