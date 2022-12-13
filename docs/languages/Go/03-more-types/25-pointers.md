---
id: pointers
title: Pointers
sidebar_label: Pointers
sidebar_position: 25
---

## Introduction

When a variable gets instantiated, its value is stored in one or more cells of the computer’s main memory. A memory cell can be identified by a unique number called address.

## Creating a pointer variable

In Go, we can create a variable that holds the address of another variable. This variable is called a pointer variable.

Here, we create an integer variable named `a`, and a variable `p` of type "pointer to int".

```go
var a int = 1
var p *int  // The type of p is "pointer to int"
```

## nil: The zero value for pointer types

At this point, variable `p` does not point to any variable yet. In the lecture on variables, we learned that every type has a well-defined zero value. For pointer types, this zero value is called `nil`. Go has the keyword `nil` to express this zero value.

```go
fmt.Println("The zero value of p is:", p)
```

Output:

```go
The zero value of p is: <nil>
```

## The address operator

To assign the address of variable `a` to `p`, we can use the address operator, which is simply an ampersand (`&`) prepended to the variable's name.

When we print `p`, we get the address of `a` as a result.

```go
fmt.Println("p's value is a's address:", p)
```

Output:

```
p's value is a's address: 0xc42007e178
```

### Not everything has an address
Note that the address operator only works with operands that are, well, addressable. Variables are addressable, and we will learn about other addressable entities in later lectures.

Two examples of non-addressable entities are literal values and constants. Trying to get the address of one of these triggers an error at compile time.

```go
// A literal has no address.
p = &123  // error: "cannot take the address of 123"
// A constant has no address.
const c int = 123
p = &c  // error: "cannot take the address of c"
```

## Retrieving a value through a pointer

Since `p` holds the address of `a`, we can use this address to retrieve the value stored in `a`. This operation is called __pointer indirection__. Technically, we simply prepend an asterisk to the pointer name: `*p` - this yields the value of `a`.

Output:

```go
*p yields a's value: 1
```

We just have retrieved the value of `a` without using `a`'s name at all.

The pointer `p` is not bound to `a` for the rest of its lifetime. We can reassign the pointer to another variable, say, `b`.

```go
b := 2
p = &b
fmt.Println("Now p contains b's address:", p, "and *p is:", *p)
```

Output:

```go
Now p contains b's address: 0xc42007e1a0 and *p is: 2
```

Be aware that pointer indirection fails if the pointer has the value `nil`. This makes sense because no target address exists to fetch a value from. Applying pointer indirection to a nil pointer causes a panic at runtime.

```go
var p *int  // p is nil by default
fmt.Println("A nil pointer:", p, *p)  // "panic: runtime error: invalid memory address or nil pointer dereference"
```

### Pointer "side effect"

When two pointers point to the same variable and the variable is updated through one pointer, then accessing the variable through the other pointer retrieves the updated value.

```go
p1 := &a
p2 := &a
fmt.Println("p1:", p1)
fmt.Println("p2:", p2)
fmt.Println("*p2 is:", *p2)
*p1 = 4
fmt.Println("*p2 is:", *p2)
```

This might seem obvious - it’s the same variable, after all -, but in later lectures we will see situations where pointers are hidden within data types, and these are the situations where reasoning about the effects of updating a value can get tricky. We will discuss this in more detail in the lectures about more complex data types.

## Pointer comparisons

Pointers can be compared to each other at two levels:

+ Wheter the addresses they store are equal, and wheter the values stored at these addresses are equal

```go
    fmt.Println("p1 == p2:", p1 == p2)
    fmt.Println("*p1 == *p2:", *p1 == *p2)
```

+ A pointer can also be compared to `nil`. This is a rather common test before attempting to apply a pointer indirection.

```go
    p = nil
    fmt.Println("p == nil:", p == nil)
    p = &a
    if p != nil {
        *p = 5
    }
```

:::info Remember
A pointer holds the memory address of a value.

The type `*T` is a pointer to a `T` value. Its zero value is nil.

```go
var p *int
```

The & operator generates a pointer to its operand.

```go
i := 42
p = &i
```

The * operator denotes the pointer's underlying value.

```go
fmt.Println(*p) // read i through the pointer p
*p = 21         // set i through the pointer p
```


This is known as "dereferencing" or "indirecting".

Unlike C, Go has no pointer arithmetic.
:::

## The new function

The built-in function `new` instantiates a new variable of a given type and returns a pointer to that variable.

```go
p = new(int)
fmt.Println("p points to an unnamed int:", p, *p)
*p = 6
fmt.Println("The unnamed int has been changed to:", *p)
```

In this case, the new variable has no name, and therefore the pointer is the only way to access the variable.

## No pointer arithmetic

Pointers in Go do not allow any kind of pointer arithmetic, and the automatic memory management prevents pointers from becoming stale.

```go
p = &a + 64  // error: "invalid operation: &a + 64 (mismatched types *int and int)"
```

## Pointers have a static type

+ A pointer always has a specific type that is bound to the type of the variable that the pointer points to.

+ A pointer of type `*int` can only point to variables of type `int`, and a pointer of type `*string` can only point to string variables.

There is no void pointer type like in C or C++.

```go
var a int64
p := &a  // p's type is now *int64
var b int32
p = &b  // fails
var s string
var t *string = &s
```

## nil has no type

The keyword `nil` represents no specific type of nil pointer. You can assign it to any pointer type...

```go
var s *string
var f *float64
s = nil
f = nil
```

…but you cannot use nil in a short declaration:

```go
p := nil  // what is the type of p now?
```

## Summary
+ Pointer variables contain the address of a memory location that holds data of a given type.
+ Pointer indirection yields the value stored at this address.
+ The zero value of a pointer of any type is nil.
+ Attempting to access the value of a nil pointer causes a runtime panic.
+ Pointers can be compared to each other and to nil.
+ Automatic memory management and the lack of pointer arithmetic make pointers safe to use.