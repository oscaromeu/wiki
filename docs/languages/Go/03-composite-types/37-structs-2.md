---
id: structs-2
title: Structs embedding and anonymous fields
sidebar_label: Structs embedding and anonymous fields
sidebar_position: 37
hide_title: true
draft: false
---

## Structs embedding and anonymous fields

Structs can be composed by adding one struct as field to another struct. Let's do this with the `Planet` and `CelestialBody` types.

```go
type CelestialBody struct {
    Name           string
    Mass           int64
    Diameter       int64
    Gravity        float64
    RotationPeriod time.Duration
}
```

`Planet` has all these fields, too, so we can remove them and add a `CelestialBody` field instead. 

```go
type Planet struct {
    HeavenlyBody     CelestialBody
    HasAtmosphere    bool
    HasMagneticField bool
    Satellites       []string
    next, previous   *Planet
}
```

Now we can access fields like `Name` through the `HeavenlyBody` field. 

```go
func main() {
    var p Planet
    p.HeavenlyBody.Name = "Venus"
    p.HeavenlyBody.Mass = 4.87e15
    p.HeavenlyBody.Diameter = 12104
    fmt.Println("Accessing a struct field's field:", p.HeavenlyBody.Name)
}
```

But as you can see, this can quickly become quite verbose. Luckily, Go provides a shortcut named __struct embedding.__

A named `struct` type can be embedded into another `struct` type without using a field name. 

```go
type Planet struct {
    CelestialBody    // Anonymous field: No name, only a type
    HasAtmosphere    bool
    HasMagneticField bool
    Satellites       []string
    next, previous   *Planet
}
```

This kind of field is called an __anonymous field__. Now we can access the fields of the anonmyous `CelestialBody` field as if they were fields of the `Planet` struct. 

```go
type Planet struct {
    CelestialBody    // Anonymous field: No name, only a type
    HasAtmosphere    bool
    HasMagneticField bool
    Satellites       []string
    next, previous   *Planet
}
func main() {
    var p Planet
    p.Name = "Venus" 
    fmt.Println("Accessing an anonymous field's field:", p.Name)
}
```

Note that the `CelestialBody` field is not entirely anonymous; we still can refer to it through the type name.

```go
p.CelestialBody.Name = 8.9
fmt.Println("p2.CelestialBody.Name is", p.CelestialBody.Name)
```

It is entirely possible that a struct contains a field of the same name as the field of an embedded struct. In this case, the shortcut to the embedded struct’s field does not work anymore, and we have to use the full access path.

```go
type Planet struct {
    Name             string
    CelestialBody    // also contains a "Name" field
    HasAtmosphere    bool
    HasMagneticField bool
    Satellites       []string
    next, previous   *Planet
}
func main() {
    var p Planet
    p.CelestialBody.Name = "Mercury"
    p.Name = "Venus" // now refers to Planet.Name
    fmt.Println("p.Name:", p.Name)
    fmt.Println("p.CelestialBody.Name:", p.CelestialBody.Name)
}
```

One thing to be aware of is that the shortcut notation is not available for `struct` literals. 

```go
p := Planet{
    Name: "Mercury",
}
```

The reason is that in a struct literals, no evaluation of field access paths takes place. So in a struct literal, embedded structs are initialized just like normal field values. 

```go
p := Planet{
    CelestialBody: CelestialBody{
        Name: "Mercury",
        Diameter: 4879,
    },
    HasAtmosphere: true,
}
```

So far, struct embedding and anonymous fields may seem no more than a convenience for the developer, but later we will see that embedding opens new possibilities when combined with struct methods.

### Summary

+ Nested struct definitions can cause long access paths to a field.
+ Structs can be embedded into other structs as fields without a name.
+ Fields of the embedded struct can be accessed like fields of the enclosing struct, thus keeping the access path short.
+ The full access path is still available by using the type name of the embedded struct in place of the missing field name.
+ If a field name appears in both structs, the shortcut to the embedded struct’s field is not available.
+ Struct literals have no shortcut notation. Use the embedded struct’s literal representation to initialize an embedded struct’s field.