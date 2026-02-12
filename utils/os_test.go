package utils

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindProjectRoot(t *testing.T) {
	projectRoot, err := FindProjectRoot()
	assert.NoError(t, err, "Expected no error when finding project root")
	assert.NotEmpty(t, projectRoot, "Expected project root to be non-empty")
	log.Printf("Project root found at: %s", projectRoot)	
}