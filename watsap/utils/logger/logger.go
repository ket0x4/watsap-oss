package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
	"watsap/utils/config"
)

var (
	mu      sync.Mutex
	logFile *os.File
)

func InitLogger() {
	mu.Lock()
	defer mu.Unlock()

	if logFile != nil {
		logFile.Close()
		logFile = nil
	}

	if !config.DebugMode && !config.WaLogging {
		log.SetOutput(io.Discard)
		return
	}

	var writers []io.Writer

	if config.WaLogging {
		file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logFile = file
			writers = append(writers, file)
		}
	}

	if config.DebugMode {
		writers = append(writers, os.Stdout)
	}

	if len(writers) > 0 {
		log.SetOutput(io.MultiWriter(writers...))
		log.SetFlags(0)
		if config.DebugMode {
			fmt.Println("[Init] Logging started (Debug Mode)")
		}
	} else {
		log.SetOutput(io.Discard)
	}
}

func logMessage(level string, module string, format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatted := fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, level, module, msg)

	mu.Lock()
	defer mu.Unlock()

	if config.WaLogging && logFile != nil {
		logFile.WriteString(formatted + "\n")
		logFile.Sync() // Flush to disk for safe concurrent reading
	}
	if config.DebugMode {
		fmt.Println(formatted)
	}
}

func Debug(module string, format string, v ...any) {
	if config.DebugMode {
		logMessage("DEBUG", module, format, v...)
	}
}

func Info(module string, format string, v ...any) {
	logMessage("INFO", module, format, v...)
}

func Warn(module string, format string, v ...any) {
	logMessage("WARN", module, format, v...)
}

func Error(module string, format string, v ...any) {
	logMessage("ERROR", module, format, v...)
}

// SafeReadLogFile returns the contents of log.w safely without Windows file sharing violations
func SafeReadLogFile() ([]byte, error) {
	mu.Lock()
	defer mu.Unlock()

	if logFile != nil {
		logFile.Sync()
	}
	return os.ReadFile(config.LogFile)
}
