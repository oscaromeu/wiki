---
id: structs-3
title: Structs field tags and JSON
sidebar_label: Structs field tags and JSON
sidebar_position: 38
hide_title: true
draft: false
---

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