package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/nacl/secretbox"

	"github.com/howeyc/gopass"
	"github.com/pkg/errors"
)

type EnvVar struct {
	Name    string `json:"name"`
	Exposed string `json:"exposed"`
	Value   string `json:"value"`
	Path    string `json:"-"`
}

type EnvBox struct {
	System
	// Config
}

func NewEnvBox() (*EnvBox, error) {
	return &EnvBox{
		System: &DefaultSystem{},
	}, nil
}

func (box *EnvBox) AddVariable(name, exposed, file string) error {

	var err error

	key, err := box.ReadKey()
	if err != nil {
		return errors.Wrap(err, "unable to read key")
	}

	// check for duplicate name
	vars, err := box.LoadEnvVars(key)
	if err != nil {
		return errors.Wrap(err, "unable to load vars")
	}

	if _, ok := vars[name]; ok {
		return fmt.Errorf("var %s already exists", name)
	}

	if len(exposed) == 0 {
		exposed = name
	}

	var value string
	if len(file) > 0 {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return errors.Wrap(err, "error reading file")
		}
		value = strings.TrimSpace(string(data))
	} else {
		value, err = box.PromptForValue()
		if err != nil {
			return errors.Wrap(err, "error reading value")
		}
	}

	message, err := json.Marshal(EnvVar{Name: name, Exposed: exposed, Value: value})
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

	var fname [24]byte
	if _, err := io.ReadFull(rand.Reader, fname[:]); err != nil {
		return errors.Wrap(err, "unable to read random")
	}

	dataPath, err := box.DataPath()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(dataPath, fmt.Sprintf("%s.envenc", hex.EncodeToString(fname[:]))), out, 0600)
}

func (box *EnvBox) keyPath() (string, error) {
	dataPath, err := box.DataPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(dataPath, "secret.key"), nil
}

func (box *EnvBox) ReadKey() (string, error) {
	keyPath, err := box.keyPath()
	if err != nil {
		return "", errors.Wrap(err, "unable to get key path")
	}

	if !box.FileExists(keyPath) {
		key, err := box.PromptForKey()
		if err != nil {
			return "", errors.Wrap(err, "unable prompt for key")
		}

		err = box.StoreKey(key)
		if err != nil {
			return "", errors.Wrap(err, "unable to set key")
		}
	}

	data, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return "", errors.Wrap(err, "unable to read keypath")
	}

	key := strings.TrimSpace(string(data))

	return key, nil
}

func (box *EnvBox) PromptForKey() (string, error) {
	fmt.Printf("enter key: ")

	key, err := gopass.GetPasswdMasked()
	if err != nil {
		if err == gopass.ErrInterrupted {
			return "", fmt.Errorf("interrupted")
		} else {
			return "", errors.Wrap(err, "unable to prompt for key")
		}
	}

	// TODO: check that key is valid

	return string(key), nil
}

func (box *EnvBox) PromptForValue() (string, error) {
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

func (box *EnvBox) StoreKey(key string) error {
	keyPath, err := box.keyPath()
	if err != nil {
		return errors.Wrap(err, "unable to get key path")
	}

	return ioutil.WriteFile(keyPath, []byte(key), 0600)
}

func (box *EnvBox) LoadEnvVars(key string) (map[string]EnvVar, error) {
	vars := make(map[string]EnvVar)

	dataPath, err := box.DataPath()
	if err != nil {
		return vars, errors.Wrap(err, "unable to get data path")
	}
	files, err := ioutil.ReadDir(dataPath)
	if err != nil {
		return vars, errors.Wrap(err, "unable to read directory")
	}

	var keyBytes [32]byte
	copy(keyBytes[:], []byte(key)[:32])

	for _, info := range files {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".envenc") {
			fileName := filepath.Join(dataPath, info.Name())
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

func (box *EnvBox) ListVariables() error {
	key, err := box.ReadKey()
	if err != nil {
		return errors.Wrap(err, "unable to read key")
	}

	vars, err := box.LoadEnvVars(key)
	if err != nil {
		return errors.Wrap(err, "unable to load vars")
	}

	for name, envVar := range vars {
		// TODO: figure out a better way to list these
		fmt.Print(name)
		if envVar.Exposed != envVar.Name {
			fmt.Printf("(%s)", envVar.Exposed)
		}
		fmt.Printf("=%s", envVar.Value)
		fmt.Println()
	}

	return nil
}

func (box *EnvBox) GenerateNewKey(set bool) error {
	var pass [32]byte
	if _, err := io.ReadFull(rand.Reader, pass[:]); err != nil {
		return errors.Wrap(err, "unable to read random")
	}

	key := hex.EncodeToString(pass[:])
	fmt.Println(key)

	if set {
		// TODO: warn when overriding existing key
		return box.StoreKey(key)
	}
	return nil
}

func (box *EnvBox) PromptAndStoreKey() error {
	key, err := box.PromptForKey()
	if err != nil {
		return errors.Wrap(err, "unable to prompt for key")
	}

	return box.StoreKey(key)
}

func (box *EnvBox) RemoveVariable(name string) error {
	key, err := box.ReadKey()
	if err != nil {
		return errors.Wrap(err, "unable to read key")
	}

	vars, err := box.LoadEnvVars(key)
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

func (box *EnvBox) RunCommandWithEnv(varNames, command []string) error {
	key, err := box.ReadKey()
	if err != nil {
		return errors.Wrap(err, "unable to read key")
	}

	hostEnv := os.Environ()
	vars, err := box.LoadEnvVars(key)
	if err != nil {
		return errors.Wrap(err, "unable to load env vars")
	}
	for _, varName := range varNames {
		if envVar, ok := vars[varName]; ok {
			hostEnv = append(hostEnv, fmt.Sprintf("%s=%s", envVar.Exposed, envVar.Value))
		} else {
			fmt.Fprintf(os.Stdout, "unable to find %s\n", varName)
		}
	}

	return box.ExecCommandWithEnv(command[0], command[1:], hostEnv)
}
