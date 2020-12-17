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
	fmt.Printf("Active before: %d\n", len(points))
	fmt.Printf("********\n")
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

func PrintPoints(allPoints map[Point]bool) {
	minX := 9999999999
	maxX := 0
	minY := 9999999999
	maxY := 0
	minZ := 9999999999
	maxZ := 0
	for k := range allPoints {
		if k.X < minX {
			minX = k.X
		} else if k.X > maxX {
			maxX = k.X
		}
		if k.Y < minY {
			minY = k.Y
		} else if k.Y > maxY {
			maxY = k.Y
		}
		if k.Z < minZ {
			minZ = k.Z
		} else if k.Z > maxZ {
			maxZ = k.Z
		}
	}
	var sb strings.Builder
	for k := minZ; k <= maxZ; k++ {
		sb.WriteString(fmt.Sprintf("\n**z = %d **\n", k))
		for i := minX; i <= maxX; i++ {
			for j := minY; j <= maxY; j++ {
				p, _ := allPoints[Point{X: i, Y: j, Z: k}]
				if p {
					sb.WriteString("#")
				} else {
					sb.WriteString(".")
				}
			}
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("\n**z = %d **\n", k))
	}
	fmt.Printf(sb.String())
	//fmt.Printf("min x: %d, min y: %d, min z: %d\n", minX, minY, minZ)
	//fmt.Printf("max x: %d, max y: %d, max z: %d\n", maxX, maxY, maxZ)
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
		_, ok := allPoints[p]
		numberOfActiveNeighbors := len(allNeighbors) - totalInactiveNeighbors
		if numberOfActiveNeighbors == 3 && ok == false {
			//Flip from inactive to active
			result[p] = active
		} else if ok == true && (numberOfActiveNeighbors == 2 || numberOfActiveNeighbors == 3) {
			//do nothing
			result[p] = active
		}
	}

	return result
}
func neighbors(allPoints map[Point]bool, p Point, enableFourthDimension bool) ([]Point, int) {
	var result []Point
	var inactivePoints = 0
	for w := -1; w < 2; w++ {
		if !enableFourthDimension {
			w = 0
		}
		for k := -1; k < 2; k++ {
			for i := -1; i < 2; i++ {
				for j := -1; j < 2; j++ {
					checks++
					if i == 0 && j == 0 && k == 0 && w == 0 {
						continue
					}
					candidatePoint := Point{
						X: p.X + i,
						Y: p.Y + j,
						Z: p.Z + k,
						W: p.W + w,
					}
					result = append(result, candidatePoint)
					if _, ok := allPoints[candidatePoint]; !ok {
						inactivePoints++
					}
				}
			}
		}
		if !enableFourthDimension {
			break
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
