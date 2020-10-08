package main

import (
	"log"
	"os"
)

// LogHeimdall represents a LogHeimdall Object
type LogHeimdall struct {
	Logger *log.Logger
}

// NewLoggerHeimdall returns a new LoggerHeimdall
func NewLoggerHeimdall(serviceName string) *LogHeimdall {
	Logger := &LogHeimdall{Logger: log.New(os.Stdout, serviceName+" - ", log.LstdFlags|log.Lshortfile)}
	return Logger
}
