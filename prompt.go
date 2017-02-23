package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/pkg/errors"
)

type Prompter interface {
	PromptMasked(string) (string, error)
	PromptFor(string) (string, error)
}

type DefaultPrompter struct{}

func (dp DefaultPrompter) PromptMasked(prompt string) (string, error) {
	fmt.Printf(prompt)

	val, err := gopass.GetPasswdMasked()
	if err != nil {
		if err == gopass.ErrInterrupted {
			return "", fmt.Errorf("interrupted")
		} else {
			return "", errors.Wrap(err, "unable to prompt")
		}
	}

	return string(val), nil
}

func (db DefaultPrompter) PromptFor(prompt string) (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", errors.Wrap(err, "failed to open /dev/tty")
	}

	fmt.Fprintf(tty, prompt)
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
