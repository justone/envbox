package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/nacl/secretbox"

	"github.com/pkg/errors"
)

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Path  string `json:"-"`
}

func ReadKey() (string, error) {
	data, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return "", errors.Wrap(err, "unable to read keypath")
	}

	key := strings.TrimSpace(string(data))

	return key, nil
}

func PromptForValue() (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", errors.Wrap(err, "failed to open /dev/tty")
	}

	fmt.Fprintf(tty, "value: ")
	value, err := bufio.NewReader(tty).ReadString('\n')
	if err != nil {
		return "", errors.Wrap(err, "unable to read value")
	}

	err = tty.Close()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(value), nil
}

func AddVariable(key, name, value string) error {
	fmt.Printf("adding variable %s=%s (key: %s)\n", name, value, key)

	message, err := json.Marshal(EnvVar{Name: name, Value: value})
	if err != nil {
		return err
	}

	var keyBytes [32]byte
	copy(keyBytes[:], []byte(key)[:32])

	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	out := make([]byte, 24)
	copy(out, nonce[:])

	out = secretbox.Seal(out, message, &nonce, &keyBytes)

	// fmt.Println(out)

	return ioutil.WriteFile(fmt.Sprintf("%s.envenc", name), out, 0600)
}

func RemoveVariable(key, name string) error {
	fmt.Printf("removing variable %s (key: %s)\n", name, key)

	vars, err := LoadEnvVars(key)
	if err != nil {
		return errors.Wrap(err, "unable to load env vars")
	}

	if envVar, ok := vars[name]; ok {
		err = os.Remove(envVar.Path)
		if err != nil {
			return errors.Wrap(err, "unable to remove file")
		}
	} else {
		return fmt.Errorf("variable %s not found", name)
	}
	return nil
}

func RunCommandWithEnv(key string, varNames, command []string) error {

	fmt.Printf("running %v with vars %v (key: %s)\n", command, varNames, key)

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	hostEnv := os.Environ()
	vars, err := LoadEnvVars(key)
	if err != nil {
		return errors.Wrap(err, "unable to load env vars")
	}
	for _, varName := range varNames {
		if envVar, ok := vars[varName]; ok {
			hostEnv = append(hostEnv, fmt.Sprintf("%s=%s", varName, envVar.Value))
		} else {
			// TODO: handle variable not found
		}
	}
	cmd.Env = hostEnv

	return cmd.Run()
}

func LoadEnvVars(key string) (map[string]EnvVar, error) {
	vars := make(map[string]EnvVar)

	files, err := ioutil.ReadDir(secretPath)
	if err != nil {
		return vars, errors.Wrap(err, "unable to read directory")
	}

	var keyBytes [32]byte
	copy(keyBytes[:], []byte(key)[:32])

	for _, info := range files {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".envenc") {
			fileName := filepath.Join(secretPath, info.Name())
			// fmt.Println("Loading file", fileName)

			data, err := ioutil.ReadFile(fileName)
			if err != nil {
				return vars, errors.Wrap(err, "unable to read file")
			}

			nonce := new([24]byte)
			copy(nonce[:], data[:24])

			if message, ok := secretbox.Open(nil, data[24:], nonce, &keyBytes); ok {
				// fmt.Println(string(message))

				var envVar EnvVar
				err := json.Unmarshal(message, &envVar)
				if err != nil {
					// ignore
				}
				envVar.Path = fileName

				vars[envVar.Name] = envVar
			} else {
				// ignore
			}

		}
	}

	return vars, nil
}

// var pass [32]byte
// if _, err := io.ReadFull(rand.Reader, pass[:]); err != nil {
// 	panic(err)
// }
