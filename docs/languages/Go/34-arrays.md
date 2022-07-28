---
id: arrays
title: Arrays
sidebar_label: Arrays
sidebar_position: 34
draft: false
---

Arrays is collection of values of the same type which are stored contiguously in memory. The type `[n]T` is an array of `n` values of type `T`. An array type in Go has two parts:

+ The type of the values that can be inserted, and
+ The maximum number of values that the array can hold.


_The size of an array is indeed part of its type_ which means that you can neither resize an array nor assign an array of different size to an array variable. The point of fixed-sized arrays is that you have control over memory usage. No unexpected memory allocation and no garbage collection happen behind the scenes. 

## Declaring arrays and accessing their elements

The following `planets` array contains exactly eight elements

```go
var planets [8]string
```

Every element of an array has the same type. In this case, `planets` is an array of strings. Individual elements of an array can be accessed by using square brackets `[]` with an index that begins at `0`. 

```go
var planets [8]string

planets[0] = "Mercury"
planets[1] = "Venus"
planets[2] = "Earth"

earth := planets[2]
fmt.Println(earth)
fmt.Println("%#v\n",planets)
fmt.Println(len(planets))
```

Even though only three planets have been assigned, the `planets`array has eight elements. The length of an array can be determined with the built-in `len` function. The other elements contain the zero value for their type, an empty string:

:::info
`%#v` prints information about the variable along with its contents. Try `%v` without the hash sign (`#`) for comparison.)
:::

:::warning
The Go compiler will report an error if it detects access to an element outside the range of the array. 

```go
var planets [8]string

planets[8] = "Pluto"
pluto := planets[8] // Invalid array index 8; out of bounds for 8 element array
```

If the Go compiler is unable to detect the error the program may _panic_ while it's running.

:::


## Initialize arrays with composite literals

A _composite literal_ is a concise syntax to initialize any composite type with the values we want. Rather than declare an array and assign elements one by one, Go's composite literal syntax will declare and initialize an array in a single step.

To initialize an array we can use array literals:

```go
var a [3]int = [3]int{3,6,99}
b := [2]string{"yes", "no"}
```

A composite literal consist of:

+ A type declaration, and
+ A list of initial values within curly braces.

In this case, the size declaration is not necessary, as the length is determined by the list of literals, and can be replaced with an ellipsis

```go
a := [...]int{1,2,3,4,5}
```

If we only want to initialize a few of the values, we can use this syntax:


```go
a := [10]string{0: "First", 9: "Last"}
fmt.Printf("%#v\n", a)
```

This will only set elements 0 and 9 to specific values, and leave all other elements initialized to their zero value.



### Composite literals with line breaks need a final comma

We can divide a composite literal into multiple lines. However, in this case, the final value also requires a comma like all other values:

```go
a = [3]int{
    8,
    16,
    32, // <-- Comma required!
}
```

Otherwise, the compiler detects a possible end of statement at the wrong place and errors out. 
The positive side of this is that you can shuffle those lines around without worrying about adding or removing the final comma.

## Iterating through arrays

We can iterate through each element of an array using the `for` loop or the `range` built-in function.

```go
dwarfs := [5]string{"Ceres", "Pluto", "Haumea", "Makemake", "Eris"}

for i :=0; i < len(dwarfs); i++ {
    dwarf := dwarfs[i]
    fmt.Println(i, dwarf)
}
```

The `range` keyword provides an index and value for each element of an array with less code and less chance for mistakes

```go
dwarfs := [5]string{"Ceres", "Pluto", "Haumea", "Makemake", "Eris"}
for i, dwarf := range dwarfs {
    fmt.Println(i, dwarf)
}
```

:::info
Remember that we can use the blank identifier [underscode] if we don't need the index variable provided by `range`
:::

## Arrays are copied

Assigning an array to a new variable or passing it to a function makes a complete copy of its contents.

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/cuXncjRhhrD">https://go.dev/play/p/cuXncjRhhrD</a></b></figcaption>

```go
	planets := [...]string{
		"Mercury",
		"Venus",
		"Earth",
		"Mars",
		"Jupiter",
		"Saturn",
		"Uranus",
		"Neptune",
	}
	planetsMarkII := planets
	planets[2] = "whoops"
	fmt.Println(planets)
	fmt.Println(planetsMarkII)
```

Arrays are values, and functions pass by value, which means if we pass an array to a function, the function receives a copy of the array, rather than a pointer to it. 

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/v3BHQAVKK1j">https://go.dev/play/p/v3BHQAVKK1j</a></b></figcaption>

```go
package main

import "fmt"

// terraform accomplishes nothing
func terraform(planets [8]string) {
	for i := range planets {
		planets[i] = "New " + planets[i]
	}
}
func main() {
	planets := [...]string{
		"Mercury",
		"Venus",
		"Earth",
		"Mars",
		"Jupiter",
		"Saturn",
		"Uranus",
		"Neptune",
	}
	terraform(planets)
	fmt.Println(planets)
}
```

The `terraform` function is operating on a copy of the `planets` array, so the modifications don't affect `planets` in the `main` function.

Also, it's important to recognize that the length of an array is part of its type. The type `[8]string` and type `[5]string` are both collections of strings, but they're two different types. The Go compiler will report an error when attempting to pass an array of a different length.

## Arrays of arrays

We can declare multidimensional arrays like so

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/yD_1tMN58dr">https://go.dev/play/p/yD_1tMN58dr</a></b></figcaption>

```go
	var board [8][8]string
	board[0][0] = "r"
	board[0][7] = "r"
	for column := range board[1] {
		board[1][column] = "p"
	}
	fmt.Print(board)
```

## Comparing arrays

If the element type of an array is a comparable type, arrays of the same type can be compared as a whole. Comparasion then happens element-wise

```go
a1 := [3]string{"one", "two", "three"}
a2 := [3]string{"one", "two", "three"}
fmt.Println("a1 == a2:", a1 == a2)
```

## Summary

+ Arrays are ordered collections of elements with a fixed length.
+ Composite literals provide a convenient means to initialize arrays.
+ The `range` keyword can iterate over arrays.
+ When accessing elements of an array, we stay inside its boundaries.
+ Arrays are copied when assigned or passed to functions.