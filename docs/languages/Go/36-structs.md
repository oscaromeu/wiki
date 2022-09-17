---
id: structs
title: Structs
sidebar_label: Structs
sidebar_position: 36
hide_title: false
draft: false
---

## Introduction

A `struct` is an aggregate data type. It is composed of zero or more fields. A field is like a variable: it represents a value of a given type. The struct data type declaration consists of: 

+ the keyword `struct`
+ and a list of field declaration within curly braces.

The field declaration syntax is the same as for variable declarations; that is, a name followed by a type.

```go
type Planet struct {
    Name             string
    Mass             int64
    Diameter         int
    Gravity          float64
    RotationPeriod   time.Duration
    HasAtmosphere    bool
    HasMagneticField bool
    Satellites       []string
    next, previous   *Planet
}
```

Fields can have any type, except for the struct type itself. Doing so would trigger an infinite recursion when instantiating a variable of that type.

So if you need to include a field of the same type as the struct itself, use a pointer instead.

Using this type declaration, we now can instantiate variables of type Planet.

```go
var earth, jupiter Planet
```



## Initialization

As with all types in Go, structs are initialized to their zero value. The zero value of a struct type is a struct with all fields being set to their respective zero value.

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/YcWU_LWKP8S">https://go.dev/play/p/YcWU_LWKP8S</a></b></figcaption>

```go
package main

import (
	"fmt"
	"time"
)

type Planet struct {
	Name             string
	Mass             int64
	Diameter         int
	Gravity          float64
	RotationPeriod   time.Duration
	HasAtmosphere    bool
	HasMagneticField bool
	Satellites       []string
	next, previous   *Planet
}

func main() {
	var earth Planet
	fmt.Printf("Zero value of earth: %v\n", earth)
}
```

## Initialize structures with composite literals

Composite literals for initializing structures come in two different forms:

:::info
A _composite literal_ is a concise syntax to initialize any composite type with the values you want. Rather than declare an array and assign elements one by one, Go's composite literal syntax will declare and initialize an array in a single step.
:::

### Using field-value pairs

```go
mars := Planet{
    Name:           "Mars",
    Diameter:       6792,
    RotationPeriod: 24.7 * 60 * time.Minute,
    HasAtmosphere:  true,
    Satellites:     []string{"Phobos", "Deimos"},
    Mass:           642e15, // in millon metric tons (1t == 1000kg)
    previous:       &earth,
    next:           &jupiter, // Remember the final comma
}
```
Names and values are separated by a colon, and name-value pairs are separated by a comma and can be spread across multiple lines. (Just remember to set the final comma when using the multi-line syntax.) This form tolerates change and will continue to work correctly even if fields are added to the structure or if fields are reordered. Omitted fields take their default value.



### Specifying only value pairs

In this form the composite literal doesn't specify field names. Instead, a value must be provided for each field in the same order in which they're listed in the structure definition. This form works best for types that are stable and only have a few fields.

```go
    // Possible but unclear and error prone
    mars = Planet{
        "Mars",
        642e15,
        6792,
        3.7,
        24.7 * 60 * time.Minute,
        true,
        false,
        []string{"Phobos", "Deimos"},
        &earth,
        &jupiter,
    }
```

## Access

Struct fields can be accessed through the dot notation, just like package functions, types and variables are accessed. This applies to both reading from and writing to a struct field.

```go
mars.Gravity = 3.7
fmt.Println("Dot notation:", mars.Gravity)
```

Dot access even works if the variable is a pointer to a struct. Unlike in some other languages, no special syntax is required for dereferencing the pointer first.

```go
var pmars = &mars
fmt.Println("Dot notation with pointer:", pmars.Gravity)
// same effect: fmt.Println((*pmars).Gravity)
```

## Visibility

In the lecture on packages, we learned that package-level identifiers are exported if they start with an uppercase letter, and internal otherwise. The same principle applies to struct fields.

+ A struct field with a capital letter is exported and therefore accessible to users of the struct outside the package in which the struct is defined.

+ All field names that do not start with a capital letter denote an internal field. Internal fields can only be accessed from within the same package.

Taking our Planet struct as an example, all fields except for the fields previous and next start with an uppercase letter and hence are exported. Hence if the Planet struct was part of a library package, then any code that imports this package can access all fields except for the previous and next fields.

## Comparasion

Structs are comparable if each of their fields are comparable. Our Planet struct is not comparable because slices are not comparable.

Trying to compare two planets already fails at compile-time, thanks to static typing.

```go
fmt.Println("Is Mars the same as Jupiter?", mars == jupiter)
```

Output:

```go
invalid operation: mars == jupiter (struct containing []string cannot be compared)
```

Let’s create a comparable struct by generalizing our Planet struct to a struct representing an arbitrary celestial body.

A celestial body can be a star, a planet, a moon, an asteroid, etc., so some of the properties do not make sense and need to be removed. A moon, for example, has no satellites. So we can reduce the Planet struct to a more general struct named CelestialBody.

```go
type CelestialBody struct {
    Name           string
    Mass           int64
    Diameter       int64
    Gravity        float64
    RotationPeriod time.Duration
}
```

And now we can compare two celestial bodies. For example, let's create sun and moon with different names, and the comparison compiles now and returns `false`.

```go
var sun, moon CelestialBody
sun.Name = "Sun"
moon.Name = "Moon"
fmt.Println("Are sun and moon the same?", sun == moon)
```

## Passing to and returning from functions

Like any other type, structs can be passed to functions, and functions can return structures. Here, we pass a `Planet` structure to function `hasSatellites()` to determine if a planet has moons.

```go
func hasSatellites(p Planet) bool {
    return len(p.Satellites) > 0
}
// in main():
fmt.Println("Does Mars have satellites?", hasSatellites(mars))
```

Remember that Go has pass-by-value semantics. This means the struct inside the function is a copy of the struct that the caller passed to the function.

In this example, function `uppercase` takes a `Planet` value, changes the name to all-caps, and returns the value.

After calling uppercase, we can see that the returned struct has an all-caps name now, but the original mars struct still has the old name.

```go
func uppercase(p Planet) Planet {
    p.Name = strings.ToUpper(p.Name)
    return p
}
// in main():
fmt.Println(uppercase(mars).Name, "is uppercase of", mars.Name)
```

Note, by the way, that we can use dot access on a function result. For large structs, pass-by-value semantics may become inefficient. Here, pointers come in handy. Just keep in mind that the function then can have side effects.

The function rename intentionally makes use of this side effect by updating the name of the original struct variable through the pointer.

```go
func rename(p *Planet, n string) {
    p.Name = n
}
// in main():
rename(&mars, "Ἄρεως ἀστἡρ")
fmt.Println("An ancient Greek name for mars is", mars.Name)
```

After calling `rename()`, variable mars now has a different name.

## Shallow copies and shallow comparisons
As you know from slices and pointers, the pass-by-value semantics and the comparison operation work at the top level only. Pointers aren’t followed when data with pointers is copied or compared.

All this also applies to structs.

## Summary

+ Structs aggregate multiple values into a single entity.
+ A struct field is accessed through the dot notation.
+ Fields are exported if their name starts with an uppercase letter.
+ Exported fields are accessible from outside the package.
+ Internal fields are only accessible whithin the package.
+ Structs are comparable if all of their fields are comparable.
+ Structs can be passed to and returned by functions.
+ Use a pointer to a struct to make function calls more efficient, but be aware of possible side effects.


https://asankov.dev/blog/2022/01/29/different-ways-to-initialize-go-structs/