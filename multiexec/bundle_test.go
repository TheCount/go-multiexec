package multiexec

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	// simpleMessage is a message printed by the simpleTestProgram.
	simpleMessage = "Hello, world!"

	// panicMessage is the value the panickingTestProgram raises a panic with.
	panicMessage = "Goodbye, world!"
)

// invalidConfig is an invalid configuration.
var invalidConfig = Config{
	PanicRecovery: 42,
}

// simpleTestProgram is a simple program which exits normally.
func simpleTestProgram(*Context) {
	println(simpleMessage)
}

// panickyTestProgram is a program which immediately starts panicking.
func panickyTestProgram(*Context) {
	panic(panicMessage)
}

// TestEmptyBundle tests waiting on an empty bundle.
func TestEmptyBundle(t *testing.T) {
	b := NewBundle()
	exits := b.Wait()
	if len(exits) > 0 {
		t.Error("Expected zero exits on empty bundle")
	}
}

// TestNilProgram tests adding a nil program to a bundle.
func TestNilProgram(t *testing.T) {
	b := NewBundle()
	if err := b.AddProgram("foo", nil, DefaultConfig()); err == nil {
		t.Error("Expected error on adding a nil program")
	}
}

// TestInvalidConfig tests adding a program with an invalid configuration to
// a bundle.
func TestInvalidConfig(t *testing.T) {
	b := NewBundle()
	if err := b.AddProgram("foo", simpleTestProgram, invalidConfig); err == nil {
		t.Error("Expected error on program with invalid configuration")
	}
}

// TestBadBasename tests adding a program with a bad basename.
func TestBadBasename(t *testing.T) {
	os.Args[0] = "bad"
	b := NewBundle()
	if err := b.AddProgram(
		"foo", simpleTestProgram, DefaultConfig(),
	); err != nil {
		t.Fatalf("Unable to add simple test program in bad basename test: %s", err)
	}
	if b.numRunning != 0 {
		t.Error("Bad basename program was not ignored")
	}
}

// TestEmptyBasename tests adding a program when someone has tampered with the
// program name and set it to an empty string.
func TestEmptyBasename(t *testing.T) {
	os.Args[0] = "" // Nyahaha
	b := NewBundle()
	if err := b.AddProgram(
		"foo", simpleTestProgram, DefaultConfig(),
	); err != nil {
		t.Fatalf(
			"Unable to add simple test program in empty basename test: %s", err,
		)
	}
	if b.numRunning != 0 {
		t.Error("Empty basename program was not ignored")
	}
}

// TestSimpleProgram tests running a simple program.
func TestSimpleProgram(t *testing.T) {
	const name = "foo"
	os.Args[0] = filepath.Join("/usr/bin", name)
	b := NewBundle()
	if err := b.AddProgram(name, simpleTestProgram, DefaultConfig()); err != nil {
		t.Fatalf("Unable to add simple test program: %s", err)
	}
	exits := b.Wait()
	if len(exits) != 1 {
		t.Fatal("Expected one exit for simple program test")
	}
	exit := exits[0]
	if exit.Context.Name != name {
		t.Errorf("Bad exit name in simple program test: %s", exit.Context.Name)
	}
	if exit.Reason != nil {
		t.Errorf("Simple program exited with reason: %s", exit.Reason)
	}
}

// TestPanickyProgram tests running a program which panics.
func TestPanickyProgram(t *testing.T) {
	const name = "bar"
	os.Args[0] = filepath.Join("/usr/bin", name)
	b := NewBundle()
	if err := b.AddProgram(
		name, panickyTestProgram, DefaultConfig(),
	); err != nil {
		t.Fatalf("Unable to add panicky test program: %s", err)
	}
	exits := b.Wait()
	if len(exits) != 1 {
		t.Fatal("Expected one exit for panicky program test")
	}
	exit := exits[0]
	if exit.Context.Name != name {
		t.Errorf("Bad exit name in panicky program test: %s", exit.Context.Name)
	}
	if exit.Reason != panicMessage {
		t.Errorf("Panicky program exited with wrong reason: %s", exit.Reason)
	}
}
