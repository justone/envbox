package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// testBoxUtils is a handy handle to all of the testing implementations of
// interfaces needed by EnvBox.
type testBoxUtils struct {
	*testSystem
	// more to come
}

// cleanup will clean up any testing resources created
func (tau testBoxUtils) cleanup() {
	tau.testSystem.cleanup()
}

// newTestBox creates a new EnvBox struct with all of the sub-interfaces
// replaced with doppelgangers that make testing easy.
func newTestBox() (*EnvBox, *testBoxUtils) {
	box, _ := NewEnvBox()

	tu := &testBoxUtils{
		testSystem: newTestSystem(),
	}

	box.System = tu.testSystem

	return box, tu
}

// testSystem is a testing implementation of the System interface.
type testSystem struct {
	homePath string
	Env      map[string]string
}

func newTestSystem() *testSystem {
	tempdir, _ := ioutil.TempDir("", "envboxapi")

	return &testSystem{
		homePath: tempdir,
		Env:      map[string]string{"HOME": tempdir},
	}
}

func (ts testSystem) Getenv(key string) string {
	val, ok := ts.Env[key]
	if ok {
		return val
	}

	return ""
}

func (ts *testSystem) Setenv(key, value string) {
	ts.Env[key] = value
}

func (ts *testSystem) cleanup() {
	os.RemoveAll(ts.homePath)
}

func (ts testSystem) FileExists(localPath string) bool {
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (ts testSystem) DataPath() (string, error) {
	var envboxPath string
	if xdgDataHome := ts.Getenv("XDG_DATA_HOME"); len(xdgDataHome) > 0 {
		envboxPath = filepath.Join(xdgDataHome, "envbox")
	} else {
		var home string
		if home = ts.Getenv("HOME"); len(home) == 0 {
			return "", fmt.Errorf("$HOME environment variable not found")
		}
		envboxPath = filepath.Join(home, ".local", "share", "envbox")
	}
	os.MkdirAll(envboxPath, 0755)

	return envboxPath, nil
}

func (ts testSystem) ExecCommandWithEnv(command string, args []string, extraEnv []string) error {
	return nil
}
