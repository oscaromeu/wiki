---
id: functions_1
title: Functions 1 - Declaring Functions
sidebar_label: Functions 1 - Declaring Functions
sidebar_position: 26
---


## The basic function definition

Let's start with the simplest form of a function

```go
func f() {
    fmt.Println("A simple function")
}
func main() {
    f()
}
```

## Function parameters

This function neither takes nor returns any parameter. Now let's add parameters. A parameter consist of a name followed by a type. This is just declaring a variable but without the `var` keyword.

```go
func f(n int) {
    fmt.Println("A function with a parameter:", n)
}
func main() {
    f(37)
}
```

Multiple parameters are separated by a comma

```go
func f(n int, m int, s string) {
    fmt.Println("A function with parameters:", n, m, s)
}
func main() {
    f(37, 29, "ok")
}
```

Subsequent parameters of the same type can be grouped together.

```go
func f(n, m int, s string) {
    fmt.Println("A function with grouped parameters:", n, m, s)
}
func main() {
    f(37, 29, "ok")
}
```
