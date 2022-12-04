package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Range struct {
	Min int
	Max int
}

func (r Range) String() string {
	return fmt.Sprintf("Min: %d, Max: %d", r.Min, r.Max)
}

func (r Range) Contains(other Range) bool {

	if other.Min == r.Min {
		return true
	} else if other.Min < r.Min {
		return other.Max >= r.Max
	}
	return r.Max >= other.Max
}
func GenerateOverlap(a Range, b Range) *Range {
	/*
		A - 1 - - 4 -
		B - - 2 - - 5
		C - 1 - - 4 -

		A - - - 3 4 -
		B - - 2 - - 5
		C - - - 3 4 -

		A - 1 - - 4 -
		B - - 2 - - 5
		C - - 2 - 4 -

		A - - - - 4 5
		B 1 - 3 - - -
		C - - - - - -
	*/
	// This logic is simpler but doesn't generate the overlap
	// if a.Min <= b.Min && a.Max >= b.Max {
	// 	return &Range{}
	// } else if b.Min <= a.Min && b.Max >= a.Max {
	// 	return &Range{}
	// }

	if a.Min <= b.Min {
		min := b.Min
		if a.Max >= b.Min {
			max := b.Max
			if a.Max <= b.Max {
				max = a.Max
			}
			return &Range{Min: min, Max: max}
		}
	} else {
		min := a.Min
		max := a.Max
		if b.Max < a.Min {
			return nil
		} else if a.Max >= b.Max {
			max = b.Max
		}
		return &Range{Min: min, Max: max}
	}
	return nil
}

func parsePair(l string) Range {
	parts := strings.Split(l, "-")
	r := Range{}
	r.Min, _ = strconv.Atoi(parts[0])
	r.Max, _ = strconv.Atoi(parts[1])
	return r
}
func parseLine(l string) (Range, Range) {
	parts := strings.Split(l, ",")

	return parsePair(parts[0]), parsePair(parts[1])

}

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	part1 := 0
	part2 := 0
	for _, l := range trimmedInput {
		a, b := parseLine(l)
		if a.Contains(b) {
			part1 += 1
		}
		if c := GenerateOverlap(a, b); c != nil {
			part2 += 1
		}
	}
	fmt.Println("Part1", part1)
	fmt.Println("Part2", part2)

}
