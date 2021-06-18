package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripV1(t *testing.T) {
	tests := map[string]string{
		"/v1/cf/8cb5a2ea-d20a-4ea0-815b-742075dc92ba/secret":  "/cf/8cb5a2ea-d20a-4ea0-815b-742075dc92ba/secret",
		"/v1/cf/51536c9b-f91c-402a-87f5-406258c792df/transit": "/cf/51536c9b-f91c-402a-87f5-406258c792df/transit",
	}
	for input, output := range tests {
		assert.Equal(t, output, stripFirst(input))
	}
}
