package multiexec

// Context is the context within which a Program is started. A context handles
// access of programs to per-process resources since all programs bundled in a
// binary share the same process.
type Context struct {
	// Config is the configuration with which this context was created and the
	// program was started.
	Config Config

	// Name is the name with which the program was started.
	Name string
}
