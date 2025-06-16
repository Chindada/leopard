//go:build !windows

package command

import "os/exec"

func NewCMD(command string, arg ...string) *exec.Cmd {
	return exec.Command(command, arg...)
}
