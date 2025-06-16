//go:build windows

package command

import (
	"os/exec"
	"syscall"
)

func NewCMD(command string, arg ...string) *exec.Cmd {
	c := exec.Command(command, arg...)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return c
}
