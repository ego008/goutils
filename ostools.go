package goutils

import (
	"os/exec"
)

// CmdExists Golang check if command exists
func CmdExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
