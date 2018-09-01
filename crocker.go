package main

import (
	"fmt"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/justone/crocker"
)

var url = "https://github.com/justone/envbox"
var helperNotFound = fmt.Errorf("crocker could not find helper")

// GetCredHelperKey tries to retrieve the key from the docker-credential-* set
// of utilities, as discovered by crocker (https://github.com/justone/crocker).
func GetCredHelperKey() (string, error) {
	cr, err := crocker.NewWithStrategy(crocker.MemThenStockStrategy{})
	if err != nil {
		return "", helperNotFound
	}

	// logrus.Debugf("found cred helper instance", cr)
	creds, err := cr.Get(url)
	if err != nil {
		return "", err
	}

	// logrus.Debugf("found cred helper creds", creds)
	return creds.Secret, nil
}

// StoreCredHelperKey tries to store the key using the docker-credential-* set
// of utilities, as discovered by crocker (https://github.com/justone/crocker).
func StoreCredHelperKey(keys string) error {
	cr, err := crocker.NewWithStrategy(crocker.MemThenStockStrategy{})

	if err != nil {
		return helperNotFound
	}

	// logrus.Debugf("found cred helper instance", cr)
	creds := &credentials.Credentials{url, "key", keys}
	err = cr.Store(creds)
	if err != nil {
		return err
	}

	// logrus.Debugf("stored cred helper creds", creds)
	return nil
}

// ClearCredHelperKey tries to clear the key using the docker-credential-* set
// of utilities, as discovered by crocker (https://github.com/justone/crocker).
func ClearCredHelperKey() error {
	cr, err := crocker.NewWithStrategy(crocker.MemThenStockStrategy{})

	if err != nil {
		return helperNotFound
	}

	// logrus.Debugf("found cred helper instance", cr)
	err = cr.Erase(url)
	if err != nil {
		return err
	}

	// logrus.Debugf("cleared cred helper creds")
	return nil
}
