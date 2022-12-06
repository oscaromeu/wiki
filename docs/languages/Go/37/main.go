// You can edit this code!
// Click here and start typing.
package main

import "fmt"


type CelestialBody struct {
    Name           string
    Mass           int64
    Diameter       int64
    Gravity        float64
    RotationPeriod time.Duration
}

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