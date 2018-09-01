package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/secretbox"
)

type EnvVar struct {
	// Name is what the user uses to refer to the variables, defaults to the
	// same as the name of the environment variable
	Name string `json:"name"`

	// LegacyExposed is the name of the variable to expose.  Early versions of
	// envbox only supported one variable per name, and this is where it was
	// stored.  When the old format is loaded, this is moved to the Vars map
	LegacyExposed string `json:"exposed"`

	// LegacyValue is the value of the variable.  Early versions of envbox only
	// supported one variable per name, and this is where its value was stored.
	// When the old format is loaded, this is moved to the Vars map
	LegacyValue string `json:"value"`

	// Path is the name of the underlying file that the data is stored in.  It
	// isn't present in the JSON data.
	Path string `json:"-"`

	// Vars holds the key/value pairs to expose as environment variables when
	// running commands.
	Vars map[string]string
}

type EnvBox struct {
	System
	Prompter
	io.Writer
	// Config
}

func NewEnvBox() (*EnvBox, error) {
	return &EnvBox{
		System:   &DefaultSystem{},
		Prompter: &DefaultPrompter{},
		Writer:   os.Stdout,
	}, nil
}

func (box *EnvBox) AddVariable(name, exposed, file string, multiple bool) error {

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
		value, err = box.PromptFor("value: ")
		if err != nil {
			return errors.Wrap(err, "error reading value")
		}
	}

	newVars := map[string]string{exposed: value}

	if multiple {

		fmt.Fprintf(box.Writer, "enter additional variables; emtpy name to finish\n")

		for {
			varName, err := box.PromptFor("name: ")
			if err != nil {
				return errors.Wrap(err, "error reading name")
			}
			if len(varName) == 0 {
				break
			}

			varValue, err := box.PromptFor("value: ")
			if err != nil {
				return errors.Wrap(err, "error reading value")
			}

			newVars[varName] = varValue
		}
	}

	message, err := json.Marshal(EnvVar{Name: name, Vars: newVars})
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

	var key string

	if helperKey, _ := GetCredHelperKey(); len(helperKey) > 0 {
		logrus.Debugf("found cred helper key, using that")
		key = helperKey
	} else if pathKeyData, err := ioutil.ReadFile(keyPath); err == nil {
		logrus.Debugf("found file key, using that")
		key = strings.TrimSpace(string(pathKeyData))
	}

	if len(key) == 0 {
		promptedKey, err := box.PromptForKey()
		if err != nil {
			return "", errors.Wrap(err, "unable prompt for key")
		}

		err = box.StoreKey(promptedKey)
		if err != nil {
			return "", errors.Wrap(err, "unable to set key")
		}

		key = promptedKey
	}

	return key, nil
}

func (box *EnvBox) PromptForKey() (string, error) {

	key, err := box.PromptMasked("enter key: ")
	if err != nil {
		return "", errors.Wrap(err, "unable to prompt for key")
	}

	// TODO: check that key is valid

	return key, nil
}

func (box *EnvBox) StoreKey(key string) error {

	err := StoreCredHelperKey(key)
	if err == nil {
		logrus.Debugf("helper key stored")
		return nil
	} else if err != nil && err != helperNotFound {
		return errors.Wrap(err, "unable to set with helper")
	}

	logrus.Debugf("falling back on path based storage")
	keyPath, err := box.keyPath()
	if err != nil {
		return errors.Wrap(err, "unable to get key path")
	}

	return ioutil.WriteFile(keyPath, []byte(key), 0600)
}

func (box *EnvBox) ClearKey() error {
	err := ClearCredHelperKey()
	if err == nil {
		logrus.Debugf("helper key cleared")
		return nil
	} else if err != nil && err != helperNotFound {
		return errors.Wrap(err, "unable to clear cred helper key")
	}

	logrus.Debugf("falling back on path based storage")
	keyPath, err := box.keyPath()
	if err != nil {
		return errors.Wrap(err, "unable to get key path")
	}
	return os.Remove(keyPath)
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

			data, err := ioutil.ReadFile(fileName)
			if err != nil {
				return vars, errors.Wrap(err, "unable to read file")
			}

			nonce := new([24]byte)
			copy(nonce[:], data[:24])

			if message, ok := secretbox.Open(nil, data[24:], nonce, &keyBytes); ok {

				var envVar EnvVar
				err := json.Unmarshal(message, &envVar)
				if err != nil {
					// ignore
				}
				envVar.Path = fileName

				if len(envVar.LegacyExposed) > 0 {
					envVar.Vars = map[string]string{envVar.LegacyExposed: envVar.LegacyValue}
				}

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
		fmt.Fprintf(box.Writer, name)
		fmt.Fprintf(box.Writer, ": ")

		varNames := []string{}
		for k, _ := range envVar.Vars {
			varNames = append(varNames, fmt.Sprintf("%s", k))
		}

		fmt.Fprintf(box.Writer, "%s", strings.Join(varNames, ", "))
		fmt.Fprintf(box.Writer, "\n")
	}

	return nil
}

func (box *EnvBox) GenerateNewKey(set bool) error {
	var pass [32]byte
	if _, err := io.ReadFull(rand.Reader, pass[:]); err != nil {
		return errors.Wrap(err, "unable to read random")
	}

	key := hex.EncodeToString(pass[:])
	fmt.Fprintf(box.Writer, "%s\n", key)

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

func (box *EnvBox) ShowKey() error {
	key, err := box.ReadKey()
	if err != nil {
		return errors.Wrap(err, "unable to read key")
	}

	fmt.Fprintf(box.Writer, "%s\n", key)

	return nil
}

func (box *EnvBox) withFoundKey(name string, fun func(EnvVar) error) error {
	key, err := box.ReadKey()
	if err != nil {
		return errors.Wrap(err, "unable to read key")
	}

	vars, err := box.LoadEnvVars(key)
	if err != nil {
		return errors.Wrap(err, "unable to load env vars")
	}

	if envVar, ok := vars[name]; ok {
		return fun(envVar)
	} else {
		return fmt.Errorf("variable %s not found", name)
	}
	return nil
}

func (box *EnvBox) ShowVariable(name string) error {
	return box.withFoundKey(name, func(envVar EnvVar) error {
		fmt.Fprintf(box.Writer, "name: %s\n", envVar.Name)
		fmt.Fprintf(box.Writer, "vars:\n")
		for k, v := range envVar.Vars {
			fmt.Fprintf(box.Writer, "  %s: %s\n", k, v)
		}
		return nil
	})
}

func (box *EnvBox) ExportVariable(name string) error {
	return box.withFoundKey(name, func(envVar EnvVar) error {
		for k, v := range envVar.Vars {
			// TODO: better value escaping
			fmt.Fprintf(box.Writer, "export %s=%q\n", k, v)
		}
		return nil
	})
}

func (box *EnvBox) RemoveVariable(name string) error {
	return box.withFoundKey(name, func(envVar EnvVar) error {
		err := os.Remove(envVar.Path)
		if err != nil {
			return errors.Wrap(err, "unable to remove file")
		}
		return nil
	})
}

func (box *EnvBox) RunCommandWithEnv(varNames, command []string) error {
	key, err := box.ReadKey()
	if err != nil {
		return errors.Wrap(err, "unable to read key")
	}

	vars, err := box.LoadEnvVars(key)
	if err != nil {
		return errors.Wrap(err, "unable to load env vars")
	}

	var exposeVars []EnvVar
	for _, varName := range varNames {
		if envVar, ok := vars[varName]; ok {
			exposeVars = append(exposeVars, envVar)
		} else {
			fmt.Fprintf(os.Stdout, "unable to find %s\n", varName)
		}
	}

	var useEnv []string

	hostEnv := os.Environ()
	for _, hostVar := range hostEnv {
		conflictFound := false
		for _, expVar := range exposeVars {
			for exposed, _ := range expVar.Vars {
				if strings.HasPrefix(hostVar, fmt.Sprintf("%s=", exposed)) {
					conflictFound = true
				}
			}
		}

		if !conflictFound {
			useEnv = append(useEnv, hostVar)
		}
	}

	for _, expVar := range exposeVars {
		for exposed, value := range expVar.Vars {
			useEnv = append(useEnv, fmt.Sprintf("%s=%s", exposed, value))
		}
	}

	return box.ExecCommandWithEnv(command[0], command[1:], useEnv)
}
