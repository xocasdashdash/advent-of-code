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
	Coord int
	Prev  *Knot
	Next  *Knot
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

func DistanceGreaterThan(coord1 Coord, coord2 Coord, value int) bool {
	x := coord1.X - coord2.X
	y := coord1.Y - coord2.Y
	fmt.Println("X distance", x, "Y Distance", y)
	return (x > value || x < -1*value) || (y > value || y < -1*value)
}
func UpdateCoord(coord Coord, d Direction, value int) Coord {
	return Coord{
		X: coord.X + d.X*value,
		Y: coord.Y + d.Y*value,
	}
}
func WalkTheRopeP1(start Coord, instructions []Instruction, ropeLength int) ([]Coord, map[Coord]struct{}, map[Coord]struct{}) {
	headPositions := make([]Coord, 0)
	tailPos := start
	headPos := start
	headPositions = append(headPositions, tailPos)
	tailUniquePositions := make(map[Coord]struct{})
	headUniquePositions := make(map[Coord]struct{})
	headUniquePositions[start] = struct{}{}
	tailUniquePositions[start] = struct{}{}
	for _, ins := range instructions {
		for step := 0; step < ins.Value; step++ {
			headPos = UpdateCoord(headPos, ins.Dir, 1)
			headUniquePositions[headPos] = struct{}{}
			headPositions = append(headPositions, headPos)
			if DistanceGreaterThan(headPos, tailPos, ropeLength) {
				tailPos = headPositions[len(headPositions)-1-ropeLength]
				tailUniquePositions[tailPos] = struct{}{}
				fmt.Println("I", ins, "Pos", tailPos)
			}
		}
	}
	return headPositions, tailUniquePositions, headUniquePositions
}
func PrintMap(min int, max int, visited map[Coord]struct{}, char string) {

	theMap := make([][]string, max)
	for i := min; i < max; i++ {
		theMap[i] = make([]string, max)
		for j := min; j < max; j++ {
			c := Coord{j, i}
			if _, ok := visited[c]; ok {
				theMap[i][j] = char
			} else {
				theMap[i][j] = "*"
			}
		}
	}
	theMap[0][0] = "s"
	for i := max - 1; i >= min; i-- {
		for j := min; j < max; j++ {
			fmt.Printf(" %s ", theMap[i][j])
		}
		fmt.Printf("\n")
	}

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
	_, p1, _ := WalkTheRopeP1(Coord{0, 0}, ins, 1)
	fmt.Println("Part1", len(p1))
	_, p2, _ := WalkTheRopeP1(Coord{0, 0}, ins, 10)
	fmt.Println("Part2", len(p2))
	fmt.Println("Took(part1)", time.Since(part1T))
	fmt.Println("Took", time.Since(start))

}
