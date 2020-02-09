# Multiple go programs in one binary

![](https://github.com/TheCount/go-multiexec/workflows/CI/badge.svg)
[![Documentation](https://godoc.org/github.com/TheCount/go-multiexec/multiexec?status.svg)](https://godoc.org/github.com/TheCount/go-multiexec/multiexec)

## Example

```golang
package main

import(
  "log"
  "os"

  "somepackage/cmd/a" // program entry points of the form
  "somepackage/cmd/b" // func(*multiexec.Context) in external packages.

  "github.com/TheCount/go-multiexec/multiexec"
)

func main() {
  // Create a program bundle and add program entry points from packages a and b.
  cfg := multiexec.DefaultConfig()
  bundle := multiexec.NewBundle()
  if err := bundle.AddProgram("a", a.EntryPoint, cfg); err != nil {
    log.Fatal(err)
  }
  if err := bundle.AddProgram("b", b.EntryPoint, cfg); err != nil {
    log.Fatal(err)
  }
  // a or b are only run if this binary was called as "a" or "b"

  // Wait for execution to finish
  exits := bundle.Wait()
  if len(exits) == 0 {
    // Nothing was run, binary was not executed as "a" or "b".
    log.Fatalf("Unknown command: %s", os.Args[0])
  }
  exit := exits[0]
  log.Printf("Command '%s' finished", exit.Context.Name)
  if exit.Reason == nil {
    log.Print("Program ended normally")
  } else {
    log.Printf("Program ended with error: %s", exit.Reason)
  }
}
```
