package util

import (
	"fmt"
	"os"
	"syscall"
)

const (
	UID = 1000
	GUID = 1000
)

func StartProc(args ...string) int {
    // The Credential fields are used to set UID, GID and attitional GIDS of the process
    // You need to run the program as  root to do this
	var cred =  &syscall.Credential{ UID, GUID, []uint32{}, true }
    // the Noctty flag is used to detach the process from parent tty
    var sysproc = &syscall.SysProcAttr{  Credential:cred, Noctty:true }
    var attr = os.ProcAttr{
        Dir: ".",
        Env: os.Environ(),
        Files: []*os.File{
            os.Stdin,
            nil,
            nil,
        },
            Sys:sysproc,

    }
    process, err := os.StartProcess(args[0], args, &attr)
	pid := process.Pid
    if err == nil {
        // It is not clear from docs, but Realease actually detaches the process
        err = process.Release();
        if err != nil {
            fmt.Println(err.Error())
        }
    } else {
        fmt.Println(err.Error())
    }
	return pid
}
