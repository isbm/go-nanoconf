package nanoconf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoot(t *testing.T) {
	assert.Equal(t, "root-value",
		NewConfig("example.conf").Root().String("key", ""),
		"Result is not as expected")
}

func TestFindFirst(t *testing.T) {
	assert.Equal(t, "section-value",
		NewConfig("example.conf").Find("section").String("key", ""),
		"Result is not as expected")
}

func TestFindPath(t *testing.T) {
	assert.Equal(t, "inner-section-value",
		NewConfig("example.conf").Find("section:inner-section").String("key", ""),
		"Result is not as expected")
}
