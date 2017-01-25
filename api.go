package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/nacl/secretbox"

	"github.com/pkg/errors"
)

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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
	fmt.Printf("TODO: adding variable %s=%s (key: %s)\n", name, value, key)

	message, err := json.Marshal(EnvVar{name, value})
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

	fmt.Println(out)

	return ioutil.WriteFile(fmt.Sprintf("%s.enc", name), out, 0600)
}

func RunCommandWithEnv(key string, vars, cmd []string) error {

	fmt.Printf("TODO: running %v with vars %v (key: %s)\n", cmd, vars, key)

	return nil
}

// var pass [32]byte
// if _, err := io.ReadFull(rand.Reader, pass[:]); err != nil {
// 	panic(err)
// }
