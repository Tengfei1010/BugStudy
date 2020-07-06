package analyzer

import (
	"bytes"
	"github.com/google/logger"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func RunCmd(command string, args string) (string, string) {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command(command, args)
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	if err != nil {
		logger.Fatalf("cmd.Start() failed with running %s'%s'\n", err)
	}
	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()
	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()
	err = cmd.Wait()
	if err != nil {
		logger.Infof("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		logger.Infof("failed to capture stdout or stderr\n")
	}
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())

	return outStr, errStr
}

func RunModuleCmd(pathToDir string, command string, arg ...string) (string, string) {
	// if the pathToDir has vendor and vendor has files, return
	if command == "mod" {
		vendorFile := filepath.Join(pathToDir, "vendor")
		if _, err := os.Stat(vendorFile); os.IsExist(err) {
			return "success", ""
		}
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command(command, arg ...)
	cmd.Dir = pathToDir
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	if err != nil {
		logger.Fatalf("cmd.Start() failed with '%s'\n", err)
	}
	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()
	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()
	err = cmd.Wait()
	if err != nil {
		logger.Infof("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		logger.Infof("failed to capture stdout or stderr\n")
	}
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())

	return outStr, errStr
}
