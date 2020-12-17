package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")
var iterations = flag.Int("n", 6, "Number of iterations to run")

const active = true
const inactive = false

type Point struct {
	X, Y, Z, W int
}

var checks = 0

var neighborDelta []Point

func neighborMatrix(enableFourthDimension bool) []Point {
	if len(neighborDelta) != 0 {
		return neighborDelta
	}
	wMin := -1
	wMax := 2
	if !enableFourthDimension {
		wMin = 0
		wMax = 1
	}
	for w := wMin; w < wMax; w++ {
		for k := -1; k < 2; k++ {
			for i := -1; i < 2; i++ {
				for j := -1; j < 2; j++ {
					if i == 0 && j == 0 && k == 0 && w == 0 {
						continue
					}
					neighborDelta = append(neighborDelta, Point{
						X: i,
						Y: j,
						Z: k,
						W: w,
					})
				}
			}
		}
	}
	return neighborDelta

}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	points := parse(trimmedInput)
	t := time.Now()
	for i := 0; i < *iterations; i++ {
		points = evolve(points, false)
		fmt.Printf("Active after iteration %d: %d\n", i+1, len(points))
	}
	fmt.Printf("Took %s|", time.Since(t))
	fmt.Printf("Checks for part 1: %d|", checks)
	fmt.Printf("Active points after %d: %d\n", *iterations, len(points))
	neighborDelta = make([]Point, 0, 80)
	checks = 0
	points = parse(trimmedInput)
	t = time.Now()
	for i := 0; i < *iterations; i++ {
		points = evolve(points, true)
		fmt.Printf("Active after iteration %d: %d\n", i+1, len(points))
	}

	fmt.Printf("Took %s|", time.Since(t))
	fmt.Printf("Checks for part 2: %d|", checks)
	fmt.Printf("Active points after %d: %d", *iterations, len(points))

}

func evolve(allPoints map[Point]bool, enableFourthDimension bool) map[Point]bool {

	result := make(map[Point]bool, len(allPoints))
	var candidateToSwitch []Point
	visitedPoints := make(map[Point]bool)
	//We solve the current universe
	for p := range allPoints {
		allNeighbors := neighbors(allPoints, p, visitedPoints, enableFourthDimension)
		for k := range allNeighbors {
			visitedPoints[allNeighbors[k]] = true
		}
		candidateToSwitch = append(candidateToSwitch, allNeighbors...)
	}
	for _, p := range candidateToSwitch {
		activeNeighbors := activeNeighbors(allPoints, p, enableFourthDimension)
		if activeNeighbors != 2 && activeNeighbors != 3 {
			continue
		}
		if allPoints[p] == true {
			result[p] = active
			//keep active
		} else if activeNeighbors == 3 {
			//Flip from inactive to active
			result[p] = active
		}
	}
	return result
}

func activeNeighbors(allPoints map[Point]bool, p Point, enableFourthDimension bool) int {
	var activePoints int
	for _, n := range neighborMatrix(enableFourthDimension) {
		candidatePoint := Point{
			X: p.X + n.X,
			Y: p.Y + n.Y,
			Z: p.Z + n.Z,
			W: p.W + n.W,
		}
		if _, ok := allPoints[candidatePoint]; ok {
			activePoints++
		}
	}
	return activePoints
}
func neighbors(allPoints map[Point]bool, p Point, visitedPoints map[Point]bool, enableFourthDimension bool) []Point {
	var result []Point
	checks++
	for _, n := range neighborMatrix(enableFourthDimension) {
		candidatePoint := Point{
			X: p.X + n.X,
			Y: p.Y + n.Y,
			Z: p.Z + n.Z,
			W: p.W + n.W,
		}
		if _, ok := visitedPoints[candidatePoint]; ok {
			continue
		}
		result = append(result, candidatePoint)
	}
	return result
}
func parse(s []string) map[Point]bool {
	result := make(map[Point]bool, len(s)*len(s))
	for x, row := range s {
		for y := range row {
			if row[y] == '#' {
				p := Point{
					X: x, Y: y, Z: 0,
				}
				result[p] = active
			}
		}
	}
	return result
}
