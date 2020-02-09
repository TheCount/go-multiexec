package multiexec

import (
	"testing"
)

// TestDefaultConfig tests whether the default configuration validates.
func TestDefaultConfig(t *testing.T) {
	if err := DefaultConfig().Validate(); err != nil {
		t.Errorf("Default configuration does not validate: %s", err)
	}
}
