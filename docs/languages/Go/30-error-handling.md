---
id: error_handling
title: Error handling
sidebar_label: Error handling
sidebar_position: 30
draft: true
---

## Creating an error

A function that detects a failure can create a value of type `error` and return this value to the caller. In its simplest form, an error object just contains an error message. A call to `fmt.Errorf()` creates this kind of error.

```go
func verify(i int) error {
    if i < 0 || i > 10 {
        return fmt.Errorf("verify: %d is outside the allowed range (0..10)", i)
    }
    return nil
}
```

If no error ocurred, the function can pass back a `nil` value. An `error` object is not a pointer type; rather it is an interface type, but interfaces can also become `nil`.


By convention, if a function already has one or more return values, the error value is the last one in the list of return values.

The caller then can inspect the returned error and act accordingly.

```go
err := verify(12)
if err != nil {
    // process the error
}
```

## Error handling strategies

### Strategie 1: Propagate the error to the next caller

When a function receives an error that it cannot handle by itself, it can pass the error back to its caller. However, ensure to never pass the naked error up the call chain. Instead, add relevant context information to the error. Go makes this easy: All you need to do is create another error, and in the format string, use the verb `"%w"` for including the original error. 

Then pass the `err` variable as an argument matching the position of the `%w` verb. In this example, function `propagate()` may get an error back from calling function `verify()`. In this case, `propagate()` creates a new error that includes the current value of `i` as the context information, and the original error.

```go
func propagate(i int) error {
    if err := verify(i); err != nil {
        return fmt.Errorf("propagate(%d): %w", i, err)
    }
    return nil
}
```

As a tip, if you avoid any newline character in the error messages, the final concatenated error message consists of one single line. This allows tools like grep to easily retrieve a complete error message from a log file, rather than only one part of it.

Note how this code example uses the "if action; check error" idiom that we discussed in the lecture on the if statement.

### Strategy 2: Retry

Some errors are only transient, like a failing network connection. In this case, the caller may attempt one or more retries before declaring failure.

```go
func retry(i int) error {
    err := propagate(i)
    if err != nil {
        err = propagate(i / 2)
        if err != nil {
            return fmt.Errorf("retry: %s", err)
        }
    }
    return nil
}
```

Our simple demo code just tries to call the failed function again but this time with half the value than before.

Real-life examples of retry strategies are:

+ to wait for a defined time before retrying,
+ to try a different port, file name, IP address, etc., or
+ to fall back to default values.

### Strategy 3: Log the error and continue

Sometimes, an error might be insignificant enough to not justify disturbing the normal program flow. In this case, the handling function might just log the incident and continue.

```go
func onlyLog(i int) {
    if err := retry(i); err != nil {
        log.Println("onlyLog:", err)
    }
}
```

The log package from the standard library provides some convenient logging functions. By default, these functions prepend time and date to the log messages. You can change this through the functions `log.SetPrefix()` and `log.SetFlags()`.

```go
log.SetPrefix("")
log.SetPrefix("module x: ")

log.SetFlags(0)
log.SetFlags(log.Lshortfile | logLdate)
```

If you are writing a library, be aware that libraries usually should not roll their own logging. Instead, they should always return any errors to the caller. The main application should be the only instance that implements a logging strategy. This is the only way to ensure a consistent logging strategy throughout an application.

### Strategy 4: Log and exit

Severe errors may require to stop the application immediately, to prevent further damage. For this situation, the log package provides the methods `log.Fatal()`, `log.Fatalln()`, and `log.Fatalf()` that act like their `fmt.Print/ln/f` counterparts except that they immediately exit the running application.

```go
func logAndExit(i int) {
    err := retry(i)
    if err != nil {
        log.Fatalln("exit:", err)
    }
}
```

## Panic and Recover

In some cases, a function is expected to succeed given that all preconditions are met. Such a function may only fail in one of two cases.

+ Case one: something unpredictable happens, like running out of memory. Usually, the function cannot properly recover from such an event.

+ Case two: the caller passed invalid arguments to the function or called the function in a context where it clearly must not have called that function. For example, if a function expects to receive a valid pointer, the caller must not pass a nil pointer instead.

In this case, the failing function can raise a panic.

```go
func unexpectedError(p *int) {
    if p == nil {
        panic("p must not be nil")
    }
}

func main() {
    unexpectedError(nil)
}
```

Go's runtime also raises panics if it detects illegal actions like dereferencing a nil pointer or accessing an array beyond its bounds.

```go
func unexpectedError(p *int) {
    *p = 1
}

func main() {
    unexpectedError(nil)
}
```

```go
func main() {
    var a [10]int
    fmt.Println(a[20])
}
```

### What happens when a panic is raised?

First, the current function executes any deferred calls (see the lecture on functions about deferred calls) and then exits immediately. The calling function also just executes any deferred calls and exits. This continues until either the program exits or one of the deferred calls in the call chain successfully recovers from the panic.

When a panic terminates a program, a panic message is printed out, along with a call stack for each goroutine. The call stack allows tracking down the cause of the panic. The topmost entry of the stack usually lists the function that panicked

```bash
panic: p must not be nil

goroutine 1 [running]:
main.unexpectedError(0x0)
    /Users/foo/temp/panic.go:5 +0x78
main.main()
    /Users/foo/temp/panic.go:10 +0x2a
exit status 2
```

## Recover

A deferred call may call the built-in function `recover()` to stop the panic. Then it may attempt to restore a stable, defined state, and continue execution from there.

```go
func main() {
    defer func() {
        res := recover()
        if res != nil {
            fmt.Println("Recovered from a panic:", res)
        }
    }()

    unexpectedError()
}
```

In most cases, however, this is not possible, so the usual option is to let the program exit. This is kind of a “controlled crash”, and usually it is better to let a program crash early than to continue execution and to risk processing broken data.

## Summary

+ In Go, error handling is part of the normal program flow.
+ Expected errors are simply returned to the caller.
+ The caller should check each returned error value before calling another function.
+ Error handling strategies: retry, propagate, log and continue, log and exit
+ Unexpected errors and very serious failures can trigger a panic.
+ A panic leads to a controlled crash.
+ Deferred calls can attempt to recover from a panic.
