package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func SetupLog() {
	if !WaLogging {
		return
	}
	file, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Logging started")
}

func Logger(msg, logtype string) {
	if !WaLogging {
		return
	}
	currentTime := time.Now().Format(time.RFC3339)
	_, file, line, _ := runtime.Caller(1)
	shortFile := filepath.Base(file)
	formattedMsg := fmt.Sprintf("%s %s:%d %s", currentTime, shortFile, line, msg)

	switch logtype {
	case "fatal":
		log.Fatal(formattedMsg)
	case "panic":
		log.Panic(formattedMsg)
	case "error":
		log.Println("ERROR:", formattedMsg)
	default:
		log.Println("INFO:", formattedMsg)
	}
}
