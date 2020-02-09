package multiexec

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRestartSimple tests restarting a simple counting program.
func TestRestartSimple(t *testing.T) {
	const numRestarts = 2
	const name = "foo"
	os.Args[0] = filepath.Join("/usr/bin", name)
	cfg := DefaultConfig()
	cfg.Restart = RestartNTimes(numRestarts, 0)
	b := NewBundle()
	starts := 0
	if err := b.AddProgram(name, func(*Context) {
		starts++
	}, cfg); err != nil {
		t.Fatalf("Unable to add simple counting program: %s", err)
	}
	b.Wait()
	if starts != numRestarts+1 {
		t.Errorf("Expected %d restarts, got %d", numRestarts, starts-1)
	}
}

// TestRestartPanic tests restarting a panicky counting program.
func TestRestartPanic(t *testing.T) {
	const numRestarts = 2
	const name = "bar"
	os.Args[0] = filepath.Join("/usr/bin", name)
	cfg := DefaultConfig()
	cfg.Restart = RestartNTimes(numRestarts, 0)
	b := NewBundle()
	starts := 0
	if err := b.AddProgram(name, func(*Context) {
		starts++
		panic(panicMessage)
	}, cfg); err != nil {
		t.Fatalf("Unable to add panicky counting program: %s", err)
	}
	b.Wait()
	if starts != 1 {
		t.Errorf("Expected no restarts, got %d", starts-1)
	}
}
