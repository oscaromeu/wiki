// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"time"
)


type CelestialBody struct {
    Name           string
    Mass           int64
    Diameter       int64
    Gravity        float64
    RotationPeriod time.Duration
}

type Planet struct {
    //Name             string
    CelestialBody    // also contains a "Name" field
    HasAtmosphere    bool
    HasMagneticField bool
    Satellites       []string
    next, previous   *Planet
}


func main() {
    p := Planet{
		//Name: "Mercury",
		// In struct literal, embedded structs are initialized just like normal field values
		CelestialBody: CelestialBody{
			Name: "Venus",
			Diameter: 4879,
		},
		HasAtmosphere: true,
	}

	fmt.Println( 
	//p.Name, 
	p.Name, p.CelestialBody.Diameter)
}