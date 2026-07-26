package config

import (
	"io"
	"log"
	"os"
)

func SetupLog() {
	if !DebugMode && !WaLogging {
		log.SetOutput(io.Discard)
		return
	}

	if WaLogging {
		file, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.SetOutput(io.Discard)
			return
		}
		log.SetOutput(file)
		log.SetFlags(0)
		log.Println("Logging started")
	} else if DebugMode {
		log.SetOutput(os.Stdout)
		log.SetFlags(log.LstdFlags)
		log.Println("Logging started (Debug Mode)")
	} else {
		log.SetOutput(io.Discard)
	}
}
