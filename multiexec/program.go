package multiexec

// Program describes a program to be run.
// This is just like the main function of a normal go binary, except that it
// has an additional parameter for the context within which the program is
// started.
type Program func(*Context)
