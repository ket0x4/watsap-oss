package secure

import (
	"fmt"
	"os"
	"watsap/utils/config"
)

func Imha() {
	fmt.Println("Self-destructing...")
	if err := os.Remove(config.WaDir); err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
