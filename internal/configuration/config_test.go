package configuration

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	err := LoadConfig()
	assert.NoError(t, err, "Expected no error when loading config")
	log.Printf("Configuration loaded successfully")
}