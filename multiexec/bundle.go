package multiexec

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ProgramExit describes an exiting program.
type ProgramExit struct {
	// Context is the context of the exiting program.
	Context *Context

	// Reason is the reason the program exited.
	// If the program returned normally, this is nil.
	// If the program panicked, the reason holds the return value of recover().
	Reason interface{}
}

// programExitChan is a channel type through which a program exit can be
// reported.
type programExitChan chan ProgramExit

// Bundle represents multiple program executions.
type Bundle struct {
	// args is a copy of os.Args at bundle creation.
	args []string

	// numRunning is the number of currently running programs.
	numRunning int

	// exit is the channel through which program exits are reported.
	exit programExitChan
}

// NewBundle creates a new program bundle.
func NewBundle() *Bundle {
	args := make([]string, len(os.Args))
	copy(args, os.Args)
	return &Bundle{
		args: args,
		exit: make(programExitChan),
	}
}

// AddProgram adds the specified program with the specified configuration
// under the specified name to this bundle.
// If name does not match the call basename, AddProgram does nothing.
// Otherwise, the program is started in a separate goroutine.
// AddProgram is not concurrency-safe and should only be called from the
// goroutine this bundle was created in, usually the main goroutine.
func (b *Bundle) AddProgram(name string, prog Program, cfg Config) error {
	if prog == nil {
		return fmt.Errorf("Program supplied for '%s' is nil", name)
	}
	if err := cfg.Validate(); err != nil {
		return err
	}
	if b.basename() != name {
		return nil
	}
	ctx := &Context{
		Config:  cfg,
		Name:    name,
		program: prog,
	}
	b.startProgram(ctx, 0)
	return nil
}

// basename returns the basename of the file name with which the binary in which
// this bundle was created was called (based on os.Args at bundle creation).
func (b *Bundle) basename() string {
	result := filepath.Base(b.args[0])
	if result == "." || result == string(filepath.Separator) {
		// This should not happen unless someone altered os.Args
		result = ""
	}
	return result
}

// start program starts the specified program with the specified context in
// a separate goroutine after the specified delay.
func (b *Bundle) startProgram(ctx *Context, delay time.Duration) {
	b.numRunning++
	go func() {
		if ctx.Config.PanicRecovery != PanicCrash {
			defer func() {
				if r := recover(); r != nil {
					b.exit <- ProgramExit{
						Context: ctx,
						Reason:  r,
					}
				}
			}()
		}
		time.Sleep(delay)
		ctx.program(ctx)
		b.exit <- ProgramExit{
			Context: ctx,
		}
	}()
}

// Wait blocks until all programs in this bundle have finished executing.
// It returns the program exits, which may be nil or empty if no programs
// were running. Note that programs ignored by AddProgram (e. g., due to
// basename mismatch) will not have a corresponding program exit.
// Wait is not concurrency-safe and should only be called from the
// goroutine this bundle was created in, usually the main goroutine.
// Once Wait has returned, this bundle can be reused by adding new programs
// and calling wait again.
func (b *Bundle) Wait() []ProgramExit {
	result := make([]ProgramExit, 0, b.numRunning)
	for b.numRunning > 0 {
		exit := <-b.exit
		b.numRunning--
		result = append(result, exit)
	}
	return result
}
