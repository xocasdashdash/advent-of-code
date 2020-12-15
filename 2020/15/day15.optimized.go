package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")
var t = flag.Int("t", 2020, "Number of turns to play")
var debug = flag.Bool("d", false, "In order to enable debugging")

func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")

	fmt.Printf("Input: %s\n", trimmedInput)
	numbers := parse(trimmedInput)
	if *debug {
		fmt.Printf("Input numbers: %+v\n", numbers)
	}
	//This is a map of the number and an array
	//of length 2. The first position holds the last
	//time the number was said out loud, and the second one holds
	//the number of times it was said.
	seen := make([]int, *t, *t)
	lastNumberSpoken := 0
	for i := range numbers {
		turn := i + 1
		lastNumberSpoken = numbers[i]
		seen[lastNumberSpoken] = turn
	}
	numberToSay := 0
	lastNumber := 0
	var j int
	lastNumber = 0
	//var ok bool
	var n int
	maxNumber := 0
	t0 := time.Now()
	for j = len(numbers); j < *t; j++ {
		lastNumber = numberToSay

		if n = seen[lastNumber]; n != 0 {
			numberToSay = j + 1 - n
		} else {
			numberToSay = 0
		}
		seen[lastNumber] = j + 1
		if lastNumber > maxNumber {
			maxNumber = lastNumber
		}
	}
	fmt.Printf("Took %s\n", time.Since(t0))
	fmt.Printf("Max number: %d", maxNumber)
	fmt.Printf("Number at %d: %d \n", *t, lastNumber)

}
func printNumbers(n []int) {
	for k, v := range n {
		fmt.Printf("Turn: %d => %d\n", k+1, v)
	}
}
func parse(s []string) []int {
	var result []int
	for i := range s {
		for _, candidateNumber := range strings.Split(s[i], ",") {
			n, _ := strconv.Atoi(candidateNumber)
			result = append(result, n)
		}
	}
	return result

}
