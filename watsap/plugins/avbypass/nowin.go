//go:build !windows

package avbypass

import "log"

func Main() {
	log.Println("Running on non-Windows OS. Defender bypass not applicable.")
}
