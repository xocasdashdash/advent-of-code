package main

import (
	"bufio"
	"fmt"
	"os"
)

func localizer(b []byte, minRange int, maxRange int, maxChar byte, minChar byte) int {

	fmt.Printf("Parsing %s\n", string(b))
	var letter byte
	for _, letter = range b {
		if letter == minChar {
			maxRange = maxRange - (maxRange+1-minRange)/2
		}
		if letter == maxChar {
			minRange = minRange + (maxRange+1-minRange)/2
		}

	}
	if letter == minChar {
		return minRange
	}
	return maxRange
}
func rowParser(b []byte) int {
	minRange := 0
	maxRange := 127

	return localizer(b, minRange, maxRange, 'B', 'F')
}
func columnParser(b []byte) int {
	return localizer(b, 0, 7, 'R', 'L')
}
func main() {
	f, _ := os.Open("input")
	s := bufio.NewScanner(f)
	highestID := 0
	for s.Scan() {
		b := s.Bytes()
		r := rowParser(b[0:7])
		c := columnParser(b[7:])
		candidate := r*8 + c
		fmt.Printf("R %d - %d, %d\n", r, c, candidate)
		if candidate > highestID {
			highestID = candidate
		}
	}

	fmt.Printf("Highest id %d", highestID)
}
