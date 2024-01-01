package util

import (
	"testing"
	"os"
	"syscall"
)

func TestRunCmd(t *testing.T) {
	pid := StartProc("/bin/sleep", "5")
	p, err := os.FindProcess(pid)
	err = p.Signal(syscall.SIGCONT)
	if err != nil {
		t.Errorf("Run command error with pid %d", pid)
	}
}
