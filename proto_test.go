package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtocol(t *testing.T) {
	msg := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	cmd, err := parseCommand(msg)
	assert.Nil(t, err)
	assert.Equal(t, "", cmd)
}
