---
id: functions_4
title: Function and pointers
sidebar_label: Functions 4 - Function and pointers
sidebar_position: 29
---

## Go has pass-by-value semantics

In Go, function parameters are always passed by _value_. That is, when passing a variable to a function, then the variable's value gets copied into the function, and any change to this value inside the function has no effect on the original value. 



<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/zz-n5AtGOE4">https://go.dev/play/p/zz-n5AtGOE4</a></b></figcaption>

```go
func f(a int) {
    a = a + 1
}
func main() {
    x := 1
    f(x)           
    fmt.Println(x) 
}
```

A pointer is also passed by value; that is, when the receiving function changes the local pointer `a` to point to a different variable, this does not change the original pointer `p`, nor the value of `x` that `p` points to. 

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/P7DS-ZIkrqW">https://go.dev/play/p/P7DS-ZIkrqW</a></b></figcaption>

```go
func f(a *int) {
    b := 2
    a = &b 
}
func main() {
    x := 1
    p := &x // to make clear that we deal with two pointers
    f(p)           
    fmt.Println(x) 
}
```

## Pointer indirection

Through pointer indirection we can manipulate the original value. As long as `a` points to `x`, we can manipulate the value of `x` inside the function altough `x` lives outside of the function. 



<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/bMOx4RX6uQS">https://go.dev/play/p/bMOx4RX6uQS</a></b></figcaption>

```go
func f(a *int) {
    *a = *a + 1
}
func main() {
    x := 1
    p := &x 
    f(p)           
    fmt.Println(x) 
}
```

## Returning a pointer to a function-local variable

Go allows returning a pointer to a variable that was created inside the function.

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/2yRZ6DaWlT9">https://go.dev/play/p/2yRZ6DaWlT9</a></b></figcaption>

```go
func f() *int {
    a := 7
    return &a
}
func main() {
    p := f()           
    fmt.Println(*p) 
}
```

## Summary

+ Pass-by-value semantics: When a variable is passed to a function, the function only receives a copy of the value.
+ Even a pointer parameter is only a copy of the original pointer.
+ However, through pointer indirection, the original value can be changed from inside the function.
+ In Go it is perfectly legal that a function returns a pointer to one of its local variables. 