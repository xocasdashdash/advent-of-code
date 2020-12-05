package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func parser(b []byte, minRange int, maxRange int, maxChar byte, minChar byte) int {

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

	return parser(b, minRange, maxRange, 'B', 'F')
}
func columnParser(b []byte) int {
	return parser(b, 0, 7, 'R', 'L')
}
func main() {
	f, _ := os.Open("input")
	defer f.Close()
	s := bufio.NewScanner(f)
	seats := make([]int, 8*128, 8*128)
	firstRowLimit := 8
	lastRowLimit := 127 * 8
	for s.Scan() {
		b := s.Bytes()
		r := rowParser(b[0:7])
		c := columnParser(b[7:])
		candidate := r*8 + c
		if !(candidate >= lastRowLimit || candidate <= firstRowLimit) {
			seats[candidate] = candidate
		}
	}
	sort.Ints(seats)
	for k := 1; k < len(seats)-1; k++ {
		if seats[k] != 0 && (seats[k+1]-seats[k] != 1) {
			fmt.Printf("Previous seat %d, next seat %d.\n", seats[k], seats[k+1])
			fmt.Printf("My seat %d \n", seats[k]+1)
		}
	}
}
