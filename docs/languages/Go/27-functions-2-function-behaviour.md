---
id: functions_2
title: Function Behaviour - recursion, deferred functions, scope
sidebar_label: Functions 2 - Function Behaviour
sidebar_position: 27
---

## Recursion

A function can call itself - this is called a recursive call

```go
func rec() {
    rec()
}
```

This function would never return, so we need to add some logic to break the recursive calls at one point. 

A common practice is to include a parameter that counts downwards on each call. Then the function can stop the recursive calls if this parameter is zero. 

```go
// The classic faculty example:
// The faculty of n (or "n!") is defined as: 
// 1 if n is 0
// n * (n-1)! if n > 0

func faculty(n int) int {
    if n == 0 {
        return 1
    }
    return n * faculty(n - 1)
}
```

## Deferred function call

Go's `defer` statement schedules a function call (the deferred function) to be run immediately before the function executing the defer returns. A defer statement pushes a function call onto a list. The list of saved calls is executed after the surrounding function returns. Defer is commonly used to simplify functions that perform various clean-up actions.

For example, let’s look at a function that opens two files and copies the contents of one file to the other:

```go
func CopyFile(dstName, srcName string) (written int64, err error) {
    src, err := os.Open(srcName)
    if err != nil {
        return
    }

    dst, err := os.Create(dstName)
    if err != nil {
        return
    }

    written, err = io.Copy(dst, src)
    dst.Close()
    src.Close()
    return
}
```

This works, but there is a bug. If the call to `os.Create` fails, the function will return without closing the source file. This can be easily remedied by putting a call to `src.Close` before the second return statement, but if the function were more complex the problem might not be so easily noticed and resolved. By introducing defer statements we can ensure that the files are always closed:

```go
func CopyFile(dstName, srcName string) (written int64, err error) {
    src, err := os.Open(srcName)
    if err != nil {
        return
    }
    defer src.Close()

    dst, err := os.Create(dstName)
    if err != nil {
        return
    }
    defer dst.Close()

    return io.Copy(dst, src)
}
```

:::info
Defer statements allow us to think about closing each file right after opening it, guaranteeing that, regardless of the number of return statements in the function, the files will be closed.
:::

The behavior of defer statements is straightforward and predictable. There are three simple rules:

1. A deferred function's arguments are evaluated when the defer statement is evaluated.

In this example, the expression "i" is evaluated when the Println call is deferred. The deferred call will print "0" after the function returns.

```go
func a() {
    i := 0
    defer fmt.Println(i)
    i++
    return
}
```

2. Deferred function calls are executed in Last In First Out order after the surrounding function returns.

This function prints "3210":

```go
func b() {
    for i := 0; i < 4; i++ {
        defer fmt.Print(i)
    }
}
```

3. Deferred functions may read and assign to the returning function’s named return values.

In this example, a deferred function increments the return value i after the surrounding function returns. Thus, this function returns 2:

```go
func c() (i int) {
    defer func() { i++ }()
    return 1
}
```

This is convenient for modifying the error return value of a function; we will see an example of this shortly.

## Scope

With functions, a fundamental concept becomes relevant: The concept about the scope of variables. 

We discuss scope here as part of the lectures on functions, but in fact, there are also other levels of variable scope. But let’s start with function-level scope.

### Function-level scope

A variable declared in a function is only visible within that function. Neither the calling function nor a function called by the current function can access this variable. 

```go
func f() {
    a := 2  // only visible within f
    fmt.Println("func f: a is", a)
}

func main() {
    a := 1  // not visible within f
    fmt.Println("main: a is", a)
    f()
    fmt.Println("main: a is", a)
}
```

### Block scope

Two curly braces define a block that has its own scope. You can declare variables that have the same name as variables outside the braces. The variable inside the block then shadows the variable outside.

```go
func f() {
    a := 1
    {
        a := 2
        fmt.Println("in block: a is", a)
    }
    fmt.Println("func f: a is", a)
}

func main() {
    f()
}
```

This behavior also applies to `if`, `switch`, `select`, or `loop` blocks.

### The scope of variables in an if, switch, or loop condition

Variables that are defined in the conditional clause of an `if`, `switch`, or `loop` block already live in a new scope.

For example, the variable in this `if` condition does not exist after the `if` block.

```go
if ret := f(); ret == 0 {
    // do something
}
// ret not accessible here
```

Furthermore, in a loop body, the scope is even limited to the current iteration.

To illustrate this, consider this loop:

```go
for i := 0; i < 4; i++ {
    n := i
    fmt.Println(&i, &n) 
}
```

On each iteration, a new instance of `n` is declared, while `i` remains the same instance, as we can see from the memory locations of `i` and `n` at each iteration.

If `n` is a large data type with complex initialization, it is usually more efficient to declare `n` outside the loop, provided that the same instance of `n` can safely be reused across iterations.

```go
var n int
for i := 0; i < 4; i++ {
    n = i // Note: no := here, only =
    fmt.Println(&i, &n) 
}
```

## Summary

The takeaways:

+ Go functions support recursion.
+ The defer keyword defers a function call until the calling function exits.
+ Variables are visible only in the function or block they are declared in. 
+ Variables in a loop condition live across all loop iterations, whereas variables declared inside the loop body only live within a single iteration. 