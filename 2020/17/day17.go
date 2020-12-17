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

	result := make(map[Point]bool)
	var candidateToSwitch []Point
	for p := range allPoints {
		allNeighbors, _ := neighbors(allPoints, p, enableFourthDimension)
		candidateToSwitch = append(candidateToSwitch, allNeighbors...)
	}
	for _, p := range candidateToSwitch {
		allNeighbors, totalInactiveNeighbors := neighbors(allPoints, p, enableFourthDimension)
		if len(allNeighbors)-totalInactiveNeighbors != 2 && len(allNeighbors)-totalInactiveNeighbors != 3 {
			continue
		}
		numberOfActiveNeighbors := len(allNeighbors) - totalInactiveNeighbors
		if allPoints[p] == true {
			result[p] = active
			//keep active
		} else if numberOfActiveNeighbors == 3 {
			//Flip from inactive to active
			result[p] = active
		}
	}
	return result
}
func neighbors(allPoints map[Point]bool, p Point, enableFourthDimension bool) ([]Point, int) {
	var result []Point
	var inactivePoints = 0
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
					checks++
					candidatePoint := Point{
						X: p.X + i,
						Y: p.Y + j,
						Z: p.Z + k,
						W: p.W + w,
					}
					result = append(result, candidatePoint)
					if !allPoints[candidatePoint] {
						inactivePoints++
					}
				}
			}
		}
	}
	return result, inactivePoints
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
