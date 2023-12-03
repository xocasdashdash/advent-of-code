package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//go:embed input
var input string

//go:embed testInput
var testInput string
var testMode = flag.Bool("test", false, "Set to run using the testInput")

func main() {
	start := time.Now()
	flag.Parse()
	if *testMode {
		input = testInput
	}
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	fmt.Println("File reading took", time.Since(start))
	parts, symbols := parseInput(trimmedInput)
	maxX := len(trimmedInput[0])
	maxY := len(trimmedInput)
	start = time.Now()
	pm := PartMap(maxX, maxY, parts)
	end := time.Now()
	fmt.Println("generating map took", end.Sub(start))
	start = time.Now()
	fmt.Println("part1", FindAndAddValidParts(symbols, pm))
	end = time.Now()
	fmt.Println("part 1 took", end.Sub(start))
	start = time.Now()
	fmt.Println("part2", Part2(symbols, pm))
	end = time.Now()
	fmt.Println("part 2 took", end.Sub(start))

}

func Part2(symbols []Symbol, pm map[Coord][]Part) int {
	gearRatio := 0
	for _, s := range symbols {
		foundNeighbors := make(map[string]Part, 0)
	findingNeighbors:
		for _, neighborCoord := range s.NeighborCoords() {
			if parts, ok := pm[neighborCoord]; ok {
				for _, p := range parts {
					foundNeighbors[p.PartId()] = p
					// If we found more than two neighbors already we can quit
					if len(foundNeighbors) > 2 {
						fmt.Println("More than two!")
						break findingNeighbors
					}
				}
			}
		}
		// fmt.Printf("Found %d neighbors for symbol %s at %s\n", len(foundNeighbors), s.Value, s.Location)
		if len(foundNeighbors) == 2 {
			symbolGearRatio := 1
			for _, p := range foundNeighbors {
				symbolGearRatio *= p.PartNumber
			}
			// fmt.Println("symbol", s.Value, "gear ratio", symbolGearRatio, foundNeighbors)
			gearRatio += symbolGearRatio
		}
	}
	return gearRatio
}

func FindAndAddValidParts(symbols []Symbol, pm map[Coord][]Part) int {
	validParts := make(map[string]Part, 0)
	for _, s := range symbols {
		for _, neighborCoord := range s.NeighborCoords() {
			if parts, ok := pm[neighborCoord]; ok {
				for _, p := range parts {
					validParts[p.PartId()] = p
				}
			}
		}
	}
	part1 := 0
	for _, p := range validParts {
		part1 += p.PartNumber
	}
	return part1
}

func PartMap(maxX, maxY int, parts []Part) map[Coord][]Part {
	result := make(map[Coord][]Part)
	for _, p := range parts {
		for i := p.start.X; i <= p.end.X; i++ {
			currentCoord := Coord{X: i, Y: p.start.Y}
			listOfParts, ok := result[currentCoord]
			if !ok {
				listOfParts = make([]Part, 0)
			}
			listOfParts = append(listOfParts, p)
			result[currentCoord] = listOfParts
		}
	}

	return result
}

func parseInput(input []string) ([]Part, []Symbol) {

	parts := make([]Part, 0, 10)
	symbols := make([]Symbol, 0, 10)
	for y, line := range input {

		splitLine := strings.Split(line, "")
		for x := 0; x < len(splitLine); x++ {

			start := x
			var end Coord
		findingParts:
			for x < len(splitLine) {
				runes := []rune(splitLine[x])
				if runes[0] == '.' {
					break findingParts
				} else if unicode.IsDigit(runes[0]) {
					end = Coord{X: x, Y: y}
					x++
				} else {
					// is a symbol
					symbols = append(symbols, Symbol{Value: splitLine[x], Location: Coord{X: x, Y: y}})
					break findingParts
				}
			}
			if start == x {
				continue
			}
			partNumber, _ := strconv.Atoi(strings.Join(splitLine[start:x], ""))
			parts = append(parts, Part{
				start:           Coord{X: start, Y: y},
				end:             end,
				PartNumber:      partNumber,
				SymbolNeighbors: make([]Symbol, 0),
			})
		}
	}
	return parts, symbols
}

type Coord struct {
	X int
	Y int
}

func (c Coord) String() string {
	return fmt.Sprintf("x:%d,y:%d", c.X, c.Y)
}

type Part struct {
	start           Coord
	end             Coord
	PartNumber      int
	SymbolNeighbors []Symbol
}

func (p Part) PartId() string {
	return fmt.Sprintf("%d-%s-%s", p.PartNumber, p.start, p.end)
}

type Symbol struct {
	Value    string
	Location Coord
}

func (s Symbol) NeighborCoords() []Coord {

	result := make([]Coord, 0)
	modifiers := []int{-1, 0, 1}

	for _, xmod := range modifiers {
		for _, ymod := range modifiers {
			if xmod == 0 && ymod == 0 {
				continue
			}
			result = append(result, Coord{X: s.Location.X + xmod, Y: s.Location.Y + ymod})
		}
	}
	return result
}
