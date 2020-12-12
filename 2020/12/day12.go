package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"testing"
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

type Coord struct {
	x, y int
}
type ShipPosition struct {
	p Coord
	//Waypoint is relative to the ship's Position
	Waypoint Coord
}

func (p ShipPosition) String() string {
	return fmt.Sprintf("X: %+.3d, Y: %+.3d, Waypoint: x:%d, y:%d", p.p.x, p.p.y, p.Waypoint.x, p.Waypoint.y)
}

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

func main() {
	flag.Parse()

	input, _ := ioutil.ReadFile(*inputFile)

	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	instructions := parseInstructions(trimmedInput)

	//Part1
	shipPos := ShipPosition{
		Coord{0, 0}, Coord{1, 0},
	}
	for _, i := range instructions {
		newPosition := EvolveP1(shipPos, i)
		//fmt.Printf("Position: %+v from %+v after applying %+v\n", newPosition, shipPos, i)
		shipPos = newPosition
	}
	fmt.Printf("Part1: %d\n", int(math.Abs(float64(shipPos.p.x))+math.Abs(float64(shipPos.p.y))))

	//Part2
	shipPos = ShipPosition{
		Coord{0, 0}, Coord{10, 1},
	}
	for _, i := range instructions {
		newPosition := EvolveP2(shipPos, i)
		//fmt.Printf("Position: %+v from %+v after applying %+v\n", newPosition, shipPos, i)
		shipPos = newPosition
	}

	fmt.Printf("Part2: %d", int(math.Abs(float64(shipPos.p.x))+math.Abs(float64(shipPos.p.y))))
}

func TestRotateWaypoint(t *testing.T) {
	type rotTest struct {
		o, r Coord
		rot  int
	}
	rt := []rotTest{
		{
			Coord{3, 0}, Coord{0, -3}, 90,
		},
		{
			Coord{1, 1}, Coord{1, -1}, 90,
		},
		{
			Coord{1, 0}, Coord{0, -1}, 90,
		},
		{
			Coord{10, 4}, Coord{10, 4}, 360,
		},
		{
			Coord{10, 4}, Coord{4, -10}, 90,
		},
	}
	for _, r := range rt {
		c := RotateWaypoint(r.o, Instruction{Direction: "R", Change: r.rot})
		if c != r.r {
			t.Errorf("Bad rotation %d . Got %+v, expected %+v from origin: %+v\n", r.rot, c, r.r, r.o)
			t.Fail()
		}
	}
}

func RotateWaypoint(o Coord, i Instruction) Coord {
	if i.Direction == "L" {
		i.Change *= -1
	}
	radiantChange := float64(i.Change) * math.Pi / 180
	cosChange := int(math.Cos(radiantChange))
	sinChange := int(math.Sin(radiantChange))
	newX := o.x*cosChange + o.y*sinChange
	newY := -1*o.x*sinChange + o.y*cosChange
	return Coord{
		x: newX,
		y: newY,
	}
}
func EvolveP2(ship ShipPosition, i Instruction) ShipPosition {

	switch i.Direction {
	case "L", "R":
		ship.Waypoint = RotateWaypoint(ship.Waypoint, i)
		//fmt.Printf("New waypoint: %+v after rotation\n", ship.Waypoint)
	case "F":
		ship.p = Coord{
			x: ship.p.x + i.Change*ship.Waypoint.x,
			y: ship.p.y + i.Change*ship.Waypoint.y,
		}
		//fmt.Printf("Moved ship to position: %+v\n", ship.p)
	default:
		ship.Waypoint.y = ship.Waypoint.y + yMultipliers[i.Direction]*i.Change
		ship.Waypoint.x = ship.Waypoint.x + xMultipliers[i.Direction]*i.Change
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
	}
	c.p.x = c.p.x + c.Waypoint.x*i.Change
	c.p.y = c.p.y + c.Waypoint.y*i.Change
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
