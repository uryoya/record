package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"time"
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
	Rtime  int64              `json:"rtime"`
	Rusage RusageJson         `json:"rusage"`
}

func runCommand(name string, args ...string) (*RecordJson, error) {
	cmd := exec.Command(name, args...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	start := time.Now().UnixNano()
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	recStdout, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	recStderr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return nil, err
	}
	cmd.Wait()
	end := time.Now().UnixNano()
	status := cmd.ProcessState.Sys().(syscall.WaitStatus)
	rusage := cmd.ProcessState.SysUsage().(*syscall.Rusage)
	record := RecordJson{
		string(recStdout),
		string(recStderr),
		status,
		(end - start) / 1000, // micro second
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
	return &record, nil
}

func main() {
	record, err := runCommand(os.Args[1], os.Args[2:]...)
	if err != nil {
		panic(err)
	}
	recordJson, _ := json.Marshal(record)
	fmt.Println(string(recordJson))
}
