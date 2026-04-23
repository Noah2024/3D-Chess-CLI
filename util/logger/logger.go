// Custom logger config
package logger

import (
	"3DChessCLI/util/color"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var logger *log.Logger
var logLevel LOGLEVEL

type LOGLEVEL int

const (
	debug = iota // 0
	info         // 1
	warn         // 2
	error        // 3
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
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetFlags(0)
	logLevel = debug
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
	if logLevel <= 0 {
		fmt := "DEBUG: " + metadata() + msg
		logger.Print(color.ColorText(fmt, color.CLR))
	}
}

func Info(msg string) {
	if logLevel <= 1 {
		fmt := "INFO: " + metadata() + msg
		logger.Print(color.ColorText(fmt, color.Blue))
	}
}

func Warn(msg string) {
	if logLevel <= 2 {
		fmt := "WARN: " + metadata() + msg
		logger.Print(color.ColorText(fmt, color.Yellow))
	}
}

func Error(msg string) {
	if logLevel <= 2 {
		fmt := "ERROR: " + metadata() + msg
		logger.Print(color.ColorText(fmt, color.Red))
	}
}

func Fatal(msg string) {
	fmt := "FATAL: " + metadata() + msg
	logger.Print(color.ColorText(fmt, color.Purple))
	//Exit
}
