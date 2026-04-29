// Custom logger config
package logger

import (
	"3DC/config"
	"3DC/util/color"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var logger *log.Logger
var LogLevel = config.LogLevel

type LOGLEVEL int

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

func init() {
	logFile, _ := os.OpenFile(config.CurrentLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

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

func Debug(msg string) {
	logger.Print("DEBUG: " + metadata() + msg)
	if LogLevel <= 0 {
		fmt.Println(color.ColorText("DEBUG: "+msg, color.CLR))
	}
}

func Info(msg string) {
	logger.Print("INFO: " + metadata() + msg)
	if LogLevel <= 1 {
		fmt.Println(color.ColorText("INFO: "+msg, color.Blue))
	}
}

func Warn(msg string) {
	logger.Print("WARN: " + metadata() + msg)
	if LogLevel <= 2 {
		fmt.Print(color.ColorText("WARN: "+msg, color.Yellow))
	}
}

func Error(msg string) {
	logger.Print("ERROR: " + metadata() + msg)
	if LogLevel <= 3 {
		fmt.Print(color.ColorText("ERROR: "+msg, color.Red))
	}
}

func Fatal(msg string) {
	logger.Print("FATAL: " + metadata() + msg)
	fmt.Println(color.ColorText("ERROR: "+msg, color.Red))
	os.Exit(1)
}
