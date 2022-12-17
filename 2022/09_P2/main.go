package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"
	"time"
)

//go:embed input
var input string

//go:embed testInput
var testInput string

//go:embed testInput2
var testInput2 string

var testMode = flag.Bool("test", false, "Set to run using the testInput")
var part2 = flag.Bool("part2", false, "Set to run using the testInput2")

type Coord struct {
	X, Y int
}
type Direction Coord
type Instruction struct {
	Value int
	Dir   Direction
	Alias string
}

type Knot struct {
	Index   int
	Coord   Coord
	Visited map[Coord]struct{}
	Prev    *Knot
	Next    *Knot
}

func (i Instruction) String() string {
	return fmt.Sprintf("%s %d", i.Alias, i.Value)
}

func parseLines(s []string) []Instruction {
	result := make([]Instruction, 0)
	for _, l := range s {
		i := Instruction{}
		var dir string
		var value int

		toks, err := fmt.Sscanf(l, "%s %d", &dir, &value)
		if toks != 2 || err != nil {
			fmt.Printf("Tokens: %d\n", toks)
			panic(err)
		}
		i.Value = value
		i.Alias = dir
		switch dir {
		case "U":
			i.Dir = Direction{Y: 1, X: 0}
		case "D":
			i.Dir = Direction{Y: -1, X: 0}
		case "L":
			i.Dir = Direction{Y: 0, X: -1}
		case "R":
			i.Dir = Direction{Y: 0, X: 1}
		default:
			panic(fmt.Sprintf("bad value %s", dir))
		}
		result = append(result, i)
	}
	return result
}

func NextStep(head Coord, tail Coord) Coord {
	x := head.X - tail.X
	y := head.Y - tail.Y
	distance := Direction{X: x, Y: y}

	switch distance {
	case Direction{X: 2, Y: 0}:
		return Coord{X: tail.X + 1, Y: tail.Y}
	case Direction{X: -2, Y: 0}:
		return Coord{X: tail.X - 1, Y: tail.Y}
	case Direction{X: 0, Y: 2}:
		return Coord{X: tail.X, Y: tail.Y + 1}
	case Direction{X: 0, Y: -2}:
		return Coord{X: tail.X, Y: tail.Y - 1}

	// Move upper left diagonal
	case Direction{X: -1, Y: 2}:
		fallthrough
	case Direction{X: -2, Y: 1}:
		fallthrough
	case Direction{X: -2, Y: 2}:
		return Coord{X: tail.X - 1, Y: tail.Y + 1}

		// Move upper right diagonal
	case Direction{X: 2, Y: 2}:
		fallthrough
	case Direction{X: 2, Y: 1}:
		fallthrough
	case Direction{X: 1, Y: 2}:
		return Coord{X: tail.X + 1, Y: tail.Y + 1}

		// Move bottom right diagonal
	case Direction{X: 2, Y: -2}:
		fallthrough
	case Direction{X: 2, Y: -1}:
		fallthrough
	case Direction{X: 1, Y: -2}:
		return Coord{X: tail.X + 1, Y: tail.Y - 1}

		// Move bottom left diagonal
	case Direction{X: -2, Y: -2}:
		fallthrough
	// Move bottom left diagonal
	case Direction{X: -2, Y: -1}:
		fallthrough
	case Direction{X: -1, Y: -2}:
		return Coord{X: tail.X - 1, Y: tail.Y - 1}

	// Distance is 1 or less, keep as is.
	default:
		// if math.Abs(float64(distance.X)) > 2.0 || math.Abs(float64(distance.Y)) > 2.0 {
		// 	runtime.Breakpoint()
		// 	panic("something is off")
		// }
		return Coord{X: tail.X, Y: tail.Y}
	}

}
func UpdateCoord(coord Coord, d Direction, value int) Coord {
	return Coord{
		X: coord.X + d.X*value,
		Y: coord.Y + d.Y*value,
	}
}
func WalkTheRope(start Coord, instructions []Instruction, numberOfKnots int) map[Coord]struct{} {

	currentKnot := new(Knot)
	startKnot := &Knot{
		Index: 0,
		Coord: start,
		Prev:  nil,
		Next:  nil,
		Visited: map[Coord]struct{}{
			start: {},
		},
	}
	currentKnot = startKnot
	var nextKnot *Knot
	knots := make([]*Knot, numberOfKnots, numberOfKnots)
	for i := 0; i < numberOfKnots; i++ {
		nextKnot = &Knot{
			Index: currentKnot.Index + 1,
			Coord: start,
			Prev:  currentKnot,
			Visited: map[Coord]struct{}{
				start: {},
			},
		}
		currentKnot.Next = nextKnot
		currentKnot = nextKnot
		knots[i] = nextKnot
	}
	headKnot := startKnot
	tailKnot := currentKnot

	for _, ins := range instructions {
		for step := 0; step < ins.Value; step++ {
			headKnot.Coord = UpdateCoord(headKnot.Coord, ins.Dir, 1)
			nextKnot := headKnot.Next
			currentCoord := headKnot.Coord
			for nextKnot != nil {

				prevCoord := nextKnot.Coord
				nextKnot.Coord = NextStep(currentCoord, nextKnot.Coord)
				// If part of the chain does not move
				// we can quit faster
				if prevCoord == nextKnot.Coord {
					break
				}
				nextKnot.Visited[nextKnot.Coord] = struct{}{}
				currentCoord = nextKnot.Coord
				nextKnot = nextKnot.Next
			}
		}
		PrintMap(headKnot, ins)
		time.Sleep(10 * time.Millisecond)
	}
	return tailKnot.Visited
}
func PrintMap(head *Knot, instruction Instruction) {

	// Set some sane defaults
	minX, minY := -10, 10
	maxX, maxY := -10, 10
	current := head.Next
	for current != nil {

		if current.Coord.X < minX {
			minX = current.Coord.X
		} else if current.Coord.X > maxX {
			maxX = current.Coord.X
		}
		if current.Coord.Y < minY {
			minY = current.Coord.Y

		} else if current.Coord.Y > maxY {
			maxY = current.Coord.Y
		}
		current = current.Next
	}
	// minX = -13
	// maxX = 17
	// minY = -5
	// maxY = 16
	positions := make(map[Coord]int)
	current = head.Next
	for current != nil {
		positions[current.Coord] = current.Index
		current = current.Next
	}
	positions[head.Coord] = head.Index
	fmt.Printf("\n---%s %d---\n", instruction.Alias, instruction.Value)
	fmt.Printf("      ")
	for i := minX; i <= maxX; i++ {
		fmt.Printf("%+03d", i)
	}
	fmt.Printf("\n")
	for j := maxY; j >= minY; j-- {
		fmt.Printf("%+03d - ", j)
		for i := minX; i <= maxX; i++ {
			char := fmt.Sprintf("% 3s", "*")
			if index, ok := positions[Coord{X: i, Y: j}]; ok {
				char = fmt.Sprintf("% 3d", index)
			}
			if i == 0 && j == 0 {
				char = fmt.Sprintf("% 3s", "s")
			}
			fmt.Printf("%s", char)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n---------\n")
}
func main() {
	start := time.Now()
	flag.Parse()
	if *testMode {
		if *part2 {
			input = testInput2
		} else {
			input = testInput
		}
	}

	parsingT := time.Now()
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	ins := parseLines(trimmedInput)
	fmt.Println("Took(parsing)", time.Since(parsingT))
	part1T := time.Now()
	// p1 := WalkTheRope(Coord{0, 0}, ins, 1)
	// fmt.Println("Part1", len(p1))
	p2 := WalkTheRope(Coord{0, 0}, ins, 9)
	fmt.Println("Part2", len(p2))
	fmt.Println("Took(part1)", time.Since(part1T))
	fmt.Println("Took", time.Since(start))

}
