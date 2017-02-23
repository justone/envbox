package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreKey(t *testing.T) {
	assert := assert.New(t)

	box, tu := newTestBox()
	defer tu.cleanup()

	assert.Nil(box.StoreKey("testkey"))
	fileData, _ := ioutil.ReadFile(filepath.Join(tu.testSystem.homePath, ".local/share/envbox/secret.key"))
	assert.Equal(string(fileData), "testkey")
}
