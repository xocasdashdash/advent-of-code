package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

type Instruction struct {
	Direction string
	Change    int
}

var positions = []string{"E", "S", "W", "N"}
var positionsMap = map[string]int{
	"E": 0,
	"S": 1,
	"W": 2,
	"N": 3,
}
var xMultipliers = map[string]int{
	"E": 1,
	"W": -1,
	"S": 0,
	"N": 0,
}
var yMultipliers = map[string]int{
	"E": 0,
	"W": 0,
	"S": -1,
	"N": 1,
}

type ShipPosition struct {
	p complex128
	//Waypoint is relative to the ship's Position
	Waypoint complex128
}

func (p ShipPosition) String() string {
	return fmt.Sprintf("P: %v, Waypoint: %v", p.p, p.Waypoint)
}

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

func main() {
	flag.Parse()

	input, _ := ioutil.ReadFile(*inputFile)

	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	instructions := parseInstructions(trimmedInput)

	//Part1
	shipPos := ShipPosition{
		0, 1,
	}
	for _, i := range instructions {
		newPosition := EvolveP1(shipPos, i)
		//fmt.Printf("Position: %+v from %+v after applying %+v\n", newPosition, shipPos, i)
		shipPos = newPosition
	}
	fmt.Printf("Part1: %+v, %d\n", shipPos, int(math.Abs(real(shipPos.p))+math.Abs(imag(shipPos.p))))

	//Part2
	shipPos = ShipPosition{
		0, 10 + 1i,
	}
	for _, i := range instructions {
		newPosition := EvolveP2(shipPos, i)
		//fmt.Printf("Position: %+v from %+v after applying %+v\n", newPosition, shipPos, i)
		shipPos = newPosition
	}

	fmt.Printf("Part2: %+v, %d\n", shipPos, int(math.Abs(real(shipPos.p))+math.Abs(imag(shipPos.p))))
}

func RotateWaypoint(o complex128, i Instruction) complex128 {
	var r complex128
	r = o
	if i.Direction == "R" {
		//If we're going clockwise we need to remap this to going counter clockwise
		i.Change = 360 - i.Change
	}
	for j := 0; j < i.Change/90; j++ {
		//We need to do math like this to avoid issues with floats
		r *= 1i
	}
	//fmt.Printf("S: %+v, R: %+v, Angle: %+v, D: %s\n", o, r, i.Change, i.Direction)
	return r
}
func EvolveP2(ship ShipPosition, i Instruction) ShipPosition {

	switch i.Direction {
	case "L", "R":
		ship.Waypoint = RotateWaypoint(ship.Waypoint, i)
		//fmt.Printf("New waypoint: %+v after rotation\n", ship.Waypoint)
	case "F":
		ship.p += complex(float64(i.Change), 0) * ship.Waypoint
	default:
		ship.Waypoint += complex(float64(xMultipliers[i.Direction]*i.Change), float64(yMultipliers[i.Direction]*i.Change))
		//fmt.Printf("New waypoint: %+v after moving it\n", ship.Waypoint)
	}
	return ship

}

//EvolveP1 Returns the map after applying the instruction
func EvolveP1(c ShipPosition, i Instruction) ShipPosition {
	switch i.Direction {
	case "L", "R":
		c.Waypoint = RotateWaypoint(c.Waypoint, i)
		return c
	case "F":
		c.p += c.Waypoint * complex(float64(i.Change), 0)
	default:
		c.p += complex(float64(xMultipliers[i.Direction]*i.Change), float64(yMultipliers[i.Direction]*i.Change))
	}

	//fmt.Printf("Change X: %d, Change y: %d\n", xMultiplier*i.Change, yMultiplier*i.Change)
	return c
}
func parseInstructions(directions []string) []Instruction {
	var ret []Instruction
	for d := range directions {
		dir := directions[d]
		var r Instruction
		r.Direction = fmt.Sprintf("%c", dir[0])
		change, _ := strconv.Atoi(fmt.Sprintf("%s", dir[1:]))
		r.Change = change
		ret = append(ret, r)
	}
	return ret
}
