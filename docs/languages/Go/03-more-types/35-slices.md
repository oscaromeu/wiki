---
id: slices
title: Slices
sidebar_label: Slices
sidebar_position: 35
hide_title: false
draft: false
---

## Introduction

An array has a fixed size. A slice, on the other hand, is a dynamically-sized, flexible view into the elements of an array. In practice, slices are much more common than arrays.

The type `[]T` is a slice with elements of type `T`.

A slice is formed by specifying two indices, a low and high bound, separated by a colon:

```
a[low : high]
```

This selects a half-open range which includes the first element, but excludes the last one. The following expression creates a slice which includes elements 1 through 3 of `primes`

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/VME38Ao4-Q5">https://go.dev/play/p/VME38Ao4-Q5</a></b></figcaption>

```go
primes := [6]int{2, 3, 5, 7, 11, 13}
var s []int = primes[1:4]
```

## Slices are like references to arrays

A slice does not store any data, it just describes a section of an underlying array. Changing the elements of a slice modifies the corresponding elements of its underlying array. Other slices that share the same underlying array will see those changes.

```go
package main

import "fmt"

func main() {
	names := [4]string{
		"John",
		"Paul",
		"George",
		"Ringo",
	}
	fmt.Println(names)

	a := names[0:2]
	b := names[1:3]
	fmt.Println(a, b)

	b[0] = "XXX"
	fmt.Println(a, b)
	fmt.Println(names)
}
```

## Slice literals

A slice literal is like an array literal without the length. This is an array literal:

```go
[3]bool{true, true, false}
```

And this creates the same array as above, then builds a slice that references it:

```go
[]bool{true, true, false}
```

## Creating a slice from scratch

There are three ways of actually instantiating a slice:

+ Using an slice literal

```go
s = []int{}  // s == []
s = []int{1, 2, 3}
```

+ A slice can be created with the built-in function called `make`, which has the signature,

```go
func make([]T, len, cap) []T
```

where T stands for the element type of the slice to be created. The `make` function takes a type, a length, and an optional capacity. When called, `make` allocates an array and returns a slice that refers to that array. 

+ The third way is to call `append()` on a nil slice

```go
append(s, 1, 2, 3)
```

## Slice defaults

When slicing, we may omit the high or low bounds to use their defaults instead.
The default is zero for the low bound and the length of the slice for the high bound. Consider the following array

```go
var a [10]int
```

these slice expressions are equivalent:

```go
a[0:10]
a[:10]
a[0:]
a[:]
```
## Slice internals

A slice is a descriptor of an array segment. It consists of :

+ The length of the segment
+ The capacity (the maximum length of the segment)
+ A pointer to the data in the underlying array.

![Slice structure](./img/35-01-slice-struct.png)

A variable `s` created with `make([]byte,5)`, is structured like this:

![Slice structure 2](./img/35-02-slice-struct.png)


+ The length of a slice is the number of elements it contains. 

+ The capacity of a slice is the number of elements in the underlying array, counting from the first element in the slice.

+ The length and capacity of a slice `s` can be obtained using the expressions `len(s)` and `cap(s)`.

Now, observe the changes in the slice data structure and their relation to the underlying array as we slice `s`:

```go
s = s[2:4]
```

![Slice structure 3](./img/35-03-slice-struct.png)

Slicing does not copy the slice's data instead it creates a new slice value that points to the original array. This makes slice operations as efficient as manipulating array indices. Therefore, modifying the elements (not the slice itself) of a re-slice modifies the elements of the original slice:

```go
d := []byte{'r', 'o', 'a', 'd'}
e := d[2:]
// e == []byte{'a', 'd'}
e[1] = 'm'
// e == []byte{'a', 'm'}
// d == []byte{'r', 'o', 'a', 'm'}
```

Earlier we sliced s to a length shorter than its capacity. We can grow s to its capacity by slicing it again:

```go
s = s[:cap(s)]
```

![Slice structure 4](./img/35-04-slice-struct.png)

A slice cannot be grown beyond its capacity. Attempting to do so will cause a runtime panic, just as when indexing outside the bounds of a slice or array. Similarly, slices cannot be re-sliced below zero to access earlier elements in the array.

## Growing slices (the copy and append functions)

To increase the capacity of a slice one must create a new, larger slice and copy the contents of the original slice into it. This technique is how dynamic array implementations from other languages work behind the scenes. The next example doubles the capacity of `s` by making a new slice, `t`, copying the contents of `s` into `t`, and then assigning the slice value `t` to `s`:

```go
t := make([]byte, len(s), (cap(s)+1)*2) // +1 in case cap(s) == 0
for i := range s {
        t[i] = s[i]
}
s = t
```

### Copy function

The looping piece of this common operation is made easier by the built-in `copy` function. As the name suggests, copy copies data from a source slice to a destination slice. It returns the number of elements copied.

```go
func copy(dst, src []T) int
```

The `copy` function supports copying between slices of different lengths (it will copy only up to the smaller number of elements). In addition, `copy` can handle source and destination slices that share the same underlying array, handling overlapping slices correctly.

Using `copy`, we can simplify the code snippet above:

```go
t := make([]byte, len(s), (cap(s)+1)*2)
copy(t, s)
s = t
```

### Append function

A common operation is to append data to the end of a slice. This function appends byte elements to a slice of bytes, growing the slice if necessary, and returns the updated slice value:

```go
func AppendByte(slice []byte, data ...byte) []byte {
    m := len(slice)
    n := m + len(data)
    if n > cap(slice) { // if necessary, reallocate
        // allocate double what's needed, for future growth.
        newSlice := make([]byte, (n+1)*2)
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[0:n]
    copy(slice[m:n], data)
    return slice
}
```

One could use `AppendByte` like this:

```go
p := []byte{2, 3, 5}
p = AppendByte(p, 7, 11, 13)
// p == []byte{2, 3, 5, 7, 11, 13}
```

But most programs don't need complete control, so Go provides a built-in append function that's good for most purposes; it has the signature

```go
func append(s []T, x ...T) []T
```

The `append` function appends the elements `x` to the end of the slice `s`, and grows the slice if a greater capacity is needed.

```go
a := [8]int{} // new, zero-filled array
s := a[:7]
fmt.Printf("&a: %p\n", &a[0])
fmt.Printf("len: %2d, cap: %2d, &s: %p s: %v\n", len(s), cap(s), &s[0], s)
for i := 1; i <= 4; i++ {
    s = append(s, i)
    fmt.Printf("len: %2d, cap: %2d, &s: %p s: %v\n", len(s), cap(s), &s[0], s)
}
```

After the second iteration, the slice has reached its capacity, and so the next `append()` allocates a new array and copies the contents of s over to the new array, as we can tell from the output: 

```
&a: 0xc000126000
len:  7, cap:  8, &s: 0xc000126000 s: [0 0 0 0 0 0 0]
len:  8, cap:  8, &s: 0xc000126000 s: [0 0 0 0 0 0 0 1]
len:  9, cap: 16, &s: 0xc00012e000 s: [0 0 0 0 0 0 0 1 2]
len: 10, cap: 16, &s: 0xc00012e000 s: [0 0 0 0 0 0 0 1 2 3]
len: 11, cap: 16, &s: 0xc00012e000 s: [0 0 0 0 0 0 0 1 2 3 4]
```

+ The capacity has risen from 8 to 16, and
+ the first element of `s` is now at a different adress.

To append one slice to another, use `...` to expand the second argument to a list of arguments.

```go
a := []string{"John", "Paul"}
b := []string{"George", "Ringo", "Pete"}
a = append(a, b...) // equivalent to "append(a, b[0], b[1], b[2])"
// a == []string{"John", "Paul", "George", "Ringo", "Pete"}
```

Since the zero value of a slice (`nil`) acts like a zero-length slice, you can declare a slice variable and then append to it in a loop:

```go
// Filter returns a new slice holding only
// the elements of s that satisfy fn()
func Filter(s []int, fn func(int) bool) []int {
    var p []int // == nil
    for _, v := range s {
        if fn(v) {
            p = append(p, v)
        }
    }
    return p
}
```

### Re-slice the slice

There is another option to extend the slice up to its capacity. 

```go
a := [8]int{}
s := a[3:6]    // len(s) == 3, cap(s) == 5
s = s[:1]      // len(s) == 2, cap(s) == 5
s = s[:cap(s)] // len(s) == 5, cap(s) == 5
```
### A possible "gotcha"

As mentioned earlier, re-slicing a slice doesn't make a copy of the underlying array. The full array will be kept in memory until it is no longer referenced. Occasionally this can cause the program to hold all the data in memory when only a small piece of it is needed.

For example, this `FindDigits` function loads a file into memory and searches it for the first group of consecutive numeric digits, returning them as a new slice.

```go
var digitRegexp = regexp.MustCompile("[0-9]+")

func FindDigits(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    return digitRegexp.Find(b)
}
```

This code behaves as advertised, but the returned `[]byte` points into an array containing the entire file. Since the slice references the original array, as long as the slice is kept around the garbage collector can’t release the array; the few useful bytes of the file keep the entire contents in memory.

To fix this problem one can copy the interesting data to a new slice before returning it:

```go
func CopyDigits(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    b = digitRegexp.Find(b)
    c := make([]byte, len(b))
    copy(c, b)
    return c
}
```

A more concise version of this function could be constructed by using append. This is left as an exercise for the reader.

## Nil slices

The zero value of a slice is `nil`. A nil slice has a length and capacity of 0 and has no underlying array. Operations like `len()` or a `range` loop, work on `nil` slices

+ Assigning to a `nil` slice is not possible and fails with a panic:

  ```go
  var s []int
  s[0] = 1234  // panic: runtime error: index out of range
  ```

+ Create a `nil` slice an empty slice:

  ```go
  var s []int   // s == nil
  s := []int{}  // s == []
  ```

To check whether a slice is empty, always use the `len()` function. 

```go
len(s) == 0
```

The other option, testing for `s == nil`, can fail because non-nil slices can also have zero length. 


<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/BnJyOUt5eCb">https://go.dev/play/p/BnJyOUt5eCb</a></b></figcaption>

```go
package main

import "fmt"

func main() {
	var s []int

    for i, v := range s {
        fmt.Println(i,v)
    }

	fmt.Println(s, len(s), cap(s))
	if s == nil {
		fmt.Println("nil!")
	}
}
```

## Range 

The range form of the `for` loop iterates over a slice or map. When ranging over a slice, two values are returned for each iteration. The first is the index, and the second is a copy of the element at that index. 

```go
package main

import "fmt"

var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}

func main() {
	for i, v := range pow {
		fmt.Printf("2**%d = %d\n", i, v)
	}
}
```

We can skip the index or value by assigning to `_`. 

```go
for i, _ := range pow
for _, value := range pow
```

If you only want the index, you can omit the second variable.

```go
for i := range pow
```

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/aHNwwoVbOcG">https://go.dev/play/p/aHNwwoVbOcG</a></b></figcaption>

```go
package main

import "fmt"

func main() {
	pow := make([]int, 10)
	for i := range pow {
		pow[i] = 1 << uint(i) // == 2**i
	}
	for _, value := range pow {
		fmt.Printf("%d\n", value)
	}
}
```

## Summary

+ A slice represents a section of an underlying array.
+ Slices starts with an initial length and can expand until the end of the array.
+ The `append()` function can expand a slice, but may allocate a new array under the hood.
+ We can create a slice using an slice literal or the function `make()`. Both functions will create an unnmaed array for the slice automatically. 


## Exercises

TODO