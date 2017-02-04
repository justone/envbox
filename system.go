package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
)

type System interface {
	Getenv(string) string
	FileExists(string) bool
	DataPath() (string, error)
	ExecCommandWithEnv(string, []string, []string) error
}

type DefaultSystem struct{}

func (ds DefaultSystem) Getenv(key string) string {
	return os.Getenv(key)
}

func (ds DefaultSystem) FileExists(localPath string) bool {
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (ds DefaultSystem) DataPath() (string, error) {
	var envboxPath string
	if xdgDataHome := ds.Getenv("XDG_DATA_HOME"); len(xdgDataHome) > 0 {
		envboxPath = filepath.Join(xdgDataHome, "envbox")
	} else {
		var home string
		if home = ds.Getenv("HOME"); len(home) == 0 {
			return "", fmt.Errorf("$HOME environment variable not found")
		}
		envboxPath = filepath.Join(home, ".local", "share", "envbox")
	}
	os.MkdirAll(envboxPath, 0755)

	return envboxPath, nil
}

func (ds DefaultSystem) ExecCommandWithEnv(command string, args []string, extraEnv []string) error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command(command, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = extraEnv

		if err := cmd.Run(); err != nil {
			if eerr, ok := err.(*exec.ExitError); ok {
				os.Exit(eerr.Sys().(syscall.WaitStatus).ExitStatus())
			}
		}

		os.Exit(0)
		return nil
	} else {
		// adapted from https://gobyexample.com/execing-processes
		fullPath, err := exec.LookPath(command)
		if err != nil {
			return err
		}
		return syscall.Exec(fullPath, append([]string{filepath.Base(command)}, args...), extraEnv)
		// end adapted from
	}

	return nil
}
