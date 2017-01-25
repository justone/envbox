package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

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

	return nil
}
