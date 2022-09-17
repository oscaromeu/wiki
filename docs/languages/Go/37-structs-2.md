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

Now we can access fields like `Mass` through the `HeavenlyBody` field. 

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
    fmt.Println("Accessing an anonymous field's field:", p.HeavenlyBody.Name)
}
```

Note that the `CelestialBody` field is not entirely anonymous; we still can refer to it through the type name.

```go
p.CelestialBody.Gravity = 8.9
fmt.Println("p2.CelestialBody.Gravity is", p.CelestialBody.Gravity)
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

## Struct field tags and JSON

Assume you have this `weatherData` struct

```go
type weatherData struct {
    LocationName  string  
    Weather       string
    Temperature   int  
    Celsius       bool 
    TempForecast []int 
}
```

and you want to generate JSON data from a `weatherData` variable. The package `encoding/json` from the standard library makes this easy as pie:

```go
weather := weatherData{
    LocationName: "Zzyzx",
    Weather:      "sunny",
    Temperature:  80,
    Celsius:      false,
    TempForecast: []int{ 27, 25, 28 },
}
data, err := json.MarshalIndent(data, &weather)
```

And we get a nicely formatted JSON output:

```json
{
  "LocationName": "Zzyzx",
  "Weather": "sunny",
  "Temperature": 80,
  "Celsius": false,
  "TempForecast": [
    27,
    25,
    28
  ]
}
```

However, there are a few problems here:

+ The JSON fields should have lowercase names
+ The field `LocationName` should have the name location in the JSON data
+ `TempForecast` should contain an underscore in JSON, but this is not in line with the Go style guide.
+ And finally, the "celsius" field is optional and thus can be omitted from the JSON data if it carries the default value.

So you don't want to rename the struct fields, but how can we then create the desired JSON format?
Luckily there is a solution that can't be easier: Field tags.

```go
type weatherData struct {
    LocationName string `json:"location"`
    Weather      string `json:"weather"`
    Temperature  int    `json:"temp"`
    Celsius      bool   `json:"celsius,omitempty"`
    TempForecast []int  `json:"temp_forecast"`
}
```

A field tag is simply a string literal that can be added to a struct field. These field tags contain metadata about the fields. In case of JSON, a field tag can specify the JSON name of the field and other properties; for example, whether to omit the field entirely when it is empty. 

We use raw strings enclosed in backticks here, because the strings contain double quotes, and this saves us from having to escape each double quote with a backslash.

The `json` package processes all field tags starting with "json:" when marshalling the data.

Now the output meets our requirements:

```json
{
  "location": "Zzyzx",
  "weather": "sunny",
  "temp": 80,
  "temp_forecast": [
    27,
    25,
    28
  ]
}
```

Field tags are also useful in other contexts; for example, when accessing the tables and fields of a SQL database.

You can even write your own field tag parser; however, this requires the use of Reflection techniques that are explained later, and therefore we will not go into details here.

### Summary

+ Struct fields can have tags.
+ A tag is a string literal that is appended to a field declaration.
+ A tag contains metadata about the tagged field.
+ Some packages of the standard library can parse these tags to get more information about how to process the fields.