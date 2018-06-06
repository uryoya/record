package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

type RusageJson struct {
	Utime   int64 `json:"utime"`
	Stime   int64 `json:"stime"`
	Maxrss  int64 `json:"maxrss"`
	Minflt  int64 `json:"minflt"`
	Majflt  int64 `json:"majflt"`
	Inblock int64 `json:"inblock"`
	Oublock int64 `json:"oublock"`
	Nvcsw   int64 `json:"nvcsw"`
	Nivcsw  int64 `json:"nivcsw"`
}

type RecordJson struct {
	Stdout string             `json:"stdout"`
	Stderr string             `json:"stderr"`
	Status syscall.WaitStatus `json:"status"`
	Rusage RusageJson         `json:"rusage"`
}

func main() {
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	recStdout, err := ioutil.ReadAll(stdout)
	if err != nil {
		panic(err)
	}
	recStderr, err := ioutil.ReadAll(stderr)
	if err != nil {
		panic(err)
	}
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitErr.Sys()
		}
	}
	status := cmd.ProcessState.Sys().(syscall.WaitStatus)
	rusage := cmd.ProcessState.SysUsage().(*syscall.Rusage)
	record := RecordJson{
		string(recStdout),
		string(recStderr),
		status,
		RusageJson{
			rusage.Utime.Usec,
			rusage.Stime.Usec,
			rusage.Maxrss,
			rusage.Minflt,
			rusage.Majflt,
			rusage.Inblock,
			rusage.Oublock,
			rusage.Nvcsw,
			rusage.Nivcsw,
		},
	}
	recordJson, _ := json.Marshal(record)
	fmt.Println(string(recordJson))
}
