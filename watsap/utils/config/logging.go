package config

import (
	"log"
	"os"
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
	log.SetFlags(0)
	log.Println("Logging started")
}

func Logger(msg, logtype string) {
	if !WaLogging {
		return
	}
	switch logtype {
	case "fatal":
		log.Fatal(msg)
	case "panic":
		log.Panic(msg)
	case "error":
		log.Println("ERROR:", msg)
	default:
		log.Println("INFO:", msg)
	}
}
