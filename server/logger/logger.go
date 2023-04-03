package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

const (
	LOG_LEVEL_DEBUG = iota
	LOG_LEVEL_INFO
	LOG_LEVEL_WARN
	LOG_LEVEL_ERROR
)

var logLevel int = LOG_LEVEL_DEBUG

func SetLogLevel(level int) { logLevel = level }

func writeLog(level, log string) {
	_, file, line, _ := runtime.Caller(2)

	file = filepath.Base(file)
	now := time.Now()

	fmt.Printf("[%s][%04d-%02d-%02d %02d:%02d:%02d.%06d][%s:%d] %s\n", level,
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000,
		file, line, log)
}

func DEBUG(format string, args ...any) {
	if logLevel > LOG_LEVEL_DEBUG {
		return
	}
	writeLog("DEBUG", fmt.Sprintf(format, args...))
}

func INFO(format string, args ...any) {
	if logLevel > LOG_LEVEL_INFO {
		return
	}
	writeLog("INFO", fmt.Sprintf(format, args...))
}

func WARN(format string, args ...any) {
	if logLevel > LOG_LEVEL_WARN {
		return
	}
	writeLog("WARN", fmt.Sprintf(format, args...))
}

func ERROR(format string, args ...any) {
	if logLevel > LOG_LEVEL_ERROR {
		return
	}
	writeLog("ERROR", fmt.Sprintf(format, args...))
}
