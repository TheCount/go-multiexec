package multiexec

import (
	"fmt"
)

// PanicRecovery is an enumeration type for the different types of panic
// recovery.
type PanicRecovery uint8

// PanicRecovery types
const (
	// PanicEnd means that if a program panics, it ends without affecting other
	// programs. The fact that a panic occurred is logged.
	// This does not work with panicking goroutines created by the program's main
	// goroutine.
	PanicEnd PanicRecovery = iota

	// PanicCrash means that if a program panics, the entire process dies.
	PanicCrash

	// PanicRestart means that if a program panics, it is nevertheless restarted
	// according to its restart policy.
	PanicRestart

	// numPanicRecoveryTypes is the number of panic recovery types.
	numPanicRecoveryTypes
)

// Config is the configuration describing how a program should be started.
type Config struct {
	// PanicRecovery describes what should be done when the program panics.
	PanicRecovery PanicRecovery

	// Restart determines how the program should be restarted when it finishes.
	// If Restart is nil, the program is not restarted after it finishes.
	Restart RestartPolicy
}

// DefaultConfig returns a default configuration for a program to be run at most
// once, ending it if it panics.
func DefaultConfig() Config {
	return Config{
		PanicRecovery: PanicEnd,
	}
}

// Validate validates this configuration. It returns the first error found
// in the configuration, or nil if the configuration is without error.
func (cfg Config) Validate() error {
	if cfg.PanicRecovery >= numPanicRecoveryTypes {
		return fmt.Errorf("Invalid panic recovery type %d", cfg.PanicRecovery)
	}
	return nil
}
