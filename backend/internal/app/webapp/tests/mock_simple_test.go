package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleTest(t *testing.T) {
	// This is a simple test to verify that the test environment works
	assert.True(t, true, "True should be true")

	// Test some basic time handling
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	assert.True(t, tomorrow.After(now), "Tomorrow should be after today")
}

func TestFormatting(t *testing.T) {
	// Test string formatting
	name := "User"
	formatted := "Hello, " + name
	assert.Equal(t, "Hello, User", formatted)
}
