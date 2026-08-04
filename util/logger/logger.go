// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// Custom logger configuration
// The logger package provides a custom logging solution for the application, allowing for different log levels (debug, info, warn, error, fatal) and output formatting.
// It supports writing logs to a file and optionally printing them to the console with color-coded messages based on severity.
// The logger captures metadata such as timestamps and file locations for better traceability of log messages.
package logger

import (
	"3DC/config"
	"3DC/util/color"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// logger allows the stdout of the log to be changed for testing purposes, and allows the log to be written to a file for later review
var logger *log.Logger

// Currelt level of logging, can be set in config.go
var LogLevel = config.LogLevel

type LOGLEVEL int

// Log levels for the logger, with increasing severity from debug to fatal.
const (
	debug = iota // 0
	info         // 1
	warn         // 2
	err          // 3
	fatal        // 4
)

// Stolen from go documentation
const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

// Setting up io writer so that the logger can be used later in test cases, or redirected for other purposes
var output io.Writer = os.Stdout

// Sets the output destination for the logger. This allows for redirecting log output to different destinations, such as files or buffers, for testing or other purposes.
func SetOutput(w io.Writer) {
	output = w
}

// Ensures that the log file can write as soon as initlaized
func init() {
	logFile, err := os.OpenFile(config.CurrentLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o755)
	if err != nil {
		Fatal(fmt.Sprintf("Could not find or load log file for reason %s", err))
	}

	logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetFlags(0)
}

// For now im accepting the time cost of usign runtime.Caller
// In the future I can scope these calls only to fatal and error calls
// But for simplicity sake right now, im sticking with this
func metadata() string {
	now := time.Now().Format("2006/01/02 15:04:05")
	_, file, line, _ := runtime.Caller(2)
	fileLine := fmt.Sprintf("%s:%d", filepath.Base(file), line)
	return now + " " + fileLine + " "
}

// Debug logs a debug message with metadata. It prints the message to the log file and, if the log level is set to debug or lower, also prints it to the console in NO color.
func Debug(msg string) {
	logger.Print("DEBUG: " + metadata() + msg)
	if LogLevel <= 0 {
		fmt.Fprintln(output, color.ColorText("DEBUG: "+msg, color.CLR))
	}
}

// Info logs an informational message with metadata. It prints the message to the log file and, if the log level is set to info or lower, also prints it to the console in a BLUE color.
func Info(msg string) {
	logger.Print("INFO: " + metadata() + msg)
	if LogLevel <= 1 {
		fmt.Fprintln(output, color.ColorText("INFO: "+msg, color.Blue))
	}
}

// Warn logs a warning message with metadata. It prints the message to the log file and, if the log level is set to warn or lower, also prints it to the console in a YELLOW color.
func Warn(msg string) {
	logger.Print("WARN: " + metadata() + msg)
	if LogLevel <= 2 {
		fmt.Fprint(output, color.ColorText("WARN: "+msg, color.Yellow))
	}
}

// Error logs an error message with metadata. It prints the message to the log file and, if the log level is set to error or lower, also prints it to the console in a RED color.
func Error(msg string) {
	logger.Print("ERROR: " + metadata() + msg)
	if LogLevel <= 3 {
		fmt.Fprint(output, color.ColorText("ERROR: "+msg, color.Red))
	}
	// os.Exit(1)
}

// Fatal logs a fatal error message with metadata. It prints the message to the log file and, if the log level is set to fatal or lower, also prints it to the console in a PURPLE color. The function is intended to be used for critical errors that require immediate attention.
func Fatal(msg string) {
	logger.Print("FATAL: " + metadata() + msg)
	fmt.Fprintln(output, color.ColorText("FATAL: "+msg, color.Purple))
	os.Exit(1)
}
