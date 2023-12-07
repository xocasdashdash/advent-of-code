package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
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
	fmt.Println("Took", time.Since(start))

	sm := parseInput(trimmedInput)
	fmt.Println("seed map", sm)
	locations := make([]int, len(sm.Seeds))
	for seedIndex, s := range sm.Seeds {
		locations[seedIndex] = s
		for _, mapping := range sm.Mappings {
			for _, mappingLine := range mapping.mappingLines {
				if mappingLine.Matches(locations[seedIndex], locations[seedIndex]) {
					locations[seedIndex] = mappingLine.Destination(locations[seedIndex])
					break
				}
			}
		}
	}
	minIndex := 0
	for currentIndex, loc := range locations {
		if loc <= locations[minIndex] {
			minIndex = currentIndex
		}
	}
	fmt.Println("part1", locations[minIndex])
}

type Mapping struct {
	name         string
	mappingLines []MappingLine
}
type SeedMap struct {
	Seeds    []int
	Mappings []Mapping
}

func parseInput(input []string) SeedMap {
	sm := SeedMap{}
	linePointer := 0
	seedParts := strings.Split(input[linePointer], ":")

	sm.Seeds = make([]int, 0)
	for _, seed := range strings.Split(strings.TrimSpace(seedParts[1]), " ") {
		seedInt, _ := strconv.Atoi(seed)
		sm.Seeds = append(sm.Seeds, seedInt)
	}
	linePointer++
	sm.Mappings = make([]Mapping, 0)
	for linePointer < len(input) {
		if strings.TrimSpace(input[linePointer]) == "" {
			linePointer++
			continue
		}
		m := Mapping{}
		if strings.Contains(input[linePointer], "map") {
			m.name = strings.Split(input[linePointer], " ")[0]
		}
		// fmt.Println("name", m.name, "line", linePointer)
		m.mappingLines = make([]MappingLine, 0)
		linePointer++
		for linePointer < len(input) && strings.TrimSpace(input[linePointer]) != "" {
			parts := strings.Split(input[linePointer], " ")
			if len(parts) != 3 {
				panic("bad line" + input[linePointer])
			}
			mappingLine := MappingLine{}
			mappingLine.DestinationRangeStart, _ = strconv.Atoi(parts[0])
			mappingLine.SourceRangeStart, _ = strconv.Atoi(parts[1])
			mappingLine.Length, _ = strconv.Atoi(parts[2])
			m.mappingLines = append(m.mappingLines, mappingLine)
			linePointer++
			// fmt.Println("line pointer", linePointer, "len(input)", len(input))
		}
		sm.Mappings = append(sm.Mappings, m)
	}
	return sm
}

type MappingLine struct {
	DestinationRangeStart int
	SourceRangeStart      int
	Length                int
}

func (m MappingLine) Matches(startRange, endRange int) bool {
	diff := startRange + (startRange - endRange) - m.SourceRangeStart
	if diff < 0 {
		return false
	}
	if diff > m.Length {
		return false
	}
	return true
}
func (m MappingLine) Destination(source int) int {

	diff := source - m.SourceRangeStart

	return m.DestinationRangeStart + diff

}
