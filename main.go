package main

import (
	"os/exec"
	"fmt"
	"syscall"
	"io"
	"os"
)

type RecordJson struct {
	Stdout string
	Stderr string
	Status int
	Rusage syscall.Rusage
}

func main() {
	cmd := exec.Command("ls", "-la")
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, stdout)
	io.Copy(os.Stderr, stderr)
	if err := cmd.Wait(); err != nil {
		panic(err)
	}
	rusage := cmd.ProcessState.SysUsage().(*syscall.Rusage)
	fmt.Println(rusage.Maxrss)
}