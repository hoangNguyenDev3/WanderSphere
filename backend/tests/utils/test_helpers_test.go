package utils

import (
	"testing"
)

// TestUtilsPackage is a simple test to verify the utils package
// This eliminates the "[no test files]" warning during test runs
func TestUtilsPackage(t *testing.T) {
	// Simple test to ensure the package compiles correctly
	client, err := NewAPIClient()
	if err != nil {
		t.Skipf("Could not create API client in test environment: %v", err)
	}

	if client == nil {
		t.Error("NewAPIClient should not return nil client")
	}
}
